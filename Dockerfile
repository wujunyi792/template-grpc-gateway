FROM golang:latest AS builder

COPY . /build

WORKDIR /build

RUN set -ex \
    && GO111MODULE=auto CGO_ENABLED=0 go build -ldflags "-s -w -extldflags '-static' -X 'pinnacle-primary-be/core/version.SysVersion=$(git show -s --format=%h)'" -o App

FROM alpine:latest

WORKDIR /Serve
RUN mkdir "config"

COPY --from=builder /build/App ./App

RUN ls -R

RUN  echo 'http://mirrors.ustc.edu.cn/alpine/v3.9/main' > /etc/apk/repositories \
    && echo 'http://mirrors.ustc.edu.cn/alpine/v3.9/community' >>/etc/apk/repositories \
    && apk update && apk add tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

ENTRYPOINT [ "/Serve/App", "server" ]