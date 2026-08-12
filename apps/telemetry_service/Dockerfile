FROM golang:1.22-alpine AS builder

WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /telemetry-service .

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /telemetry-service /telemetry-service
EXPOSE 8082
ENTRYPOINT ["/telemetry-service"]
