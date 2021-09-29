FROM golang:1.17 as builder

ENV GO111MODULE=on \
    GOPROXY=https://proxy.golang.org \
    CGO_ENABLED=0

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -ldflags="-w -s" -o starcharts

FROM alpine:3.14

COPY --from=builder /build/starcharts /starcharts

ENTRYPOINT [ "/starcharts" ]
