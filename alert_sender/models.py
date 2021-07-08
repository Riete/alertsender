from django.db import models


class DingTalkNotify(models.Model):
    env = models.CharField('环境', max_length=100, unique=True)
    url = models.CharField('地址', max_length=200)

    class Meta:
        db_table = 'ding_talk_notify'
        verbose_name = '钉钉告警'
        verbose_name_plural = verbose_name

    def __str__(self):
        return self.env
