# syntax=docker/dockerfile:1
# https://docs.docker.com/language/golang/build-images/

FROM golang:1.20 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go test
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GIN_MODE=release go build -o ./go-rest-sample
RUN addgroup --system client && adduser --system --group client


#FROM gcr.io/distroless/base
FROM alpine:3.18
RUN apk update && apk add bash && apk add util-linux
WORKDIR /app
COPY --from=build-stage /etc/passwd /etc/passwd
COPY --from=build-stage /app/go-rest-sample .
USER client
EXPOSE 8080
ENTRYPOINT ["./go-rest-sample"]