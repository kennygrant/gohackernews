// A simple news site inspired by hacker news, written in Go
package main

import (
	"fmt"

	// Import the fragmenta command line tool so that if you go-get this repo, you get the tool
	_ "github.com/fragmenta/fragmenta"

	"github.com/fragmenta/server"

	"github.com/kennygrant/gohackernews/src/app"
)

func main() {

	// Setup server
	server, err := server.New()
	if err != nil {
		fmt.Printf("Error creating server %s", err)
		return
	}

	app.Setup(server)

	// Inform user of server setup
	server.Logf("#info Starting server in %s mode on port %d", server.Mode(), server.Port())

	// Start the server
	err = server.Start()
	if err != nil {
		server.Fatalf("Error starting server %s", err)
	}

}
