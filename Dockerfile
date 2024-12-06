ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN update-ca-certificates
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app .


FROM debian:bookworm

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /run-app /usr/local/bin/
CMD ["run-app"]
