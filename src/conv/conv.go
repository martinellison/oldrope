// conv.go
package main

import "log"

func main() {
	log.Print("start")
	lineChan = make(chan scanLine, 1)
	go getLines("/home/martin/git/twine/test.data")
	linesDone = make(chan int)
	go getTokens()
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

type token struct {
	theType    tokenType
	text       string
	lineNumber int
}
type tokenType int

const (
	eofTokenType tokenType = iota
)

var tokenChan chan token

func getTokens() {
	for {
		scanState.line = <-lineChan
		if scanState.line.eof {
		//	tokenChan <- token{theType: eofTokenType}			
	linesDone <- 1
			break
		}
		log.Printf("%d: %s", scanState.line.number, scanState.line.text)
		scanState.lineLen = len(scanState.line.text)
		scanState.charPos = 0
		for scanState.charPos < scanState.lineLen {
			log.Printf("scan char is: '%s'",scanChars(1))
			
			scanState.charPos ++
		}
	}
	log.Print("all lines read")
}
var scanState struct {
	lineLen int
	charPos int
	line scanLine
}
func scanChars (leng int)(chars string) {
	useLen:=leng
	if leng+scanState.charPos>scanState.lineLen{useLen=scanState.lineLen-scanState.charPos}
	return scanState.line.text[scanState.charPos:useLen]
}
