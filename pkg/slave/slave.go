package master

import (
	"gitlab.com/oivoodoo/webhooks/pkg/batcher"
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
	v1 "gitlab.com/oivoodoo/webhooks/pkg/master/v1"
)

type slave struct {
	*batcher.Batcher
}

func (s *slave) Receive(webhook *webhooks.Webhook) error {
	return v1.Receive(s.Batcher, webhook)
}

func Create() *slave {
	return &slave{
		Batcher: batcher.New(),
	}
}
