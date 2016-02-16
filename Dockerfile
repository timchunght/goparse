# golang image where workspace (GOPATH) configured at /go.
FROM golang:latest

# Copy the local package files to the container’s workspace.
ADD . /go/src/goparse
WORKDIR /go/src/goparse
# Build the golang-docker command inside the container.
RUN go install goparse

# Run the golang-docker command when the container starts.

# http server listens on port 8080.
