FROM golang:1.18 AS build

ADD . /api
WORKDIR /api
RUN go build -o bin/main cmd/main.go

FROM ubuntu:20.04

WORKDIR /usr/src/app

COPY --from=build /api/bin/main .
COPY --from=build /api/config.json .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /api/kit-lokle.crt .
COPY --from=build /api/kit-lokle.key .

EXPOSE 3001
CMD ./main
