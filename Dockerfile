FROM golang:1.17.12

WORKDIR /src
COPY . .

ENV GO111MODULE=on

RUN go build -o /bin/action

ENTRYPOINT ["/bin/action"]
