package webhooks

type Webhook struct {
	Body     []byte
	Checksum string
}

func New(body []byte) *Webhook {
	return &Webhook{
		Body: body,
	}
}
