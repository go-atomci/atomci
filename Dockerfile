## build
FROM golang:1.15-buster AS build-env

ADD . /go/src/atomci

WORKDIR /go/src/atomci

RUN make build

## run
FROM alpine:3.9

LABEL maintainer="Colynn Liu <colynn.liu@gmail.com>"

ADD conf /atomci/conf

RUN mkdir -p /atomci && mkdir -p /atomci/log

WORKDIR /atomci

COPY --from=build-env /go/src/atomci/atomci /atomci

ENV PATH $PATH:/atomci

EXPOSE 8080
CMD ["./atomci"]
