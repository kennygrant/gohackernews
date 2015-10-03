package app

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/fragmenta/assets"
	"github.com/fragmenta/query"
	"github.com/fragmenta/router"
	"github.com/fragmenta/server"
	"github.com/fragmenta/server/log"
	"github.com/fragmenta/view"
	"github.com/fragmenta/view/helpers"

	"github.com/kennygrant/gohackernews/src/lib/authorise"
)

// appAssets holds a reference to our assets for use in asset setup
var appAssets *assets.Collection

// Setup sets up our application
func Setup(server *server.Server) {

	// Setup log
	server.Logger = log.New(server.Config("log"), server.Production())

	// Set up our assets
	setupAssets(server)

	// Setup our view templates
	setupView(server)

	// Setup our database
	setupDatabase(server)

	// Routing
	router, err := router.New(server.Logger, server)
	if err != nil {
		server.Fatalf("Error creating router %s", err)
	}

	// Setup our authentication and authorisation
	authorise.Setup(server)

	// Add a prefilter to store the current user on the context, so that we only fetch it once
	// We use this below in Resource, and also in views to determine current user attributes
	router.AddFilter(authorise.CurrentUserFilter)

	// Setup our router and handlers
	setupRoutes(router)

}

// Compile or copy in our assets from src into the public assets folder, for use by the app
func setupAssets(server *server.Server) {
	defer server.Timef("#info Finished loading assets in %s", time.Now())

	// Compilation of assets is done on deploy
	// We just load them here
	assetsCompiled := server.ConfigBool("assets_compiled")
	appAssets = assets.New(assetsCompiled)

	// Load asset details from json file on each run
	err := appAssets.Load()
	if err != nil {
		// Compile assets for the first time
		server.Logf("#info Compiling assets")
		err := appAssets.Compile("src", "public")
		if err != nil {
			server.Fatalf("#error compiling assets %s", err)
		}
	}

	// Set up helpers which are aware of fingerprinted assets
	// These behave differently depending on the compile flag above
	// when compile is set to no, they use precompiled assets
	// otherwise they serve all files in a group separately
	view.Helpers["style"] = appAssets.StyleLink
	view.Helpers["script"] = appAssets.ScriptLink

}

func setupView(server *server.Server) {
	defer server.Timef("#info Finished loading templates in %s", time.Now())

	// A very limited translation - would prefer to use editable.js instead and offer proper editing
	// with styles, and proper sanitization...
	view.Helpers["markup"] = markup
	view.Helpers["timeago"] = timeago

	view.Production = server.Production()
	err := view.LoadTemplates()
	if err != nil {
		server.Fatalf("Error reading templates %s", err)
	}

}

func markup(s string) template.HTML {
	// Nasty find/replace
	s = strings.Replace(s, "\n", "</p><p>", -1)

	return helpers.Sanitize(s)
}

func timeago(d time.Time) string {

	duration := time.Since(d)
	hours := duration / time.Hour

	switch {
	case duration < time.Minute:
		return fmt.Sprintf("%d seconds ago", duration/time.Second)
	case duration < time.Hour:
		return fmt.Sprintf("%d minutes ago", duration/time.Minute)
	case duration < time.Hour*24:
		unit := "hour"
		if hours > 1 {
			unit = "hours"
		}
		return fmt.Sprintf("%d %s ago", hours, unit)
	default:
		unit := "day"
		if hours > 48 {
			unit = "days"
		}
		return fmt.Sprintf("%d %s ago", hours/24, unit)
	}

}

// Setup db - at present query pkg manages this...
func setupDatabase(server *server.Server) {
	defer server.Timef("#info Finished opening in %s database %s for user %s", time.Now(), server.Config("db"), server.Config("db_user"))

	config := server.Configuration()
	options := map[string]string{
		"adapter":  config["db_adapter"],
		"user":     config["db_user"],
		"password": config["db_pass"],
		"db":       config["db"],
	}

	// Ask query to open the database
	err := query.OpenDatabase(options)

	if err != nil {
		server.Fatalf("Error reading database %s", err)
	}

}
