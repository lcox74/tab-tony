# Builder
FROM golang:1.18-alpine AS builder

RUN apk update && apk add git gcc musl-dev

# Set working directory
COPY ./ /src/
WORKDIR /src

# Get packages and Build
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /bin/backend .

# Runner
FROM alpine:latest as deploy
LABEL description="The Aussie BroadWAN Bot Backend" Version="0.0.1"

# Copy the binary from the builder
COPY --from=builder /bin/backend /bin/backend

# Expose port 3000 to the outside world
EXPOSE 3000
RUN chmod +x /bin/backend
ENTRYPOINT ["/bin/backend"]