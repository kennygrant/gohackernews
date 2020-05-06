package app

import (
	"github.com/fragmenta/mux"
	"github.com/fragmenta/mux/middleware/gzip"

	//	"github.com/fragmenta/mux/middleware/secure"
	"github.com/fragmenta/server/log"

	// Resource Actions
	commentactions "github.com/kennygrant/gohackernews/src/comments/actions"
	"github.com/kennygrant/gohackernews/src/lib/session"
	storyactions "github.com/kennygrant/gohackernews/src/stories/actions"
	stripeactions "github.com/kennygrant/gohackernews/src/stripe/actions"
	useractions "github.com/kennygrant/gohackernews/src/users/actions"
)

// SetupRoutes creates a new router and adds the routes for this app to it.
func SetupRoutes() *mux.Mux {

	router := mux.New()
	mux.SetDefault(router)

	// Add the home page route
	router.Get("/", storyactions.HandleHome)

	// Add a route to handle static files
	router.Get("/favicon.ico", fileHandler)
	router.Get("/files/{path:.*}", fileHandler)
	router.Get("/assets/{path:.*}", fileHandler)

	// Resource Routes

	router.Get("/stripe/pay", stripeactions.HandleShowPay)
	router.Get("/stripe/thanks", stripeactions.HandleShowPayThanks)
	router.Get("/stripe/cancel", stripeactions.HandleShowPayCancel)

	// Add story routes
	router.Get("/go-jobs", storyactions.HandleJobs)
	router.Get("/index{format:(.xml)?}", storyactions.HandleIndex)
	router.Get("/stories/create", storyactions.HandleCreateShow)
	router.Post("/stories/create", storyactions.HandleCreate)
	router.Get("/stories/code{format:(.xml)?}", storyactions.HandleListCode)
	router.Get("/stories/upvoted{format:(.xml)?}", storyactions.HandleListUpvoted)
	router.Get("/stories/{id:[0-9]+}/update", storyactions.HandleUpdateShow)
	router.Post("/stories/{id:[0-9]+}/update", storyactions.HandleUpdate)
	router.Post("/stories/{id:[0-9]+}/destroy", storyactions.HandleDestroy)
	router.Post("/stories/{id:[0-9]+}/upvote", storyactions.HandleUpvote)
	router.Post("/stories/{id:[0-9]+}/downvote", storyactions.HandleDownvote)
	router.Post("/stories/{id:[0-9]+}/flag", storyactions.HandleFlag)
	router.Get("/stories/{id:[0-9]+}", storyactions.HandleShow)
	router.Get("/stories{format:(.xml)?}", storyactions.HandleIndex)
	router.Get("/sitemap.xml", storyactions.HandleSiteMap)

	router.Get("/comments", commentactions.HandleIndex)
	router.Get("/comments/create", commentactions.HandleCreateShow)
	router.Post("/comments/create", commentactions.HandleCreate)
	router.Get("/comments/{id:[0-9]+}/update", commentactions.HandleUpdateShow)
	router.Post("/comments/{id:[0-9]+}/update", commentactions.HandleUpdate)
	router.Post("/comments/{id:[0-9]+}/destroy", commentactions.HandleDestroy)
	router.Post("/comments/{id:[0-9]+}/upvote", commentactions.HandleUpvote)
	router.Post("/comments/{id:[0-9]+}/downvote", commentactions.HandleDownvote)
	router.Post("/comments/{id:[0-9]+}/flag", commentactions.HandleFlag)
	router.Get("/comments/{id:[0-9]+}", commentactions.HandleShow)

	router.Get("/users", useractions.HandleIndex)
	router.Get("/users/create", useractions.HandleCreateShow)
	router.Post("/users/create", useractions.HandleCreate)
	router.Get("/users/{id:[0-9]+}/update", useractions.HandleUpdateShow)
	router.Post("/users/{id:[0-9]+}/update", useractions.HandleUpdate)
	router.Post("/users/{id:[0-9]+}/destroy", useractions.HandleDestroy)
	router.Get("/users/{id:[0-9]+}", useractions.HandleShow)
	router.Get("/u/{name:.*}", useractions.HandleShowName)
	router.Get("/users/login", useractions.HandleLoginShow)
	router.Post("/users/login", useractions.HandleLogin)
	router.Post("/users/logout", useractions.HandleLogout)

	// Set the default file handler
	router.FileHandler = fileHandler
	router.ErrorHandler = errHandler

	// Add middleware
	router.AddMiddleware(log.Middleware)
	router.AddMiddleware(session.Middleware)
	router.AddMiddleware(gzip.Middleware)
	//	router.AddMiddleware(secure.Middleware)

	return router
}
