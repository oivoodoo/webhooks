package batcher

import (
	"crypto/sha256"
	"sync"
	"time"

	"github.com/ztrue/tracerr"
	"gitlab.com/oivoodoo/webhooks/pkg"
	"gitlab.com/oivoodoo/webhooks/pkg/db"
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
)

type Batcher struct {
	collector []*webhooks.Webhook
	locker    *sync.Mutex
	die       chan bool
}

func New() *Batcher {
	batcher := &Batcher{
		locker:    &sync.Mutex{},
		collector: []*webhooks.Webhook{},
		die:       pkg.App.SubscribeDie(),
	}
	batcher.start()
	return batcher
}

func (b *Batcher) sync() error {
	b.locker.Lock()
	{
		println("[batcher] begin inserting to database", len(b.collector))
		repo := pkg.App.DB.(*db.DB).WebhookRepo

		err := repo.BatchInsert(b.collector)
		if err != nil {
			return tracerr.Wrap(err)
		}
		b.collector = []*webhooks.Webhook{}
		println("[batcher] done inserting to database collector: ", len(b.collector), " repo:", len(repo.(*webhooks.MockWebhookRepo).Data))
	}
	b.locker.Unlock()

	return nil
}

func (b *Batcher) Push(webhook *webhooks.Webhook) error {
	b.locker.Lock()
	{
		if b.unique(webhook) {
			b.collector = append(b.collector, webhook)
		}
	}
	b.locker.Unlock()

	return nil
}

func (b *Batcher) unique(webhook *webhooks.Webhook) bool {
	bs := sha256.Sum256(webhook.Body)
	webhook.Checksum = string(bs[:])

	for _, w := range b.collector {
		if w.Checksum == "" {
			bs = sha256.Sum256(w.Body)
			w.Checksum = string(bs[:])
		}

		if w.Checksum == webhook.Checksum {
			return false
		}
	}

	return true
}

func (b *Batcher) start() {
	go func() {
		for {
			select {
			case <-b.die:
				println("[batcher] begin die because of app exit and do sync() before to exit")
				if err := b.sync(); err != nil {
					println(err.Error())
				}
				println("[batcher] done die because of app exit and do sync() before to exit")
				break
			case <-time.After(time.Duration(pkg.App.Config.SYNC_DATABASE_SECONDS_WINDOW) * time.Second):
				if err := b.sync(); err != nil {
					// TODO: add errors channel to output it
					println(err.Error())
				}
			}
		}
	}()
}
