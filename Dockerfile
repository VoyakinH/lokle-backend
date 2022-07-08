FROM golang:1.18 AS build

ADD . /api
WORKDIR /api
RUN go build main.go

FROM ubuntu:20.04

WORKDIR /usr/src/app

COPY . .
COPY --from=build /api/main/ .

EXPOSE 3001
CMD ./main
