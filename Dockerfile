# syntax=docker/dockerfile:1

FROM golang

WORKDIR /go/src/multicast
COPY . .


RUN go get -d -v ./... \
    && go install -v github.com/SimoBenny8/MulticastSDCCProject/cmd/multicast

CMD [ "multicast" ]