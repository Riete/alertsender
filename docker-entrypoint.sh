#!/bin/bash

cd /usr/src/app

[ ! -f data/db.sqlite3 ] && cp db.sqlite3 data/db.sqlite3
python manage.py migrate
gunicorn -b 0.0.0.0:80 --access-logfile - alertsender.wsgi
