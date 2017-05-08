package cmd

import "github.com/op/go-logging"

const (
	normalFormat = `%{color}%{time:15:04:05.000} %{level:-7s} ▶ %{message}%{color:reset}`
	debugFormat  = `%{color}%{time:15:04:05.000} %{level} ¶ %{shortfile} ▶ %{message}%{color:reset}`
)

var (
	log = logging.MustGetLogger(appName)
)

func init() {
	logging.SetLevel(logging.INFO, appName)
	logging.SetFormatter(logging.MustStringFormatter(normalFormat))
}

func enableDebug(debug bool) {
	if debug {
		logging.SetLevel(logging.DEBUG, appName)
		logging.SetFormatter(logging.MustStringFormatter(debugFormat))
	}
}
