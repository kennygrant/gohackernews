# Golang Hacker News
A responsive version of hacker news written in golang. 

## Gettting Started

Go get this app:

    go get -u github.com/kennygrant/gohackernews

Then to build and run the server locally, get the fragmenta command line tool:

    go get -u github.com/fragmenta/fragmenta

Change dir and run the server to bootstrap (requires a Postgresql install at present):

    cd $GOPATH/src/github.com/kennygrant/gohackernews
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
lib is used to store utility packages which can be used by several parts of the app. Some examples of libraries are included, but unused in this example application. 

#### The src/lib/templates folder
Templates for generating new resources are stored in here and used by fragmenta generate to generate a new resource package, containing assets, code and views.  
