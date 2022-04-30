FROM golang:1.17.9-alpine3.15 as builder

RUN apk add --update make \
    && rm -rf /tmp/* \
    && rm -rf /var/cache/apk/*

RUN mkdir /build
WORKDIR /build

COPY . .

RUN GOOS=linux GOARCH=amd64 make build

FROM alpine:3.9

COPY --from=builder /build/docker/app/files /
COPY --from=builder /build/artifacts /app/

WORKDIR /app

RUN apk --no-cache add tzdata bash \
    && chmod +x /docker/bin/* \
    && chmod +x bin \
    && rm -rf /tmp/* \
    && rm -rf /var/cache/apk/*

ENTRYPOINT ["/docker/bin/entrypoint.sh"]
