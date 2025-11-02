FROM golang:1.24-alpine as builder

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests TO DO
#RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN go build -o ./auth-service ./cmd/auth-service

# Start fresh from a smaller image
FROM alpine:3.9
RUN apk add ca-certificates

COPY --from=builder /app/auth-service /app/auth-service

# This container exposes port 8080 to the outside world
EXPOSE 9999

# Run the binary program produced by `go install`
CMD ["/app/auth-service"]