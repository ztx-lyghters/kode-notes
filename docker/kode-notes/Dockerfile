FROM golang:1.21.1-alpine3.18

RUN go version

RUN apk update
RUN apk upgrade
RUN apk add postgresql15-client

RUN mkdir -p /usr/local/bin
COPY ./run /usr/local/bin/
COPY ./kode-notes /usr/local/bin/

RUN chmod +x /usr/local/bin/run
RUN chmod +x /usr/local/bin/kode-notes

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/run"]
