package app

import (
	"os"
	"time"

	// psql driver - we only use a psql db at the moment
	_ "github.com/lib/pq"

	"github.com/fragmenta/query"
	"github.com/fragmenta/server/config"
	"github.com/fragmenta/server/log"
)

// SetupDatabase sets up the db with query given our server config.
func SetupDatabase() {
	defer log.Time(time.Now(), log.V{"msg": "Finished opening database", "db": config.Get("db"), "user": config.Get("db_user")})

	options := map[string]string{
		"adapter":  config.Get("db_adapter"),
		"user":     config.Get("db_user"),
		"password": config.Get("db_pass"),
		"db":       config.Get("db"),
	}

	// Optionally Support remote databases
	if len(config.Get("db_host")) > 0 {
		options["host"] = config.Get("db_host")
	}
	if len(config.Get("db_port")) > 0 {
		options["port"] = config.Get("db_port")
	}
	if len(config.Get("db_params")) > 0 {
		options["params"] = config.Get("db_params")
	}

	// Ask query to open the database
	err := query.OpenDatabase(options)

	if err != nil {
		log.Fatal(log.V{"msg": "unable to read database", "db": config.Get("db"), "error": err})
		os.Exit(1)
	}

}
