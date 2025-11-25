# -------- Stage 1: Build --------
FROM golang:1.24.9-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./

RUN apk add --no-cache make

RUN make test

RUN make build

# -------- Stage 2: Production image --------

FROM alpine:latest

COPY --from=builder /app/bin/pr-reviewer-service .
COPY --from=builder /app/configs ./configs

EXPOSE 8080

ENTRYPOINT ["./pr-reviewer-service"]

