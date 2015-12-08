// Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.

// tokens.go (tokeniser (scanner) for input)
package main

import (
	"fmt"
	"log"
	"os"
)

/* token describes a token from the input */
type token struct {
	theType    tokenType
	text       string
	lineNumber int
}

/* */ type tokenType int

const (
	eofTokenType tokenType = iota
	jsCodeTokenType
	jsExprTokenType
	textTokenType
	numberTokenType
	identTokenType
	delimTokenType
	htmlTokenType
)

/* tokenChan is a channel of tokens */
var tokenChan chan token

/* a scanLineState is the current state of the tokeniser */
type scanLineState int

/* possible states of the tokeniser */
const (
	textState scanLineState = iota
	commentState
	jsCodeState
	jsExprState
	convState
	htmlState
)

func (theTokeniser *tokeniser) init() { theTokeniser.state = textState }

/* getTokens is the tokeniser. It divides the input text into tokens. */
func (theTokeniser *tokeniser) getTokens() {
	defer func() {
		rec := recover()
		if rec == nil {
			return
		}
		os.Stderr.WriteString(fmt.Sprintf("getTokens: internal error: %v", rec))
	}()
	theTokeniser.setTokenState(textState)
	for {
		theTokeniser.line = <-lineChan
		logfIfLogging("%d inline: '%s' (%d) eof: %t", theTokeniser.line.number, theTokeniser.line.text, len(theTokeniser.line.text), theTokeniser.line.eof)
		theTokeniser.lineLen = len(theTokeniser.line.text)
		theTokeniser.charPos = 0
		for theTokeniser.charPos < theTokeniser.lineLen {
			//	log.Printf("%d: scan char is: '%s'", theTokeniser.charPos, scanChars(1))
			switch theTokeniser.state {
			case textState:
				switch theTokeniser.scanChars(2) {
				case "/*":
					theTokeniser.emitToken(textTokenType)
					theTokeniser.state = commentState
					theTokeniser.charPos += 2
				case "$/":
					theTokeniser.emitToken(textTokenType)
					theTokeniser.state = jsCodeState
					theTokeniser.charPos += 2
				case "$(":
					theTokeniser.emitToken(textTokenType)
					theTokeniser.state = jsExprState
					theTokeniser.charPos += 2
				case "$[":
					theTokeniser.emitToken(textTokenType)
					theTokeniser.state = convState
					theTokeniser.charPos += 2
				case "$<":
					theTokeniser.emitToken(textTokenType)
					theTokeniser.state = htmlState
					theTokeniser.charPos += 2
				default:
					theTokeniser.charToToken()
				}
			case commentState:
				if theTokeniser.scanChars(2) == "*/" {
					//if logging {if logging {log.Print("end of comment")
					theTokeniser.setTokenState(textState)
					theTokeniser.charPos += 2
				} else {
					theTokeniser.charPos++
				}
			case jsCodeState:
				if theTokeniser.scanChars(2) == "/$" {
					theTokeniser.emitToken(jsCodeTokenType)
					theTokeniser.state = textState
					theTokeniser.charPos += 2
				} else {
					theTokeniser.charToToken()
				}
			case jsExprState:
				if theTokeniser.scanChars(2) == ")$" {
					theTokeniser.emitToken(jsExprTokenType)
					theTokeniser.state = textState
					theTokeniser.charPos += 2
				} else {
					theTokeniser.charToToken()
				}
			case convState:
				//	log.Printf("scanning directive: %s...", scanChars(5))
				if theTokeniser.scanChars(2) == "]$" {
					//emitToken(convTokenType)
					theTokeniser.state = textState
					theTokeniser.charPos += 2
				} else {
					theTokeniser.getConvToken()
				}
			case htmlState:
				if theTokeniser.scanChars(2) == ">$" {
					theTokeniser.emitToken(htmlTokenType)
					theTokeniser.state = textState
					theTokeniser.charPos += 2
				} else {
					theTokeniser.charToToken()
				}
			default:
				reportError("internal error, unknown scan state", theTokeniser.line.number)
				log.Fatalf("unknown scan state: %d", theTokeniser.state)
				theTokeniser.charPos++
			}
		}
		logIfLogging("end of input line")
		if theTokeniser.line.eof {
			logIfLogging("in lines at eof, emitting eof token")
			break
		}
	}
	switch theTokeniser.state {
	case textState:
		theTokeniser.emitToken(textTokenType)
		logIfLogging("in text at end")
	case commentState:
		reportError("in comment at end", theTokeniser.line.number)
	case jsCodeState:
		theTokeniser.emitToken(jsCodeTokenType)
		reportError("in code at end", theTokeniser.line.number)
	case jsExprState:
		theTokeniser.emitToken(jsExprTokenType)
		reportError("in expression at end", theTokeniser.line.number)
	case convState:
		reportError("in directive at end", theTokeniser.line.number)
	case htmlState:
		theTokeniser.emitToken(htmlTokenType)
		reportError("in html at end", theTokeniser.line.number)
	}
	logIfLogging("emitting end of file token")
	theTokeniser.emitToken(eofTokenType)
	linesDone <- 1
	logIfLogging("all lines read")
}

