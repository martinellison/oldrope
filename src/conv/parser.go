// parser.go
package main

import (
	"fmt"
	"log"
)

//func dumpTokens() {
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
				reportError(fmt.Sprintf("unknown identifier: %s", theParser.theCurrentToken.text), theParser.theCurrentToken.lineNumber)
				getToken()
			}
		} else {
			if !isAllSpace(tokText()) {
				reportError(fmt.Sprintf("not in directive: %s", tokText()), theParser.theCurrentToken.lineNumber)
			}
			getToken()
		}
	}
}
func parsePage() {
	expectIdent("page")
	pageName := tokText()
	if logging {
		log.Printf("page name: '%s'", pageName)
	}
	getToken()
	thePage := page{local: make([]string, 0), theFragmentsByName: make(map[string]*fragment, 0), theName: pageName}
	thePage.theFragments = parseBody([]string{"page"})
	if thePage.theName != "" {
		thePageSet[thePage.theName] = thePage
	}
	for _, fr := range thePage.theFragments {
		if fr.name != "" {
			thePage.theFragmentsByName[fr.name] = fr
		}
	}
}
func parseBody(stopIdents []string) (theFragments []*fragment) {
	if logging {
		log.Printf("parsing body")
	}
	theFragments = make([]*fragment, 0)
	for !stopped(stopIdents) {
		theFragment := &fragment{theFragments: make([]*fragment, 0), name: "", text: "", auxName: "", actionFragments: make([]*fragment, 0)}
		if logging {
			log.Printf("getting fragment with token: %s", tokText())
		}
		switch tokTyp() {
		case identTokenType:
			switch tokText() {
			case "link":
				theFragment.theFragType = linkFragType
				expectIdent("link")
				if tokTyp() == identTokenType {
					if logging {
						log.Printf("link name: %s", tokText())
					}
					theFragment.name = tokText()
					getToken()
				}
				theFragment.theFragments = parseBody([]string{"end", "page", "goto", "act"})
				if tokIsIdent("goto") {
					expectIdent("goto")
					if logging {
						log.Printf("goto target: %s", tokText())
					}
					theFragment.auxName = tokText()
					getToken()
				} else if tokIsIdent("end") {
					expectIdent("end")
				} else if tokIsIdent("act") {
					expectIdent("act")
					theFragment.actionFragments = parseBody([]string{"end"})
					expectIdent("end")
				}
			case "div", "span":
				switch tokText() {
				case "div":
					theFragment.theFragType = divFragType
				case "span":
					theFragment.theFragType = spanFragType
				}
				getToken()
				if logging {
					log.Printf("div/span name: %s", tokText())
				}
				theFragment.name = tokText()
				getToken() // div/span name
				theFragment.actionFragments = parseBody([]string{"end"})
				expectIdent("end")
			case "include":
				expectIdent("include")
				if logging {
					log.Printf("include target: %s", tokText())
				}
				theFragment.theFragType = includeFragType
				theFragment.auxName = tokText()
				getToken() //target
			default:
				reportError(fmt.Sprintf("unknown directive: %s", tokText()), theParser.theCurrentToken.lineNumber)
				getToken()
				getToken()
			}
		case jsCodeTokenType:
			theFragment.theFragType = jsCodeFragType
			theFragment.text = tokText()
			getToken()
		case jsExprTokenType:
			theFragment.theFragType = jsExprFragType
			theFragment.text = tokText()
			getToken()
		case textTokenType:
			theFragment.theFragType = textFragType
			theFragment.text = tokText()
			getToken()
		case htmlTokenType:
			theFragment.theFragType = htmlFragType
			theFragment.text = tokText()
			getToken()
		default:
			reportError(fmt.Sprintf("wrong kind of token: %s", tokText()), theParser.theCurrentToken.lineNumber)
			getToken()
		}
		theFragments = append(theFragments, theFragment)
	}
	if logging {
		log.Print("body done")
	}
	return
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
	if logging {
		log.Printf("parsing token: '%s'", theParser.theCurrentToken.text)
	}
	if tokTyp() == eofTokenType && logging {
		log.Printf("end of input")
	}
}
func tokTyp() tokenType { return theParser.theCurrentToken.theType }
func tokText() string   { return theParser.theCurrentToken.text }
func tokIsIdent(id string) bool {
	if tokTyp() != identTokenType {
		return false
	}
	return theParser.theCurrentToken.text == id
}
func expectIdent(id string) {
	if tokIsIdent(id) {
		getToken()
	} else {
		reportError(fmt.Sprintf("expected %s but not found", id), theParser.theCurrentToken.lineNumber)
	}
}

type pageSet map[string]page
type page struct {
	local              []string
	theFragments       []*fragment
	theFragmentsByName map[string]*fragment
	theName            string
}
type fragType int

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

func (theFragType fragType) String() string {
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

type fragment struct {
	theFragType     fragType
	name            string
	text            string
	theFragments    []*fragment
	actionFragments []*fragment
	auxName         string
}

var thePageSet pageSet

var startPageName string

func init() {
	thePageSet = make(map[string]page, 0)
	startPageName = "start"
}
func dumpPages() {
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
func dumpFragment(fr *fragment, indent string) {
	log.Printf("%s%s (%s): %s\n", indent, fr.name, fr.theFragType, fr.text)
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
