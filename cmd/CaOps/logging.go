package main

import (
	"log/syslog"

	"github.com/sirupsen/logrus"
	lrSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

func init() {
	// TODO make level and syslog server configurable
	logrus.SetLevel(logrus.DebugLevel)
	hook, err := lrSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
	if err != nil {
		logrus.Error("Unable to connect to local syslog daemon")
	} else {
		logrus.AddHook(hook)
	}
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true, DisableTimestamp: true})
}
