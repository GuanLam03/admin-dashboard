FROM golang:1.24

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy project files
COPY . .

# Build Go binary
RUN go build -o app main.go

# Expose port (Render uses $PORT)
EXPOSE 3000

# Run binary (server only)
CMD ["./app"]
