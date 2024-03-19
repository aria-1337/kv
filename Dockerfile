FROM ubuntu:22.04

# Args (ex: docker run --env PORT=3001 LEVEL_DB_PATH=awesome -t kv:latest)
ARG PORT
ENV PORT=$PORT

ARG LEVEL_DB_PATH
ENV LEVEL_DB_PATH=$LEVEL_DB_PATH

EXPOSE $PORT

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
  apt-get -y --no-install-recommends install \
    build-essential \
    curl \
    golang \
    ca-certificates \
    git && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN update-ca-certificates

WORKDIR /
ENV GOPATH /go
ENV PATH ${PATH}:/dist

WORKDIR /dist/src
WORKDIR /

COPY build.sh /dist
COPY ./go.mod /dist/src
COPY ./go.sum /dist/src
COPY src/*.go /dist/src
WORKDIR /dist

CMD ./build.sh -port=$PORT -leveldbPath=$LEVEL_DB_PATH
