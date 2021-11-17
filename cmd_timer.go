package main

import (
	"fmt"
	"github.com/jinzhu/now"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func cmdTimerStart() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "Start the timer",
		Action: func(c *cli.Context) error {
			startTimer()
			return nil
		},
	}
}

func cmdTimerStop() *cli.Command {
	return &cli.Command{
		Name:  "stop",
		Usage: "Stop the timer",
		Action: func(c *cli.Context) error {
			stopTimer()
			return nil
		},
	}
}

func cmdTimerNew() *cli.Command {
	return &cli.Command{
		Name:  "new",
		Usage: "Create new date entry",
		Action: func(c *cli.Context) error {
			startNewTimer()
			return nil
		},
	}
}

func getCurrentEntry(entries []DateEntry) *DateEntry {
	if len(entries) == 0 {
		return nil
	}

	last := entries[len(entries)-1]

	if last.StartedAt != nil {
		return &last
	}

	return nil
}

func writeToFile(filename string, entries []DateEntry) {
	text := ""
	for i, entry := range entries {
		if i == len(entries)-1 {
			text += serializeDateEntry(entry)
		} else {
			text += fmt.Sprintf("%s\n", serializeDateEntry(entry))
		}
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	if _, err := f.WriteString(text); err != nil {
		log.Fatal(err)
	}
}

func startNewTimer() {
	filename, err := findConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	entries, err := readConfig(filename)
	if err != nil {
		log.Fatal(err)
	}

	entry := getCurrentEntry(entries)
	if entry == nil {
		entries = append(entries, createEntry())
		writeToFile(filename, entries)
		return
	}

	stopCurrentEntry(entries, entry)
	entry.StartedAt = nil
	entry.StoppedAt = nil
	entries[len(entries)-1] = *entry
	entries = append(entries, createEntry())
	writeToFile(filename, entries)
	fmt.Println("Created new entry and started timer")
}

func startTimer() {
	filename, err := findConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	entries, err := readConfig(filename)
	if err != nil {
		log.Fatal(err)
	}

	timeNow := Now()
	entry := getCurrentEntry(entries)

	// No current entry found
	if entry == nil {
		newEntry := createEntry()
		entries = append(entries, newEntry)
		writeToFile(filename, entries)
		fmt.Println("Timer started. Current time:", newEntry.Time)
		return
	}

	// Current entry not stopped
	if entry.StoppedAt == nil {
		fmt.Println("Timer already started")
		os.Exit(1)
	}

	if timeNow.After(now.With(*entry.StartedAt).EndOfDay()) {
		stopCurrentEntry(entries, entry)
		entry.StartedAt = nil
		entry.StoppedAt = nil
		entries[len(entries)-1] = *entry
		newEntry := createEntry()
		entries = append(entries, newEntry)
		entry = &newEntry
	} else {
		entry.StoppedAt = nil
		entry.StartedAt = &timeNow
		entries[len(entries)-1] = *entry
	}

	writeToFile(filename, entries)
	fmt.Println("Timer started. Current time:", entry.Time)
}

func stopTimer() {
	filename, err := findConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	entries, err := readConfig(filename)
	if err != nil {
		log.Fatal(err)
	}

	entry := getCurrentEntry(entries)
	if entry == nil {
		log.Fatal("Timer not started")
	}

	if entry.StoppedAt != nil {
		fmt.Println("Timer already stopped")
		os.Exit(1)
	}

	stopCurrentEntry(entries, entry)
	fmt.Println("Timer stopped. Current time:", entry.Time)

	writeToFile(filename, entries)
}

func stopCurrentEntry(entries []DateEntry, current *DateEntry) {
	timeNow := Now()
	current.Time = getCurrentTime(*current)
	current.StoppedAt = &timeNow
	entries[len(entries)-1] = *current
}
