package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"strings"
)

func cmdCsv() *cli.Command {
	return &cli.Command{
		Name:  "csv",
		Usage: "Print entries as csv",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "project",
				Aliases: []string{"p"},
				Value:   "",
				Usage:   "Filter by project, default all",
			},
		},
		Action: func(c *cli.Context) error {
			printCsv(c.String("project"))
			return nil
		},
	}
}

func printCsv(project string) {
	filename, err := findConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	entries, err := readConfig(filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Datum;Beschreibung;Aufwand")
	for _, entry := range entries {
		if (project != "") {
			if (entry.Project != project) {
				continue
			}
		}

		fmt.Print(entry.Date)
		fmt.Print(";")
		fmt.Print(strings.Trim(strings.ReplaceAll(entry.Description, "\n", ", "), ", "))
		fmt.Print(";")
		fmt.Print(getCurrentTime(entry))
		fmt.Println("")
	}
}
