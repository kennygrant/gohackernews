package app

import (
	"github.com/fragmenta/router"
	"github.com/kennygrant/gohackernews/src/comments/actions"
	"github.com/kennygrant/gohackernews/src/stories/actions"
	"github.com/kennygrant/gohackernews/src/users/actions"
)

// Define routes for this app
func setupRoutes(r *router.Router) {

	// Set the default file handler
	r.FileHandler = fileHandler
	r.ErrorHandler = errHandler

	// Add the home page route
	r.Add("/", storyactions.HandleHome)
	r.Add("/index{format:(.xml)?}", storyactions.HandleIndex)
	r.Add("/stories/create", storyactions.HandleCreateShow)
	r.Add("/stories/create", storyactions.HandleCreate).Post()
	r.Add("/stories/code{format:(.xml)?}", storyactions.HandleCode)
	r.Add("/stories/upvoted{format:(.xml)?}", storyactions.HandleUpvoted)
	r.Add("/stories/{id:[0-9]+}/update", storyactions.HandleUpdateShow)
	r.Add("/stories/{id:[0-9]+}/update", storyactions.HandleUpdate).Post()
	r.Add("/stories/{id:[0-9]+}/destroy", storyactions.HandleDestroy).Post()
	r.Add("/stories/{id:[0-9]+}/upvote", storyactions.HandleUpvote).Post()
	r.Add("/stories/{id:[0-9]+}/downvote", storyactions.HandleDownvote).Post()
	r.Add("/stories/{id:[0-9]+}/flag", storyactions.HandleFlag).Post()
	r.Add("/stories/{id:[0-9]+}", storyactions.HandleShow)
	r.Add("/stories{format:(.xml)?}", storyactions.HandleIndex)

	r.Add("/comments", commentactions.HandleIndex)
	r.Add("/comments/create", commentactions.HandleCreateShow)
	r.Add("/comments/create", commentactions.HandleCreate).Post()
	r.Add("/comments/{id:[0-9]+}/update", commentactions.HandleUpdateShow)
	r.Add("/comments/{id:[0-9]+}/update", commentactions.HandleUpdate).Post()
	r.Add("/comments/{id:[0-9]+}/destroy", commentactions.HandleDestroy).Post()
	r.Add("/comments/{id:[0-9]+}/upvote", commentactions.HandleUpvote).Post()
	r.Add("/comments/{id:[0-9]+}/downvote", commentactions.HandleDownvote).Post()
	r.Add("/comments/{id:[0-9]+}/flag", commentactions.HandleFlag).Post()
	r.Add("/comments/{id:[0-9]+}", commentactions.HandleShow)

	r.Add("/users", useractions.HandleIndex)
	r.Add("/users/create", useractions.HandleCreateShow)
	r.Add("/users/create", useractions.HandleCreate).Post()
	r.Add("/users/{id:[0-9]+}/update", useractions.HandleUpdateShow)
	r.Add("/users/{id:[0-9]+}/update", useractions.HandleUpdate).Post()
	r.Add("/users/{id:[0-9]+}/destroy", useractions.HandleDestroy).Post()
	r.Add("/users/{id:[0-9]+}", useractions.HandleShow)
	r.Add("/users/login", useractions.HandleLoginShow)
	r.Add("/users/login", useractions.HandleLogin).Post()
	r.Add("/users/logout", useractions.HandleLogout).Post()

	// Add a files route to handle static images under files
	// - nginx deals with this in production - perhaps only do this in dev?
	r.Add("/files/{path:.*}", fileHandler)
	r.Add("/favicon.ico", fileHandler)

}
