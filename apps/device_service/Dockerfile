FROM golang:1.22-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /device-service .

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /device-service /device-service
EXPOSE 8083
ENTRYPOINT ["/device-service"]
