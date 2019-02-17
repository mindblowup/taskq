FROM golang:alpine

RUN apk update && apk add git

RUN go get github.com/mindblowup/taskq

EXPOSE 8001

ENTRYPOINT ["taskq"]

WORKDIR /root/