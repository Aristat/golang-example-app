![golang logo](golang_logo.png)

# Golang Example Application


# Table of Contents

- [Overview](#overview)
- [Package list](#package-list)
- [Installing](#installing)
  * [Local environment](#local-environment)
  * [Docker environment](#docker-enviroment)
- [Run services](#run-services)
  * [Start in local](#start-in-local-machine)
  * [Start in docker](#start-in-docker)
- [Getting started](#getting-started)  
  * [Jaeger](#jaeger)
  * [Oauth2 client](#oauth2-client)
  * [Http with gRPC](#http-example-with-grpc)
  * [Graphql with gRPC](#graphql-example-with-grpc)
- [Deprecate version](#deprecated-version)  
- [Testing](#testing)

# Overview

This is an example golang application.
Commands list:
1. Daemon - main service
2. Product service - service that returns Product, an example of gRPC client/server interaction
3. Health check service - this service is needed to show how convenient to understand on which of the services an error occurred in jaeger
4. Oauth client - this service is needed to show a simple example of http client and server, for example oauth2 server
5. Migrate - commands for migration
6. Jwt - commands for generate jwt token

# Package list

Packages which use in this example project

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
13. [jaeger](https://github.com/uber/jaeger-client-go) - Jaeger Bindings for Go OpenTracing API
14. [casbin](https://github.com/casbin/casbin) - Supports access control
15. [dataloaden](https://github.com/vektah/dataloaden) - DataLoader for graphql
16. [nats](https://github.com/nats-io/nats.go) - Golang client for NATS, the cloud native messaging system
17. [nats-streaming](https://github.com/nats-io/stan.go) - NATS Streaming System

# Installing

Install the Golang and GO environment

```$xslt
https://golang.org/doc/install
```

Clone repository

```$xslt
git clone git@github.com:Aristat/golang-example-app.git (go get)
```

## Local environment

Install Golang packages without modules

```$xslt
make install
```

Install database

```$xslt
make createdb
```

Sql migrations

```$xslt
sql-migrate up
```

Install Golang dependencies

```$xslt
make dependencies
```

Generate artifacts(binary files and configs)

```$xslt
make build
```

## Docker environment


Install Golang packages without modules

```$xslt
make install
```

Generate artifacts(binary files and configs)

```$xslt
GOOS=linux GOARCH=amd64 make build
```

Generate docker image

```$xslt
DOCKER_IMAGE=golang-example-app TAG=development make docker-image
```

# Run services

## Start in local machine

Up jaeger in docker-compose or disable Jaeger(and rebuild binary file) in `resources/configs/*.yaml`

```$xslt
docker-compose up jaeger
```

Start daemon (main service)

```$xslt
make start
```

or

```$xslt
./artifacts/bin daemon -c ./artifacts/configs/development.yaml -d
```

Start client 

```$xslt
./artifacts/bin oauth-client -c ./artifacts/configs/development.yaml -d
```

Start product service 

```$xslt
./artifacts/bin product-service -c ./artifacts/configs/development.yaml -d
```

Start health-check service

```$xslt
./artifacts/bin health-check -c ./artifacts/configs/development.yaml -d
```

## Start in docker

```$xslt
docker-compose rm # Remove previous containers
GOOS=linux GOARCH=amd64 make build # Generate binary file
REMOVE_CONTAINERS=on DOCKER_IMAGE=golang-example-app TAG=development make docker-image # Generate new docker image
docker-compose up
```

# Getting Started

## Jaeger

```$xslt
http://localhost:16686
```

## Oauth2 client

```$xslt
http://localhost:9094/login
```

## Http example with gRPC

```$xslt
http://localhost:9096/products_grpc
```

## Http example with Nats Streaming

```$xslt
http://localhost:9096/products_nats
```

## Http example with graceful shutdown(long request, you can check the server shutdown)

```$xslt
http://localhost:9096/products_slowly
```

## Graphql example with gRPC

Graphql [client](https://github.com/prisma-labs/graphql-playground) for testing. End-point `http://localhost:9096/query`.

Generate JWT for Graphql authorization

```$xslt
./artifacts/bin jwt token --key='./artifacts/keys/local/private_key.pem' --fields='{"sub": "owner", "iss": "test-service"}'
```

Set JWT token in headers

```$xslt
{
  "Authorization": "bearer token"
}
```

Example query

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

Example query with data loader

```
query allProducts {
  products {
    list {
      products {
        id
        productItems {
          id
          name
        }
      }
    }
  }
}
```

Example mutation

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

# Deprecated version

[Deprecated version](https://github.com/Aristat/golang-example-app/tree/gin-example)

Usage only `gin` package and oauth2 client/server mechanic

# Testing
```
export APP_WD=go_path to project_path/rsources or project_path/artifacts - needed for load templates
➜  make test
```
