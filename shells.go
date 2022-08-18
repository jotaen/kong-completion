package kongcompletion

import (
	"github.com/pkg/errors"
)

type shell struct {
	// name is the name of the shell
	name string

	// initCode is a template for the shell initialization code.
	initCode *template

	// dynamicInitCode is a template for the command that prints the shell initialization code.
	dynamicInitCode *template

	// initFilePath is the path of the shell’s default init file, e.g. ~/.bashrc
	initFilePath string
}

var shells = map[string]shell{
	bash.name: bash,
	zsh.name:  zsh,
	fish.name: fish,
}

func newShellFromString(shellName string) (shell, error) {
	sh, ok := shells[shellName]
	if !ok {
		return shell{}, errors.New("")
	}
	return sh, nil
}

var bash = shell{
	name:            "bash",
	initCode:        tmpl(`complete -o default -o bashdefault -C {{.BinPath}} {{.BinName}}`),
	dynamicInitCode: tmpl(`source <({{.BinName}} {{.SubCmdName}} -c bash)`),
	initFilePath:    "~/.bashrc",
}

var zsh = shell{
	name: "zsh",
	initCode: tmpl(`autoload -U +X bashcompinit && bashcompinit
complete -o default -o bashdefault -C {{.BinPath}} {{.BinName}}`),
	dynamicInitCode: tmpl(`source <({{.BinName}} {{.SubCmdName}} -c zsh)`),
	initFilePath:    "~/.zshrc",
}

var fish = shell{
	name: "fish",
	initCode: tmpl(`function __complete_{{.BinName}}
    set -lx COMP_LINE (commandline -cp)
    test -z (commandline -ct)
    and set COMP_LINE "$COMP_LINE "
    {{.BinPath}}
end
complete -f -c {{.BinName}} -a "(__complete_{{.BinName}})"`),
	dynamicInitCode: tmpl(`{{.BinName}} {{.SubCmdName}} -c fish | source`),
	initFilePath:    "~/.config/fish/config.fish",
}
