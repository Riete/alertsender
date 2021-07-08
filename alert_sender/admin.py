from django.contrib import admin
from alert_sender.models import DingTalkNotify


@admin.register(DingTalkNotify)
class DingTalkUrlAdmin(admin.ModelAdmin):
    list_display = ('env', 'url')
    search_fields = list_display
    list_filter = list_display


admin.site.site_header = '告警管理后台'
admin.site.site_title = '告警管理后台'
