FROM golang:1.24-alpine

RUN apk add --no-cache postgresql-client

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY wait-for-postgres.sh /wait-for-postgres.sh
RUN chmod +x /wait-for-postgres.sh

RUN go build -o main ./cmd/main.go

EXPOSE 8000

CMD ["/wait-for-postgres.sh", "./main"]