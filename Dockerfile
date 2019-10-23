FROM golang:1.13.3-alpine3.10 AS build-env
RUN apk --no-cache add make git
ADD . /work
WORKDIR /work
RUN make build

FROM alpine:3.10
RUN apk --no-cache add ca-certificates && update-ca-certificates
COPY --from=build-env /work/bin/myaws /usr/local/bin/myaws
