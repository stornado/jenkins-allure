package jenkinsallure

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type JenkinsJob struct {
	JobName        string   `yaml:"name"`
	EmailReceivers []string `yaml:"receivers,omitempty"`
}

type JenkinsConfig struct {
	Server   string       `yaml:"server"`
	Username string       `yaml:"username"`
	Password string       `yaml:"password,omitempty"`
	Jobs     []JenkinsJob `yaml:"jobs,omitempty"`
}

type EmailConfig struct {
	Sender    string   `yaml:"sender"`
	Host      string   `yaml:"host"`
	Port      int      `yaml:"port"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password,omitempty"`
	Receivers []string `yaml:"receivers,omitempty"`
}

type JenkinsAllureConfig struct {
	Jenkins JenkinsConfig `yaml:"jenkins"`
	Email   EmailConfig   `yaml:"email"`
}

func (config *JenkinsAllureConfig) ParseConfig(cfgFilepath string) error {
	data, err := ioutil.ReadFile(cfgFilepath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &config)
}
