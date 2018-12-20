# BUILD BINARY FIRST WITH OFFICIAL GOLANG IMAGE
FROM golang:alpine as builder

# Set working directory and copy code into working directory
WORKDIR /go/src/github.com/edwintcloud/shippy/user-cli
COPY . .

# Fetch dependencies using govendor
RUN apk update && apk add git
RUN go get -u github.com/kardianos/govendor
RUN govendor sync

# Build binary
RUN go build

# SETUP OUR ALPINE ENV FOR RUNNING BUILT EXECUTABLE
FROM alpine:latest

# Create app directory and set it as working directory
RUN mkdir /app
WORKDIR /app

# Pull binary from container we labeled as builder
COPY --from=builder /go/src/github.com/edwintcloud/shippy/user-cli/user-cli .

# Run the binary
CMD ["./user-cli"]