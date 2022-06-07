FROM golang:alpine as builder
WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -o main

FROM alpine:latest
COPY --from=builder /app/main .
CMD ["./main"]