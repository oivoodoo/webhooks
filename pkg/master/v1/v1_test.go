package v1

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/oivoodoo/webhooks/pkg"
	"gitlab.com/oivoodoo/webhooks/pkg/batcher"
	"gitlab.com/oivoodoo/webhooks/pkg/db"
	"gitlab.com/oivoodoo/webhooks/pkg/db/webhooks"
	"gitlab.com/oivoodoo/webhooks/pkg/router"
)

type mock_master struct {
	*batcher.Batcher
}

func (m *mock_master) Receive(webhook *webhooks.Webhook) error {
	return Receive(m.Batcher, webhook)
}

var server *echo.Echo
var repo *webhooks.MockWebhookRepo

func setup() {
	repo = &webhooks.MockWebhookRepo{}

	pkg.App = pkg.New()
	pkg.App.DB = &db.DB{
		WebhookRepo: repo,
	}

	master := &mock_master{
		Batcher: batcher.New(),
	}
	server = echo.New()
	server.POST("/webhooks", router.CreateWebhook(master))
}

const msg = `{"key":"value"}`

func TestCreateWebooksAndSync(t *testing.T) {
	setup()

	assert.Len(t, repo.Data, 0)

	request := func(msg string) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/webhooks", bytes.NewBufferString(msg))
		req.Header.Set("Content-Type", "application/json")
		server.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, router.OK_RESP, w.Body.String())
	}

	request(strings.Replace(msg, "value", "value1", -1))
	assert.Len(t, repo.Data, 0)
	request(strings.Replace(msg, "value", "value2", -1))
	assert.Len(t, repo.Data, 0)
	request(strings.Replace(msg, "value", "value3", -1))
	assert.Len(t, repo.Data, 0)
	// should skip it because of not unique
	request(strings.Replace(msg, "value", "value3", -1))
	assert.Len(t, repo.Data, 0)

	time.Sleep(time.Duration(pkg.App.Config.SYNC_DATABASE_SECONDS_WINDOW+1) * time.Second)

	assert.Len(t, repo.Data, 3)

	pkg.App.Die()
}
