package run

import (
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/fastly/cli/pkg/argparser"
	"github.com/fastly/cli/pkg/commands/profile"
	fsterr "github.com/fastly/cli/pkg/errors"
	"github.com/fastly/cli/pkg/global"
)

// RootCommand is the parent command for all subcommands in this package.
// It should be installed under the primary root command.
type RootCommand struct {
	argparser.Base
	tokenTTL  time.Duration
	command   string
	arguments []string
}

// CommandName is the string to be used to invoke this command
const CommandName = "run"

// NewRootCommand returns a new command registered in the parent.
func NewRootCommand(parent argparser.Registerer, g *global.Data) *RootCommand {
	var c RootCommand
	c.Globals = g
	c.CmdClause = parent.Command(CommandName, "Execute an external command using a token from a profile")
	c.CmdClause.Flag("token-ttl", "Amount of time for which the token must be valid (in seconds 's', minutes 'm', or hours 'h')").Default(profile.DefaultTokenTTL.String()).DurationVar(&c.tokenTTL)
	c.CmdClause.Arg("command", "Command to be executed").StringVar(&c.command)
	c.CmdClause.Arg("arguments", "Arguments to be passed to the command").StringsVar(&c.arguments)
	return &c
}

// Exec implements the command interface.
func (c *RootCommand) Exec(_ io.Reader, out io.Writer) error {

	cmd := exec.Command(c.command, c.arguments...)
	cmd.Env = append(cmd.Environ(), "FASTLY_API_KEY=foo", "FASTLY_API_TOKEN=bar")
	// the external command's input, output, and error output
	// should be passed through unmodified
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// if the command was executed, but reported an error,
		// we should not generate an error message, but
		// instead cause the same exit code to be returned by
		// this program
		if exiterr, ok := err.(*exec.ExitError); ok {
			return fsterr.PassthroughExitError{
				ExitCode: exiterr.ExitCode(),
				Err:      err,
			}
		}

		return err
	}

	return nil
}
