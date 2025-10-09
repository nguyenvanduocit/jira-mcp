FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Use build arguments for cross-compilation
ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o jira-mcp .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/jira-mcp .

# Expose port for SSE server (optional)
EXPOSE 8080

ENTRYPOINT ["/app/jira-mcp"]

CMD []