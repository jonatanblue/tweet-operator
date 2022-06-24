# syntax=docker/dockerfile:1
FROM golang:1.18-alpine AS build

RUN apk add git

COPY go.mod go.sum main.go /app/
COPY pkg /app/pkg

WORKDIR /app/

RUN go get

RUN go build -o /bin/app

FROM alpine:latest
COPY --from=build /bin/app /bin/app
ENTRYPOINT [ "/bin/app" ]
