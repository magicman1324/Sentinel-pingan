package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

)

// DingTalkCard renders a Markdown alert card and POSTs to the webhook URL.
type DingTalkCard struct {
	WebhookURL string
	HTTPClient *http.Client
}

func NewDingTalkCard(webhookURL string) *DingTalkCard {
	return &DingTalkCard{
		WebhookURL: webhookURL,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}
}

type dingtalkMessage struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
}

func (c *DingTalkCard) Send(hostname, severity, metric, message string, value, threshold float64, alertID int64) error {
	title := fmt.Sprintf("%s — %s", severity, hostname)

	text := fmt.Sprintf(`### %s

| 字段 | 值 |
|------|-----|
| **主机** | %s |
| **指标** | %s |
| **当前值** | %.2f |
| **阈值** | %.2f |
| **规则** | %s |
| **时间** | %s |

> [ 查看Dashboard](%s) | [ 一键处理](%s)`,
		title,
		hostname,
		metric,
		value,
		threshold,
		message,
		time.Now().Format("15:04:05"),
		"https://monitor.pingan.com/dashboard",
		"https://monitor.pingan.com/ack/"+fmt.Sprintf("%d", alertID),
	)

	msg := dingtalkMessage{MsgType: "markdown"}
	msg.Markdown.Title = title
	msg.Markdown.Text = text

	body, _ := json.Marshal(msg)
	resp, err := c.HTTPClient.Post(c.WebhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
