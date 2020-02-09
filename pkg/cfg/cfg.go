package cfg

import (
	"github.com/kelseyhightower/envconfig"
)

// Configuration contains the basic settings
// that should be loaded from env variables.
type Configuration struct {
	Debug bool

	// Example:
	//   export SERVER_PORT=5000
	PORT string `default:"5001"`

	// Example:
	//	 export SERVER_MASTER_HOST=master.probackup.io
	MASTER_HOST string `default:"127.0.0.1"`
	SLAVE_HOST  string `default:"127.0.0.1"`

	BEGIN_TIME_SECONDS_WINDOW    int `default:"240"` // 4 minutes
	END_TIME_SECONDS_WINDOW      int `default:"120"` // 2 minutes
	EXPIRE_TIME_SECONDS_AGO      int `default:"60"`  // 1 minute
	SYNC_DATABASE_SECONDS_WINDOW int `default:"5"`   // each 5 seconds sync with database collector using batch inserts`

	// Amazon AWS
	AWS_REGION            string `envconfig:"AWS_REGION" default:"us-east-1"`
	AWS_ACCESS_KEY_ID     string `envconfig:"AWS_ACCESS_KEY_ID" required:"true"`
	AWS_SECRET_ACCESS_KEY string `envconfig:"AWS_SECRET_ACCESS_KEY" required:"true"`
}

func Create() *Configuration {
	println("[cfg] begin loading configuration...")

	var config Configuration

	err := envconfig.Process("server", &config)
	if err != nil {
		panic(err.Error())
	}

	println("[cfg] PORT=", config.PORT)
	println("[cfg] MASTER_HOST=", config.MASTER_HOST)
	println("[cfg] SLAVE_HOST=", config.SLAVE_HOST)
	println("[cfg] AWS_ACCESS_KEY_ID=", config.AWS_ACCESS_KEY_ID)

	println("[cfg] done loading configuration...")

	return &config
}
