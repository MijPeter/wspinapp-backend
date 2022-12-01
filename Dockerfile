FROM golang:1.19.3

# Set the Current Working Directory inside the container
WORKDIR /wspinapp

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./bin/wspinapp .


# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["./bin/wspinapp"]