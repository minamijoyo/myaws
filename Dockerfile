FROM golang:1.17.3-alpine3.14 AS build-env
RUN apk --no-cache add make git
WORKDIR /work

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM alpine:3.14
RUN apk --no-cache add ca-certificates && update-ca-certificates
COPY --from=build-env /work/bin/myaws /usr/local/bin/myaws
