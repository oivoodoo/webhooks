package router

import (
	"io/ioutil"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ztrue/tracerr"
	"gitlab.com/oivoodoo/webhooks/pkg"
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
)

// master / slave should have the same implementation
// for handling messages
//
// master should have separate background worker to make sync up
// with slave. in case if slave has more data than master it would require
// to push the difference to master.
// slave should receive acknowledge that the messages where delivered
// otherwise we could have the same situation when we lost the messages
//
// the same idea used to use in rabbit mq with acknowledge pattern.
type ServerHandler interface {
	Receive(webhook Webhook) error
}

type Router struct {
	handler ServerHandler
}

func Create(handler ServerHandler) *Router {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("/v1")
	{
		v1.POST("/webhooks", createWebhook(handler))
	}

	// Start server
	go func() {
		config := pkg.App.Config

		e.Logger.Fatal(e.Start(":" + config.PORT))
	}()

	return &Router{
		handler: handler,
	}
}

func createWebhook(handler ServerHandler) error {
	return func(e *echo.Context) {
		req := e.Request()

		if b, err := ioutil.ReadAll(req.Body); err != nil {
			return tracerr.Wrap(err)
		} else {
			webhook := webhooks.New(b)

			if err := handler.Receive(webhook); err != nil {
				return tracerr.Wrap(err)
			}
		}

		return nil
	}
}
