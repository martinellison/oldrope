// conv.go
package main

import (
	"html"
	"log"
	"os"
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
	makeTemplate()
	makeGenData()
	expandTemplate(os.Stdout)
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
func makeGenData() {
	theOutData.Pages = make([]*outPage, 0, len(thePageSet))
	for _, page := range thePageSet {
		outPage := &outPage{Name: page.theName, SetLines: make([]string, 0, 0), FixLines: make([]fix, 0, 0), RedisplayLines: make([]string, 0, 0)}
		outPage.codePage(&page)
		theOutData.Pages = append(theOutData.Pages, outPage)
	}
}
func (theOutPage *outPage) codePage(thePage *page) {
	for _, theFragment := range thePage.theFragments {
		theOutPage.codeFragment(theFragment)
	}
}
func (theOutPage *outPage) codeFragment(theFragment *fragment) {
	comprText := compress(theFragment.text)
	escapeText := html.EscapeString(comprText)
	switch theFragment.theFragType {
	case spanFragType:
		spanText := "parts.push('<span>" + escapeText + "</span>');"
		theOutPage.SetLines = append(theOutPage.SetLines, spanText)
	case divFragType:
		divText := "parts.push('<div>" + escapeText + "</div>');"
		theOutPage.SetLines = append(theOutPage.SetLines, divText)
	case paraFragType:
		paraText := "parts.push('<p>" + escapeText + "</p>');"
		theOutPage.SetLines = append(theOutPage.SetLines, paraText)
	case jsCodeFragType:
		theOutPage.SetLines = append(theOutPage.SetLines, comprText)
	case jsExprFragType:
		exprText := "parts.push(" + comprText + ");"
		theOutPage.SetLines = append(theOutPage.SetLines, exprText)
	case textFragType:
		textText := "parts.push('" + escapeText + "');"
		theOutPage.SetLines = append(theOutPage.SetLines, textText)
	case linkFragType:
	default:
	}
	//	fmt.Printf("%s%s (%s): %s\n", theFragment.name, theFragment.theFragType, theFragment.text)
	if theFragment.auxName != "" {
		//		fmt.Printf("%sgoto:%s\n", indent, theFragment.auxName)
	}
	for _, theFragment := range theFragment.theFragments {
		theOutPage.codeFragment(theFragment)
	}
	for _, theFragment := range theFragment.actionFragments {
		theOutPage.codeFragment(theFragment)
	}
}
