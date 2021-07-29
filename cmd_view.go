package main

import (
  "fmt"
  "log"
  "github.com/urfave/cli/v2"
)

func CmdView() *cli.Command {
  return &cli.Command{
    Name: "view",
    Usage: "Print entries",
    Action:  func(c *cli.Context) error {
      PrintView()
      return nil
    },
  }
}

func PrintView() {
  filename, err := findConfigFile()
  if err != nil {
    log.Fatal(err)
  }

  entries, err := readConfig(filename)
  if err != nil {
    log.Fatal(err)
  }

  for _, entry := range entries {
    entry.Time = getCurrentTime(entry)
    fmt.Println(serializeDateEntry(entry))
  }
}
