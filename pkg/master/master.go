package master

import (
	"gitlab.com/oivoodoo/webhooks/pkg/batcher"
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
	v1 "gitlab.com/oivoodoo/webhooks/pkg/master/v1"
)

type master struct {
	*batcher.Batcher
}

func (m *master) Receive(webhook *webhooks.Webhook) error {
	return v1.Receive(m.Batcher, webhook)
}

func Create() *master {
	return &master{
		Batcher: batcher.New(),
	}
}
