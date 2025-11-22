FROM golang:1.24

WORKDIR /app

# Install git only if needed — safe to keep, small package
# RUN apk add --no-cache git

# Copy Go module files first (build cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy project
COPY . .

EXPOSE 3000

CMD sh -c " \
  go run . artisan key:generate && \
  go run . artisan migrate && \
  go run . artisan db:seed && \
  go run main.go \
"
