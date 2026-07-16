FROM golang:1.25.5-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/todoapp/exe /app/cmd/todoapp

FROM alpine:3.23.5
WORKDIR /app
COPY --from=builder /app/cmd/todoapp/exe /app
CMD ["/app/exe"]
