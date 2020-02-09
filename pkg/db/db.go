package db

import (
	"gitlab.com/oivoodoo/webhooks/pkg/db/dynamo"
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
)

type DB struct {
	WebhookRepo webhooks.Repo
}

func Create() *DB {
	conn := dynamo.Connect()

	db := &DB{
		WebhookRepo: &webhooks.DynamoRepo{Conn: conn},
	}

	return db
}
