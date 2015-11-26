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

func (theOutFrag *outPage) codeFragment(theFragment *fragment) {
	theOutFrag = new(outPage)
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
		theOutFrag.addLine(spanText, set)
		theOutFrag.addLine("ld.s"+fragName+"=false;", set)
		theOutFrag.addLine("if (ld.s"+fragName+") {parts=[];", false)
		subset = false
	case divFragType:
		divText := "parts.push('<div" + fragIdAttr + ">" + escapeText + "');"
		theOutFrag.addLine(divText, set)
		subset = false
	case paraFragType:
		paraText := "parts.push('<p" + fragIdAttr + ">" + escapeText + "');"
		theOutFrag.addLine(paraText, set)
	case jsCodeFragType:
		theOutFrag.addLine(comprText, set)
	case jsExprFragType:
		exprText := "parts.push(" + comprText + ");"
		theOutFrag.addLine(exprText, set)
	case textFragType:
		if !allSpace.MatchString(escapeText) {
			textText := "parts.push('" + escapeText + "');"
			theOutFrag.addLine(textText, set)
		}
	case linkFragType:
		textText := "parts.push('<a" + fragIdAttr + ">');"
		theOutFrag.addLine(textText, set)
		code := ""
		if theFragment.auxName != "" {
			code = " setPage('" + theFragment.auxName + "'); displayPage();"
		} else {
			code = " ld.s" + fragName + "=true; displayPage();"
		}
		linkFix := fix{Name: fragName, Code: code}
		theOutFrag.FixLines = append(theOutPage.FixLines, linkFix)
	default:
	}
	for _, theFragment := range theFragment.theFragments {
		theOutFrag.codeFragment(theFragment, set)
	}
	for _, theFragment := range theFragment.actionFragments {
		theOutFrag.codeFragment(theFragment, subset)
	}

	switch theFragment.theFragType {
	case spanFragType:
		theOutFrag.addLine("parts.push('</span>');", set)
		theOutFrag.addLine("setHtml('"+fragName+"',parts.join(','));}", false)
	case divFragType:
		theOutFrag.addLine("parts.push('</div>');", set)
	case paraFragType:
		theOutFrag.addLine("parts.push('</p>');", set)
	case linkFragType:
		theOutFrag.addLine("parts.push('</a>');", set)
	default:
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