/* getConvToken gets a token */
func (theTokeniser *tokeniser) getConvToken() {
	theTokeniser.examineNextChar()
	//log.Printf("looking for conv token, examining from: %c at: %d type: %s", theTokeniser.nextChar, theTokeniser.charPos, theTokeniser.nextCharType)
	switch theTokeniser.nextCharType {
	case identCharType:
		for theTokeniser.nextCharType == identCharType || theTokeniser.nextCharType == digitCharType {
			theTokeniser.charToToken()
		}
		theTokeniser.emitToken(identTokenType)
	case digitCharType:
		for theTokeniser.nextCharType == digitCharType {
			theTokeniser.charToToken()
		}
		theTokeniser.emitToken(numberTokenType)
	case spaceCharType:
		for theTokeniser.nextCharType == spaceCharType {
			theTokeniser.charPos++
			theTokeniser.examineNextChar()
		}
	case specialCharType:
		theTokeniser.charToToken()
		theTokeniser.emitToken(delimTokenType)
	case otherCharType:
		for theTokeniser.nextCharType == otherCharType {
			theTokeniser.charToToken()
		}
		theTokeniser.emitToken(delimTokenType)
	case stopCharType:
		log.Printf("stop char detected %c in %s...", theTokeniser.nextChar, theTokeniser.scanChars(5))
		theTokeniser.charPos++
		theTokeniser.examineNextChar()
	default:
		log.Printf("unknown detected %c", theTokeniser.nextChar)
		theTokeniser.charPos++
		theTokeniser.examineNextChar()
	}
}

type tokeniser struct {
	lineLen      int
	charPos      int
	line         scanLine
	state        scanLineState
	tokenText    []byte
	nextChar     byte
	nextCharType charType
	more         bool
}

/* */ func (theTokeniser *tokeniser) scanChars(leng int) (chars string) {
	useLen := leng
	if leng+theTokeniser.charPos > theTokeniser.lineLen {
		useLen = theTokeniser.lineLen - theTokeniser.charPos
		//log.Printf("only returning %d chars", useLen)
	}
	return theTokeniser.line.text[theTokeniser.charPos : theTokeniser.charPos+useLen]
}

/* */ func (theTokeniser *tokeniser) charToToken() {
	currentChar := theTokeniser.line.text[theTokeniser.charPos]
	theTokeniser.tokenText = append(theTokeniser.tokenText, currentChar)
	//	if theTokeniser.charPos >= theTokeniser.lineLen {
	//		if logging {if logging {log.Print("ran off end of line!")
	//	}
	theTokeniser.charPos++
	theTokeniser.examineNextChar()
}

/* */ func (theTokeniser *tokeniser) examineNextChar() {
	if theTokeniser.charPos >= theTokeniser.lineLen {
		//if logging {if logging {log.Print("at end of line")
		theTokeniser.nextCharType = stopCharType
		theTokeniser.more = false
		return
	}
	theTokeniser.nextChar = theTokeniser.line.text[theTokeniser.charPos]
	theTokeniser.nextCharType = charTypes[theTokeniser.nextChar]
}

/* */ func (theTokeniser *tokeniser) setTokenState(newState scanLineState) {
	theTokeniser.state = newState
	theTokeniser.tokenText = make([]byte, 0)
}

/* */ func (theTokeniser *tokeniser) emitToken(theTokenType tokenType) {
	tokenText := string(theTokeniser.tokenText)
	if tokenText == "" && theTokenType != eofTokenType {
		return
	}
	logfIfLogging("emitting token: %s", tokenText)
	theToken := token{theType: theTokenType, text: tokenText, lineNumber: theTokeniser.line.number}
	theTokeniser.tokenText = make([]byte, 0)
	tokenChan <- theToken
}

/* */ type charType int

const (
	identCharType charType = iota
	digitCharType
	spaceCharType
	specialCharType
	stopCharType
	otherCharType
)

/* */ func (theCharType charType) String() string {
	switch theCharType {
	case identCharType:
		return "ident"
	case digitCharType:
		return "digit"
	case spaceCharType:
		return "spaceC"
	case specialCharType:
		return "special"
	case stopCharType:
		return "stop"
	case otherCharType:
		return "other"
	default:
		return "(??)"
	}
}

/* */ var charTypes []charType

/* */ func init() {
	charTypes = make([]charType, 256)
	for index := range charTypes {
		charTypes[index] = otherCharType
	}
	for index := 'a'; index <= 'z'; index++ {
		charTypes[index] = identCharType
	}
	for index := 'A'; index <= 'Z'; index++ {
		charTypes[index] = identCharType
	}
	for index := '0'; index <= '9'; index++ {
		charTypes[index] = digitCharType
	}
	for _, char := range " \n\t\r" {
		charTypes[char] = spaceCharType
	}
	for _, char := range "()" {
		charTypes[char] = specialCharType
	}
	charTypes['_'] = identCharType
	charTypes[']'] = stopCharType
}

// This file is part of OldRope. OldRope is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. OldRope is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with OldRope. If not, see <http://www.gnu.org/licenses/>.
