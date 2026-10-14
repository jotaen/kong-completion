package kongcompletion

import (
	"errors"
	"strings"
	gotemplate "text/template"
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
	name: "bash",
	// bash parses the -C command again when it completes, so BinPath needs to be quoted twice.
	initCode: tmpl(`complete{{if .UseShellDefault}} -o default -o bashdefault{{ end }} -C {{quote (quote .BinPath)}} {{.BinName}}`,
		gotemplate.FuncMap{"quote": shquote}),
	configFileCode: tmpl(`source <({{.BinName}} {{.SubCmdName}} -c bash)`, nil),
	initFilePath:   "~/.bashrc",
}

var zsh = shell{
	name: "zsh",
	// zsh parses the -C command again when it completes, so BinPath needs to be quoted twice.
	initCode: tmpl(`#compdef {{.BinName}}
autoload -U +X bashcompinit && bashcompinit
complete{{if .UseShellDefault}} -o default -o bashdefault{{ end }} -C {{quote (quote .BinPath)}} {{.BinName}}`,
		gotemplate.FuncMap{"quote": shquote}),
	configFileCode: tmpl(`source <({{.BinName}} {{.SubCmdName}} -c zsh)`, nil),
	initFilePath:   "~/.zshrc",
}

var fish = shell{
	name: "fish",
	initCode: tmpl(`function __complete_{{.BinName}}
    set -lx COMP_LINE (commandline -cp)
    test -z (commandline -ct)
    and set COMP_LINE "$COMP_LINE "
    {{quote .BinPath}}
end
complete -f -c {{.BinName}} -a "(__complete_{{.BinName}})"`,
		gotemplate.FuncMap{
			"quote": func(s string) string {
				if s != "" && strings.Trim(s, safeChars) == "" {
					return s
				}
				return "'" + strings.NewReplacer(`\`, `\\`, "'", `\'`).Replace(s) + "'"
			},
		}),
	configFileCode: tmpl(`{{.BinName}} {{.SubCmdName}} -c fish | source`, nil),
	initFilePath:   "~/.config/fish/config.fish",
}

// safeChars never need quoting in a shell word.
const safeChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@%+=:,./_-"

func shquote(s string) string {
	if s != "" && strings.Trim(s, safeChars) == "" {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
