package v1

import (
	"gitlab.com/oivoodoo/webhooks/pkg/batcher"
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
)

func Receive(batcher *batcher.Batcher, webhook *webhooks.Webhook) error {
	return batcher.Push(webhook)
}
