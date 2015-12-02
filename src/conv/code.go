// code.go
package main

import (
	"fmt"
	"html"
	"log"
	"regexp"
	"strings"
)

/* outData contains all the data required for generating pages */
type outData struct {
	Pages []*outPage
}

/* outPage contains the data for generating a single page */
type outPage struct {
	Name      string
	InitLines []string
	SetLines  []string
	FixLines  []string
}

/* makeGenData builds theOutData with page data */
func makeGenData() {
	theOutData.Pages = make([]*outPage, 0, len(thePageSet))
	for _, page := range thePageSet {
		outPage := makeOutPage(page.theName)
		outPage.codePage(&page)
		theOutData.Pages = append(theOutData.Pages, outPage)
	}
}

/* makeOutPage creates and initialises a new outPage*/ func makeOutPage(theName string) (theOutPage *outPage) {
	theOutPage = &outPage{Name: theName, SetLines: make([]string, 0, 0), FixLines: make([]string, 0, 0)}
	return
}

/* */ func (theOutPage *outPage) codePage(thePage *page) {
	for _, theFragment := range thePage.theFragments {
		theOutFragment := theFragment.code()
		theOutPage.InitLines = append(theOutPage.InitLines, theOutFragment.InitLines.lines()...)
		theOutPage.SetLines = append(theOutPage.SetLines, theOutFragment.SetLines.lines()...)
		theOutPage.FixLines = append(theOutPage.FixLines, theOutFragment.FixLines.lines()...)
	}
}

/* */ var autoLink int

/* */ var allSpace *regexp.Regexp

/* */ func init() { allSpace = regexp.MustCompile(`^[\s]*$`) }

/* */ func isAllSpace(s string) bool { return allSpace.MatchString(s) }

/* */ func escapeAll(inString string) string {
	outSegs := make([]string, 0, 0)
	for _, rune := range inString {
		outSegs = append(outSegs, fmt.Sprintf("&#%d;", rune))
	}
	return strings.Join(outSegs, "")
}

/* */ func (theFragment *fragment) code() (theOutFragment *outFragment) {
	theOutFragment = new(outFragment)
	theOutFragment.init()
	comprText := compress(theFragment.text)
	escapeText := ""
	if hashText {
		escapeText = escapeAll(comprText)
	} else {
		escapeText = html.EscapeString(comprText)
	}
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
		theOutFragment.SetLines.addStr(comprText)
	case jsExprFragType:
		theOutFragment.SetLines.addStr("parts.push(" + comprText + ");")
	case textFragType:
		if !isAllSpace(escapeText) {
			theOutFragment.SetLines.addStrClose("parts.push(\"", escapeText, "\");")
		}
	case linkFragType:
		if theFragment.auxName == "" {
			theOutFragment.includeOutSubfragment(outOnPageLink(fragName, theFragment.theFragments))
			theOutFragment.includeOutSubfragment(outBlock(fragName, "span", theFragment.actionFragments))
		} else {
			theOutFragment.includeOutSubfragment(outOffPageLink(fragName, theFragment.auxName, theFragment.theFragments))
		}
	case htmlFragType:
		theOutFragment.SetLines.addStrClose("parts.push(\"", "<"+comprText+">", "\");")
	case includeFragType:
		theOutFragment.InitLines.addStr("pages." + theFragment.auxName + ".init(parts);")
		theOutFragment.SetLines.addStr("pages." + theFragment.auxName + ".display(parts);")
		theOutFragment.FixLines.addStr("pages." + theFragment.auxName + ".fix();")
	default:
	}
	return
}

/* */ func outOffPageLink(id string, targetPage string, subFragments []*fragment) (theOutFragment *outFragment) {
	theOutFragment = new(outFragment)
	theOutFragment.init()
	theOutFragment.SetLines.addStrClose("parts.push(\"", "<a id='"+id+"'>", "\");")
	fixLine := "setClick('" + id + "', function(){setPage('" + targetPage + "');});"
	theOutFragment.FixLines.addStr(fixLine)
	theOutFragment.includeSubfragments(subFragments)
	theOutFragment.SetLines.addStrClose("parts.push(\"", "</a>", "\");")
	return
}

