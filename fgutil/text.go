package fgutil

import (
	"bufio"
	"io"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

//RenderTemplate renders the specified template
func RenderTemplate(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": Capitalize})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func Capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

func printUsage(w io.Writer, template string, data interface{}) {
	bw := bufio.NewWriter(w)
	RenderTemplate(bw, template, data)
	bw.Flush()
}



func IsStringInSlice(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}