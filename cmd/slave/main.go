package main

import (
	"gitlab.com/oivoodoo/webhooks/pkg"
	"gitlab.com/oivoodoo/webhooks/pkg/cfg"
	"gitlab.com/oivoodoo/webhooks/pkg/db"
	"gitlab.com/oivoodoo/webhooks/pkg/router"
	"gitlab.com/oivoodoo/webhooks/pkg/slave"
)

func main() {
	app := start()
	<-app.Die
}

func start() *pkg.AppStr {
	slave := slave.Create()

	app := pkg.New()
	pkg.App = app

	app.Config = cfg.Create()
	app.DB = db.Create()
	app.Router = router.Create(slave)

	return app
}
