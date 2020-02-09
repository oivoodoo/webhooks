package webhooks

type MockWebhookRepo struct {
	Data []*Webhook
}

func (m *MockWebhookRepo) BatchInsert(webhooks []*Webhook) error {
	for _, webhook := range webhooks {
		m.Data = append(m.Data, webhook)
	}
	return nil
}
