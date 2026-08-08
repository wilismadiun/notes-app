FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o notes-app .


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/notes-app .

EXPOSE 3000

CMD ["./notes-app"]