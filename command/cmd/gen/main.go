package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
	"text/template"

	"github.com/sethvargo/ratchet/command"
)

func main() {
	if err := realMain(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func realMain() error {
	topLevelHelp := buildTopLevelHelp()

	tmpl, err := template.New("gen").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var w bytes.Buffer
	if err := tmpl.Execute(&w, &Params{
		TopLevelHelp: topLevelHelp,
	}); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	root := folderRoot()
	pth := path.Join(root, "command_gen.go")
	if err := os.WriteFile(pth, w.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write generated file: %w", err)
	}

	return nil
}

func folderRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	for i := 0; i < 3; i++ {
		filename = path.Dir(filename)
	}
	return filename
}

func buildTopLevelHelp() string {
	longest := 0
	names := make([]string, 0, len(command.Commands))
	for name := range command.Commands {
		names = append(names, name)
		if l := len(name); l > longest {
			longest = l
		}
	}
	sort.Strings(names)

	var w strings.Builder
	fmt.Fprint(&w, "Usage: ratchet COMMAND\n\n")
	for _, name := range names {
		cmd := command.Commands[name]

		padding := longest - len(name) + 4
		fmt.Fprint(&w, "  ")
		fmt.Fprint(&w, name)
		for i := 0; i < padding; i++ {
			fmt.Fprint(&w, " ")
		}
		fmt.Fprint(&w, cmd.Desc())
		fmt.Fprintln(&w)
	}

	return w.String()
}

type Params struct {
	TopLevelHelp string
}

const tmpl = `// Code generated by "command/cmd/gen". DO NOT EDIT.

package command

const topLevelHelp = ` + "`" + `{{.TopLevelHelp}}` + "`" + `
`
