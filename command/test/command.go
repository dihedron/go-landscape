package test

import (
	"log/slog"

	"github.com/dihedron/landscape/command/base"
)

// Version is the command that prints information about the application
// or plugin to the console; it support both compact and verbose mode.
type Test struct {
	base.AuthenticatedCommand
}

// Execute is the real implementation of the Version command.
func (cmd *Test) Execute(args []string) error {
	slog.Debug("running test command")
	slog.Debug("email", "value", cmd.Email)
	slog.Debug("password", "value", cmd.Password)
	slog.Debug("account", "value", cmd.Account)
	slog.Debug("key", "value", cmd.Key)
	slog.Debug("secret", "value", cmd.Secret)
	slog.Debug("expiry", "value", cmd.Expiry)
	slog.Debug("endpoint", "value", cmd.Endpoint)
	slog.Debug("command done")
	return nil
}
