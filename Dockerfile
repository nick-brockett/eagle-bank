FROM golang:1.23-alpine  AS build

# set working directory
WORKDIR /app


COPY go.mod ./
COPY go.sum ./

# install dependencies
RUN go mod download

# copy source code
COPY . .

RUN go build -o main ./cmd/server/main.go

EXPOSE 8080

CMD ["./main"]