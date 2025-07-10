# Build Stage
FROM golang:1.20-alpine as builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o app .

# Final Stage
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/app .
COPY --from=builder /app/db.json .

EXPOSE 8080
CMD ["./app"]
