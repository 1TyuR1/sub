# build
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd/app

# run
FROM gcr.io/distroless/base-debian12
ENV TZ=UTC
COPY --from=builder /bin/app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
