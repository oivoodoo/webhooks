package router

import (
	"encoding/json"
	"io/ioutil"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ztrue/tracerr"
	"gitlab.com/oivoodoo/webhooks/pkg"
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
	"gitlab.com/oivoodoo/webhooks/pkg/slave"
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
	Receive(webhook *webhooks.Webhook) error
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
		v1.POST("/webhooks", CreateWebhook(handler))

		if s, ok := handler.(*slave.Slave); ok {
			v1.POST("/sync", SyncSlave(s))
		}
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

const OK_RESP = "RECEIVED"
const BAD_RESP = "BAD-RECEIVED"

func CreateWebhook(handler ServerHandler) func(echo.Context) error {
	return func(ctx echo.Context) error {
		req := ctx.Request()

		if b, err := ioutil.ReadAll(req.Body); err != nil {
			return tracerr.Wrap(err)
		} else {
			webhook := webhooks.New(b)

			if err := handler.Receive(webhook); err != nil {
				return tracerr.Wrap(err)
			}
		}

		return ctx.String(200, OK_RESP)
	}
}

func SyncSlave(handler ServerHandler) func(echo.Context) error {
	return func(ctx echo.Context) error {
		req := ctx.Request()

		if b, err := ioutil.ReadAll(req.Body); err != nil {
			return tracerr.Wrap(err)
		} else {
			type reqData struct {
				Checksum string
			}

			var data []reqData

			if err := json.Unmarshal(b, &data); err != nil {
				return ctx.String(400, BAD_RESP)
			}

			var checksums []string
			for _, str := range data {
				checksums = append(checksums, str.Checksum)
			}

			slave := handler.(*slave.Slave)
			if err := slave.Sync(slave.History, checksums); err != nil {
				return tracerr.Wrap(err)
			}
		}

		return ctx.String(200, OK_RESP)
	}
}
