package base

import (
	"log/slog"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

type Command struct {
	// CPUProfile sets the (optional) path of the file for CPU profiling info.
	CPUProfile *string `short:"C" long:"cpu-profile" description:"The (optional) path where the CPU profiler will store its data." optional:"yes"`
	// MemProfile sets the (optional) path of the file for memory profiling info.
	MemProfile *string `short:"M" long:"mem-profile" description:"The (optional) path where the memory profiler will store its data." optional:"yes"`
}

func (cmd *Command) ProfileCPU() *Closer {
	var f *os.File
	if cmd.CPUProfile != nil {
		var err error
		f, err = os.Create(*cmd.CPUProfile)
		if err != nil {
			slog.Error("could not create CPU profile", "file", cmd.CPUProfile, "error", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			slog.Error("could not start CPU profiler", "error", err)
		}
	}
	return &Closer{
		file: f,
	}
}

func (cmd *Command) ProfileMemory() {
	if cmd.MemProfile != nil {
		f, err := os.Create(*cmd.MemProfile)
		if err != nil {
			slog.Error("could not create memory profile", "file", cmd.MemProfile, "error", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			slog.Error("could not write memory profile", "error", err)
		}
	}
}

type Closer struct {
	file *os.File
}

func (c *Closer) Close() {
	if c.file != nil {
		pprof.StopCPUProfile()
		c.file.Close()
	}
}

// AuthenticatedCommand is a command that requires authentication.
type AuthenticatedCommand struct {
	// Command is the base command.
	Command
	// Email is the email to use for authentication.
	Email string `short:"e" long:"email" description:"The email to use for authentication." required:"true" env:"LANDSCAPE_EMAIL"`
	// Password is the password to use for authentication.
	Password string `short:"p" long:"password" description:"The password to use for authentication." required:"true" env:"LANDSCAPE_PASSWORD"`
	// Account is the account to use for authentication.
	Account *string `short:"a" long:"account" description:"The account to use for authentication." optional:"true" env:"LANDSCAPE_ACCOUNT"`
	// Key is the access key to use for authentication.
	Key string `short:"k" long:"access-key" description:"The access key to use for authentication." required:"true" env:"LANDSCAPE_KEY"`
	// Secret is the secret key to use for authentication.
	Secret string `short:"s" long:"secret-key" description:"The secret key to use for authentication." required:"true" env:"LANDSCAPE_SECRET"`
	// Expiry is the expiry time for the authentication token.
	Expiry *time.Duration `short:"t" long:"token-expiry" description:"The expiry time for the authentication token." optional:"true" env:"LANDSCAPE_EXPIRY" default:"24h"`
	// Endpoint is the Landscape API endpoint.
	Endpoint string `short:"x" long:"endpoint" description:"The Landscape API endpoint." required:"true" env:"LANDSCAPE_ENDPOINT"`
}
