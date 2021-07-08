# run
```
docker run -d -p 80:80 -v /app/docker_data/alertsender:/usr/src/app/data --name alertsender registry.cn-hangzhou.aliyuncs.com/riet/alertsender:latest
```

# alertmanager config
```
global:
  resolve_timeout: 1m

receivers:
- name: 'default'
  webhook_configs:
  - url: 'http://172.21.54.60:28888/alert-receiver/default/'
    send_resolved: true


route:
  receiver: default
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 2m
  repeat_interval: 3h
```