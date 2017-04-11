package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/op/go-logging"
	"gopkg.in/yaml.v2"
)

type Location struct {
	Desciption      string `yaml:"description"`
	FileNameRegexp  string `yaml:"filename_re"`
	fileNameRe      *regexp.Regexp
	URI             string `yaml:"uri"`
	SavePath        string `yaml:"save_path"`
	BashExecTimeout int64  `yaml:"bash_exec_timeout"`
	BashExec        string `yaml:"bash_exec"`
	Sync            bool   `yaml:"sync"`
}

type Config struct {
	Listen    string      `yaml:"listen"`
	LogFile   string      `yaml:"log_file"`
	LogLevel  string      `yaml:"log_level"`
	Locations []*Location `yaml:"locations"`
}

func loadConfig(configPath string) (*Config, error) {
	var config Config
	configYAML, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("Can't open config file. %s", err)
	}
	err = yaml.Unmarshal([]byte(configYAML), &config)
	if err != nil {
		return nil, fmt.Errorf("Can't parse config file [%s] [%s]", configPath, err)
	}
	for _, l := range config.Locations {
		if len(l.FileNameRegexp) > 0 {
			var err error
			if l.fileNameRe, err = regexp.Compile(l.FileNameRegexp); err != nil {
				return nil, fmt.Errorf("Bad regexp [%s]", l.FileNameRegexp)
			}
		} else {
			l.fileNameRe, _ = regexp.Compile(".*")
		}
		if l.BashExecTimeout < 1 {
			l.BashExecTimeout = 30
		}
	}
	return &config, nil
}

func newLog(logFile, level string) (*logging.Logger, error) {
	logLevel, err := logging.LogLevel(level)
	if err != nil {
		logLevel = logging.DEBUG
	}
	var logBackend *logging.LogBackend
	if logFile == "stdout" || logFile == "" {
		logBackend = logging.NewLogBackend(os.Stdout, "", 0)
		logBackend.Color = true
	} else {
		logFileName := logFile
		logFile, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("Can't open log file %s: %s", logFileName, err.Error())
		}
		logBackend = logging.NewLogBackend(logFile, "", 0)
	}
	logging.SetFormatter(logging.MustStringFormatter("%{time:2006-01-02 15:04:05}\t%{level}\t%{message}"))
	logger := logging.MustGetLogger("module")
	leveledLogBackend := logging.AddModuleLevel(logBackend)
	leveledLogBackend.SetLevel(logLevel, "module")
	logger.SetBackend(leveledLogBackend)
	return logger, nil
}
