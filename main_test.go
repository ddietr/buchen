package main

import (
	"testing"
	"time"
)

func TestGetCurrentTimeFloat64(t *testing.T) {
	Now = func() time.Time {
		return time.Date(2021, 10, 13, 0, 0, 0, 0, time.UTC)
	}

	entry := createEntry()

	startedAt := Now().Add(time.Minute * -90)
	entry.StartedAt = &startedAt
	result, _ := getCurrentTimeFloat64(entry)
	if result != 1.5 {
		t.Fail()
	}

	startedAt = Now().Add(time.Minute * -15)
	entry.StartedAt = &startedAt
	result, _ = getCurrentTimeFloat64(entry)
	if result != 0.25 {
		t.Fail()
	}

	entry.Time = "0,5"
	result, _ = getCurrentTimeFloat64(entry)
	if result != 0.75 {
		t.Fail()
	}

	entry.StoppedAt = &startedAt
	result, _ = getCurrentTimeFloat64(entry)
	if result != 0.5 {
		t.Fail()
	}
}
