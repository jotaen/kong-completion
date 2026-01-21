// This code is copied over from kong to allow performing interpolation here.
// (The original function is not exported.)
// See https://github.com/alecthomas/kong/blob/master/interpolate.go
//
// Copyright (C) 2018 Alec Thomas (https://github.com/alecthomas/kong)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is furnished to do
// so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package kongcompletion

import (
	"fmt"
	"regexp"

	"github.com/alecthomas/kong"
)

var interpolationRegex = regexp.MustCompile(`(\$\$)|((?:\${([[:alpha:]_][[:word:]]*))(?:=([^}]+))?})|(\$)|([^$]+)`)

// Interpolate variables from vars into s for substrings in the form ${var} or ${var=default}.
func interpolate(s string, vars kong.Vars, updatedVars map[string]string) (string, error) {
	out := ""
	matches := interpolationRegex.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return s, nil
	}
	for key, val := range updatedVars {
		if vars[key] != val {
			vars = vars.CloneWith(updatedVars)
			break
		}
	}
	for _, match := range matches {
		if dollar := match[1]; dollar != "" {
			out += "$"
		} else if name := match[3]; name != "" {
			value, ok := vars[name]
			if !ok {
				// No default value.
				if match[4] == "" {
					return "", fmt.Errorf("undefined variable ${%s}", name)
				}
				value = match[4]
			}
			out += value
		} else {
			out += match[0]
		}
	}
	return out, nil
}
