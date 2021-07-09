package main

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/gin-gonic/gin"
	"github.com/riete/go-tools/notify"
	"log"
	"strings"
	"time"
)

const DingTalk = "dingtalk"

func ParseConfig(channel string) map[string]string {
	cfg, err := goconfig.LoadConfigFile("config.ini")
	if err != nil {
		log.Fatalln("load config file failed, config.ini is not exists")
	}
	conf, err := cfg.GetSection(channel)
	if err != nil {
		log.Fatalln("parse config file config.ini failed")
	}
	return conf
}

func GetDingTalkChannel(name string) (webhook, secret string) {
	conf := ParseConfig(DingTalk)
	channel := strings.Split(conf[name], "|")
	return channel[0], channel[1]
}

type Alert struct {
	Status string `json:"status"`
	Alerts []struct {
		Labels struct {
			Alertname string `json:"alertname"`
			Severity  string `json:"severity"`
		} `json:"labels"`
		Annotations struct {
			Message string `json:"message"`
		} `json:"annotations"`
		StartsAt string `json:"startsAt"`
		EndsAt   string `json:"endsAt"`
	} `json:"alerts"`
}

func (alert Alert) ConvertUtcToLocal(utcTime string) string {
	var timeLayoutStr = "2006-01-02T15:04:05Z"
	utcStr := strings.Split(utcTime, ".")[0]
	if !strings.HasSuffix(utcStr, "Z") {
		utcStr += "Z"
	}
	location, _ := time.LoadLocation("Asia/Shanghai")
	localTime, _ := time.ParseInLocation(timeLayoutStr, utcStr, location)
	return localTime.Format("2006-01-02 15:04:05")
}

func (alert Alert) FStatus() string {
	if alert.Status == "firing" {
		return "故障"
	}
	return "恢复"
}

func (alert Alert) StartsAt() string {
	startsAt := alert.ConvertUtcToLocal(alert.Alerts[0].StartsAt)
	return fmt.Sprintf("\n\n**故障时间:**\n\n\n- %s", startsAt)
}

func (alert Alert) EndsAt() string {
	endsAt := alert.ConvertUtcToLocal(alert.Alerts[0].EndsAt)
	return fmt.Sprintf("\n\n**恢复时间:**\n\n\n- %s", endsAt)
}

func (alert Alert) AlertsAt() string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	localTime := time.Now().In(loc)
	return fmt.Sprintf("\n\n**告警时间:**\n\n\n- %s", localTime.Format("2006-01-02 15:04:05"))
}

func (alert Alert) Alertname() string {
	return alert.Alerts[0].Labels.Alertname
}

func (alert Alert) Severity() string {
	return alert.Alerts[0].Labels.Severity
}

func (alert Alert) Message() [][]string {
	var messages [][]string
	i := 1
	for {
		if i*5 >= len(alert.Alerts) {
			start := (i - 1) * 5
			end := len(alert.Alerts)
			var message []string
			for _, v := range alert.Alerts[start:end] {
				message = append(message, fmt.Sprintf("\n- %s\n\n&nbsp;\n\n", v.Annotations.Message))
			}
			messages = append(messages, message)
			break
		} else {
			start := (i - 1) * 5
			end := i * 5
			var message []string
			for _, v := range alert.Alerts[start:end] {
				message = append(message, fmt.Sprintf("\n- %s\n\n&nbsp;\n\n", v.Annotations.Message))
			}
			messages = append(messages, message)
			i += 1
		}
	}
	return messages
}

func (alert Alert) SendDingTalk(webhook, secret string) {
	alertname, status, severity := alert.Alertname(), alert.FStatus(), alert.Severity()
	title := fmt.Sprintf("%s\n\n**[%s] [%s]**\n\n&nbsp;\n\n---", alertname, status, severity)
	for _, v := range alert.Message() {
		message := fmt.Sprintf("**告警内容:**\n\n%s\n\n---\n\n", strings.Join(v, ""))
		if status == "恢复" {
			message += alert.StartsAt() + "\n\n---" + alert.EndsAt()
		} else {
			message += alert.StartsAt() + "\n\n---" + alert.AlertsAt()
		}
		level := strings.ToLower(severity)
		if level == "p0" || level == "p1" || level == "p2" {
			log.Println(notify.SendDingTalkMarkdown(title, message, webhook, secret, true))
		} else {
			log.Println(notify.SendDingTalkMarkdown(title, message, webhook, secret, false))
		}
	}
}

func main() {
	route := gin.Default()
	route.POST("/alert-receiver/:name", func(c *gin.Context) {
		alert := Alert{}
		if err := c.ShouldBindJSON(&alert); err != nil {
			c.JSON(400, err.Error())
			return
		}
		name := c.Param("name")
		webhook, secret := GetDingTalkChannel(name)
		alert.SendDingTalk(webhook, secret)
		c.JSON(200, "ok")
	})
	route.Run(":8000")
}
