# Start from the latest golang base image
FROM golang:1.12-alpine as builder

RUN apk add --no-cache tzdata
RUN apk add git
# This is an opt in feature in Go 1.1
ENV GO111MODULE=on

WORKDIR /app

## Copy go mod and sum files into the gatekeeper directory in the container.
COPY go.mod ./

RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

WORKDIR /app

## Build the Go app
RUN GOOS=linux go build -o app .
#
# New stage.
FROM alpine

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app .

# Command to run the executable
CMD ["./app"]