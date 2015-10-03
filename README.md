# Golang Hacker News
A responsive version of hacker news written in golang. 

## Gettting Started

Get the fragmenta tool

    go get github.com/fragmenta/fragmenta

Go get this app:

    go get -u github.com/kennygrant/gohackernews

to run the server locally:

    cd to/app    

Run the server locally to bootstrap:

    fragmenta server


TODO: Some folders are not in version control -  (bin,secrets,log,db/backup) - add these with .keep files.

Performance test on blitz.io : https://www.blitz.io/report/2b7c5290b6f9436c3c50a090e1559576 

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
