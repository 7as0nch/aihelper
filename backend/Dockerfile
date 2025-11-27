FROM golang:1.24 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn make build

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 make build

FROM debian:stable-slim

# RUN apt-get update && apt-get install -y --no-install-recommends \
# 		ca-certificates  \
#         netbase \
#         && rm -rf /var/lib/apt/lists/ \
#         && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  \
#         && echo Asia/Shanghai > /etc/timezone \
#         && apt-get autoremove -y && apt-get autoclean -y

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  \
        && echo Asia/Shanghai > /etc/timezone \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app

WORKDIR /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf
VOLUME /app/logs

CMD ["./backend", "-conf", "/data/conf/config.yaml"]
