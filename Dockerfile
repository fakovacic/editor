FROM golang:1.21-alpine as builder

RUN apk update && apk add --no-cache git ca-certificates tzdata build-base openssh-client

WORKDIR /src

COPY go.* .

RUN go mod download

COPY . .

# Build the Go app
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o ./bin/app ./cmd/app/main.go

FROM scratch AS final

# Import the time zone files
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# Import the CA certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Import the compiled go executable
COPY --from=builder /src/bin/app /app 

WORKDIR /

ENTRYPOINT ["/app"]

EXPOSE 8080
EXPOSE 8081