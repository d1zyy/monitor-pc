FROM golang:1.26 AS builder 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY  . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o monitor ./cmd/monitor

FROM alpine:3.23

WORKDIR /app

RUN addgroup -S monitor && adduser -S monitor -G monitor

RUN apk add --no-cache curl

COPY --from=builder /app/monitor ./monitor

RUN chown -R monitor:monitor /app

USER monitor

EXPOSE 8080

CMD ["./monitor"]