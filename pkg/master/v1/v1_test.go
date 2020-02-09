package v1

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
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
	return Receive(m, webhook)
}

var server *echo.Echo

func setup() {
	pkg.App = pkg.New()
	pkg.App.DB = db.Create()

	server = echo.New()
	server.POST("/webhooks", router.CreateWebhook(&mock_master{}))
}

const msg = `{"key":"value"}`

func TestCreateWebhook(t *testing.T) {
	setup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/webhooks", bytes.NewBufferString(msg))
	req.Header.Set("Content-Type", "application/json")
	server.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, router.OK_RESP, w.Body.String())
}
