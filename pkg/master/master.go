package master

import (
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
)

type master struct {
}

func (master) Receive(webhook *webhooks.Webhook) error {
	return nil
}

func Create() *master {
	return &master{}
}
