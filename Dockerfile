FROM golang:1.22-alpine3.20 AS build-env
RUN apk --no-cache add make git

# A workaround for a permission issue of git.
# Since UIDs are different between host and container,
# the .git directory is untrusted by default.
# We need to allow it explicitly.
# https://github.com/actions/checkout/issues/760
RUN git config --global --add safe.directory /work

WORKDIR /work

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM alpine:3.20
RUN apk --no-cache add ca-certificates && update-ca-certificates
COPY --from=build-env /work/bin/myaws /usr/local/bin/myaws
