package main

import (
	"gitlab.com/oivoodoo/webhooks/pkg"
	"gitlab.com/oivoodoo/webhooks/pkg/db"
	"gitlab.com/oivoodoo/webhooks/pkg/master"
	"gitlab.com/oivoodoo/webhooks/pkg/router"
)

func main() {
	app := start()
	<-app.Done
}

func start() *pkg.AppStr {
	app := pkg.New()
	pkg.App = app

	app.DB = db.Create()
	app.Router = router.Create(master.Create())

	return app
}
