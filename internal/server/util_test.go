package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetNextRoundedTimeWithin(t *testing.T) {
	format := "2006-01-02 15:04:05"
	start, error := time.Parse(format, "2017-09-13 11:06:15")
	assert.Nil(t, error)
	next := getNextRoundedTimeWithin(start, 5*time.Second)
	assert.Equal(t, "2017-09-13 11:06:20", next.Format(format))
	next = getNextRoundedTimeWithin(start, 5*time.Minute)
	assert.Equal(t, "2017-09-13 11:10:00", next.Format(format))
	next = getNextRoundedTimeWithin(start, 1*time.Minute)
	assert.Equal(t, "2017-09-13 11:07:00", next.Format(format))
}
