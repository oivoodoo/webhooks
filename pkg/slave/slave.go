package slave

type slave struct {
}

func (slave) Receive(body []byte) error {
	return nil
}

func Create() *slave {
	return &slave{}
}
