FROM golang:1.22-alpine AS builder

WORKDIR /src
COPY go.mod main.go main_test.go ./
RUN go test ./... && CGO_ENABLED=0 go build -o /temperature-api .

FROM alpine:3.20
COPY --from=builder /temperature-api /temperature-api
EXPOSE 8081
ENTRYPOINT ["/temperature-api"]
