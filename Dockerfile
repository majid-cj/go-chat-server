FROM golang:bookworm

RUN apt update && apt upgrade -y && apt install -y git make openssh-client

RUN mkdir -p /app

WORKDIR /app

ADD . /app

EXPOSE 8080

RUN go mod tidy

RUN curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
