// A simple news site inspired by hacker news, written in Go
package main

import (
	"fmt"

	"github.com/fragmenta/server"

	"github.com/kennygrant/gohackernews/src/app"
)

func main() {

	// If we have no config, bootstrap first by generating config/migrations
	if app.RequiresBootStrap() {
		err := app.Bootstrap()
		if err != nil {
			fmt.Printf("Error bootstrapping server %s\n", err)
			return
		}
	}

	// Setup server
	server, err := server.New()
	if err != nil {
		fmt.Printf("Error creating server %s", err)
		return
	}

	app.Setup(server)

	// Inform user of server setup
	server.Logf("#info Starting server in %s mode on port %d", server.Mode(), server.Port())

	// In production, server
	if server.Production() {

		// Redirect all port 80 traffic to our canonical url
		server.StartRedirectAll(80, "https://golangnews.com")

		// Start the server using cert and key locally held
		// later change to autocert
		err = server.StartTLS()
		if err != nil {
			server.Fatalf("Error starting server %s", err)
		}

		/*
			// If in production, serve over tls with autocerts from let's encrypt
			err = server.StartTLSAutocert(server.Config("mail_from"), server.Config("autocert_domains"))
			if err != nil {
				server.Fatalf("Error starting server %s", err)
			}
		*/

	} else {
		// In development just serve with http on local port
		err = server.Start()
		if err != nil {
			server.Fatalf("Error starting server %s", err)
		}
	}

}
