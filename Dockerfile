FROM golang:latest

RUN mkdir /app

ADD . /app/

WORKDIR /app

VOLUME /app/dist
CMD ["go", "build", "-o", "/app/dist/service", "./main.go"]