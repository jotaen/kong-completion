package kongcompletion

import (
	"bytes"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
	"github.com/riywo/loginshell"
	"os"
	"path/filepath"
)

// Completion is a kong subcommand that prints out the shell code for
// initializing tab completions in various shells. It also educates the
// user what to do with the printed code.
type Completion struct {
	Bash bool `name:"bash" help:"Print the code for bash"`
	Zsh  bool `name:"zsh" help:"Print the code for zsh"`
	Fish bool `name:"fish" help:"Print the code for fish"`
}

// Help is a predefined kong method for printing the help text.
func (c *Completion) Help() string {
	sh, err := detectShell()
	if err != nil {
		sh = bash
	}

	return `
The output of this command is code for instructing your shell to use tab completions for this program. After you have executed that code in your shell, the tab completions are available.

In order to persist the completions beyond your current shell session, put the code into your shell initialization file, e.g. in ` + sh.initFilePath + `

If no flag for the shell type is specified, it tries to detect your current login shell automatically.
`
}

// Run is a predefined kong method that contains the command’s main procedure.
func (c *Completion) Run(ctx *kong.Context) error {
	binInfo, err := determineBinaryInfo(ctx)
	if err != nil {
		return err
	}

	sh, err := (func() (shell, error) {
		if c.Bash {
			return bash, nil
		} else if c.Zsh {
			return zsh, nil
		} else if c.Fish {
			return fish, nil
		}
		return detectShell()
	})()
	if err != nil {
		return err
	}

	output := expandInitTemplate(binInfo, sh) + "\n"
	_, err = fmt.Fprint(ctx.Stdout, output)
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
func determineBinaryInfo(ctx *kong.Context) (binaryInfo, error) {
	bin, err := os.Executable()
	if err != nil {
		return binaryInfo{}, errors.Wrapf(err, "couldn't determine absolute path to binary")
	}
	bin, err = filepath.Abs(bin)
	if err != nil {
		return binaryInfo{}, errors.Wrapf(err, "couldn't determine absolute path to binary")
	}
	return binaryInfo{ctx.Model.Name, bin}, nil

}

// expandInitTemplate resolves the init template of the given shell.
func expandInitTemplate(bi binaryInfo, sh shell) string {
	result := &bytes.Buffer{}
	err := sh.initCode.Execute(result, bi)
	if err != nil {
		panic(err)
	}
	return result.String()
}
