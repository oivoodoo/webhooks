package slave

import (
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
)

type slave struct {
}

func (slave) Receive(webhook *webhooks.Webhook) error {
	return nil
}

func Create() *slave {
	return &slave{}
}
