FROM golang:alpine

RUN apk add --update git
RUN go get -u github.com/valyala/fasthttp
RUN go get -u github.com/hashicorp/consul/api
RUN go get -u github.com/satori/uuid

RUN mkdir /app
WORKDIR /app
ADD weatherservice/ /app/

CMD ["/app/main"]

RUN go build -o main .
