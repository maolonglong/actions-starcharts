FROM golang:1.19.1-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOPROXY=https://proxy.golang.org

WORKDIR /src
COPY . .

RUN go build -o /bin/action -ldflags "-w -s"

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /bin/action /bin/action

ENTRYPOINT ["/bin/action"]
