package kongcompletion

import (
	"bytes"
	gotemplate "text/template"
)

type binaryInfo struct {
	BinName    string
	BinPath    string
	SubCmdName string
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

func (bi binaryInfo) fill(t *template) string {
	result := &bytes.Buffer{}
	err := (*gotemplate.Template)(t).Execute(result, bi)
	if err != nil {
		panic(err)
	}
	return result.String()
}
