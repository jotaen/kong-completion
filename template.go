package kongcompletion

import (
	"bytes"
	gotemplate "text/template"
)

type templateData struct {
	BinName         string // The canonical name of the binary, e.g. `greet`
	BinPath         string // The full path to the binary, e.g. `/usr/bin/greet`
	SubCmdName      string // The name of the invoked subcommand, e.g. `completion` for `greet completion`.
	UseShellDefault bool   // Whether to fall back to default shell completions.
}

type template gotemplate.Template

// tmpl compiles a template from the given input string. It panics if the text is malformed.
func tmpl(tmpl string) *template {
	t, err := gotemplate.New("").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	return (*template)(t)
}

func (bi templateData) fill(t *template) string {
	result := &bytes.Buffer{}
	err := (*gotemplate.Template)(t).Execute(result, bi)
	if err != nil {
		panic(err)
	}
	return result.String()
}
