package webhooks

type Webhook struct {
	Body     []byte `json:"-",bson:"-"`
	Checksum string
	// TODO: timestamp int
	// History should clear using timestamp by time window
	// Slave should search for differences using timestamp by time window
}

func New(body []byte) *Webhook {
	return &Webhook{
		Body: body,
	}
}
