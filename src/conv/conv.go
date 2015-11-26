// conv.go
package main

import (
	"fmt"
	"html"
	"log"
	"os"
	"regexp"
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
	file, err := os.Create("test/testres.js")
	if err != nil {
		log.Fatal(err)
	}
	genStart(file)
	expandTemplate(file)
	genEnd(file)
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
		theOutPage.codeFragment(theFragment, true)
	}
}

var autoLink int
var allSpace *regexp.Regexp

func init() { allSpace = regexp.MustCompile(`^[\s]*$`) }

func (theOutPage *outPage) codeFragment(theFragment *fragment, set bool) {
	comprText := compress(theFragment.text)
	escapeText := html.EscapeString(comprText)
	fragName := theFragment.name
	if fragName == "" {
		autoLink++
		fragName = fmt.Sprintf("z%d", autoLink)
	}
	fragIdAttr := " id=\"" + fragName + "\" "
	subset := true
	switch theFragment.theFragType {
	case spanFragType:
		spanText := "parts.push('<span" + fragIdAttr + ">" + escapeText + "');"
		theOutPage.addLine(spanText, set)
		theOutPage.addLine("ld.s"+fragName+"=false;", set)
		theOutPage.addLine("if (ld.s"+fragName+") {parts=[];", false)
		subset = false
	case divFragType:
		divText := "parts.push('<div" + fragIdAttr + ">" + escapeText + "');"
		theOutPage.addLine(divText, set)
		subset = false
	case paraFragType:
		paraText := "parts.push('<p" + fragIdAttr + ">" + escapeText + "');"
		theOutPage.addLine(paraText, set)
	case jsCodeFragType:
		theOutPage.addLine(comprText, set)
	case jsExprFragType:
		exprText := "parts.push(" + comprText + ");"
		theOutPage.addLine(exprText, set)
	case textFragType:
		if !allSpace.MatchString(escapeText) {
			textText := "parts.push('" + escapeText + "');"
			theOutPage.addLine(textText, set)
		}
	case linkFragType:
		textText := "parts.push('<a" + fragIdAttr + ">');"
		theOutPage.addLine(textText, set)
		code := ""
		if theFragment.auxName != "" {
			code = " setPage('" + theFragment.auxName + "'); displayPage();"
		} else {
			code = " ld.s" + fragName + "=true; displayPage();"
		}
		linkFix := fix{Name: fragName, Code: code}
		theOutPage.FixLines = append(theOutPage.FixLines, linkFix)
	default:
	}
	//	fmt.Printf("%s%s (%s): %s\n", theFragment.name, theFragment.theFragType, theFragment.text)
	//if theFragment.auxName != "" {
	//		fmt.Printf("%sgoto:%s\n", indent, theFragment.auxName)
	//}
	for _, theFragment := range theFragment.theFragments {
		theOutPage.codeFragment(theFragment, set)
	}
	for _, theFragment := range theFragment.actionFragments {
		theOutPage.codeFragment(theFragment, subset)
	}

	switch theFragment.theFragType {
	case spanFragType:
		theOutPage.addLine("parts.push('</span>');", set)
		theOutPage.addLine("$('#"+fragName+"').html(parts.join(','));}", false)
	case divFragType:
		theOutPage.addLine("parts.push('</div>');", set)
	case paraFragType:
		theOutPage.addLine("parts.push('</p>');", set)
	case linkFragType:
		theOutPage.addLine("parts.push('</a>');", set)
	default:
	}
}
func (theOutPage *outPage) addLine(line string, set bool) {
	if set {
		theOutPage.SetLines = append(theOutPage.SetLines, line)
	} else {
		theOutPage.RedisplayLines = append(theOutPage.RedisplayLines, line)
	}
}
