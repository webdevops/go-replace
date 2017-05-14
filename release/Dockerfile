FROM alpine

RUN apk --no-cache add --virtual .gocrond-deps \
        ca-certificates  \
        wget \
    && GOREPLACE_VERSION=1.1.2 \
    && wget -O /usr/local/bin/go-replace https://github.com/webdevops/go-replace/releases/download/$GOREPLACE_VERSION/gr-64-linux \
    && chmod +x /usr/local/bin/go-replace \
    && apk del .gocrond-deps

CMD ["go-replace"]
