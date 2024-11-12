FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o execute-command .

FROM alpine:3.16

WORKDIR /app

COPY --from=builder /app/execute-command .

EXPOSE 8080

CMD ["./execute-command"]