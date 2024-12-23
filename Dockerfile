## Build Stage ##
# First pull Golang image
FROM golang:1.23-alpine as build-env

ENV APP_NAME go-store-gateway
ENV CMD_PATH cmd/main.go

COPY . /app
WORKDIR /app

# Build application
RUN CGO_ENABLED=0 go build -v -o /$APP_NAME $CMD_PATH

## Run Stage ##
FROM alpine:3.14

ENV APP_NAME go-store-gateway
COPY --from=build-env /$APP_NAME .
# Copy config folder
# COPY <src-path> <destination-path>
COPY --from=build-env /app/config config

EXPOSE 8081

CMD ./$APP_NAME