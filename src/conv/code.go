// code.go
package main

import (
	"fmt"
	"html"
	"regexp"
)

type outData struct {
	Pages []*outPage
}
type outPage struct {
	Name      string
	InitLines []string
	SetLines  []string
	FixLines  []string
	//	RedisplayLines []string
	//	Refixes        []string
}

/*type fix struct {
	Name string
	Code string
}*/
func makeGenData() {
	theOutData.Pages = make([]*outPage, 0, len(thePageSet))
	for _, page := range thePageSet {
		outPage := makeOutPage(page.theName)
		outPage.codePage(&page)
		theOutData.Pages = append(theOutData.Pages, outPage)
	}
}
func makeOutPage(theName string) (theOutPage *outPage) {
	theOutPage = &outPage{Name: theName, SetLines: make([]string, 0, 0), FixLines: make([]string, 0, 0)}
	return
}
func (theOutPage *outPage) codePage(thePage *page) {
	for _, theFragment := range thePage.theFragments {
		//subOutFrag := makeOutPage(theFragment.name)
		//subOutFrag.codeFragment(theFragment, true)
		theOutFragment := theFragment.code()
		theOutPage.InitLines = append(theOutPage.InitLines, theOutFragment.InitLines...)
		theOutPage.SetLines = append(theOutPage.SetLines, theOutFragment.SetLines...)
		theOutPage.FixLines = append(theOutPage.FixLines, theOutFragment.FixLines...)
		//	theOutPage.RedisplayLines = append(theOutPage.RedisplayLines, subOutFrag.RedisplayLines...)
		//	theOutPage.Refixes = append(theOutPage.Refixes, subOutFrag.Refixes...)
	}
}

var autoLink int
var allSpace *regexp.Regexp

func init() { allSpace = regexp.MustCompile(`^[\s]*$`) }

type outFragment struct {
	InitLines []string
	SetLines  []string
	FixLines  []string
}

//func (theOutFrag *outPage) codeFragment(theFragment *fragment, topLevel bool) {
//	//theOutFrag = new(outPage)
//	comprText := compress(theFragment.text)
//	escapeText := html.EscapeString(comprText)
//	fragName := theFragment.name
//	if fragName == "" {
//		autoLink++
//		fragName = fmt.Sprintf("z%d", autoLink)
//	}
//	theOutFrag.Name = fragName
//	fragNameExtend := fragName + "-x"
//	fragIdAttr := " id=\"" + fragName + "\" "
//	fragIdExtendAttr := " id=\"" + fragNameExtend + "\" "
//	mainLT := setLineType
//	fixLT := fixLineType
//	if !topLevel {
//		mainLT = redisplayLineType
//		fixLT = refixLineType
//	}
//	//	subset := true
//	switch theFragment.theFragType {
//	case spanFragType, divFragType:
//		if topLevel {
//			tag := ifThenElse(theFragment.theFragType == spanFragType, "span", "div")
//			spanText := "parts.push('<" + tag + fragIdAttr + ">" + escapeText + "');"
//			theOutFrag.addLine(spanText, setLineType)
//			theOutFrag.addLine("ld.s"+fragName+"=false;", setLineType)
//			theOutFrag.addLine("if (ld.s"+fragName+") {parts=[];", redisplayLineType)
//		} else {
//			theOutFrag.addLine("if (ld.s"+fragName+") {parts=[];", redisplayLineType)
//		}
//	case paraFragType:
//		paraText := "parts.push('<p" + fragIdAttr + ">" + escapeText + "');"
//		theOutFrag.addLine(paraText, mainLT)
//	case jsCodeFragType:
//		theOutFrag.addLine(comprText, mainLT)
//	case jsExprFragType:
//		exprText := "parts.push(" + comprText + ");"
//		theOutFrag.addLine(exprText, mainLT)
//	case textFragType:
//		if !allSpace.MatchString(escapeText) {
//			textText := "parts.push('" + escapeText + "');"
//			theOutFrag.addLine(textText, mainLT)
//		}
//	case htmlFragType:
//		if !allSpace.MatchString(comprText) {
//			htmlText := "parts.push('<" + comprText + ">');"
//			theOutFrag.addLine(htmlText, mainLT)
//		}
//	case linkFragType:
//		textText := "parts.push('<a" + fragIdExtendAttr + ">');"
//		theOutFrag.addLine(textText, mainLT)
//		code := ""
//		if theFragment.auxName != "" {
//			code = " setPage('" + theFragment.auxName + "'); displayPage();"
//		} else {
//			code = " ld.s" + fragName + "=true; displayPage();"
//		}
//		//	if topLevel {
//		linkFix := "setClick('" + fragNameExtend + "', function(){" + code + "});"
//		theOutFrag.addLine(linkFix, fixLT)
//		//	} else {
//		//		theOutFrag.addLine("if (ld.s"+fragName+") {", redisplayLineType)
//		//	}
//		if theFragment.auxName == "" {
//			theOutFrag.addLine("if (ld.s"+fragName+") {", mainLT)
//			theOutFrag.addLine("if (ld.s"+fragName+") {", fixLT)
//			theOutFrag.addLine("/*"+theFragment.name+"*/", fixLT)
//		}
//	case includeFragType:
//		theOutFrag.addLine("pages."+theFragment.auxName+".set(parts);", setLineType)
//		theOutFrag.addLine("pages."+theFragment.auxName+".fix();", fixLineType)
//		theOutFrag.addLine("pages."+theFragment.auxName+".redisplay();", redisplayLineType)
//		theOutFrag.addLine("pages."+theFragment.auxName+".refix();", refixLineType)
//	default:
//	}
//	for _, theFragment := range theFragment.theFragments {
//		subOutFrag := makeOutPage(theFragment.name)
//		subOutFrag.codeFragment(theFragment, false)
//		theOutFrag.SetLines = append(theOutFrag.SetLines, subOutFrag.SetLines...)
//		theOutFrag.FixLines = append(theOutFrag.FixLines, subOutFrag.FixLines...)
//		//		theOutFrag.RedisplayLines = append(theOutFrag.RedisplayLines, subOutFrag.RedisplayLines...)
//		//		theOutFrag.Refixes = append(theOutFrag.Refixes, subOutFrag.Refixes...)
//	}
//	for _, theFragment := range theFragment.actionFragments {
//		subOutFrag := makeOutPage(theFragment.name)
//		subOutFrag.codeFragment(theFragment, false)
//		//		theOutFrag.RedisplayLines = append(theOutFrag.RedisplayLines, subOutFrag.RedisplayLines...)
//		//		theOutFrag.Refixes = append(theOutFrag.Refixes, subOutFrag.Refixes...)
//		//		theOutFrag.RedisplayLines = append(theOutFrag.RedisplayLines, subOutFrag.RedisplayLines...)
//		//		theOutFrag.Refixes = append(theOutFrag.Refixes, subOutFrag.Refixes...)
//	}

