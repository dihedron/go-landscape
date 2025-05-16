package command

import (
	"github.com/dihedron/landscape/command/test"
)

// Commands is the set of root command groups.
type Commands struct {
	// Test prints the command arguments and exits.
	//lint:ignore SA5008 commands can have multiple aliases
	Test test.Test `command:"test" alias:"tst" alias:"t" description:"Show the command arguments and exit."`
}
