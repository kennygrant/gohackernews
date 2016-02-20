# Golang News
A responsive news website for links about the Go programming language, written in Go. You can see a this website in action at this link:

[http://golangnews.com](http://golangnews.com)

## Gettting Started

The app requires postgresql just now to bootstrap locally (not Mysql). So make sure you have psql installed. The bootstrap process will create a database and settings for you, but you'll need to promote the first user to admin in order to use the site locally.

Go get this app:

    go get -u github.com/kennygrant/gohackernews

Then to build and run the server locally, as you'd expect:

    go run server.go

or get the fragmenta command line tool (for things like migrations, deploy etc) and run it with that:

    go get -u github.com/fragmenta/fragmenta
    fragmenta server



## App Structure

#### server.go
This is the entry point main() for the application. It includes packages within src and starts a server. 

#### The src folder
This holds the website assets, actions and views - the meat of the app. 

#### The src/app folder
This contains general app files, resources like pages or users should go in a separate pkg.

#### The src/users folder
This contains files related to users on the website.

#### The src/stories folder
This contains files related to stories on the website.

#### The src/comments folder
This contains files related to comments on the website.

#### The src/lib folder
lib is used to store utility packages which can be used by several parts of the app.

#### The src/lib/templates folder
Templates for generating new resources are stored in here and used by fragmenta generate to generate a new resource package, containing assets, code and views for the resource.  
