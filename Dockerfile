FROM registry.cn-hangzhou.aliyuncs.com/riet/python:3.7.4
ENV DJANGO_SETTINGS_MODULE alertsender.settings
WORKDIR /usr/src/app
ADD . .
RUN pip install --no-cache-dir -r requirements.txt -i https://mirrors.aliyun.com/pypi/simple/ && mkdir data
EXPOSE 80
CMD cp db.sqlite3 data/db.sqlite3 && python manage.py migrate && gunicorn -b 0.0.0.0:80 --access-logfile - alertsender.wsgi
