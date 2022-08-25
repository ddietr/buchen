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

	total := 0.0
	for _, k := range keys {
		entry := data[k]
		f, _ := getCurrentTimeFloat64(entry)
		total += f
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
	fmt.Println("Total hours:", strings.Replace(fmt.Sprintf("%.2f", total), ".", ",", 1))
	fmt.Println("Remaining:", strings.Replace(fmt.Sprintf("%.2f", 40-total), ".", ",", 1))
}

func toFromTo(entry DateEntry) string {
	f, _ := getCurrentTimeFloat64(entry)

	if f <= 6 {
		return fmt.Sprintf("8:00-%s", formatTime(calcDuration(f+8)))
	}

	breakT := 0.5
	if f > 8 {
		breakT = 1
	}

	return fmt.Sprintf(
		"8:00-%s 🍜%s",
		formatTime(calcDuration(8+f+breakT)),
		formatTime(calcDuration(breakT)),
	)
}

func calcDuration(f float64) float64 {
	floor := math.Floor(f)
	return (f-floor)/10*6 + floor
}

func toDuration(entry DateEntry) string {
	f, _ := getCurrentTimeFloat64(entry)
	return formatTime(calcDuration(f))
}

func formatTime(f float64) string {
	return strings.Replace(fmt.Sprintf("%.2f", f), ".", ":", 1)
}
