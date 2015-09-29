# fragmenta-app
A minimal website built with fragmenta, with only a home page and minimal styling.

## Gettting Started
To create a copy of this app, run:

fragmenta new $GOPATH/src/my/app/name app

Then cd to your new app and run migrations:

fragmenta migrate

Then run the server:

fragmenta



## App Structure

#### server.go
This is the main entrypoint for the application. The structure of other parts of the application is dictated by what you need from it. 

#### The src folder
This is a suggested structure for an application, the structure used is entirely up to you, if you prefer you don't have to use a src folder. 


#### The src/app folder
This contains general app files, resources like pages or users should go in a separate pkg.


#### The src/lib folder

lib is used to store utility packages which can be used by several parts of the app. Some examples of libraries are included, but unused in this example application. 

#### The src/lib/templates folder

Templates for generating new resources are stored in here and used by fragmenta generate to generate a new resource package, containing assets, code and views.  