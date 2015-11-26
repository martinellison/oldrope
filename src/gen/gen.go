// gen.go
package main

import (
	"io"
	"log"
	"os"
	"text/template"
)

func main() {
	makeTemplate()
	expandTemplate(os.Stdout)
}

var templ *template.Template

func makeTemplate() {
	var err error
	templateText := "{{range .}} {{.theName}}: , {{end}}"
	templ, err = template.New("page").Parse(templateText)
	if err != nil {
		log.Fatalf("template def error: %v", err)
	}
}
func expandTemplate(w io.Writer) {
	data := "???"
	err := templ.Execute(os.Stdout, data)
	if err != nil {
		log.Fatalf("template exp error: %v", err)
	}
}
