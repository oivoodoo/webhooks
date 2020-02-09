package cfg

import (
	"github.com/kelseyhightower/envconfig"
)

// Configuration contains the basic settings
// that should be loaded from env variables.
type Configuration struct {
	Debug bool

	PORT string `default:5000`

	// PostgreSQL connecting settings
	MASTER_HOST string `default:"127.0.0.1"`
	SLAVE_HOST  string `default:"127.0.0.1"`

	BEGIN_TIME_SECONDS_WINDOW int `default:240` // 4 minutes
	END_TIME_SECONDS_WINDOW   int `default:120` // 2 minutes

	EXPIRE_TIME_SECONDS_AGO int `default:60` // 1 minutes
}

func Create() *Configuration {
	var config *Configuration

	err := envconfig.Process("server", config)

	if err != nil {
		panic(err.Error())
	}

	return config
}
