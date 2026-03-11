# syntax=docker/dockerfile:1
FROM golang:latest AS build-stage
ARG TARGETOS
ARG TARGETARCH
RUN mkdir -p -m 0700 /root/.ssh && \
    echo 'Host *' > /root/.ssh/config && \
    echo '    StrictHostKeyChecking no' >> /root/.ssh/config && \
    echo '    UserKnownHostsFile=/dev/null' >> /root/.ssh/config && \
    chmod 0644 /root/.ssh/config
WORKDIR /build
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=ssh \
    GOPROXY=direct GOSUMDB=off GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 \
    go build -ldflags "-s -w" -trimpath -o certbot

FROM scratch
COPY --from=build-stage --chmod=0755 /build/certbot /usr/sbin/certbot
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
ENTRYPOINT [ "certbot" ]