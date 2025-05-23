FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o tupload

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/tupload .
COPY --from=builder /app/config/config.yaml ./config/config.yaml

RUN mkdir -p /app/uploads

EXPOSE 6060

CMD ["./tupload"] 
