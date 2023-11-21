FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache 'git=~2'

# Install dependencies
ENV GO111MODULE=on
WORKDIR $GOPATH/src/packages/goginapp/
COPY . .

# Fetch dependencies.
# Using go get.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/main .

ENV PORT 8080
ENV GIN_MODE release
EXPOSE 8080

WORKDIR /$GOPATH/src/packages/goginapp/

# Run the Go Gin binary.
ENTRYPOINT ["go", "run", "main.go"]