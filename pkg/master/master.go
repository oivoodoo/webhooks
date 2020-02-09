package master

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gojektech/heimdall/httpclient"
	"gitlab.com/oivoodoo/webhooks/pkg"
	"gitlab.com/oivoodoo/webhooks/pkg/batcher"
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
	v1 "gitlab.com/oivoodoo/webhooks/pkg/master/v1"
)

type master struct {
	*batcher.Batcher

	locker    *sync.Mutex
	checksums []string

	die chan bool
}

func (m *master) Receive(webhook *webhooks.Webhook) error {
	return v1.Receive(m.Batcher, webhook)
}

func Create() *master {
	m := &master{
		Batcher:   batcher.New(),
		locker:    &sync.Mutex{},
		checksums: []string{},
		die:       pkg.App.SubscribeDie(),
	}

	m.checksumCollector()
	m.slaveSync()

	return m
}

func (m *master) checksumCollector() {
	go func() {
		for {
			select {
			case checksum := <-m.Batcher.ChecksumChan:
				m.locker.Lock()
				m.checksums = append(m.checksums, checksum)
				m.locker.Unlock()
			}
		}
	}()
}

func (m *master) slaveSync() {
	go func() {
		for {
			select {
			case <-m.die:
				return
			case <-time.After(time.Duration(pkg.App.Config.SYNC_SLAVE_SECONDS_WINDOW) * time.Second):
				m.locker.Lock()

				client := httpclient.NewClient(
					httpclient.WithHTTPTimeout(5*time.Second),
					httpclient.WithRetryCount(2),
				)

				headers := http.Header{}
				headers.Set("Content-Type", "application/json")

				body := strings.Join(m.checksums, ",")
				resp, err := client.Post(
					pkg.App.Config.SLAVE_HOST+"/v1/sync",
					bytes.NewBufferString(body),
					headers)
				if err != nil {
					println(err.Error())
					m.locker.Unlock()
					continue
				}

				_, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					println(err.Error())
					m.locker.Unlock()
					continue
				}

				// TODO: read of received webhooks in case slave has more messages
				// - m.Batcher.Push(for each new message)

				m.checksums = []string{}
				m.locker.Unlock()
			}
		}
	}()
}
