FROM golang:1.14-alpine3.11

ARG goproxy


# Download packages from aliyun mirrors
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk --update add --no-cache ca-certificates tzdata git
RUN wget https://github.com/buger/goreplay/releases/download/v1.2.0/gor_v1.2.0_x64.tar.gz \
  && tar xzf gor_v1.2.0_x64.tar.gz -C  /bin
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








COPY --from=0 /bin/gor /go/bin/check-proxy /go/bin/checksever /bin/




WORKDIR /
## 心跳
#CMD sudo nohup /bin/gor  --input-raw :8090  --input-raw-track-response  \
# --output-http-track-response \
# --output-tcp localhost:28020 \
# --output-http $logserver &
## 拦截
#CMD sudo nohup /bin/gor  --input-raw-track-response  \
#  --output-http-track-response \
#  --input-tcp localhost:28020 \
#  --output-http $compareserver \
#  --middleware "/bin/check-proxy http://$logserver" &



ENTRYPOINT ["/bin/checksever"]