/* */ func outOnPageLink(id string, subFragments []*fragment) (theOutFragment *outFragment) {
	theOutFragment = new(outFragment)
	theOutFragment.init()
	theOutFragment.InitLines.addStr("df." + id + "=false;")
	theOutFragment.SetLines.addStrClose("parts.push(\"", "<a id='"+id+"'>", "\");")
	theOutFragment.FixLines.addStr("setClick('" + id + "', function(){df." + id + "=true; displayPage();});")
	theOutFragment.includeSubfragments(subFragments)
	theOutFragment.SetLines.addStrClose("parts.push(\"", "</a>", "\");")
	return
}

/* */ func outBlock(id string, tag string, subFragments []*fragment) (theOutFragment *outFragment) {
	theOutFragment = new(outFragment)
	theOutFragment.init()
	theOutFragment.SetLines.addStr("if (df." + id + ") {parts.push(\"<" + tag + ">\");")
	theOutFragment.FixLines.addStr("if (df." + id + ") {")
	theOutFragment.includeSubfragments(subFragments)
	theOutFragment.SetLines.addStrClose("parts.push(\"", "</"+tag+">", "\");")
	theOutFragment.SetLines.addStr("}")
	theOutFragment.FixLines.addStr("}")
	return
}

/* */ type outFragment struct {
	InitLines outLineSet
	SetLines  outLineSet
	FixLines  outLineSet
}

/* */ func (theOutFragment *outFragment) init() {
	theOutFragment.InitLines.init()
	theOutFragment.SetLines.init()
	theOutFragment.FixLines.init()
}

/* */ func (theOutFragment *outFragment) includeSubfragments(subFragments []*fragment) {
	for _, theSubFragment := range subFragments {
		theOutSubfragment := theSubFragment.code()
		theOutFragment.includeOutSubfragment(theOutSubfragment)
	}
}

/* */ func (theOutFragment *outFragment) includeOutSubfragment(theOutSubfragment *outFragment) {
	theOutFragment.InitLines.includeSubLineSet(&theOutSubfragment.InitLines)
	theOutFragment.SetLines.includeSubLineSet(&theOutSubfragment.SetLines)
	theOutFragment.FixLines.includeSubLineSet(&theOutSubfragment.FixLines)
}

/* */ type outLineSet struct {
	theStrings []string
	theCloser  string
}

/* */ func (theOutLineSet *outLineSet) init() {
	theOutLineSet.theStrings = make([]string, 0, 0)
	theOutLineSet.theCloser = ""
}

/* */ func (theOutLineSet *outLineSet) addStr(theString string) {
	theOutLineSet.close()
	theOutLineSet.theStrings = append(theOutLineSet.theStrings, theString)
}

/* */ func (theOutLineSet *outLineSet) close() {
	if theOutLineSet.theCloser == "" {
		return
	}
	theOutLineSet.theStrings = append(theOutLineSet.theStrings, theOutLineSet.theCloser)
	theOutLineSet.theCloser = ""
}

/* */ func (theOutLineSet *outLineSet) addStrClose(theStarter, theString, theCloser string) {
	if logging {
		log.Printf("asc '%s'/'%s'=%t ", theCloser, theOutLineSet.theCloser, theOutLineSet.theCloser == theCloser)
	}
	if theOutLineSet.theCloser == theCloser {
		theOutLineSet.theStrings = append(theOutLineSet.theStrings, theString)
		return
	}
	theOutLineSet.close()
	theOutLineSet.theStrings = append(theOutLineSet.theStrings, theStarter, theString)
	theOutLineSet.theCloser = theCloser
}

/* */ func (theOutLineSet *outLineSet) lines() []string {
	theOutLineSet.close()
	return theOutLineSet.theStrings
}

/* */ func (theOutLineSet *outLineSet) includeSubLineSet(theSubOutLineSet *outLineSet) {
	theSubOutLineSet.close()
	theOutLineSet.close()
	theOutLineSet.theStrings = append(theOutLineSet.theStrings, theSubOutLineSet.theStrings...)
}

///* */ type fragCode struct {
//}
