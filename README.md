# Golang Example Application

## NOTE! [Deprecated version](https://github.com/Aristat/golang-example-app/tree/gin-example)

In the deprecated version using only `gin` package.

# Getting started

## Package list, which using in this example project

1. [sql-migrate](https://github.com/rubenv/sql-migrate) - sql migrations
2. [wire](https://github.com/google/wire) - dependency Injection
3. [viper](https://github.com/spf13/viper) - environment configuration
4. [cobra](https://github.com/spf13/cobra) - create commands
5. [cast](https://github.com/spf13/cast) - easy casting from one type to another
6. [gorm](https://github.com/jinzhu/gorm) - database ORM
7. [zap](https://github.com/uber-go/zap) - logger
8. [mux](https://github.com/gorilla/mux) - http router
9. [oauth2](https://github.com/go-oauth2/oauth2) - simple oauth2 server
10. [gqlgen](https://github.com/99designs/gqlgen) - graphql server library
11. [protobuf](https://github.com/golang/protobuf) - Google's data interchange format
12. [grpc](google.golang.org/grpc) - RPC framework
13. [jaeger](http://github.com/uber/jaeger-client-go) - Jaeger Bindings for Go OpenTracing API

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

## Jaeger Tracing


##### Run this command

```
 ➜ docker-compose up jaeger
 ➜ http://localhost:16686
```

##### or disable Jaeger in

```
resources/configs/*.yaml
```

##  Start

#### Oauth2
 
#### Run

```
➜ cd artifacts/
➜ ./bin daemon -c ./configs/local.yaml -d
➜ ./bin client
➜ http://localhost:9094/login
```

#### Graphql 

##### Run
```
➜ cd artifacts/
➜ ./bin daemon -c ./configs/local.yaml -d
➜ http://localhost:9096/query
```

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

#### Grpc Testing

##### Run
```
➜ cd artifacts/
➜ ./bin daemon -c ./configs/local.yaml -d
➜ ./bin product_service
➜ http://localhost:9096/products
```

## Testing
```
➜  make test
```
