FROM golang:1.25.9

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy

RUN go build -o app ./cmd/main.go

CMD ["./app"]