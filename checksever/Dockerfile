FROM golang:1.14-alpine3.11

ARG goproxy


# Download packages from aliyun mirrors
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk --update add --no-cache ca-certificates tzdata git

# build bat for healthcheck
RUN GO111MODULE=on GOPROXY=$goproxy GOOS=linux CGO_ENABLED=0 go get -v github.com/astaxie/bat

RUN git clone https://github.com/SeigeC/check-proxy.git \
  && cd check-proxy/ \
  && GO111MODULE=on GOPROXY=$goproxy GOOS=linux CGO_ENABLED=0 go install \
  && cd ./checksever/ \
  && GO111MODULE=on GOPROXY=$goproxy GOOS=linux CGO_ENABLED=0 go install



FROM scratch

ARG logserver
ARG compareserver








COPY --from=0 /go/bin/check-proxy /go/bin/checksever /bin/




WORKDIR /


ENTRYPOINT ["/bin/checksever"]