# Golang Gin OAuth2 Example

> ### Golang/Gin/Oauth2 codebase containing examples

# How it works

```
.
├── common
│   ├── bcrypt.go         //generate passwords
│   ├── config.go         //config for app
│   ├── database.go       //DB connect
│   ├── env.go            //environment for services
│   ├── error_messages.go //errors
│   ├── interfaces.go     //interfaces for app
│   └── session.go        //session
├── db        // migrations for DB
├── oauth
│   ├── clients.go        //clients data(in memory)
│   ├── routers.go        //business logic and gin routers
│   └── server.go         //oauth2 server
├── routers   // init gin
├── templates
└── users
    ├── models.go       //data models and DB operation
    ├── routers.go      //business logic and gin routers
    └── serializers.go  //json
```


# Getting started

## Install the Golang and GO environment

https://golang.org/doc/install

## Install glide

https://github.com/Masterminds/glide

## Install sql-migrate

```
go get -v github.com/rubenv/sql-migrate/...
```

## Clone repository

```
git clone git@github.com:Aristat/golang-gin-oauth2-example-app.git (go get)
```

## Start

```
run server
➜  glide install
➜  sql-migrate up (create database)
➜  go run *.go

run client
➜  cd client
➜  go run client.go

➜  http://localhost:9094/login
```

## Testing
```
➜  sql-migrate up -env="test"
➜  go test -v ./...
```