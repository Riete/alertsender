import requests

def send_dingding(url: str, title: str, message: str, is_at_all: bool=False):
    data = {
        "msgtype": "markdown",
        "markdown": {
            'title': '{0}'.format(title),
            'text': '### {0}\n\n{1}'.format(title, message)
        },
        "at": {
            "atMobiles": [],
            "isAtAll": is_at_all
        },
    }
    requests.post(url, json=data)
