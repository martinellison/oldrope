// code_test.go
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCode1(t *testing.T) {
	assert := assert.New(t)
	theLineItemSet := make(lineItemSet, 0, 0)
	theLineItemSet.addStr("alpha")
	theLineItemSet.addStrPush("beta")
	res := collapse(theLineItemSet)
	assert.Equal("alphaparts.push(\"beta\");", res)
}
func TestCode2(t *testing.T) {
	assert := assert.New(t)
	theLineItemSet := make(lineItemSet, 0, 0)
	theSubLineItemSet := make(lineItemSet, 0, 0)
	theLineItemSet.addStr("alpha();")
	theSubLineItemSet.addStrPush("beta")
	theLineItemSet.includeSubLineSet(&theSubLineItemSet)
	res := collapse(theLineItemSet)
	assert.Equal("alpha();parts.push(\"beta\");", res)
}

func TestAllSpace1(t *testing.T) {
	assert := assert.New(t)
	assert.True(isAllSpace(""))
	assert.True(isAllSpace(" "))
	assert.True(isAllSpace("\n"))
	assert.True(isAllSpace("\t"))
	assert.True(isAllSpace("  "))
	assert.False(isAllSpace("a"))
	assert.False(isAllSpace("aaaa"))
}

func TestEscapeAll1(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("&#32;", escapeAll(" "))
	assert.Equal("&#97;&#98;&#99;", escapeAll("abc"))
}

func TestCodeOutPage1(t *testing.T) {
	assert := assert.New(t)
	op := makeOutPage("fred")
	assert.Equal("fred", op.Name)
	assert.Equal("", op.InitLines)
	assert.Equal("", op.SetLines)
	assert.Equal("", op.FixLines)
}
func TestCodeOutOffPageLink1(t *testing.T) {
	assert := assert.New(t)
	theOutFragment := outOffPageLink("fred", "bill", make([]*fragment, 0, 0))
	assert.Equal(0, len(theOutFragment.InitLines))
	assert.Equal("", collapse(theOutFragment.InitLines))
	assert.Equal(2, len(theOutFragment.SetLines))
	assert.Equal("parts.push(\"<a id='fred'></a>\");", collapse(theOutFragment.SetLines))
	assert.Equal(1, len(theOutFragment.FixLines))
	assert.Equal("setClick('fred', function(){setPage('bill');});", collapse(theOutFragment.FixLines))
}
func TestCodeOutOnPageLink1(t *testing.T) {
	assert := assert.New(t)
	theOutFragment := outOnPageLink("fred", make([]*fragment, 0, 0))
	assert.Equal(1, len(theOutFragment.InitLines))
	assert.Equal("df.fred=false;", collapse(theOutFragment.InitLines))
	assert.Equal(2, len(theOutFragment.SetLines))
	assert.Equal("parts.push(\"<a id='fred'></a>\");", collapse(theOutFragment.SetLines))
	assert.Equal(1, len(theOutFragment.FixLines))
	assert.Equal("setClick('fred', function(){df.fred=true; displayPage();});", collapse(theOutFragment.FixLines))
}
func TestCodeOutBlock1(t *testing.T) {
	assert := assert.New(t)
	theOutFragment := outBlock("fred", "span", make([]*fragment, 0, 0))
	assert.Equal(0, len(theOutFragment.InitLines))
	assert.Equal("", collapse(theOutFragment.InitLines))
	assert.Equal(4, len(theOutFragment.SetLines))
	assert.Equal("if (df.fred) {parts.push(\"<span></span>\");}", collapse(theOutFragment.SetLines))
	assert.Equal(2, len(theOutFragment.FixLines))
	assert.Equal("if (df.fred) {}", collapse(theOutFragment.FixLines))
}

