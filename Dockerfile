FROM registry.cn-hangzhou.aliyuncs.com/riet/python:3.7.4
ENV DJANGO_SETTINGS_MODULE alertsender.settings
WORKDIR /usr/src/app
ADD . .
RUN pip install --no-cache-dir -r requirements.txt -i https://mirrors.aliyun.com/pypi/simple/
EXPOSE 80
CMD python manage.py migrate && gunicorn -b 0.0.0.0:80 --access-logfile - alertsender.wsgi
