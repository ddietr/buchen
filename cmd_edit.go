package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
)

func cmdEdit() *cli.Command {
	return &cli.Command{
		Name:    "edit",
		Aliases: []string{"e"},
		Usage:   "Open entries file in EDITOR",
		Action: func(c *cli.Context) error {
			openEditor()
			return nil
		},
	}
}

func openEditor() {
	editorVar := os.Getenv("EDITOR")
	if "" == editorVar {
		editorVar = "vi"
	}

	editor, err := exec.LookPath(editorVar)
	if err != nil {
		log.Fatal(err)
	}

	filename, err := findConfigFile()
	cmd := exec.Command(editor, "+99999", filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
