# Golang Example Application

## NOTE! [Deprecated version](https://github.com/Aristat/golang-example-app/tree/gin-example)

In the deprecated version using only `gin` package.

# Getting started

## Package list, which using in this example project

1. [sql-migrate](https://github.com/rubenv/sql-migrate) - sql migrations
2. [wire](https://github.com/google/wire) - dependency Injection
3. [viper](https://github.com/spf13/viper) - environment configuration
4. [cobra](https://github.com/spf13/cobra) - create commands
5. [gorm](https://github.com/jinzhu/gorm) - database ORM
6. [zap](https://github.com/uber-go/zap) - logger
7. [mux](https://github.com/gorilla/mux) - http router
8. [oauth2](https://github.com/go-oauth2/oauth2) - simple oauth2 server

## Install the Golang and GO environment

https://golang.org/doc/install

## Install sql-migrate/wire

```
go get -v github.com/rubenv/sql-migrate/...
go get -v github.com/google/wire/cmd/wire
```

## Clone repository

```
git clone git@github.com:Aristat/golang-example-app.git (go get)
```

## Start

```
➜ make createdb
➜ sql-migrate up (create database)
➜ make vendor
➜ make build
➜ cd artifacts/
➜ ./bin daemon -c ./configs/local.yaml -d (run server with debug mod)
➜ ./bin client

➜  http://localhost:9094/login
```

## Testing
```
➜  sql-migrate up -env="test"
➜  make test
```
