# Use the official Golang image
FROM golang:1.22 AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY ./src ./src

# Build the Go app
RUN cd src && GOARCH=amd64 GOOS=linux go build -o /app/main .

# Start a new stage from scratch
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/main .

# Command to run the executable
COPY config ./config
COPY abi ./abi
COPY migrations ./migrations

CMD ["./main"]
