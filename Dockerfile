# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /go/src/multicast
COPY . .


RUN go get ./... & \
    go install -v MulticastSDCCProject/cmd/multicast

CMD [ "multicast" ]