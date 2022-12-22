FROM golang:1-alpine3.12 AS build

WORKDIR /app

COPY ../ .
RUN go mod download
RUN go build main.go


FROM alpine:3.12

WORKDIR /app

COPY --from=build /app/main ./

ENTRYPOINT ["/app/main"]