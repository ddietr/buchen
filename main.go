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
      CmdInit(),
    },
    Action:  func(c *cli.Context) error {
      printTableView(true)
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
  path := usr.HomeDir + "/" + filename
  _, err = os.Stat(path)

  if err != nil {
    return "", fmt.Errorf(
      "Could not find configuration in %s. Run 'buchen init' to create the configuration",
      path,
    )
  }

  return path, nil
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
  timeF, err := getCurrentTimeFloat64(entry)

  if err != nil {
    log.Fatal("Cannot calc current time")
  }

  return strings.Replace(fmt.Sprintf("%.2f", timeF), ".", ",", 1)
}

func getCurrentTimeFloat64(entry DateEntry) (float64, error) {
  timeStr := strings.Replace(entry.Time, ",", ".", 1)
  s, err := strconv.ParseFloat(timeStr, 64)

  if err != nil {
    return 0.0, err
  }

  if entry.StartedAt != nil && entry.StoppedAt == nil {
    now := time.Now()
    diff := now.Sub(*entry.StartedAt)
    t := s+(diff.Minutes()/60)
    return t, nil
  }

  return s, nil
}

func createEntry() DateEntry {
  now := time.Now()
  return DateEntry{
    Date: now.Format("02 Jan"),
    Time: "0,0",
    StartedAt: &now,
    Project: "...",
    Description: "- ...",
  }
}

