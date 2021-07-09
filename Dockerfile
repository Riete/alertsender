FROM registry.cn-hangzhou.aliyuncs.com/riet/golang:1.13.10 as backend
COPY . .
RUN unset GOPATH && go build -mod=vendor

FROM registry.cn-hangzhou.aliyuncs.com/riet/centos:7.4.1708-cnzone
WORKDIR /opt
COPY --from=backend /go/alertsender alertsender
COPY config.ini config.ini
EXPOSE 8000
CMD ["./alertsender"]
