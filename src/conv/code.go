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
	theOutPage = &outPage{Name: theName, SetLines: make([]string, 0, 0), FixLines: make([]string, 0, 0), RedisplayLines: make([]string, 0, 0), Refixes: make([]string, 0, 0)}
	return
}
func (theOutPage *outPage) codePage(thePage *page) {
	for _, theFragment := range thePage.theFragments {
		subOutFrag := makeOutPage(theFragment.name)
		subOutFrag.codeFragment(theFragment, true)
		theOutPage.SetLines = append(theOutPage.SetLines, subOutFrag.SetLines...)
		theOutPage.FixLines = append(theOutPage.FixLines, subOutFrag.FixLines...)
		theOutPage.RedisplayLines = append(theOutPage.RedisplayLines, subOutFrag.RedisplayLines...)
		theOutPage.Refixes = append(theOutPage.Refixes, subOutFrag.Refixes...)
	}
}

var autoLink int
var allSpace *regexp.Regexp

func init() { allSpace = regexp.MustCompile(`^[\s]*$`) }

func (theOutFrag *outPage) codeFragment(theFragment *fragment, topLevel bool) {
	//theOutFrag = new(outPage)
	comprText := compress(theFragment.text)
	escapeText := html.EscapeString(comprText)
	fragName := theFragment.name
	if fragName == "" {
		autoLink++
		fragName = fmt.Sprintf("z%d", autoLink)
	}
	theOutFrag.Name = fragName
	fragNameExtend := fragName + "-x"
	fragIdAttr := " id=\"" + fragName + "\" "
	fragIdExtendAttr := " id=\"" + fragNameExtend + "\" "
	mainLT := setLineType
	fixLT := fixLineType
	if !topLevel {
		mainLT = redisplayLineType
		fixLT = refixLineType
	}
	//	subset := true
	switch theFragment.theFragType {
	case spanFragType, divFragType:
		if topLevel {
			tag := ifThenElse(theFragment.theFragType == spanFragType, "span", "div")
			spanText := "parts.push('<" + tag + fragIdAttr + ">" + escapeText + "');"
			theOutFrag.addLine(spanText, setLineType)
			theOutFrag.addLine("ld.s"+fragName+"=false;", setLineType)
			theOutFrag.addLine("if (ld.s"+fragName+") {parts=[];", redisplayLineType)
		} else {
			theOutFrag.addLine("if (ld.s"+fragName+") {parts=[];", redisplayLineType)
		}
	case paraFragType:
		paraText := "parts.push('<p" + fragIdAttr + ">" + escapeText + "');"
		theOutFrag.addLine(paraText, mainLT)
	case jsCodeFragType:
		theOutFrag.addLine(comprText, mainLT)
	case jsExprFragType:
		exprText := "parts.push(" + comprText + ");"
		theOutFrag.addLine(exprText, mainLT)
	case textFragType:
		if !allSpace.MatchString(escapeText) {
			textText := "parts.push('" + escapeText + "');"
			theOutFrag.addLine(textText, mainLT)
		}
	case htmlFragType:
		if !allSpace.MatchString(comprText) {
			htmlText := "parts.push('<" + comprText + ">');"
			theOutFrag.addLine(htmlText, mainLT)
		}
	case linkFragType:
		textText := "parts.push('<a" + fragIdExtendAttr + ">');"
		theOutFrag.addLine(textText, mainLT)
		code := ""
		if theFragment.auxName != "" {
			code = " setPage('" + theFragment.auxName + "'); displayPage();"
		} else {
			code = " ld.s" + fragName + "=true; displayPage();"
		}
		//	if topLevel {
		linkFix := "setClick('" + fragNameExtend + "', function(){" + code + "});"
		theOutFrag.addLine(linkFix, fixLT)
		//	} else {
		//		theOutFrag.addLine("if (ld.s"+fragName+") {", redisplayLineType)
		//	}
		if theFragment.auxName == "" {
			theOutFrag.addLine("if (ld.s"+fragName+") {", mainLT)
			theOutFrag.addLine("if (ld.s"+fragName+") {", fixLT)
			theOutFrag.addLine("/*"+theFragment.name+"*/", fixLT)
		}
	case includeFragType:
		theOutFrag.addLine("pages."+theFragment.auxName+".set(parts);", setLineType)
		theOutFrag.addLine("pages."+theFragment.auxName+".fix();", fixLineType)
		theOutFrag.addLine("pages."+theFragment.auxName+".redisplay();", redisplayLineType)
		theOutFrag.addLine("pages."+theFragment.auxName+".refix();", refixLineType)
	default:
	}
	for _, theFragment := range theFragment.theFragments {
		subOutFrag := makeOutPage(theFragment.name)
		subOutFrag.codeFragment(theFragment, false)
		theOutFrag.SetLines = append(theOutFrag.SetLines, subOutFrag.SetLines...)
		theOutFrag.FixLines = append(theOutFrag.FixLines, subOutFrag.FixLines...)
		theOutFrag.RedisplayLines = append(theOutFrag.RedisplayLines, subOutFrag.RedisplayLines...)
		theOutFrag.Refixes = append(theOutFrag.Refixes, subOutFrag.Refixes...)
	}
	for _, theFragment := range theFragment.actionFragments {
		subOutFrag := makeOutPage(theFragment.name)
		subOutFrag.codeFragment(theFragment, false)
		theOutFrag.SetLines = append(theOutFrag.SetLines, subOutFrag.SetLines...)
		theOutFrag.Refixes = append(theOutFrag.Refixes, subOutFrag.Refixes...)
		theOutFrag.FixLines = append(theOutFrag.RedisplayLines, subOutFrag.FixLines...)
		theOutFrag.Refixes = append(theOutFrag.Refixes, subOutFrag.Refixes...)
	}

	switch theFragment.theFragType {
	case spanFragType, divFragType:
		if topLevel {
			tag := ifThenElse(theFragment.theFragType == spanFragType, "span", "div")
			theOutFrag.addLine("parts.push('</"+tag+">');", setLineType)
			theOutFrag.addLine("/*"+theFragment.name+"*/", redisplayLineType)
			theOutFrag.addLine("setHtml('"+fragName+"',parts.join(','));}", redisplayLineType)
		} else {
			theOutFrag.addLine("/*"+theFragment.name+"*/", setLineType)
			theOutFrag.addLine("setHtml('"+fragName+"',parts.join(','));}", setLineType)
		}
	case paraFragType:
		theOutFrag.addLine("parts.push('</p>');", mainLT)
	case linkFragType:
		//theOutFrag.addLine("}", fixLineType)
		theOutFrag.addLine("parts.push('</a>');", mainLT)
		theOutFrag.addLine("/*"+theFragment.name+"*/", redisplayLineType)
		if theFragment.auxName == "" {
			if topLevel {
				theOutFrag.addLine("setHtml('"+fragName+"',parts.join(','));}", mainLT)
			} else {
				theOutFrag.addLine("}", mainLT)
			}
			theOutFrag.addLine("/*"+theFragment.name+"*/", fixLT)
			theOutFrag.addLine("}", fixLT)
		}
	default:
	}
}

type fragCode struct {
}
type lineType int

const (
	setLineType lineType = iota
	fixLineType
	redisplayLineType
	refixLineType
)

func (theOutPage *outPage) addLine(line string, theLineType lineType) {
	switch theLineType {
	case setLineType:
		theOutPage.SetLines = append(theOutPage.SetLines, line)
	case fixLineType:
		theOutPage.FixLines = append(theOutPage.FixLines, line)
	case redisplayLineType:
		theOutPage.RedisplayLines = append(theOutPage.RedisplayLines, line)
	case refixLineType:
		theOutPage.Refixes = append(theOutPage.Refixes, line)
	}
}
