package webhooks

type MockWebhookRepo struct {
	Data []*Webhook
}

func (m *MockWebhookRepo) BatchInsert(webhooks []*Webhook) error {
	println("[repo.mock] begin inserting", len(m.Data))
	for _, webhook := range webhooks {
		println(string(webhook.Body))
		m.Data = append(m.Data, webhook)
	}
	println("[repo.mock] done inserting", len(m.Data))
	return nil
}
