// conv.go
package main

import "log"

func main() {
	log.Print("start")
	lineChan = make(chan scanLine, 1)
	go getLines("/home/martin/git/twine/test.data")
	linesDone = make(chan int)
	go getTokens()
	tokenChan = make(chan token)
	go dumpTokens()
	<-linesDone
	log.Print("all done.")
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

func dumpTokens() {
	log.Print("dumping tokens...")
	for {
		log.Print("waiting for token")
		var theToken token
		theToken = <-tokenChan
		if theToken.theType == eofTokenType {
			break
		}
		log.Printf("%d: %s", theToken.lineNumber, theToken.text)
	}
	log.Print("all tokens read")
	linesDone <- 1
}

type pageSet map[string]page
type page struct {
	local []string
	fixedParas
	dynParas
}

var startPageName string

func init() {
	pageSet = make(map[string]page, 0)
	startPageName = "start"
}
