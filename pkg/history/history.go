package history

import (
	"sync"

	"gitlab.com/oivoodoo/webhooks/pkg"
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
)

type History struct {
	locker       *sync.Mutex
	Webhooks     []*webhooks.Webhook
	WebhooksChan chan *webhooks.Webhook
	die          chan bool
}

func New(webhooksChan chan *webhooks.Webhook) *History {
	h := &History{
		locker:       &sync.Mutex{},
		Webhooks:     []*webhooks.Webhook{},
		WebhooksChan: webhooksChan,
		die:          pkg.App.SubscribeDie(),
	}

	h.StartWorker()

	return h
}

func (h *History) Lock() {
	h.locker.Lock()
}

func (h *History) Unlock() {
	h.locker.Unlock()
}

func (h *History) Clear() {
	h.Webhooks = []*webhooks.Webhook{}
}

func (h *History) StartWorker() {
	go func() {
		for {
			select {
			case <-h.die:
				return
			case webhook := <-h.WebhooksChan:
				h.locker.Lock()
				h.Webhooks = append(h.Webhooks, webhook)
				h.locker.Unlock()
			}
		}
	}()
}
