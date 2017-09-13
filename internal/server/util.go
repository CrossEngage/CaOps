package server

import (
	"time"
)

func stringListToMapKeys(list []string) map[string]bool {
	ret := make(map[string]bool)
	for _, item := range list {
		ret[item] = true
	}
	return ret
}

func getNextRoundedTimeWithin(from time.Time, duration time.Duration) time.Time {
	return from.Truncate(duration).Add(duration)
}

func getErrStr(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
