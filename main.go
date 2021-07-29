package main

import (
  "bytes"
  "strconv"
  "strings"
  "os"
  "fmt"
  "log"
  "time"
  "os/user"
  "io/ioutil"
  "gopkg.in/yaml.v2"
  "github.com/urfave/cli/v2"
)

type DateEntry struct {
  Date string
  Time string
  StartedAt *time.Time `yaml:"startedAt,omitempty"`
  StoppedAt *time.Time `yaml:"stoppedAt,omitempty"`
  Project string
  Description string
}

func main() {
  app := &cli.App{
    Name: "buchen",
    Usage: "time tracking from cli",
    Commands: []*cli.Command{
      CmdView(),
      CmdEdit(),
      CmdTimerStart(),
      CmdTimerStop(),
      CmdTimerNew(),
    },
    Action:  func(c *cli.Context) error {
      PrintView()
      return nil
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal("Error: ", err)
  }
}

func readConfig(filename string) ([]DateEntry, error) {
  buf, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }

  dec := yaml.NewDecoder(bytes.NewReader(buf))

  xs := []DateEntry{}
	var doc DateEntry
	for dec.Decode(&doc) == nil {
    xs = append(xs, doc)
	}

  return xs, nil
}

func findConfigFile() (string, error) {
  usr, err := user.Current()
  if err != nil {
    return "", err
  }

  filename := "buchen.yaml"
  _, err = os.Stat(usr.HomeDir + "/" + filename)
  if err == nil {
    return usr.HomeDir + "/" + filename, nil
  }

  return filename, nil
}

func serializeDateEntry(entry DateEntry) string {
  text := "---\n"
  text += fmt.Sprintf("date: %s\n", entry.Date)
  text += fmt.Sprintf("time: %s\n", entry.Time)
  if entry.StartedAt != nil {
    text += fmt.Sprintf("startedAt: %s\n", entry.StartedAt.Format(time.RFC3339))
  }

  if entry.StoppedAt != nil {
    text += fmt.Sprintf("stoppedAt: %s\n", entry.StoppedAt.Format(time.RFC3339))
  }

  text += fmt.Sprintf("project: %s\n", entry.Project)
  text += fmt.Sprintf("description: |\n")
  desc := strings.Split(entry.Description, "\n")
  for _, line := range desc {
    if line != "" {
      text += fmt.Sprintf("  %s\n", line)
    }
  }

  return text
}

func getCurrentTime(entry DateEntry) string {
  if entry.StartedAt == nil || entry.StoppedAt != nil {
    return entry.Time
  }

  now := time.Now()
  diff := now.Sub(*entry.StartedAt)
  entry.StoppedAt = &now
  timeStr := strings.Replace(entry.Time, ",", ".", 1)
  if s, err := strconv.ParseFloat(timeStr, 64); err == nil {
    t := s+(diff.Minutes()/60)
    return strings.Replace(fmt.Sprintf("%.2f", t), ".", ",", 1)
  }

  log.Fatal("Cannot calc current time")
  return ""
}
