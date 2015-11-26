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
	makeOutData()
	expandTemplate(os.Stdout)
}

var templ *template.Template

func makeTemplate() {
	var err error
	templateText := "{{range .Pages}} {{.Name}}: , {{end}}"
	templ, err = template.New("page").Parse(templateText)
	if err != nil {
		log.Fatalf("template def error: %v", err)
	}
}

type outData struct {
	Pages []*outPage
}
type outPage struct {
	Name string
}

var theOutData outData

func makeOutData() {
	pageZzz := &outPage{Name: "zzz"}
	pageBill := &outPage{Name: "bill"}
	theOutData = outData{Pages: []*outPage{pageZzz, pageBill}}
}
func expandTemplate(w io.Writer) {
	err := templ.Execute(os.Stdout, theOutData)
	if err != nil {
		log.Fatalf("template exp error: %v", err)
	}
}
