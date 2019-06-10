FROM alpine:3.5

ENV MYAWS_VERSION=v0.3.8

RUN apk --no-cache add curl ca-certificates && update-ca-certificates
RUN curl -fsSL https://github.com/minamijoyo/myaws/releases/download/${MYAWS_VERSION}/myaws_${MYAWS_VERSION}_linux_amd64.tar.gz \
    | tar -xzC /usr/local/bin && chmod +x /usr/local/bin/myaws
RUN apk del --purge curl
