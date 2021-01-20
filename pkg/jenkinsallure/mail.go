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
	JobName                     string
	JobReceivers                []string
	LastBuildNumber             int64
	LastBuildDuration           int64
	LastBuildResult             string
	LastBuildColor              string
	LastAllureReportURL         string
	LastAllureSummary           string
	LastAllureBehaviors         string
	LastAllureTrend             string
	LastAllureSummarySnapshot   string
	LastAllureBehaviorsSnapshot string
	LastAllureTrendSnapshot     string
}

func (report JenkinsAllureReport) ToHtml() (string, error) {
	tpl := `
<h1>{{ .JobName }}</h1>
<p>构建结果: <span class="{{ .LastBuildColor }}">{{ .LastBuildResult }}</span> <a href="{{ .LastAllureReportURL }}">点击查看详情: {{ .LastBuildNumber }}</a></p>

<h2>ALLURE REPORT SUMMARY</h2>
<a href="{{- .LastAllureReportURL -}}#graph"><img src="data:image/jpeg;base64,{{- .LastAllureSummarySnapshot -}}"/></a>

<h2>FEATURES BY STORIES</h2>
<a href="{{- .LastAllureReportURL -}}#behaviors"><img src="data:image/jpeg;base64,{{- .LastAllureBehaviorsSnapshot -}}"/></a>
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
