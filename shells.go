package kongcompletion

import "text/template"

type shell struct {
	// initCode is a Go template for the shell initialization code.
	initCode *template.Template

	// initFilePath is the path of the shell’s default init file, e.g. ~/.bashrc
	initFilePath string
}

type binaryInfo struct {
	BinName string
	BinPath string
}

var shells = map[string]shell{
	"bash": bash,
	"zsh":  zsh,
	"fish": fish,
}

var bash = shell{
	initCode:     tmpl(`complete -o default -o bashdefault -C {{.BinPath}} {{.BinName}}`),
	initFilePath: "~/.bashrc",
}

var zsh = shell{
	initCode: tmpl(`autoload -U +X bashcompinit && bashcompinit
complete -o default -o bashdefault -C {{.BinPath}} {{.BinName}}`),
	initFilePath: "~/.zshrc",
}

var fish = shell{
	initCode: tmpl(`function __complete_{{.BinName}}
    set -lx COMP_LINE (commandline -cp)
    test -z (commandline -ct)
    and set COMP_LINE "$COMP_LINE "
    {{.BinPath}}
end
complete -f -c {{.BinName}} -a "(__complete_{{.BinName}})"`),
	initFilePath: "~/.config/fish/config.fish",
}

// tmpl returns an unnamed go template from the given input string.
// It panics if the template is malformed.
func tmpl(tmpl string) *template.Template {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	return t
}
