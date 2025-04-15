FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o jira-mcp .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/jira-mcp .

# Expose port for SSE server (optional)
EXPOSE 8080

ENTRYPOINT ["/app/jira-mcp"]

CMD []