FROM golang:1.24.0 as builder

WORKDIR /app
COPY . .
RUN go mod download

RUN GOSUMDB=off CGO_ENABLED=0 GOOS=linux go build -o bin/client cmd/client/main.go

FROM alpine:3.15.0

COPY --from=builder /app/bin/client /word-of-wisdom-client

EXPOSE 8080

CMD ["/word-of-wisdom-client"]
