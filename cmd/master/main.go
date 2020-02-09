package master

import (
	"gitlab.com/oivoodoo/webhooks/pkg"
)

func main() {
	app := pkg.StartMaster()
	<-app.Die
}
