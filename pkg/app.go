package pkg

import (
	"gitlab.com/oivoodoo/webhooks/pkg/cfg"
)

type AppStr struct {
	Config *cfg.Configuration

	DB     interface{}
	Router interface{}

	Die chan bool
}

var App *AppStr

func New() *AppStr {
	return &AppStr{
		Config: cfg.Create(),
		Die:    make(chan bool),
	}
}
