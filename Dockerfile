FROM golang:1.18-alpine3.15 AS build-env

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOPROXY=https://proxy.golang.org

WORKDIR /go/src/go.chensl.me/actions-starcharts
COPY . .

RUN go build -o /go/bin/actions-starcharts -trimpath -buildvcs=false -ldflags="-s -w" .

FROM gcr.io/distroless/static

COPY --from=build-env /go/bin/actions-starcharts /

ENTRYPOINT [ "/actions-starcharts" ]
