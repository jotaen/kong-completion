package kongcompletion

import (
	"errors"
	"fmt"

	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/riywo/loginshell"
)

// Completion is a kong subcommand that prints out the shell code for
// initializing tab completion in various shells. It also educates the
// user what to do with the printed code.
type Completion struct {
	Shell             string `arg:"" help:"The name of the shell you are using" enum:"bash,zsh,fish," default:""`
	NoDefaultFileComp bool   `help:"Whether or not to default to file comparison if no completion result is available"`
	Code              bool   `short:"c" help:"Generate the initialization code"`
}

// Help is a predefined kong method for printing the help text.
func (c *Completion) Help() string {
	return `
Displays a command that you need to execute in order activate tab completion for this program.

For permanent activation (i.e. beyond the current shell session), paste the command in your shell’s init file.

If no shell is specified, it tries to detect your current login shell automatically.
`
}

// Run is a predefined kong method that contains the command’s main procedure.
func (c *Completion) Run(ctx *kong.Context) error {
	binInfo, err := determineBinaryInfo(ctx, c)
	if err != nil {
		return err
	}

	// Determine targeted shell.
	sh, err := (func() (shell, error) {
		if c.Shell != "" {
			return newShellFromString(c.Shell, c.NoDefaultFileComp)
		}
		return detectShell()
	})()
	if err != nil {
		return err
	}

	// Generate command output.
	output := (func() string {
		if c.Code {
			return binInfo.fill(sh.initCode)
		} else {
			return "" +
				"Execute the following command to activate tab completion for " + binInfo.BinName + " in " + sh.name + ":\n\n" +
				"    " + binInfo.fill(sh.dynamicInitCode) + "\n\n" +
				"Note that this only takes effect for your current shell session. For permanent activation (beyond the current shell session), you can e.g. paste this command into your " + sh.name + "’s init file, which usually is: " + sh.initFilePath
		}
	})()
	_, err = fmt.Fprint(ctx.Stdout, output+"\n")
	if err != nil {
		return err
	}

	// Instruct kong to exit, to prevent the base command to be executed
	// and potentially printing something to stdout afterwards.
	ctx.Exit(0)
	return nil
}

// detectShell tries to determine from the process environment what the user’s
// login shell is.
func detectShell() (shell, error) {
	shellName, err := loginshell.Shell()
	if err != nil {
		return shell{}, errors.New("couldn't determine user's shell")
	}
	sh, ok := shells[filepath.Base(shellName)]
	if !ok {
		return shell{}, errors.New("this shell is not supported (" + shellName + ")")
	}
	return sh, nil
}

// determineBinaryInfo tries to determine information about the current command.
func determineBinaryInfo(ctx *kong.Context, cmd *Completion) (binaryInfo, error) {
	bin, err := os.Executable()
	if err != nil {
		return binaryInfo{}, fmt.Errorf("couldn't determine absolute path to binary: %w", err)
	}
	bin, err = filepath.Abs(bin)
	if err != nil {
		return binaryInfo{}, fmt.Errorf("couldn't determine absolute path to binary: %w", err)
	}

	var opts string
	if !cmd.NoDefaultFileComp {
		opts = "-o default -o bashdefault"
	}

	return binaryInfo{ctx.Model.Name, bin, ctx.Selected().Name, opts}, nil
}
