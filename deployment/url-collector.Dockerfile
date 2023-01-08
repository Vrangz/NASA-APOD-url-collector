FROM golang:1.19.4 AS base

RUN mkdir /app

COPY ./url-collector /app/url-collector

WORKDIR /app/url-collector/

RUN go test --race /app/url-collector/...

RUN go build -o url-collector-app ./cmd/main/main.go

CMD [ "./url-collector-app" ] 
