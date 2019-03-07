# Golang Gin OAuth2 Example

> ### Golang/Gin/Oauth2 codebase containing examples

# How it works

```
.
├── cmd // start commands
├── common
│   ├── bcrypt.go         //generate passwords
│   └── database.go       //DB connect
├── db        // migrations for DB
├── oauth
│   ├── session           //init session manager
│   ├── clients.go        //clients data(in memory)
│   ├── service.go        //business logic and gin routers
│   └── server.go         //oauth2 server
├── templates
└── users
    ├── models.go       //data models and DB operation
    ├── service.go      //business logic and gin routers
    └── serializers.go  //json
```


# Getting started

## Install the Golang and GO environment

https://golang.org/doc/install

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
➜  go mod vendor
➜  sql-migrate up (create database)
➜  make run_development

run client
➜  make run_client

➜  http://localhost:9094/login
```

## Testing
```
➜  sql-migrate up -env="test"
➜  make test
```