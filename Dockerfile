FROM golang:1.9

RUN mkdir -p /app

WORKDIR /cli

RUN go get github.com/aambhaik/topview/...

# Replace the value with the registry host and port that you'd like to use
ENV TMGC_REGISTRY_HOST=$TMGC_REGISTRY_HOST
ENV TMGC_REGISTRY_PORT=1080