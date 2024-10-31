# Dockerfile
FROM golang:1.23 AS dev

# Set working directory in container
WORKDIR /app

# Install fsnotify for hot reloading and any other necessary tools
RUN go install github.com/githubnemo/CompileDaemon@latest

# Copy go.mod and go.sum files, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Expose port 50051
#EXPOSE 50051

# Start CompileDaemon for auto-recompiling
CMD ["CompileDaemon", "--build=go build -o server .", "--command=./server"]

