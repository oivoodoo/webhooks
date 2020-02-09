package master

import (
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
	v1 "gitlab.com/oivoodoo/webhooks/pkg/master/v1"
)

type master struct {
}

func (master) Receive(webhook *webhooks.Webhook) error {
	return v1.Receive(webhook)
}

func Create() *master {
	return &master{}
}
