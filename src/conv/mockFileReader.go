// mockFileReader.go
package main

import "github.com/stretchr/testify/assert"

type mockFileReader struct {
	inlines   []string
	nextIndex int
}

func makeMockFileReader(theInlines []string) *mockFileReader {
	return &mockFileReader{inlines: theInlines}
}
func (theMockFileReader *mockFileReader) getLine(assert *assert.Assertions) (theLine string, eof bool) {
	numLines := len(theMockFileReader.inlines)
	assert.True(theMockFileReader.nextIndex < numLines+1)
	if theMockFileReader.nextIndex == numLines {
		theMockFileReader.nextIndex++
		theLine = ""
		eof = true
		return
	}
	//log.Printf("using mock line %d of %d", theMockFileReader.nextIndex, numLines)
	theLine = theMockFileReader.inlines[theMockFileReader.nextIndex]
	theMockFileReader.nextIndex++
	return
}
func (theMockFileReader *mockFileReader) verify(assert *assert.Assertions) {
	assert.True(theMockFileReader.nextIndex == len(theMockFileReader.inlines)+1)
}
func sendInlines(theMockFileReader *mockFileReader, assert *assert.Assertions) {
	//logging = true
	theLine, eof := theMockFileReader.getLine(assert)
	lineNumber := 1
	for !eof {
		//log.Printf("scan line %d", lineNumber)
		lineChan <- scanLine{text: theLine, number: lineNumber, eof: false}
		theLine, eof = theMockFileReader.getLine(assert)
		lineNumber++
	}
	//log.Print("end of scan lines")
	lineChan <- scanLine{text: "", number: lineNumber, eof: true}
}
