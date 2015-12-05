// parser.go
package main

import (
	"fmt"
	"log"
	"os"
)

///* */ func dumpTokens() {
//	if logging {
//		log.Print("dumping tokens...")
//	}
//	for {
//		if logging {
//			log.Print("waiting for token")
//		}
//		var theToken token
//		theToken = <-tokenChan
//		if theToken.theType == eofTokenType {
//			break
//		}
//		if logging {
//		log.Printf("%d: token  %s", theToken.lineNumber, theToken.text)
//			}
//	}
//	if logging {
//		log.Print("all tokens read")
//	}
//	linesDone <- 1
//}

/* theParser is the parser */
var theParser parser

/* a parser contains the parse state  */
type parser struct {
	theCurrentToken token
}

/* parse runs the parse (top-level parser) */
func (theParser *parser) parse() {
	defer func() {
		rec := recover()
		if rec == nil {
			return
		}
		os.Stderr.WriteString(fmt.Sprintf("parse: internal error: %v", rec))
	}()
	theParser.getToken()
	for theParser.tokTyp() != eofTokenType {
		if theParser.tokTyp() == identTokenType {
			switch theParser.tokText() {
			case "page":
				theParser.parsePage()
			default:
				reportError(fmt.Sprintf("unknown identifier: %s", theParser.theCurrentToken.text), theParser.theCurrentToken.lineNumber)
				theParser.getToken()
			}
		} else {
			if !isAllSpace(theParser.tokText()) {
				reportError(fmt.Sprintf("not in directive: %s", theParser.tokText()), theParser.theCurrentToken.lineNumber)
			}
			theParser.getToken()
		}
	}
}

/* parsePage parses a page directive */
func (theParser *parser) parsePage() {
	theParser.expectIdent("page")
	pageName := theParser.tokText()
	if logging {
		log.Printf("page name: '%s'", pageName)
	}
	theParser.getToken()
	thePage := page{local: make([]string, 0), theFragmentsByName: make(map[string]*fragment, 0), theName: pageName}
	thePage.theFragments = theParser.parseBody([]string{"page"})
	if thePage.theName != "" {
		thePageSet[thePage.theName] = thePage
	}
	for _, fr := range thePage.theFragments {
		if fr.name != "" {
			thePage.theFragmentsByName[fr.name] = fr
		}
	}
}

/* parseBody */
func (theParser *parser) parseBody(stopIdents []string) (theFragments []*fragment) {
	if logging {
		log.Printf("parsing body")
	}
	theFragments = make([]*fragment, 0)
	for !theParser.stopped(stopIdents) {
		theFragment := &fragment{theFragments: make([]*fragment, 0), name: "", text: "", auxName: "", actionFragments: make([]*fragment, 0)}
		if logging {
			log.Printf("getting fragment with token: %s", theParser.tokText())
		}
		switch theParser.tokTyp() {
		case identTokenType:
			switch theParser.tokText() {
			case "link":
				theFragment.theFragType = linkFragType
				theParser.expectIdent("link")
				if theParser.tokTyp() == identTokenType {
					if logging {
						log.Printf("link name: %s", theParser.tokText())
					}
					theFragment.name = theParser.tokText()
					theParser.getToken()
				}
				theFragment.theFragments = theParser.parseBody([]string{"end", "page", "goto", "act"})
				if theParser.tokIsIdent("goto") {
					theParser.expectIdent("goto")
					if logging {
						log.Printf("goto target: %s", theParser.tokText())
					}
					theFragment.auxName = theParser.tokText()
					theParser.getToken()
				} else if theParser.tokIsIdent("end") {
					theParser.expectIdent("end")
				} else if theParser.tokIsIdent("act") {
					theParser.expectIdent("act")
					theFragment.actionFragments = theParser.parseBody([]string{"end"})
					theParser.expectIdent("end")
				}
			case "div", "span":
				switch theParser.tokText() {
				case "div":
					theFragment.theFragType = divFragType
				case "span":
					theFragment.theFragType = spanFragType
				}
				theParser.getToken()
				if logging {
					log.Printf("div/span name: %s", theParser.tokText())
				}
				theFragment.name = theParser.tokText()
				theParser.getToken() // div/span name
				theFragment.actionFragments = theParser.parseBody([]string{"end"})
				theParser.expectIdent("end")
			case "include":
				theParser.expectIdent("include")
				if logging {
					log.Printf("include target: %s", theParser.tokText())
				}
				theFragment.theFragType = includeFragType
				theFragment.auxName = theParser.tokText()
				theParser.getToken() //target
			default:
				reportError(fmt.Sprintf("unknown directive: %s", theParser.tokText()), theParser.theCurrentToken.lineNumber)
				theParser.getToken()
				theParser.getToken()
			}
		case jsCodeTokenType:
			theFragment.theFragType = jsCodeFragType
			theFragment.text = theParser.tokText()
			theParser.getToken()
		case jsExprTokenType:
			theFragment.theFragType = jsExprFragType
			theFragment.text = theParser.tokText()
			theParser.getToken()
		case textTokenType:
			theFragment.theFragType = textFragType
			theFragment.text = theParser.tokText()
			theParser.getToken()
		case htmlTokenType:
			theFragment.theFragType = htmlFragType
			theFragment.text = theParser.tokText()
			theParser.getToken()
		default:
			reportError(fmt.Sprintf("wrong kind of token: %s", theParser.tokText()), theParser.theCurrentToken.lineNumber)
			theParser.getToken()
		}
		theFragments = append(theFragments, theFragment)
	}
	logIfLogging("body done")
	return
}

