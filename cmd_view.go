package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"regexp"
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
	headers := []string{"Date", "Time", "Project", "Description"}

	var data = make(map[string]DateEntry)
	keys := []string{}

	for i, entry := range entries {
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

			f := a + b
			s := strings.Replace(fmt.Sprintf("%.2f", f), ".", ",", 1)

			existing.Time = s
			existing.Description += ", " + entry.Description
			matched, _ := regexp.MatchString(
				"(?:^|[, ])" + existing.Project + "(,|$)",
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
		entry.Time = getCurrentTime(entry)
		row := []string{
			entry.Date,
			entry.Time,
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
