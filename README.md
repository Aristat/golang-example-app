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

## Clone repository

```
git clone git@github.com:Aristat/golang-example-app.git (go get)
```

## Install

```
➜ make install
➜ make createdb
➜ sql-migrate up (create database)
➜ make vendor
➜ make build
```

##  Start

#### Run Server and Client(only for oauth2)

```
➜ cd artifacts/
➜ ./bin daemon -c ./configs/local.yaml -d (run server with debug mod)
➜ ./bin client (client for oauth2)
➜  http://localhost:9094/login (testing oauth2)
➜  http://localhost:9096/query (graphql route)
```

#### Graphql 

##### User query
```
query oneUser {
  users {
    one(email: "test@gmail.com") {
      email
      id
    }
  }
}
```

##### Create User mutation
```
mutation createUser {
  users {
    createUser(email: "test1@gmail.com", password: "123456789") {
      id
      email
    }
  }
}
```

## Testing
```
➜  make test
```
