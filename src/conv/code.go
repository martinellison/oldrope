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
}

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
		theOutFragment := theFragment.code()
		theOutPage.InitLines = append(theOutPage.InitLines, theOutFragment.InitLines...)
		theOutPage.SetLines = append(theOutPage.SetLines, theOutFragment.SetLines...)
		theOutPage.FixLines = append(theOutPage.FixLines, theOutFragment.FixLines...)
	}
}

var autoLink int
var allSpace *regexp.Regexp

func init()                    { allSpace = regexp.MustCompile(`^[\s]*$`) }
func isAllSpace(s string) bool { return allSpace.MatchString(s) }

type outFragment struct {
	InitLines []string
	SetLines  []string
	FixLines  []string
}

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

//type fragCode struct {
//}
