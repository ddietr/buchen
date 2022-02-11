package main

import (
	"fmt"
	"log"
	"os"
	"github.com/urfave/cli/v2"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
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
			startNewTimer(c.Args().Get(0))
			return nil
		},
	}
}

func getCurrentEntry(entries []DateEntry) (int, *DateEntry) {
	if len(entries) == 0 {
		return -1, nil
	}

	for i, entry := range entries {
		if entry.StartedAt != nil {
			return i, &entry
		}
	}

	lastIndex := len(entries)-1
	last := entries[lastIndex]

	return lastIndex, &last
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

func startNewTimer(desc string) {
	filename, err := findConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	entries, err := readConfig(filename)
	if err != nil {
		log.Fatal(err)
	}

	index, entry := getCurrentEntry(entries)
	if entry == nil {
		newEntry := createEntry()
		if desc != "" {
			newEntry.Description = desc
		}
		entries = append(entries, newEntry)
		writeToFile(filename, entries)
		fmt.Printf("🏃 Start \"%s\"\n", toInlineDescription(newEntry.Description))
		return
	}

	fmt.Printf("💤 Stopped \"%s\" at %s\n", toInlineDescription(entry.Description), entry.Time)
	stopCurrentEntry(entries, index, entry)
	entry.StartedAt = nil
	entry.StoppedAt = nil
	entries[index] = *entry
	newEntry := createEntry()
	if desc != "" {
		newEntry.Description = desc
	}
	entries = append(entries, newEntry)
	writeToFile(filename, entries)
	fmt.Printf("🏃 Start \"%s\"\n", toInlineDescription(newEntry.Description))
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
	today := timeNow.Format("02.01.2006")
	index, entry := getCurrentEntry(entries)

	// No current entry found
	if entry == nil {
		newEntry := createEntry()
		entries = append(entries, newEntry)
		writeToFile(filename, entries)
		fmt.Println("🏃 Start timer for a new day.")
		return
	}

	// Start new timer if current is from yesterday
	if today != entry.Date {
		if entry.StartedAt != nil {
			stopCurrentEntry(entries, index, entry)
			entry.StartedAt = nil
			entry.StoppedAt = nil
			entries[index] = *entry
		}

		newEntry := createEntry()
		entries = append(entries, newEntry)
		writeToFile(filename, entries)
		fmt.Println("🏃 Start timer for a new day.")
		return
	}

	entriesTodayCount := 0
	for _, e := range entries {
		if today == e.Date { entriesTodayCount++ }
	}

	if entriesTodayCount == 1 {
		if entry.StartedAt != nil && entry.StoppedAt == nil {
			fmt.Println("Error: Timer already started")
			os.Exit(1)
		}

		entry.StoppedAt = nil
		entry.StartedAt = &timeNow
		entries[index] = *entry
		writeToFile(filename, entries)
		fmt.Printf("🏃 Restart \"%s\" at %s\n", toInlineDescription(entry.Description), entry.Time)
		return
	}

	switchEntryPrompt(entries)
	writeToFile(filename, entries)
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

	index, entry := getCurrentEntry(entries)
	if entry == nil {
		log.Fatal("Timer not started")
	}

	if entry.StoppedAt != nil {
		fmt.Println("Error: Timer already stopped.")
		os.Exit(1)
	}

	stopCurrentEntry(entries, index, entry)
	fmt.Printf("💤 Stopped \"%s\" at %s\n", toInlineDescription(entry.Description), entry.Time)

	writeToFile(filename, entries)
}


func stopCurrentEntry(entries []DateEntry, index int, current *DateEntry) {
	timeNow := Now()
	current.Time = getCurrentTime(*current)
	current.StoppedAt = &timeNow
	entries[index] = *current
}

func switchEntryPrompt(entries []DateEntry) {
	today := Now().Format("02.01.2006")
	options := []string{}
	tasks := []DateEntry{}

	for _, entry := range entries {
		if today != entry.Date {
			continue
		}

		text := toInlineDescription(entry.Description)
		if entry.StartedAt != nil && entry.StoppedAt == nil {
			text = text + " 🏃"
		} else if entry.StartedAt != nil {
			text = text + " 💤"
		}

		options = append(options, text)
		tasks = append(tasks, entry)
	}

	index := 0
	prompt := &survey.Select{
		Message: "Choose task:",
		Help: "",
		Options: options,
	}

	err := survey.AskOne(prompt, &index, survey.WithValidator(survey.Required))
	if err != nil {
		if err == terminal.InterruptErr {
			fmt.Println("Switch aborted.")
			os.Exit(0)
		}

		log.Fatal(err)
	}

	timeNow := Now()
	curIndex, cur := getCurrentEntry(entries)
	selected := tasks[index]
	selectedIndex := findEntryIndex(entries, selected)
	if selectedIndex == -1 {
		fmt.Println("Error: Selected entry not found.")
		os.Exit(1)
	}

	if *cur == selected {
		if selected.StoppedAt != nil || selected.StartedAt == nil {
			fmt.Printf("🏃 Restart \"%s\" at %s\n", toInlineDescription(selected.Description), selected.Time)
			selected.StartedAt = &timeNow
			selected.StoppedAt = nil
			entries[curIndex] = selected
		} else {
			fmt.Println("Error: Entry already started.")
			os.Exit(1)
		}
	} else {
		if cur != nil {
			cur.Time = getCurrentTime(*cur)
			fmt.Printf("💤 Stopped \"%s\" at %s\n", toInlineDescription(cur.Description), cur.Time)
			cur.StoppedAt = nil
			cur.StartedAt = nil
			entries[curIndex] = *cur
		}

		fmt.Printf("🏃 Start \"%s\" at %s\n", toInlineDescription(selected.Description), selected.Time)
		selected.StartedAt = &timeNow
		selected.StoppedAt = nil
		entries[selectedIndex] = selected
	}
}

func findEntryIndex(entries []DateEntry, entry DateEntry) int {
	for k, v := range entries {
		if entry == v {
			return k
		}
	}

	return -1
}
