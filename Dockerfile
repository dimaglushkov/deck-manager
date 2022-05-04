FROM golang:1.18.1-buster

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY . .

RUN ["go", "test", "./app/..."]
RUN go build ./cmd/main.go

CMD ["./main"]