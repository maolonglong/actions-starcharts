FROM golang:1.19-alpine3.16 AS build-env

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOPROXY=https://proxy.golang.org

WORKDIR /go/src/github.com/maolonglong/actions-starcharts
COPY . .

RUN go build -o /bin/action -trimpath -buildvcs=false -ldflags="-s -w" .

FROM scratch

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /bin/action /bin/action

ENTRYPOINT ["/bin/action"]
