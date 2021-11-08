## build
FROM golang:1.15-alpine AS build-env

ADD . /go/src/atomci

WORKDIR /go/src/atomci

RUN go build ./cmd/atomci &&  go build ./cmd/cli

## run
FROM alpine:3.9

LABEL maintainer="Colynn Liu <colynn.liu@gmail.com>"

ADD conf /atomci/conf

RUN mkdir -p /atomci && mkdir -p /atomci/log

WORKDIR /atomci

COPY --from=build-env /go/src/atomci/atomci /atomci
COPY --from=build-env /go/src/atomci/cli /atomci

ENV PATH $PATH:/atomci

EXPOSE 8080
CMD ["./atomci"]
