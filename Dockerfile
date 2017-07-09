FROM golang:alpine AS buildenv

COPY . /go/src/go-replace
WORKDIR /go/src/go-replace

RUN apk --no-cache add git \
    && go get \
    && go build \
    && chmod +x go-replace \
    && ./go-replace --version

FROM golang:alpine
COPY --from=buildenv /go/src/go-replace/go-replace /usr/local/bin
CMD ["go-replace"]
