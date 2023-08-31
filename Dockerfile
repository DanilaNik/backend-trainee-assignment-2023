FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
COPY config/local.yaml /app/local.yaml

RUN go build -o ./bin/app ./cmd/app/main.go

EXPOSE 8080

CMD "./bin/app"