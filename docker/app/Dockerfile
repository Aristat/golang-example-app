FROM alpine:3.9

WORKDIR /app

COPY ./docker/app/files /
COPY ./artifacts .

RUN apk --no-cache add tzdata bash \
    && chmod +x /docker/bin/* \
    && chmod +x bin \
    && rm -rf /tmp/* \
    && rm -rf /var/cache/apk/*

ENTRYPOINT ["/docker/bin/entrypoint.sh"]
