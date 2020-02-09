package main

import (
	"gitlab.com/oivoodoo/webhooks/pkg"
	"gitlab.com/oivoodoo/webhooks/pkg/cfg"
	"gitlab.com/oivoodoo/webhooks/pkg/db"
	"gitlab.com/oivoodoo/webhooks/pkg/master"
	"gitlab.com/oivoodoo/webhooks/pkg/router"
)

func main() {
	app := start()
	<-app.Die
}

func start() *pkg.AppStr {
	master := master.Create()

	app := pkg.New()
	pkg.App = app

	app.Config = cfg.Create()
	app.DB = db.Create()
	app.Router = router.Create(master)

	return app
}
