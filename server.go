package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/acme/autocert"

	"github.com/fragmenta/server"
	"github.com/fragmenta/server/config"

	"github.com/kennygrant/gohackernews/src/app"
)

// Main entrypoint for the server which performs bootstrap, setup
// then runs the server. Most setup is delegated to the src/app pkg.
func main() {

	// Bootstrap if required (no config file found).
	if app.RequiresBootStrap() {
		err := app.Bootstrap()
		if err != nil {
			fmt.Printf("Error bootstrapping server %s\n", err)
			return
		}
	}

	// Setup our server
	server, err := SetupServer()
	if err != nil {
		fmt.Printf("server: error setting up %s\n", err)
		return
	}

	// Inform user of server setup
	server.Logf("#info Starting server in %s mode on port %d", server.Mode(), server.Port())

	// In production, server
	if server.Production() {

		// If in production, serve over tls with autocerts from let's encrypt
		err = startTLSAutocert(server)
		if err != nil {
			server.Fatalf("Error starting server %s", err)
		}

	} else {
		// In development just serve with http on local port
		err = server.Start()
		if err != nil {
			server.Fatalf("Error starting server %s", err)
		}
	}

}

// SetupServer creates a new server, and delegates setup to the app pkg.
func SetupServer() (*server.Server, error) {

	// Setup server
	s, err := server.New()
	if err != nil {
		return nil, err
	}

	// Load the appropriate config
	c := config.New()
	err = c.Load("secrets/fragmenta.json")
	if err != nil {
		return nil, err
	}
	config.Current = c

	// Check environment variable to see if we are in production mode
	if os.Getenv("FRAG_ENV") == "production" {
		config.Current.Mode = config.ModeProduction
	}

	// Call the app to perform additional setup
	app.Setup()

	return s, nil
}

// startTLSAutocert starts an https server on the given port
// by requesting certs from an ACME provider.
// The server must be on a public IP which matches the
// DNS for the domains.
func startTLSAutocert(server *server.Server) error {
	autocertDomains := strings.Split(server.Config("autocert_domains"), " ")
	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Email:      server.Config("autocert_email"),            // Email for problems with certs
		HostPolicy: autocert.HostWhitelist(autocertDomains...), // Domains to request certs for
		Cache:      autocert.DirCache("secrets"),               // Cache certs in secrets folder
	}
	// Handle all :80 traffic using autocert to allow http-01 challenge responses
	go func() {
		http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	}()

	s := server.ConfiguredTLSServer(certManager)
	return s.ListenAndServeTLS("", "")
}
