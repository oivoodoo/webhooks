package v1

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
)

func Receive(batcher *batcher.Batcher, webhook *webhooks.Webhook) error {
	return batcher.Push(webhook)
}

// Sync search for difference by checksums
// and return collection of webhook bodies
// in response for master node.
// master should push bodies to batcher to
// sync with database.
func Sync(history *history.History, checksums []string) error {
	history.Lock()
	defer history.Unlock()

	ws := search(history.Webhooks, checksums)

	if len(ws) == 0 {
		return nil
	}

	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(5*time.Second),
		httpclient.WithRetryCount(2),
	)

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	body, err := json.Marshal(ws)
	if err != nil {
		println(err.Error())
		return nil
	}

	resp, err := client.Post(
		pkg.App.Config.MASTER_HOST+"/v1/webhooks",
		bytes.NewBuffer(body),
		headers)
	if err != nil {
		println(err.Error())
		return nil
	}
	if resp.StatusCode != 200 {
		println("error on pushing webhooks to master")
		return nil
	}

	return nil
}

func search(ws []*webhooks.Webhook, checksums []string) []*webhooks.Webhook {
	// TODO: compare checksums and in case of having more items return it and push to
	// master
	var result []*webhooks.Webhook

	return result
}
