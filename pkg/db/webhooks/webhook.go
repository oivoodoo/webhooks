package webhooks

type Webhook struct {
	Body []byte
}

func New(body []byte) *Webhook {
	return &Webhook{
		Body: body,
	}
}
