package pkg

import (
	"gitlab.com/oivoodoo/webhooks/pkg/cfg"
)

type AppStr struct {
	Config *cfg.Configuration

	DB     interface{}
	Router interface{}

	DieChans []chan bool
}

var App *AppStr

func New() *AppStr {
	return &AppStr{
		Config:   cfg.Create(),
		DieChans: []chan bool{},
	}
}

func (a *AppStr) SubscribeDie() chan bool {
	ch := make(chan bool)
	a.DieChans = append(a.DieChans, ch)
	return ch
}

func (a AppStr) Die() {
	for _, ch := range a.DieChans {
		ch <- true
	}
}
