package kongcompletion

import (
	"errors"
)

type shell struct {
	// name is the name of the shell
	name string

	// initCode is the actual code that registers the completion in the shell.
	initCode *template

	// configFileCode is the code that the user is supposed to put in
	// their shell config file, e.g. ~/.bashrc
	configFileCode *template

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
	name:           "bash",
	initCode:       tmpl(`complete{{if .UseShellDefault}} -o default -o bashdefault{{ end }} -C {{.BinPath}} {{.BinName}}`),
	configFileCode: tmpl(`source <({{.BinName}} {{.SubCmdName}} -c bash)`),
	initFilePath:   "~/.bashrc",
}

var zsh = shell{
	name: "zsh",
	initCode: tmpl(`autoload -U +X bashcompinit && bashcompinit
complete{{if .UseShellDefault}} -o default -o bashdefault{{ end }} -C {{.BinPath}} {{.BinName}}`),
	configFileCode: tmpl(`source <({{.BinName}} {{.SubCmdName}} -c zsh)`),
	initFilePath:   "~/.zshrc",
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
	configFileCode: tmpl(`{{.BinName}} {{.SubCmdName}} -c fish | source`),
	initFilePath:   "~/.config/fish/config.fish",
}
