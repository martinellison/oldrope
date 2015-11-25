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
	go parse()
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

var theParser struct {
	theCurrentToken token
}

func parse() {
	getToken()
	for tokTyp() != eofTokenType {
		if tokTyp() == identTokenType {
			switch tokText() {
			case "page":
				parsePage()
			default:
				log.Printf("unknown ident: %s", theParser.theCurrentToken.text)
				getToken()
			}
		} else {
			getToken()
		}
	}
}
func parsePage() {
	expectIdent("page")
	log.Printf("page name: %s", tokText())
	getToken()
	parseBody([]string{"page"})
}
func parseBody(stopIdents []string) {
	log.Printf("parsing body")
	for !stopped(stopIdents) {
		if tokTyp() == identTokenType {
			switch tokText() {
			case "link":
				expectIdent("link")
				if tokTyp() == identTokenType {
					log.Printf("link name: %s", tokText())
					getToken()
				}
				parseBody([]string{"end", "page", "goto", "act"})
				if tokIsIdent("goto") {
					expectIdent("goto")
					log.Printf("goto target: %s", tokText())
					getToken()
				} else if tokIsIdent("end") {
					expectIdent("end")
				} else if tokIsIdent("act") {
					expectIdent("act")
					parseBody([]string{"end"})
					expectIdent("end")
				}
			case "div", "span":
				getToken()
				log.Printf("div/span name: %s", tokText())
				getToken() // div/span name
				parseBody([]string{"end"})
				expectIdent("end")
			case "include":
				expectIdent("include")
				log.Printf("include target: %s", tokText())
				getToken() //target
			default:
				log.Printf("unknown ident: %s", tokText())
				getToken()
			}
		} else {
			getToken()
		}
	}
	log.Print("body done")
}
func stopped(stopIdents []string) bool {
	if tokTyp() == eofTokenType {
		return true
	}
	if theParser.theCurrentToken.theType != identTokenType {
		return false
	}
	for _, stopIdent := range stopIdents {
		if stopIdent == theParser.theCurrentToken.text {
			return true
		}
	}
	return false
}

func getToken() {
	theParser.theCurrentToken = <-tokenChan
	if tokTyp() == eofTokenType {
		log.Printf("end of input")
	}
}
func tokTyp() tokenType { return theParser.theCurrentToken.theType }
func tokText() string   { return theParser.theCurrentToken.text }
func tokIsIdent(id string) bool {
	return tokTyp() != eofTokenType && theParser.theCurrentToken.theType == identTokenType && theParser.theCurrentToken.text == id
}
func expectIdent(id string) {
	if tokIsIdent(id) {
		getToken()
	} else {
		log.Printf("expected %s but not found!", id)
	}

}

type pageSet map[string]page
type page struct {
	local []string
	//	fixedParas
	//	dynParas
}

var thePageSet pageSet

var startPageName string

func init() {
	thePageSet = make(map[string]page, 0)
	startPageName = "start"
}