func TestCodeCode1(t *testing.T) {
	assert := assert.New(t)
	theFragment := &fragment{theFragType: textFragType, name: "fred", text: "bill", theFragments: make([]*fragment, 0, 0), actionFragments: make([]*fragment, 0, 0), auxName: ""}
	theOutFragment := theFragment.code()
	assert.Equal(0, len(theOutFragment.InitLines))
	assert.Equal(1, len(theOutFragment.SetLines))
	assert.Equal("parts.push(\"bill\");", collapse(theOutFragment.SetLines))
	assert.Equal(0, len(theOutFragment.FixLines))
}
func TestCodeCode2(t *testing.T) {
	assert := assert.New(t)
	theFragment := &fragment{theFragType: spanFragType, name: "fred", text: "bill", theFragments: make([]*fragment, 0, 0), actionFragments: make([]*fragment, 0, 0), auxName: ""}
	theOutFragment := theFragment.code()
	assert.Equal(0, len(theOutFragment.InitLines))
	assert.Equal(4, len(theOutFragment.SetLines))
	assert.Equal("if (df.fred) {parts.push(\"<span></span>\");}", collapse(theOutFragment.SetLines))
	assert.Equal(2, len(theOutFragment.FixLines))
	assert.Equal("if (df.fred) {}", collapse(theOutFragment.FixLines))
}
func TestCodeCode3(t *testing.T) {
	assert := assert.New(t)
	theFragment := &fragment{theFragType: divFragType, name: "fred", text: "bill", theFragments: make([]*fragment, 0, 0), actionFragments: make([]*fragment, 0, 0), auxName: ""}
	theOutFragment := theFragment.code()
	assert.Equal(0, len(theOutFragment.InitLines))
	assert.Equal(4, len(theOutFragment.SetLines))
	assert.Equal("if (df.fred) {parts.push(\"<div></div>\");}", collapse(theOutFragment.SetLines))
	assert.Equal(2, len(theOutFragment.FixLines))
	assert.Equal("if (df.fred) {}", collapse(theOutFragment.FixLines))
}
func TestCodeCode4(t *testing.T) {
	assert := assert.New(t)
	theFragment := &fragment{theFragType: paraFragType, name: "fred", text: "bill", theFragments: make([]*fragment, 0, 0), actionFragments: make([]*fragment, 0, 0), auxName: ""}
	theOutFragment := theFragment.code()
	assert.Equal(0, len(theOutFragment.InitLines))
	assert.Equal(4, len(theOutFragment.SetLines))
	assert.Equal("if (df.fred) {parts.push(\"<p></p>\");}", collapse(theOutFragment.SetLines))
	assert.Equal(2, len(theOutFragment.FixLines))
	assert.Equal("if (df.fred) {}", collapse(theOutFragment.FixLines))
}
func TestCodeCode5(t *testing.T) {
	assert := assert.New(t)
	theFragment := &fragment{theFragType: jsCodeFragType, name: "fred", text: "bill", theFragments: make([]*fragment, 0, 0), actionFragments: make([]*fragment, 0, 0), auxName: ""}
	theOutFragment := theFragment.code()
	assert.Equal(0, len(theOutFragment.InitLines))
	assert.Equal(1, len(theOutFragment.SetLines))
	assert.Equal("bill", collapse(theOutFragment.SetLines))
	assert.Equal(0, len(theOutFragment.FixLines))
}
func TestCodeCode6(t *testing.T) {
	assert := assert.New(t)
	theFragment := &fragment{theFragType: jsExprFragType, name: "fred", text: "bill", theFragments: make([]*fragment, 0, 0), actionFragments: make([]*fragment, 0, 0), auxName: ""}
	theOutFragment := theFragment.code()
	assert.Equal(0, len(theOutFragment.InitLines))
	assert.Equal(1, len(theOutFragment.SetLines))
	assert.Equal("parts.push(bill);", collapse(theOutFragment.SetLines))
	assert.Equal(0, len(theOutFragment.FixLines))
}
func TestCodeCode7(t *testing.T) {
	assert := assert.New(t)
	theFragment := &fragment{theFragType: linkFragType, name: "fred", text: "bill", theFragments: make([]*fragment, 0, 0), actionFragments: make([]*fragment, 0, 0), auxName: ""}
	theOutFragment := theFragment.code()
	assert.Equal(1, len(theOutFragment.InitLines))
	assert.Equal("df.fred=false;", collapse(theOutFragment.InitLines))
	assert.Equal(6, len(theOutFragment.SetLines))
	assert.Equal("parts.push(\"<a id='fred'></a>\");if (df.fred) {parts.push(\"<span></span>\");}", collapse(theOutFragment.SetLines))
	assert.Equal(3, len(theOutFragment.FixLines))
	assert.Equal("setClick('fred', function(){df.fred=true; displayPage();});if (df.fred) {}", collapse(theOutFragment.FixLines))
}
func TestCodeCode8(t *testing.T) {
	assert := assert.New(t)
	theFragment := &fragment{theFragType: htmlFragType, name: "fred", text: "bill", theFragments: make([]*fragment, 0, 0), actionFragments: make([]*fragment, 0, 0), auxName: ""}
	theOutFragment := theFragment.code()
	assert.Equal(0, len(theOutFragment.InitLines))
	assert.Equal(1, len(theOutFragment.SetLines))
	assert.Equal("parts.push(\"<bill>\");", collapse(theOutFragment.SetLines))
	assert.Equal(0, len(theOutFragment.FixLines))
}
func TestCodeCode9(t *testing.T) {
	assert := assert.New(t)
	theFragment := &fragment{theFragType: includeFragType, name: "fred", text: "bill", theFragments: make([]*fragment, 0, 0), actionFragments: make([]*fragment, 0, 0), auxName: "jane"}
	theOutFragment := theFragment.code()
	assert.Equal(1, len(theOutFragment.InitLines))
	assert.Equal("pages.jane.init(parts);", collapse(theOutFragment.InitLines))
	assert.Equal(1, len(theOutFragment.SetLines))
	assert.Equal("pages.jane.display(parts);", collapse(theOutFragment.SetLines))
	assert.Equal(1, len(theOutFragment.FixLines))
	assert.Equal("pages.jane.fix();", collapse(theOutFragment.FixLines))
}
func TestCodeCode10(t *testing.T) {
	assert := assert.New(t)
	theFragment := &fragment{theFragType: linkFragType, name: "fred", text: "bill", theFragments: make([]*fragment, 0, 0), actionFragments: make([]*fragment, 0, 0), auxName: "jane"}
	theOutFragment := theFragment.code()
	assert.Equal(0, len(theOutFragment.InitLines))
	assert.Equal(2, len(theOutFragment.SetLines))
	assert.Equal("parts.push(\"<a id='fred'></a>\");", collapse(theOutFragment.SetLines))
	assert.Equal(1, len(theOutFragment.FixLines))
	assert.Equal("setClick('fred', function(){setPage('jane');});", collapse(theOutFragment.FixLines))
}
