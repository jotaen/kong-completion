package kongcompletion

import (
	"github.com/pkg/errors"
)

type shell struct {
	// initCode is a Go template for the shell initialization code.
	initCode *template

	dynamicInitCode *template

	// initFilePath is the path of the shell’s default init file, e.g. ~/.bashrc
	initFilePath string
}

var shells = map[string]shell{
	"bash": bash,
	"zsh":  zsh,
	"fish": fish,
}

func newShellFromString(shellName string) (shell, error) {
	sh, ok := shells[shellName]
	if !ok {
		return shell{}, errors.New("")
	}
	return sh, nil
}

var bash = shell{
	initCode:        tmpl(`complete -o default -o bashdefault -C {{.BinPath}} {{.BinName}}`),
	dynamicInitCode: tmpl(`source <({{.BinName}} {{.SubCmdName}} -c bash)`),
	initFilePath:    "~/.bashrc",
}

var zsh = shell{
	initCode: tmpl(`autoload -U +X bashcompinit && bashcompinit
complete -o default -o bashdefault -C {{.BinPath}} {{.BinName}}`),
	dynamicInitCode: tmpl(`source <({{.BinName}} {{.SubCmdName}} -c zsh)`),
	initFilePath:    "~/.zshrc",
}

var fish = shell{
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
