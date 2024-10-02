FROM golang:1.22

RUN mkdir /app
ADD . /app

WORKDIR /app

COPY go.mod go.sum ./

RUN go build -o app ./cmd/gophermart/main.go

EXPOSE 8080

CMD [ "/app/app" ]