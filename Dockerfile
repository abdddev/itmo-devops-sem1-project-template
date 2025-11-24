FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o server ./cmd/server/main.go
RUN go build -o migrator ./cmd/migrator/main.go


FROM alpine:3.19 AS prod
WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .
COPY --from=builder /app/migrator .
COPY --from=builder /app/migrations ./migrations

COPY .env ./

ENTRYPOINT ["./server"]
