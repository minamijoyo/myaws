FROM golang:1.19-alpine3.16 AS build-env
RUN apk --no-cache add make git
WORKDIR /work

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM alpine:3.16
RUN apk --no-cache add ca-certificates && update-ca-certificates
COPY --from=build-env /work/bin/myaws /usr/local/bin/myaws
