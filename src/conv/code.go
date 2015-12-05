// code.go
package main

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

/* outData contains all the data required for generating pages */
type outData struct {
	Pages []*outPage
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

/* outPage contains the data for generating a single page */
type outPage struct {
	Name      string
	InitLines string
	SetLines  string
	FixLines  string
}

/* makeOutPage creates and initialises a new outPage*/
func makeOutPage(theName string) (theOutPage *outPage) {
	theOutPage = &outPage{Name: theName, InitLines: "", SetLines: "", FixLines: ""}
	return
}

/* codePage creates the code for a page */
func (theOutPage *outPage) codePage(thePage *page) {
	theInitLines := make([]*lineItem, 0, 0)
	theSetLines := make([]*lineItem, 0, 0)
	theFixLines := make([]*lineItem, 0, 0)
	for _, theFragment := range thePage.theFragments {
		theOutFragment := theFragment.code()
		theInitLines = append(theInitLines, theOutFragment.InitLines...)
		theSetLines = append(theSetLines, theOutFragment.SetLines...)
		theFixLines = append(theFixLines, theOutFragment.FixLines...)
	}
	theOutPage.InitLines = collapse(theInitLines)
	theOutPage.SetLines = collapse(theSetLines)
	theOutPage.FixLines = collapse(theFixLines)
}

/* autoLink */ var autoLink int

/* allSpace is a regular expression for a string of all white space */
var allSpace *regexp.Regexp

func init() { allSpace = regexp.MustCompile(`^[\s]*$`) }

/* isAllSpace tests whether a string is all white space  */
func isAllSpace(s string) bool { return allSpace.MatchString(s) }

/* escapeAll converts a string to HTML escape characters. All characters are converted, even a..z */
func escapeAll(inString string) string {
	outSegs := make([]string, 0, 0)
	for _, rune := range inString {
		outSegs = append(outSegs, fmt.Sprintf("&#%d;", rune))
	}
	return strings.Join(outSegs, "")
}

/* code converts a fragment to javascript code */
func (theFragment *fragment) code() (theOutFragment *outFragment) {
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
			theOutFragment.SetLines.addStrPush(escapeText)
		}
	case linkFragType:
		if theFragment.auxName == "" {
			theOutFragment.includeOutSubfragment(outOnPageLink(fragName, theFragment.theFragments))
			theOutFragment.includeOutSubfragment(outBlock(fragName, "span", theFragment.actionFragments))
		} else {
			theOutFragment.includeOutSubfragment(outOffPageLink(fragName, theFragment.auxName, theFragment.theFragments))
		}
	case htmlFragType:
		theOutFragment.SetLines.addStrPush("<" + comprText + ">")
	case includeFragType:
		theOutFragment.InitLines.addStr("pages." + theFragment.auxName + ".init(parts);")
		theOutFragment.SetLines.addStr("pages." + theFragment.auxName + ".display(parts);")
		theOutFragment.FixLines.addStr("pages." + theFragment.auxName + ".fix();")
	default:
	}
	return
}

/* outOffPageLink codes an out-of-page link */
func outOffPageLink(id string, targetPage string, subFragments []*fragment) (theOutFragment *outFragment) {
	theOutFragment = new(outFragment)
	theOutFragment.init()
	theOutFragment.SetLines.addStrPush("<a id='" + id + "'>")
	fixLine := "setClick('" + id + "', function(){setPage('" + targetPage + "');});"
	theOutFragment.FixLines.addStr(fixLine)
	theOutFragment.includeSubfragments(subFragments)
	theOutFragment.SetLines.addStrPush("</a>")
	return
}

