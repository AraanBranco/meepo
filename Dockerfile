FROM golang:1.23.4-alpine3.21 AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . /build

RUN mkdir -p /app

RUN apk add --update make

RUN cd /build && \
    make build-linux-x86_64 && \
    mv ./bin/meepo-linux-x86_64 /app/meepo && \
    mv ./config /app/config

FROM alpine

WORKDIR /app

COPY --from=build /app/meepo /app/meepo
COPY --from=build /app/config /app/config

EXPOSE 3000 8080

CMD /app/meepo