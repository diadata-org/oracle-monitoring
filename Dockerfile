FROM golang:1.19 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o oraclemonitoring .


FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /app/oraclemonitoring .
COPY --from=builder /app/internal /internal/


ENTRYPOINT ["./oraclemonitoring"]