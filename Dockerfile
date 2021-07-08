FROM registry.cn-hangzhou.aliyuncs.com/riet/python:3.7.4
ENV DJANGO_SETTINGS_MODULE alertsender.settings
WORKDIR /usr/src/app
ADD . .
RUN pip install --no-cache-dir -r requirements.txt -i https://mirrors.aliyun.com/pypi/simple/ && \
    mkdir data && \
    chmod +x docker-entrypoint.sh && \
    mv docker-entrypoint.sh /usr/local/bin
EXPOSE 80
CMD ["docker-entrypoint.sh"]
