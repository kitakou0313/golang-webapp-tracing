FROM golang:1.20-bullseye

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .

RUN go build -o service

EXPOSE 8080

CMD [ "./service" ]