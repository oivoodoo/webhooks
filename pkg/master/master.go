package master

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gojektech/heimdall/httpclient"
	"gitlab.com/oivoodoo/webhooks/pkg"
	"gitlab.com/oivoodoo/webhooks/pkg/batcher"
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
	"gitlab.com/oivoodoo/webhooks/pkg/history"
	v1 "gitlab.com/oivoodoo/webhooks/pkg/master/v1"
)

type master struct {
	*batcher.Batcher
	*history.History

	die chan bool
}

func (m *master) Receive(webhook *webhooks.Webhook) error {
	return v1.Receive(m.Batcher, webhook)
}

func Create() *master {
	b := batcher.New()
	h := history.New(b.WebhooksChan)

	m := &master{
		Batcher: b,
		History: h,
		die:     pkg.App.SubscribeDie(),
	}

	m.slaveSync()

	return m
}

func (m *master) slaveSync() {
	go func() {
		for {
			select {
			case <-m.die:
				return
			case <-time.After(time.Duration(pkg.App.Config.SYNC_SLAVE_SECONDS_WINDOW) * time.Second):
				m.History.Lock()

				client := httpclient.NewClient(
					httpclient.WithHTTPTimeout(5*time.Second),
					httpclient.WithRetryCount(2),
				)

				headers := http.Header{}
				headers.Set("Content-Type", "application/json")

				body, err := json.Marshal(m.History.Webhooks)
				if err != nil {
					println(err.Error())
					m.History.Unlock()
					continue
				}

				resp, err := client.Post(
					pkg.App.Config.SLAVE_HOST+"/v1/sync",
					bytes.NewBuffer(body),
					headers)
				if err != nil {
					println(err.Error())
					m.History.Unlock()
					continue
				}
				if resp.StatusCode != 201 {
					println("error on verifying webhooks in slave")
					m.History.Unlock()
					continue
				}

				m.History.Clear()
				m.History.Unlock()
			}
		}
	}()
}
