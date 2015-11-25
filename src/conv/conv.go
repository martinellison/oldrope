// conv.go
package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
	log.Print("start")
	lineChan = make(chan scanLine, 1)
	go getLines("/home/martin/git/twine/test.data")
	linesDone = make(chan int)
	go getTokens()
	tokenChan = make(chan token)
	go parse()
	<-linesDone
	log.Print("all done.")
	dumpPages()
	templateText := "{{range .}} {{.theName}}: , {{end}}"
	templ, err := template.New("page").Parse(templateText)
	if err != nil {
		log.Fatalf("template def error: %v", err)
	}
	data := thePageSet
	err = templ.Execute(os.Stdout, data)
	if err != nil {
		log.Fatalf("template exp error: %v", err)
	}
}

func dumpLines() {
	for {
		var line scanLine
		line = <-lineChan
		if line.eof {
			break
		}
		log.Printf("%d: %s", line.number, line.text)
	}
	log.Print("all lines read")
	linesDone <- 1
}
