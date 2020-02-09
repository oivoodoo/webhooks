package master

type master struct {
}

func (master) Receive(body []byte) error {
	return nil
}

func Create() *master {
	return &master{}
}
