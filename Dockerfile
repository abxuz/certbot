FROM golang:alpine AS build-stage
WORKDIR /build
COPY . .
RUN go build -ldflags "-s -w" -trimpath -o certbot

# Deploy the application binary into a lean image
FROM alpine:latest AS release-stage
COPY --from=build-stage /build/certbot /usr/bin/
VOLUME /data
WORKDIR /data
COPY config.yaml .
ENTRYPOINT ["certbot"]
