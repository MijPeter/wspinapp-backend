FROM golang:1.19.3

# Set the Current Working Directory inside the container
WORKDIR /wspinapp

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o /server ./cmd

# Run the binary program produced by `go install`
CMD ["/server"]
