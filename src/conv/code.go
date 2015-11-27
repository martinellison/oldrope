// code.go
package main

import (
	"fmt"
	"html"
	"regexp"
)

func makeGenData() {
	theOutData.Pages = make([]*outPage, 0, len(thePageSet))
	for _, page := range thePageSet {
		outPage := makeOutPage(page.theName)
		outPage.codePage(&page)
		theOutData.Pages = append(theOutData.Pages, outPage)
	}
}
func makeOutPage(theName string) (theOutPage *outPage) {
	theOutPage = &outPage{Name: theName, SetLines: make([]string, 0, 0), FixLines: make([]fix, 0, 0), RedisplayLines: make([]string, 0, 0)}
	return
}
func (theOutPage *outPage) codePage(thePage *page) {
	for _, theFragment := range thePage.theFragments {
		subOutFrag := makeOutPage(theFragment.name)
		subOutFrag.codeFragment(theFragment, true)
		theOutPage.SetLines = append(theOutPage.SetLines, subOutFrag.SetLines...)
		theOutPage.FixLines = append(theOutPage.FixLines, subOutFrag.FixLines...)
		theOutPage.RedisplayLines = append(theOutPage.RedisplayLines, subOutFrag.RedisplayLines...)
	}
}

var autoLink int
var allSpace *regexp.Regexp

func init() { allSpace = regexp.MustCompile(`^[\s]*$`) }

func (theOutFrag *outPage) codeFragment(theFragment *fragment, set bool) {
	//theOutFrag = new(outPage)
	comprText := compress(theFragment.text)
	escapeText := html.EscapeString(comprText)
	fragName := theFragment.name
	if fragName == "" {
		autoLink++
		fragName = fmt.Sprintf("z%d", autoLink)
	}
	theOutFrag.Name = fragName
	fragIdAttr := " id=\"" + fragName + "\" "
	//	subset := true
	switch theFragment.theFragType {
	case spanFragType:
		spanText := "parts.push('<span" + fragIdAttr + ">" + escapeText + "');"
		theOutFrag.addLine(spanText, true)
		theOutFrag.addLine("ld.s"+fragName+"=false;", true)
		theOutFrag.addLine("if (ld.s"+fragName+") {parts=[];", false)
	//	subset = false
	case divFragType:
		divText := "parts.push('<div" + fragIdAttr + ">" + escapeText + "');"
		theOutFrag.addLine(divText, true)
	//	subset = false
	case paraFragType:
		paraText := "parts.push('<p" + fragIdAttr + ">" + escapeText + "');"
		theOutFrag.addLine(paraText, true)
	case jsCodeFragType:
		theOutFrag.addLine(comprText, true)
	case jsExprFragType:
		exprText := "parts.push(" + comprText + ");"
		theOutFrag.addLine(exprText, true)
	case textFragType:
		if !allSpace.MatchString(escapeText) {
			textText := "parts.push('" + escapeText + "');"
			theOutFrag.addLine(textText, true)
		}
	case htmlFragType:
		if !allSpace.MatchString(comprText) {
			htmlText := "parts.push('<" + comprText + ">');"
			theOutFrag.addLine(htmlText, true)
		}
	case linkFragType:
		textText := "parts.push('<a" + fragIdAttr + ">');"
		theOutFrag.addLine(textText, true)
		code := ""
		if theFragment.auxName != "" {
			code = " setPage('" + theFragment.auxName + "'); displayPage();"
		} else {
			code = " ld.s" + fragName + "=true; displayPage();"
		}
		linkFix := fix{Name: fragName, Code: code}
		theOutFrag.FixLines = append(theOutFrag.FixLines, linkFix)
		//theOutFrag.FixLines = []fix{linkFix}
	default:
	}
	for _, theFragment := range theFragment.theFragments {
		subOutFrag := makeOutPage(theFragment.name)
		subOutFrag.codeFragment(theFragment, true)
		theOutFrag.SetLines = append(theOutFrag.SetLines, subOutFrag.SetLines...)
	}

	switch theFragment.theFragType {
	case spanFragType:
		theOutFrag.addLine("parts.push('</span>');", true)
		theOutFrag.addLine("setHtml('"+fragName+"',parts.join(','));}", false)
	case divFragType:
		theOutFrag.addLine("parts.push('</div>');", true)
	case paraFragType:
		theOutFrag.addLine("parts.push('</p>');", true)
	case linkFragType:
		theOutFrag.addLine("parts.push('</a>');", true)
	default:
	}
	for _, theFragment := range theFragment.actionFragments {
		subOutFrag := makeOutPage(theFragment.name)
		subOutFrag.codeFragment(theFragment, true)
		theOutFrag.FixLines = append(theOutFrag.FixLines, subOutFrag.FixLines...)
		theOutFrag.RedisplayLines = append(theOutFrag.RedisplayLines, subOutFrag.RedisplayLines...)
	}
}

type fragCode struct {
}

func (theOutPage *outPage) addLine(line string, set bool) {
	if set {
		theOutPage.SetLines = append(theOutPage.SetLines, line)
	} else {
		theOutPage.RedisplayLines = append(theOutPage.RedisplayLines, line)
	}
}
