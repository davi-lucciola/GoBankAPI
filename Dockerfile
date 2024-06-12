FROM golang:1.22.4

WORKDIR /app

COPY go.mod go.sum Makefile ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -C ./src -o ../bin/go-bank-api

EXPOSE 8084

CMD ["./bin/go-bank-api"]