//	switch theFragment.theFragType {
//	case spanFragType, divFragType:
//		if topLevel {
//			tag := ifThenElse(theFragment.theFragType == spanFragType, "span", "div")
//			theOutFrag.addLine("parts.push('</"+tag+">');", setLineType)
//			theOutFrag.addLine("/*"+theFragment.name+"*/", redisplayLineType)
//			theOutFrag.addLine("setHtml('"+fragName+"',parts.join(','));}", redisplayLineType)
//		} else {
//			theOutFrag.addLine("/*"+theFragment.name+"*/", setLineType)
//			theOutFrag.addLine("setHtml('"+fragName+"',parts.join(','));}", setLineType)
//		}
//	case paraFragType:
//		theOutFrag.addLine("parts.push('</p>');", mainLT)
//	case linkFragType:
//		//theOutFrag.addLine("}", fixLineType)
//		theOutFrag.addLine("parts.push('</a>');", mainLT)
//		theOutFrag.addLine("/*"+theFragment.name+"*/", redisplayLineType)
//		if theFragment.auxName == "" {
//			if topLevel {
//				theOutFrag.addLine("setHtml('"+fragName+"',parts.join(','));}", mainLT)
//			} else {
//				theOutFrag.addLine("}", mainLT)
//			}
//			theOutFrag.addLine("/*"+theFragment.name+"*/", fixLT)
//			theOutFrag.addLine("}", fixLT)
//		}
//	default:
//	}
//}
func (theFragment *fragment) code() (theOutFragment *outFragment) {
	theOutFragment = new(outFragment)
	comprText := compress(theFragment.text)
	escapeText := html.EscapeString(comprText)
	fragName := theFragment.name
	if fragName == "" {
		autoLink++
		fragName = fmt.Sprintf("z%d", autoLink)
	}
	switch theFragment.theFragType {
	case spanFragType:
		theOutFragment.includeOutSubfragment(outBlock(fragName, "span", theFragment.actionFragments))
	case divFragType:
		theOutFragment.includeOutSubfragment(outBlock(fragName, "div", theFragment.actionFragments))
	case paraFragType:

		theOutFragment.includeOutSubfragment(outBlock(fragName, "p", theFragment.actionFragments))
	case jsCodeFragType:
		addStr(&theOutFragment.SetLines, comprText)
	case jsExprFragType:
		addStr(&theOutFragment.SetLines, "parts.push("+comprText+");")
	case textFragType:
		addStr(&theOutFragment.SetLines, "parts.push('"+escapeText+"');")
	case linkFragType:
		if theFragment.auxName == "" {
			theOutFragment.includeOutSubfragment(outOnPageLink(fragName, theFragment.theFragments))
			theOutFragment.includeOutSubfragment(outBlock(fragName, "span", theFragment.actionFragments))
		} else {
			theOutFragment.includeOutSubfragment(outOffPageLink(fragName, theFragment.auxName, theFragment.theFragments))
		}
	case htmlFragType:
		addStr(&theOutFragment.SetLines, "parts.push('<"+comprText+">');")
	case includeFragType:
		addStr(&theOutFragment.InitLines, "pages."+theFragment.auxName+".init(parts);")
		addStr(&theOutFragment.SetLines, "pages."+theFragment.auxName+".display(parts);")
		addStr(&theOutFragment.FixLines, "pages."+theFragment.auxName+".fix();")
	default:
	}
	return
}
func outOffPageLink(id string, targetPage string, subFragments []*fragment) (theOutFragment *outFragment) {
	theOutFragment = new(outFragment)
	addStr(&theOutFragment.SetLines, "parts.push('<a id=\\'"+id+"\\'>');")
	fixLine := "setClick('" + id + "', function(){setPage('" + targetPage + "');});"
	addStr(&theOutFragment.FixLines, fixLine)
	theOutFragment.includeSubfragments(subFragments)
	addStr(&theOutFragment.SetLines, "parts.push(\"</a>\");")
	return
}
func outOnPageLink(id string, subFragments []*fragment) (theOutFragment *outFragment) {
	theOutFragment = new(outFragment)
	addStr(&theOutFragment.InitLines, "df."+id+"=false;")
	addStr(&theOutFragment.SetLines, "parts.push('<a id=\\'"+id+"\\'>');")
	addStr(&theOutFragment.FixLines, "setClick('"+id+"', function(){df."+id+"=true; displayPage();});")
	theOutFragment.includeSubfragments(subFragments)
	addStr(&theOutFragment.SetLines, "parts.push(\"</a>\");")
	return
}
func outBlock(id string, tag string, subFragments []*fragment) (theOutFragment *outFragment) {
	theOutFragment = new(outFragment)
	addStr(&theOutFragment.SetLines, "if (df."+id+") {parts.push(\"<"+tag+">\");")
	addStr(&theOutFragment.FixLines, "if (df."+id+") {")
	theOutFragment.includeSubfragments(subFragments)
	addStr(&theOutFragment.SetLines, "parts.push(\"</"+tag+">\");}")
	addStr(&theOutFragment.FixLines, "}")
	return
}
func addStr(theStrings *[]string, theString string) {
	*theStrings = append(*theStrings, theString)
}
func (theOutFragment *outFragment) includeSubfragments(subFragments []*fragment) {
	for _, theSubFragment := range subFragments {
		theOutSubfragment := theSubFragment.code()
		theOutFragment.includeOutSubfragment(theOutSubfragment)
	}
}
func (theOutFragment *outFragment) includeOutSubfragment(theOutSubfragment *outFragment) {
	theOutFragment.InitLines = append(theOutFragment.InitLines, theOutSubfragment.InitLines...)
	theOutFragment.SetLines = append(theOutFragment.SetLines, theOutSubfragment.SetLines...)
	theOutFragment.FixLines = append(theOutFragment.FixLines, theOutSubfragment.FixLines...)
}

type fragCode struct {
}

//type lineType int

//const (
//	setLineType lineType = iota
//	fixLineType
//	redisplayLineType
//	refixLineType
//)

//func (theOutPage *outPage) addLine(line string, theLineType lineType) {
//	switch theLineType {
//	case setLineType:
//		theOutPage.SetLines = append(theOutPage.SetLines, line)
//	case fixLineType:
//		theOutPage.FixLines = append(theOutPage.FixLines, line)
//		//	case redisplayLineType:
//		//		theOutPage.RedisplayLines = append(theOutPage.RedisplayLines, line)
//		//	case refixLineType:
//		//		theOutPage.Refixes = append(theOutPage.Refixes, line)
//	}
//}
