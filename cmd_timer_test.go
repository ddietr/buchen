package main

import (
	"fmt"
	"testing"
	"time"
)

func TestStopCurrentEntry(t *testing.T) {
	Now = func() time.Time {
		return time.Date(2021, 10, 19, 0, 0, 0, 0, time.UTC)
	}

	startedAt := Now().Add(time.Minute * -90)

	entry1 := createEntry()
	entry1.Date = "2021.10.17"
	entry2 := createEntry()
	entry2.Date = "2021.10.18"
	current := createEntry()
	current.Date = "2021.10.19"
	current.StartedAt = &startedAt
	current.Time = "0,75"

	entries := []DateEntry{
		entry1,
		entry2,
		current,
	}

	stopCurrentEntry(entries, 2, &current)
	stopCurrentEntry(entries, 2, &current)

	if entries[2].Date != "2021.10.19" {
		t.Fail()
	}

	if entries[2].StoppedAt == nil {
		t.Fail()
	}

	if entries[2].Time != "2,25" {
		fmt.Println(entries[2].Time)
		t.Fail()
	}
}
