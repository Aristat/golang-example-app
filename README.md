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
13. [jaeger](https://github.com/uber/jaeger-client-go) - Jaeger Bindings for Go OpenTracing API
14. [casbin](https://github.com/casbin/casbin) - Supports access control

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

##### or disable Jaeger(and rebuild binary file) in

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

#### Grpc Testing

##### Run
```
➜ cd artifacts/
➜ ./bin daemon -c ./configs/local.yaml -d
➜ ./bin product_service
➜ ./bin health-check
➜ http://localhost:9096/products
```

You with a probability of 50/50 will receive either the correct answer or an error on one of the services.
It's needed for check jaeger in http://localhost:16686

#### Graphql 

Add `Authorization` in header by Bearer JWT token

Example token:
```
eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJvd25lciIsImlzcyI6InRlc3Qtc2VydmljZSJ9.WY1deEZU5IWJ4Kx-udiKakCm0Q7PEVs6ZsUje6StJLy6gzZ8MHPI89rdsM9FkhiA1ZhjxsNGM93e2huwjRsTRhV_fIwQSFrH72M2g7c7lxh4U_q8C1OfSee2Ffy4wVh3dCQ5Nz3BKoKYVh2E1PSzMSm-3SDs6q-UTTjzRCWOORKdh9gisyhHbL8zbjHLBHSsiG1DPWin0beGSmA92cpwpaICEEK-lSNhDRlrCHYMJYAjKBphwpQY4PjMC_rKykQM_mAeKFdj4pXiReDw0QuCKXseWo_b46PO-YnukYM26fogrbwFb0bhz9FuQOuusZAz-ONmAaCVeZ_OK9nHyCuswg
```

##### Run
```
➜ cd artifacts/
➜ ./bin daemon -c ./configs/local.yaml -d
➜ http://localhost:9096/query
```

##### User query

You with a probability of 50/50 will receive the correct result. It's needed for check jaeger

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
