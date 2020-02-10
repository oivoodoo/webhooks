package slave

import (
	"gitlab.com/oivoodoo/webhooks/pkg/batcher"
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
	"gitlab.com/oivoodoo/webhooks/pkg/history"
	v1 "gitlab.com/oivoodoo/webhooks/pkg/slave/v1"
)

type Slave struct {
	*batcher.Batcher
	*history.History

	die chan bool
}

func (s *Slave) Receive(webhook *webhooks.Webhook) error {
	return v1.Receive(s.Batcher, webhook)
}

func (s *Slave) Sync(history *history.History, checksums []string) error {
	return v1.Sync(history, checksums)
}

func Create() *Slave {
	b := batcher.New()
	h := history.New(b.WebhooksChan)

	return &Slave{
		Batcher: b,
		History: h,
		die:     make(chan bool),
	}
}
