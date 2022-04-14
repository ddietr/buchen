package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func cmdView() *cli.Command {
	return &cli.Command{
		Name:  "view",
		Usage: "Print entries",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "sum",
				Value: false,
				Usage: "sum time of day",
			},
		},
		Action: func(c *cli.Context) error {
			printTableView(c.Bool("sum"))
			return nil
		},
	}
}

func printTableView(sum bool) {
	filename, err := findConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	entries, err := readConfig(filename)
	if err != nil {
		log.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	headers := []string{"Date", "Hours", "Time", "From/To", "Project", "Description"}

	var data = make(map[string]DateEntry)
	keys := []string{}

	for i, entry := range entries {
		entry.Description = toInlineDescription(entry.Description)
		if entry.StartedAt != nil && entry.StoppedAt == nil {
			entry.Description = entry.Description + " 🏃"
		} else if entry.StartedAt != nil {
			entry.Description += entry.Description + " 💤"
		}

		if !sum {
			k := strconv.Itoa(i)
			keys = append(keys, k)
			data[k] = entry
			continue
		}

		if existing, ok := data[entry.Date]; ok {
			a, err := getCurrentTimeFloat64(existing)
			if err != nil {
				log.Fatal(err)
			}

			b, err := getCurrentTimeFloat64(entry)
			if err != nil {
				log.Fatal(err)
			}

			existing.Time = strings.Replace(fmt.Sprintf("%.2f", a+b), ".", ",", 1)
			existing.Description += ", " + entry.Description

			matched, _ := regexp.MatchString(
				"(?:^|[, ])"+existing.Project+"(,|$)",
				entry.Project,
			)
			if !matched {
				existing.Project += ", " + entry.Project
			}

			data[entry.Date] = existing
			continue
		}

		keys = append(keys, entry.Date)
		data[entry.Date] = entry
	}

	for _, k := range keys {
		entry := data[k]
		row := []string{
			entry.Date,
			getCurrentTime(entry),
			toDuration(entry),
			toFromTo(entry),
			entry.Project,
			entry.Description,
		}

		table.Rich(row, []tablewriter.Colors{
			{},
			{tablewriter.Bold},
			{},
			{},
		})
	}

	table.SetHeader(headers)
	table.Render()
}

func toFromTo(entry DateEntry) string {
	f, _ := getCurrentTimeFloat64(entry)
	floor := math.Floor(f)
	time := (f-floor)/10*6 + floor

	if f <= 6 {
		return fmt.Sprintf("8:00-%s", formatTime(time+8))
	}

	breakT := 0.3
	if f > 7.5 {
		breakT = 1
	}

	to := formatTime(8 + time + breakT)
	return fmt.Sprintf("8:00-%s 🍜%s", to, formatTime(breakT))
}

func toDuration(entry DateEntry) string {
	f, _ := getCurrentTimeFloat64(entry)
	floor := math.Floor(f)
	time := (f-floor)/10*6 + floor
	return formatTime(time)
}

func formatTime(f float64) string {
	return strings.Replace(fmt.Sprintf("%.2f", f), ".", ":", 1)
}
