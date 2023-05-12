# Use the official Go image as the base image
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the Go app
RUN go build -o oracle-monitoring ./cmd/oracle-monitoring

# Start a new stage for the final image
FROM alpine:3.14

# Install ca-certificates
RUN apk --no-cache add ca-certificates

# Copy the binary from the previous stage
COPY --from=0 /app/oracle-monitoring /usr/local/bin/oracle-monitoring

# Run the oracle-monitoring binary
ENTRYPOINT ["oracle-monitoring"]
