FROM alpine:3.5

RUN apk --no-cache add curl ca-certificates && update-ca-certificates
RUN curl -fsSL https://github.com/minamijoyo/myaws/releases/download/v0.1.3/myaws_v0.1.3_linux_amd64.tar.gz \
    | tar -xzC /usr/local/bin && chmod +x /usr/local/bin/myaws
