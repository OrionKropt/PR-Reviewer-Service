# -------- Stage 1: Build --------
FROM golang:1.24.9-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o pr-reviewer-service ./cmd/pr-reviewer-service/main.go


# -------- Stage 2: Production image --------

FROM alpine:latest

COPY --from=builder /app/pr-reviewer-service .
COPY --from=builder /app/configs ./configs

EXPOSE 8080

ENTRYPOINT ["./pr-reviewer-service"]

