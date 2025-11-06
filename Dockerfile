FROM golang:1.17-alpine

RUN apk add ca-certificates git

WORKDIR /go/src/github.com/pierre-emmanuelJ/iptv-proxy
COPY . .
RUN GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown") && \
    BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") && \
    CGO_ENABLED=0 GOOS=linux go build -mod=vendor \
    -ldflags "-X github.com/pierre-emmanuelJ/iptv-proxy/cmd.GitCommit=${GIT_COMMIT} \
              -X github.com/pierre-emmanuelJ/iptv-proxy/cmd.BuildDate=${BUILD_DATE} \
              -X github.com/pierre-emmanuelJ/iptv-proxy/cmd.Version=channel-logging" \
    -a -installsuffix cgo -o iptv-proxy .

FROM alpine:3
COPY --from=0  /go/src/github.com/pierre-emmanuelJ/iptv-proxy/iptv-proxy /
ENTRYPOINT ["/iptv-proxy"]
