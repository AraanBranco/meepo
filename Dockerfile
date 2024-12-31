FROM golang:1.23.4-alpine3.21 AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . /build

RUN mkdir -p /app

RUN apk add --update make

RUN cd /build && \
    make build && \
    mv ./bin/meepow /app/meepow && \
    mv ./config /app/config

FROM alpine

WORKDIR /app

COPY --from=build /app/meepow /app/meepow
COPY --from=build /app/config /app/config

ENTRYPOINT ["/app/meepow"]