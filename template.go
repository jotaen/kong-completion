package kongcompletion

import (
	"bytes"
	"strings"
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
	t, err := gotemplate.New("").Funcs(funcs).Parse(tmpl)
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

var funcs = gotemplate.FuncMap{
	"shquote":   shquote,
	"fishquote": fishquote,
}

// safeChars never need quoting in a shell word.
const safeChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@%+=:,./_-"

// shquote quotes s as a single bash or zsh word, unless it needs no quoting.
func shquote(s string) string {
	if s != "" && strings.Trim(s, safeChars) == "" {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// fishquote quotes s as a single fish word, unless it needs no quoting.
func fishquote(s string) string {
	if s != "" && strings.Trim(s, safeChars) == "" {
		return s
	}
	return "'" + strings.NewReplacer(`\`, `\\`, "'", `\'`).Replace(s) + "'"
}
