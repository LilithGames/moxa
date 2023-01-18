FROM golang:1.19-buster

WORKDIR /root
ENV PATH="$PATH:/root"

RUN go install github.com/rakyll/hey@v0.1.4

COPY bin/moxactl .
