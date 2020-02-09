package batcher

import (
	"sync"
	"time"

	"github.com/ztrue/tracerr"
	"gitlab.com/oivoodoo/webhooks/pkg"
	"gitlab.com/oivoodoo/webhooks/pkg/db"
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
)

type Batcher struct {
	collector []*webhooks.Webhook
	locker    sync.Mutex
}

func New() *Batcher {
	batcher := &Batcher{}
	batcher.start()
	return batcher
}

func (b *Batcher) sync() error {
	b.locker.Lock()
	{
		err := pkg.App.DB.(*db.DB).WebhookRepo.BatchInsert(b.collector)
		if err != nil {
			return tracerr.Wrap(err)
		}
		b.collector = []*webhooks.Webhook{}
	}
	b.locker.Unlock()

	return nil
}

func (b *Batcher) Push(webhook *webhooks.Webhook) error {
	b.locker.Lock()
	{
		b.collector = append(b.collector, webhook)
	}
	b.locker.Unlock()

	return nil
}

func (b *Batcher) start() {
	go func() {
		for {
			select {
			case <-pkg.App.Die:
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
