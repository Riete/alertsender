import json
import datetime
import pytz
import time
from rest_framework.views import APIView
from django.http import HttpResponse
from alert_sender.models import DingTalkNotify
from alert_sender.tools import send_dingding


class AlertSend(APIView):

    def post(self, requests, *args, **kwargs):
        env = self.kwargs['env']
        alert_data = requests.data
        if isinstance(alert_data, str):
            alert_data = json.loads(alert_data)
        status = '故障' if alert_data['status'] == 'firing' else '恢复'

        alerts: list = alert_data['alerts']
        get_five_items = slice(0, 5)
        while alerts:
            five_alerts = alerts[get_five_items]
            if five_alerts:
                self.process_alerts(status, five_alerts, env)
                [alerts.remove(i) for i in five_alerts]
        return HttpResponse('ok')

    def process_alerts(self, status: str, alerts: list, env: str):
        annotations_list: list = []
        alert_name: str
        severity: str
        starts_at: str
        ends_at: str
        alerts_at: str
        robot_url: str

        if DingTalkNotify.objects.filter(env=env).exists():
            robot_url = DingTalkNotify.objects.get(env=env).url
        else:
            robot_url = DingTalkNotify.objects.get(env='default').url

        for num, alert in enumerate(alerts):
            alert_labels: dict = alert['labels']
            if num == 0:
                alert_name = alert_labels['alertname']
                severity = alert_labels['severity']
                starts_at = '\n\n**故障时间:**\n\n\n- {0}'.format(self.convert_utc_to_local(alert['startsAt']))
                ends_at = '\n\n**恢复时间:**\n\n\n- {0}'.format(self.convert_utc_to_local(alert['endsAt']))
                alerts_at = '\n\n**告警时间:**\n\n\n- {0}'.format(time.strftime('%F %T'))
            annotations: dict = alert['annotations']
            annotations = [f'\n- {v}\n\n&nbsp;\n\n' for k, v in annotations.items() if k not in ['runbook_url']]
            annotations_list.extend(annotations)

        alert_title = '{0}\n\n**[{1}] [{2}]**\n\n&nbsp;\n\n---'.format(alert_name, status, severity)
        alert_message = '**告警内容:**\n\n{0}\n\n---\n\n'.format(''.join(annotations_list))
        if status == '恢复':
            alert_message += starts_at + '\n\n---' + ends_at
        else:
            alert_message += starts_at + '\n\n---' + alerts_at
        if severity in ['disaster', 'P0', 'P1', 'P2']:
            if status != '恢复':
                pass # add another alert method here, like phone call
            send_dingding(robot_url, alert_title, alert_message, True)
        else:
            send_dingding(robot_url, alert_title, alert_message, False)

    def convert_utc_to_local(self, utc_time):
        utc_time: str = utc_time.split('.')[0]
        if not utc_time.endswith('Z'):
            utc_time += 'Z'
        utc = datetime.datetime.strptime(utc_time, '%Y-%m-%dT%H:%M:%SZ').replace(tzinfo=pytz.utc)
        local = utc.astimezone(pytz.timezone('Asia/Shanghai'))
        return local.strftime('%F %T')
