FROM golang:alpine AS build-stage
WORKDIR /build
COPY . .
RUN go build -ldflags "-s -w" -trimpath -o certbot

FROM alpine:latest AS release-stage
COPY --from=build-stage /build/certbot /usr/bin/
WORKDIR /data
ENTRYPOINT ["certbot"]