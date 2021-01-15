package main

import (
	"encoding/base64"
	"flag"
	"log"
	"time"

	"github.com/bndr/gojenkins"
	"github.com/stornado/bazinga/pkg/urlparse"
	"github.com/stornado/jenkins-allure/pkg/jenkinsallure"
)

var cfgFilepath string

func init() {
	flag.StringVar(&cfgFilepath, "f", "", "the yaml config filepath")
	flag.StringVar(&cfgFilepath, "-config", "", "the yaml config filepath")
}

func main() {
	flag.Parse()

	if cfgFilepath == "" {
		flag.Usage()
	}

	var config jenkinsallure.JenkinsAllureConfig
	config.ParseConfig(cfgFilepath)

	jenkins := gojenkins.CreateJenkins(nil, config.Jenkins.Server, config.Jenkins.Username, config.Jenkins.Password)
	_, err := jenkins.Init()

	if err != nil {
		log.Fatalln(err)
	}

	for _, _job := range config.Jenkins.Jobs {
		job, err := jenkins.GetJob(_job.JobName)
		if err != nil {
			log.Fatalln(err)
		}

		detail := job.GetDetails()

		log.Println(detail.Color, detail.DisplayName, detail.FullDisplayName, detail.FullName, detail.Name, detail.NextBuildNumber, detail.URL, detail.Views)

		last, err := job.GetLastBuild()
		log.Println(last.GetResult(), last.GetDuration(), last.GetBuildNumber(), last.GetTimestamp(), last.GetRevision(), last.GetUrl())

		// Wait for build to finish
		for last.IsRunning() {
			time.Sleep(time.Duration(last.GetDuration()*2) * time.Millisecond)
			last.Poll()
		}

		allureURL, err := urlparse.URLJoin(last.GetUrl(), "allure")
		if err != nil {
			log.Fatalln(err)
		}
		summary, behaviors, trend, err := jenkinsallure.CaptureAllureResult(allureURL, _job.JobName)
		if err != nil {
			log.Fatalln(err)
		}

		report := jenkinsallure.JenkinsAllureReport{
			JobName:                       _job.JobName,
			JobReceivers:                  _job.EmailReceivers,
			LastBuildNumber:               last.GetBuildNumber(),
			LastBuildDuration:             last.GetDuration(),
			LastBuildResult:               last.GetResult(),
			LastBuildColor:                detail.Color,
			LastAllureReportURL:           allureURL,
			LastAllureSummaryScreenshot:   base64.StdEncoding.EncodeToString(summary),
			LastAllureBehaviorsScreenshot: base64.StdEncoding.EncodeToString(behaviors),
			LastAllureTrendScreenshot:     base64.StdEncoding.EncodeToString(trend),
		}
		err = report.Report(config.Email)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("build number %d with result: %v %v %v\n", last.GetBuildNumber(), last.GetResult(), last.GetUrl(), last.GetDuration())
	}
}
