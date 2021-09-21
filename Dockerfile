# Taken from https://github.com/olliefr/docker-gs-ping/blob/main/Dockerfile
FROM golang:1.16-alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod .
# COPY go.sum .
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . /app

# Build
RUN go build -o /app/bin/naivedb-server /app/cmd/server

# This is for documentation purposes only.
# To actually open the port, runtime parameters
# must be supplied to the docker command.
EXPOSE 8080

# Run
CMD [ "/app/bin/naivedb-server" ]