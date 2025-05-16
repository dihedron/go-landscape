package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/dihedron/landscape/metadata"
	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
)

var (
	red     = color.New(color.FgRed).SprintfFunc()
	green   = color.New(color.FgGreen).SprintfFunc()
	yellow  = color.New(color.FgYellow).SprintfFunc()
	magenta = color.New(color.FgMagenta).SprintfFunc()
	cyan    = color.New(color.FgCyan).SprintfFunc()
	blue    = color.New(color.FgBlue).SprintfFunc()
)

func main() {

	if len(os.Args) == 2 && (os.Args[1] == "version" || os.Args[1] == "--version") {
		metadata.Print(os.Stdout)
		os.Exit(0)
	} else if len(os.Args) == 3 && os.Args[1] == "version" && (os.Args[2] == "--verbose" || os.Args[2] == "-v") {
		metadata.PrintFull(os.Stdout)
		os.Exit(0)
	}

	var options struct {
		Format      string  `short:"f" long:"format" choice:"json" choice:"yaml" choice:"text" choice:"template" optional:"true" default:"text"`
		Template    *string `short:"t" long:"template" optional:"true"`
		Diagnostics bool    `long:"print-diagnostics" optional:"true"`
	}

	_, err := flags.Parse(&options)
	if err != nil {
		slog.Error("error parsing command line", "error", err)
		fmt.Fprintf(os.Stderr, "Invalid command line: %v\n", err)
		os.Exit(1)
	}
}