/* stopped detects a stop token */ func (theParser *parser) stopped(stopIdents []string) bool {
	if theParser.tokTyp() == eofTokenType {
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

/* getToken gets a new token */ func (theParser *parser) getToken() {
	theParser.theCurrentToken = <-tokenChan
	if logging {
		log.Printf("parsing token: '%s'", theParser.theCurrentToken.text)
	}
	if theParser.tokTyp() == eofTokenType && logging {
		log.Printf("end of input")
	}
}

/* tokTyp returns the type of the current togken */
func (theParser *parser) tokTyp() tokenType { return theParser.theCurrentToken.theType }

/* tokText returns the text of the current token */
func (theParser *parser) tokText() string { return theParser.theCurrentToken.text }

/* tokIsIdent returns whether the current token is the specified identifier */
func (theParser *parser) tokIsIdent(id string) bool {
	if theParser.tokTyp() != identTokenType {
		return false
	}
	return theParser.theCurrentToken.text == id
}

/* expectIdent consumes the current token if it is the expected token, and otherwise raises an error*/
func (theParser *parser) expectIdent(id string) {
	if theParser.tokIsIdent(id) {
		theParser.getToken()
	} else {
		reportError(fmt.Sprintf("expected %s but not found", id), theParser.theCurrentToken.lineNumber)
	}
}

/* */ type pageSet map[string]page

/* */ type page struct {
	local              []string
	theFragments       []*fragment
	theFragmentsByName map[string]*fragment
	theName            string
}

/* */ type fragType int

const (
	spanFragType fragType = iota
	divFragType
	paraFragType
	jsCodeFragType
	jsExprFragType
	textFragType
	linkFragType
	htmlFragType
	includeFragType
)

/* */ func (theFragType fragType) String() string {
	switch theFragType {
	case spanFragType:
		return "span"
	case divFragType:
		return "div"
	case paraFragType:
		return "para"
	case jsCodeFragType:
		return "jsCode"
	case jsExprFragType:
		return "jsExpr"
	case textFragType:
		return "text"
	case linkFragType:
		return "link"
	case htmlFragType:
		return "html"
	case includeFragType:
		return "include"
	default:
		return "(??)"
	}
}

/* */ type fragment struct {
	theFragType     fragType
	name            string
	text            string
	theFragments    []*fragment
	actionFragments []*fragment
	auxName         string
}

/* */ var thePageSet pageSet

/* */ var startPageName string

/* */ func init() {
	thePageSet = make(map[string]page, 0)
	startPageName = "start"
}

/* */ func dumpPages() {
	if !logging {
		return
	}
	log.Print("pages as parsed:")
	for _, page := range thePageSet {
		log.Printf("--- page %s ---\n", page.theName)
		for _, fr := range page.theFragments {
			dumpFragment(fr, "\t")
		}
	}
}

/* */ func dumpFragment(fr *fragment, indent string) {
	log.Printf("%s%s (%s): %s\n", indent, fr.name, fr.theFragType.String(), fr.text)
	if fr.auxName != "" {
		log.Printf("%sgoto:%s\n", indent, fr.auxName)
	}
	for _, fr := range fr.theFragments {
		dumpFragment(fr, indent+"\tf: ")
	}
	for _, fr := range fr.actionFragments {
		dumpFragment(fr, indent+"a: \t")
	}
}
