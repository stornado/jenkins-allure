package jenkinsallure

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"text/template"

	"gopkg.in/gomail.v2"
)

type JenkinsBuild struct {
	BuildNumber int
	BuildStatus string
}

type JenkinsAllureReport struct {
	JobName             string
	JobReceivers        []string
	LastBuildNumber     int64
	LastBuildDuration   int64
	LastBuildResult     string
	LastBuildColor      string
	LastAllureReportURL string
	LastAllureSummary   string
	LastAllureBehaviors string
	LastAllureTrend     string
}

func (report JenkinsAllureReport) ToHtml() (string, error) {
	tpl := `
<base href="{{ .LastAllureReportURL }}">
<style>
.bar{border-radius:3px;background:#eee;display:inline-flex;width:100%}.bar,.bar__fill{height:18px;overflow:hidden}.bar__fill{display:inline-block;background:#04b;text-align:center;color:#fff;font-size:12px;line-height:18px}.bar__fill_status_failed{background:#fd5a3e}.bar__fill_status_broken{background:#ffd050}.bar__fill_status_passed{background:#97cc64}.bar__fill_status_skipped{background:#aaa}.bar__fill_status_unknown{background:#d35ebe}.island{background:#fff;border:1px solid #e5e5e5;padding:15px 15px 0}.splash{margin:auto;text-align:center}.splash__title{font-size:3.5em;line-height:1}.splash__subtitle{color:#999}.widget{margin-bottom:15px;position:relative}.widget_ghost{border:1px dashed #e5e5e5;box-shadow:none;min-height:50px}.widget_ghost>*{display:none}.widget__title{margin-top:0;margin-bottom:15px;font-weight:lighter;text-transform:uppercase}.widget__subtitle{color:#999;font-size:16px;text-transform:none}.widget__noitems{font-size:16px;text-align:center;padding:10px 15px;line-height:1.5em}.widget__flex-line{display:flex}.widget__column{width:50%}.widget__handle{display:none;position:absolute;right:15px;top:15px;color:#999;cursor:move;cursor:-webkit-grab;cursor:grab}.widget__handle:active{cursor:-webkit-grabbing;cursor:grabbing}.widget:hover .widget__handle{display:block}.widget__table{border-top:1px solid #eceff1;margin:0 -15px;word-break:break-all}.widget__table .table__row:last-child{border-bottom:0}.draggable-icon{position:absolute;width:10px;height:15px;right:0;background-image:url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI1IiBoZWlnaHQ9IjUiPjxwYXRoIGZpbGw9IiNjY2MiIGQ9Ik0wIDBoM3YzSDB6Ii8+PC9zdmc+);z-index:1}body{color:#000;font-family:Helvetica,Arial,sans-serif}*{box-sizing:border-box}body,html{height:100%;font-size:14px}#content{height:100%;min-height:100%;overflow:hidden;display:flex;flex-direction:column}#content .app{flex:1}#content>.spinner{position:absolute;top:50%;left:50%;-webkit-transform:translate(-50%,-50%);transform:translate(-50%,-50%)}.view{padding:0 15px;margin:0 auto}.view_narrow{max-width:1100px}.view-small{max-width:300px}.view-medium{max-width:600px}.view-large{max-width:1200px}.clickable{cursor:pointer}.long-line{word-break:break-word}.line-nobreak{white-space:nowrap}.preformated-text{white-space:pre-wrap;word-wrap:break-word}.line-ellipsis{white-space:nowrap;text-overflow:ellipsis;overflow:hidden}.app{background:#fff;display:flex}.app__nav{padding:0}.app__content{position:relative;flex:1;overflow:auto}.app__header{background:#fff}.error-splash{padding:10px;text-align:center}.table__head,.table__row{display:flex}.table__row{border-bottom:1px solid #eceff1;text-decoration:none;color:#000}.table__head{border-bottom:1px solid #e5e5e5;font-weight:700}.table__col{line-height:1.5em;padding:10px 15px;flex:1}.table__col_center{text-align:center}.table__col_right{text-align:right;justify-content:flex-end}.table__col_sortable{cursor:pointer;display:flex}.table__col_sortable>span{overflow:hidden;padding-right:5px}.table__col_sortable:after{flex-shrink:0;vertical-align:middle;content:" ";display:inline-block;width:12px;height:18px;background:url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI1MTIiIGhlaWdodD0iNTEyIj48cGF0aCBmaWxsPSIjYTVhNWE1IiBkPSJNMjU2IDUwbDEzMi4wMzQgMTc2SDEyMy45NjZMMjU2IDUwem0xMzIuMDM0IDIzNkgxMjMuOTY2TDI1NiA0NjJsMTMyLjAzNC0xNzZ6Ii8+PC9zdmc+) 50% no-repeat;background-size:contain}.table__col_sorted_down:after{background-image:url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI1MTIiIGhlaWdodD0iNTEyIj48cGF0aCBmaWxsPSIjYTVhNWE1IiBkPSJNMjU2IDUwbDEzMi4wMzQgMTc2SDEyMy45NjZMMjU2IDUweiIvPjxwYXRoIGQ9Ik0zODguMDM0IDI4NkgxMjMuOTY2TDI1NiA0NjJsMTMyLjAzNC0xNzZ6Ii8+PC9zdmc+)}.table__col_sorted_up:after{background-image:url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI1MTIiIGhlaWdodD0iNTEyIj48cGF0aCBkPSJNMjU2IDUwbDEzMi4wMzQgMTc2SDEyMy45NjZMMjU2IDUweiIvPjxwYXRoIGZpbGw9IiNhNWE1YTUiIGQ9Ik0zODguMDM0IDI4NkgxMjMuOTY2TDI1NiA0NjJsMTMyLjAzNC0xNzZ6Ii8+PC9zdmc+)}.table_hover .table__row:not([disabled]):hover{background:#e4edfe}.table .table__row_summary{background:#f6f5f3}.table .table__row_active{background:#fffacd}.chart{margin-right:15px;margin-bottom:15px;margin-top:15px;position:relative;flex:1}.chart__title{margin-top:0;margin-bottom:15px;font-size:18px}.chart__caption{text-anchor:middle}.chart__body>div{padding-top:50%;position:relative}.chart__svg{position:absolute;top:0;left:0;width:100%;height:100%}.chart__legend{position:absolute;height:50%;top:10%;right:10%;display:grid}.chart__legend-icon{position:relative;top:-1px;border-radius:3px;display:inline-block;vertical-align:middle;width:20px;height:16px;margin-right:5px}.chart__legend-icon_status_failed{background:#fd5a3e}.chart__legend-icon_status_broken{background:#ffd050}.chart__legend-icon_status_passed{background:#97cc64}.chart__legend-icon_status_skipped{background:#aaa}.chart__legend-icon_status_unknown{background:#d35ebe}.chart__bar{shape-rendering:crispEdges;fill:#4682b4}.chart__arc{stroke:#fff;stop-opacity:0}.chart__fill_status_failed{fill:#fd5a3e}.chart__fill_status_broken{fill:#ffd050}.chart__fill_status_passed{fill:#97cc64}.chart__fill_status_skipped{fill:#aaa}.chart__fill_status_unknown{fill:#d35ebe}.chart__axis line,.chart__axis path{shape-rendering:crispEdges;stroke:#000;fill:none}.widgets-grid{position:absolute;top:0;left:0;right:0;bottom:0;padding:15px 15px 0;max-height:100%;overflow:auto;display:flex}.widgets-grid__col{flex:1 0 0%;width:50%;min-width:300px}.widgets-grid__col+.widgets-grid__col{flex:1 0 0%;padding-left:15px}.summary-widget{padding:1em 0}.summary-widget__stats{padding:2em 0}.summary-widget__chart>div{height:100%;padding-bottom:12px}.history-trend__chart>div,.summary-widget__chart>div{position:relative;padding-top:50%}
</style>

<h1>{{ .JobName }}</h1>
<p>构建结果: <span class="{{ .LastBuildColor }}">{{ .LastBuildResult }}</span> <a href="{{ .LastAllureReportURL }}">点击查看详情: {{ .LastBuildNumber }}</a></p>

<div id="content">
<div class="app">
<div class="app__content">
<div class="widgets-grid">
<div class="widgets-grid__col">
{{ .LastAllureSummary }}

{{ .LastAllureBehaviors }}
</div>
</div>
</div>
</div>
</div>
`

	htmlTpl, err := template.New("report").Parse(tpl)
	if err != nil {
		return "", err
	}

	var doc bytes.Buffer
	err = htmlTpl.Execute(&doc, report)
	if err != nil {
		return "", err
	}

	return doc.String(), nil
}

func (report JenkinsAllureReport) Report(email EmailConfig) error {
	doc, err := report.ToHtml()
	if err != nil {
		return err
	}

	d := gomail.NewDialer(email.Host, email.Port, email.Username, email.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", email.Sender)
	m.SetHeader("To", report.JobReceivers...)
	m.SetHeader("Subject", fmt.Sprintf("【Allure】Jenkins build %s: %s # %d", report.LastBuildResult, report.JobName, report.LastBuildNumber))
	m.SetBody("text/html", doc)
	err = d.DialAndSend(m)
	if err != nil {
		return err
	}

	return nil
}
