FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY cmd .
COPY internal ./internal
COPY cmd ./cmd
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o api cmd/api/main.go

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/api .
ENTRYPOINT ["./api"]
