package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

func cmdAdd() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "Add description to current entry",
		Action: func(c *cli.Context) error {
			if c.Args().Len() > 0 {
				text := c.Args().Get(0)
				addDesc(text)
				return nil
			}

			reader := bufio.NewReader(c.App.Reader)
			text, _ := reader.ReadString('\n')
			text = strings.TrimSuffix(text, "\n")
			text = strings.TrimSpace(text)

			if len(text) == 0 {
				fmt.Println("Error: Empty text. Usage: `buchen add text`")
				os.Exit(1)
			}

			addDesc(text)
			return nil
		},
	}
}

func addDesc(text string) {
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
		fmt.Println("Error: Cound not find current entry")
		os.Exit(1)
	}

	desc := strings.TrimSuffix(entry.Description, "\n")
	desc = strings.TrimSpace(desc)
	if desc == "" || desc == "..." {
		desc = text
	} else {
		desc += "\n" + text
	}

	entry.Description = desc
	fmt.Println("New description:", toInlineDescription(entry.Description))
	entries[index] = *entry
	writeToFile(filename, entries)
}
