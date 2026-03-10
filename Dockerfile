FROM golang:1.26.1-bookworm

WORKDIR /app

COPY . .

RUN go mod tidy

CMD ["go", "run", "./cmd/app/main.go"]
