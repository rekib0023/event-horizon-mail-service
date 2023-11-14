FROM golang:1.20.4 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o event-horizon-mailsvc ./cmd

FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/event-horizon-mailsvc .

COPY --from=builder /app/internal/templates ./internal/templates

ENTRYPOINT ["./event-horizon-mailsvc"]
