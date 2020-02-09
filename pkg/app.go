package pkg

import (
	"gitlab.com/oivoodoo/webhooks/pkg/cfg"
	"gitlab.com/oivoodoo/webhooks/pkg/db"
	"gitlab.com/oivoodoo/webhooks/pkg/master"
	"gitlab.com/oivoodoo/webhooks/pkg/router"
	"gitlab.com/oivoodoo/webhooks/pkg/slave"
)

type app struct {
	Config *cfg.Configuration
	DB     *db.DB

	router *router.Router

	Die chan bool
}

var App *app

func StartMaster() *app {
	app := &app{}

	master := master.Create()

	app.Config = cfg.Create()
	app.DB = db.Create()
	app.router = router.Create(master)

	App = app

	return App
}

func StartSlave() *app {
	app := &app{}

	slave := slave.Create()

	app.Config = cfg.Create()
	app.DB = db.Create()
	app.router = router.Create(slave)

	App = app

	return App
}
