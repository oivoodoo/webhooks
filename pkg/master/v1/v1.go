package v1

import (
	"gitlab.com/oivoodoo/webhooks/pkg"
	"gitlab.com/oivoodoo/webhooks/pkg/db"
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
)

func storage() *db.DB {
	return pkg.App.DB.(*db.DB)
}

func Receive(webhook *webhooks.Webhook) error {
	storage()

	return nil
}
