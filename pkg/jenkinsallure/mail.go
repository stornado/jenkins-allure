package jenkinsallure

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"

	"gopkg.in/gomail.v2"
)

type JenkinsBuild struct {
	BuildNumber int
	BuildStatus string
}

type JenkinsAllureReport struct {
	JobName                       string
	JobReceivers                  []string
	LastBuildNumber               int64
	LastBuildDuration             int64
	LastBuildResult               string
	LastBuildColor                string
	LastAllureReportURL           string
	LastAllureSummaryScreenshot   string
	LastAllureBehaviorsScreenshot string
	LastAllureTrendScreenshot     string
}

func (report JenkinsAllureReport) ToHtml() (string, error) {
	tpl := `
<style>
.icon-health {
	width: 48px; height: 48px
}
</style>

<h1>{{ .JobName }}</h1>
<p>上次构建结果: <a href="{{ .LastAllureReportURL }}"><span class="{{ .LastBuildColor }}">{{ .LastBuildResult }}</span></a><img src="http://172.20.66.36:8080/static/80232935/images/48x48/{{- .LastBuildColor -}}.png"/></p>

<div id="allure">
<div id="summary">
<img src="data:image/png;base64,{{- .LastAllureSummaryScreenshot -}}" title="summary"/>
</div>
<div id="trend">
<img src="data:image/png;base64,{{- .LastAllureTrendScreenshot -}}" title="trend"/>
</div>
<div id="behaviors">
<img src="data:image/png;base64,{{- .LastAllureBehaviorsScreenshot -}}" title="behaviors"/>
</div>
</div>
`

	htmlTpl, err := template.New("report").Parse(tpl)
	if err != nil {
		return "", nil
	}

	var doc bytes.Buffer
	err = htmlTpl.Execute(&doc, report)
	if err != nil {
		return "", nil
	}

	return doc.String(), nil
}

func (report JenkinsAllureReport) Report(email EmailConfig) error {
	doc, err := report.ToHtml()
	if err != nil {
		log.Println(err)
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
