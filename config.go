package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SourceDir string `yaml:"sourceDir"`
	TargetDir string `yaml:"targetDir"`
	Interval  int    `yaml:"interval"`
}

func getConfig() (config Config) {
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func getLogger() *log.Logger {
	file, err := os.OpenFile("run.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	l := &log.Logger{}
	l.SetOutput(file)
	l.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return l
}
