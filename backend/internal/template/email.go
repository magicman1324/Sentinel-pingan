package template

import (
	"bytes"
	"html/template"
	"time"

	"github.com/pingan/monitor-backend/internal/model"
)

const dailyReportHTML = `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>哨兵日报</title></head>
<body style="font-family:monospace;max-width:800px;margin:0 auto;">
  <h2>哨兵 · 每日巡检报告 — {{.Date}}</h2>
  <table border="1" cellpadding="6" cellspacing="0" width="100%">
    <tr style="background:#f44336;color:#fff;">
      <th>主机</th><th>指标</th><th>当前值</th><th>阈值</th><th>等级</th><th>时间</th>
    </tr>
    {{range .Alerts}}
    <tr>
      <td>{{.Hostname}}</td>
      <td>{{.Metric}}</td>
      <td>{{printf "%.2f" .Value}}</td>
      <td>{{printf "%.2f" .Threshold}}</td>
      <td>{{.Severity}}</td>
      <td>{{.CreatedAt.Format "15:04:05"}}</td>
    </tr>
    {{end}}
  </table>
  <p style="color:#666;">哨兵监控系统自动生成 · {{.Date}}</p>
</body>
</html>`

var dailyTmpl = template.Must(template.New("daily").Parse(dailyReportHTML))

type DailyData struct {
	Date   string
	Alerts []model.Alert
}

func RenderDailyReport(alerts []model.Alert) ([]byte, error) {
	var buf bytes.Buffer
	data := DailyData{
		Date:   time.Now().Format("2006-01-02"),
		Alerts: alerts,
	}
	if err := dailyTmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
