package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os/user"
)

func cmdInit() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Init",
		Action: func(c *cli.Context) error {
			initialize()
			return nil
		},
	}
}

func initialize() {
	filename, err := findConfigFile()
	if err == nil {
		log.Fatal("Already initialized in " + filename)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dest := usr.HomeDir + "/" + "buchen.yaml"

	content := serializeDateEntry(createEntry())
	bytes := []byte(content)

	ioutil.WriteFile(dest, bytes, 0644)
	fmt.Println("Successfully initialized " + dest)
}
