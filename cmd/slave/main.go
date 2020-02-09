package main

import (
	"gitlab.com/oivoodoo/webhooks/pkg"
	"gitlab.com/oivoodoo/webhooks/pkg/db"
	"gitlab.com/oivoodoo/webhooks/pkg/router"
	"gitlab.com/oivoodoo/webhooks/pkg/slave"
)

func main() {
	app := start()
	<-app.Done
}

func start() *pkg.AppStr {
	app := pkg.New()
	pkg.App = app

	app.DB = db.Create()
	app.Router = router.Create(slave.Create())

	return app
}