/* outOnPageLink codes an on-page link */
func outOnPageLink(id string, subFragments []*fragment) (theOutFragment *outFragment) {
	theOutFragment = new(outFragment)
	theOutFragment.init()
	theOutFragment.InitLines.addStr("df." + id + "=false;")
	theOutFragment.SetLines.addStrPush("<a id='" + id + "'>")
	theOutFragment.FixLines.addStr("setClick('" + id + "', function(){df." + id + "=true; displayPage();});")
	theOutFragment.includeSubfragments(subFragments)
	theOutFragment.SetLines.addStrPush("</a>")
	return
}

/* outBlock codes  a span or div */
func outBlock(id string, tag string, subFragments []*fragment) (theOutFragment *outFragment) {
	theOutFragment = new(outFragment)
	theOutFragment.init()
	theOutFragment.SetLines.addStr("if (df." + id + ") {")
	theOutFragment.SetLines.addStrPush("<" + tag + ">")
	theOutFragment.FixLines.addStr("if (df." + id + ") {")
	theOutFragment.includeSubfragments(subFragments)
	theOutFragment.SetLines.addStrPush("</" + tag + ">")
	theOutFragment.SetLines.addStr("}")
	theOutFragment.FixLines.addStr("}")
	return
}

/* outFragment represents a fragment of code */ type outFragment struct {
	InitLines lineItemSet
	SetLines  lineItemSet
	FixLines  lineItemSet
}

/* init */ func (theOutFragment *outFragment) init() {
	theOutFragment.InitLines = make([]*lineItem, 0, 0)
	theOutFragment.SetLines = make([]*lineItem, 0, 0)
	theOutFragment.FixLines = make([]*lineItem, 0, 0)
}

/* includeSubfragments copies the line items from sub-fragments into a fragment */
func (theOutFragment *outFragment) includeSubfragments(subFragments []*fragment) {
	for _, theSubFragment := range subFragments {
		theOutSubfragment := theSubFragment.code()
		theOutFragment.includeOutSubfragment(theOutSubfragment)
	}
}

/* includeOutSubfragment copies the line items from a sub-fragment into a fragment */
func (theOutFragment *outFragment) includeOutSubfragment(theOutSubfragment *outFragment) {
	theOutFragment.InitLines.includeSubLineSet(&theOutSubfragment.InitLines)
	theOutFragment.SetLines.includeSubLineSet(&theOutSubfragment.SetLines)
	theOutFragment.FixLines.includeSubLineSet(&theOutSubfragment.FixLines)
}

/* a lineItemSet is a sequence of lineItems*/
type lineItemSet []*lineItem

/* addStr creates a line item without pushing */
func (theOutLineSet *lineItemSet) addStr(theString string) {
	*theOutLineSet = append(*theOutLineSet, &lineItem{theText: theString})
}

/* addStrPush creates a line item with pushing */
func (theOutLineSet *lineItemSet) addStrPush(theString string) {
	*theOutLineSet = append(*theOutLineSet, &lineItem{theText: theString, pushing: true})
}

/* includeSubLineSet appends a lineItemSet */
func (theOutLineSet *lineItemSet) includeSubLineSet(theSubOutLineSet *lineItemSet) {
	*theOutLineSet = append(*theOutLineSet, *theSubOutLineSet...)
}

/* collapse combines the line items into a string */
func collapse(theLineItems []*lineItem) string {
	outParts := make([]string, 0, len(theLineItems))
	pushing := false
	for _, outPart := range theLineItems {
		if !pushing && outPart.pushing {
			outParts = append(outParts, "parts.push(\"")
		} else if pushing && !outPart.pushing {
			outParts = append(outParts, "\");")
		}
		outParts = append(outParts, outPart.theText)
		pushing = outPart.pushing
	}
	if pushing {
		outParts = append(outParts, "\");")
	}
	return strings.Join(outParts, "")
}

/* a lineItem represents some text that will be output. If pushing is set, it needs to be quoted and pushed onto the parts array. */
type lineItem struct {
	theText string
	pushing bool
}

func (theLineItem lineItem) String() string {
	if theLineItem.pushing {
		return "*push* " + theLineItem.theText
	}
	return theLineItem.theText
}
