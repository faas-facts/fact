package main

import (
	prefixed "github.com/chappjc/logrus-prefix"
	"github.com/faas-facts/fact/cmd"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var log *logrus.Entry

var (
	Build string
)

func init() {
	if Build == "" {
		Build = "Debug"
	}
	logger.Formatter = new(prefixed.TextFormatter)
	logger.SetLevel(logrus.DebugLevel)
	log = logger.WithFields(logrus.Fields{
		"prefix": "req-mon",
		"build":  Build,
	})
}

func main() {
	cmd.Execute(log)
}
