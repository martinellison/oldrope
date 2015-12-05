// tokens.go
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

/* */ var tokenChan chan token

/* */ type scanLineState int

const (
	textState scanLineState = iota
	commentState
	jsCodeState
	jsExprState
	convState
	htmlState
)

/* */ func getTokens() {
	defer func() {
		rec := recover()
		if rec == nil {
			return
		}
		os.Stderr.WriteString(fmt.Sprintf("getTokens: internal error: %v", rec))
	}()
	setTokenState(textState)
	for {
		scanState.line = <-lineChan
		logfIfLogging("%d inline: '%s' (%d) eof: %t", scanState.line.number, scanState.line.text, len(scanState.line.text), scanState.line.eof)
		scanState.lineLen = len(scanState.line.text)
		scanState.charPos = 0
		for scanState.charPos < scanState.lineLen {
			//	log.Printf("%d: scan char is: '%s'", scanState.charPos, scanChars(1))
			switch scanState.state {
			case textState:
				switch scanChars(2) {
				case "/*":
					emitToken(textTokenType)
					scanState.state = commentState
					scanState.charPos += 2
				case "$/":
					emitToken(textTokenType)
					scanState.state = jsCodeState
					scanState.charPos += 2
				case "$(":
					emitToken(textTokenType)
					scanState.state = jsExprState
					scanState.charPos += 2
				case "$[":
					emitToken(textTokenType)
					scanState.state = convState
					scanState.charPos += 2
				case "$<":
					emitToken(textTokenType)
					scanState.state = htmlState
					scanState.charPos += 2
				default:
					charToToken()
				}
			case commentState:
				if scanChars(2) == "*/" {
					//if logging {if logging {log.Print("end of comment")
					setTokenState(textState)
					scanState.charPos += 2
				} else {
					scanState.charPos++
				}
			case jsCodeState:
				if scanChars(2) == "/$" {
					emitToken(jsCodeTokenType)
					scanState.state = textState
					scanState.charPos += 2
				} else {
					charToToken()
				}
			case jsExprState:
				if scanChars(2) == ")$" {
					emitToken(jsExprTokenType)
					scanState.state = textState
					scanState.charPos += 2
				} else {
					charToToken()
				}
			case convState:
				//	log.Printf("scanning directive: %s...", scanChars(5))
				if scanChars(2) == "]$" {
					//emitToken(convTokenType)
					scanState.state = textState
					scanState.charPos += 2
				} else {
					getConvToken()
				}
			case htmlState:
				if scanChars(2) == ">$" {
					emitToken(htmlTokenType)
					scanState.state = textState
					scanState.charPos += 2
				} else {
					charToToken()
				}
			default:
				reportError("internal error, unknown scan state", scanState.line.number)
				log.Fatalf("unknown scan state: %d", scanState.state)
				scanState.charPos++
			}
		}
		logIfLogging("end of input line")
		if scanState.line.eof {
			logIfLogging("in lines at eof, emitting eof token")
			break
		}
	}
	switch scanState.state {
	case textState:
		emitToken(textTokenType)
		logIfLogging("in text at end")
	case commentState:
		reportError("in comment at end", scanState.line.number)
	case jsCodeState:
		emitToken(jsCodeTokenType)
		reportError("in code at end", scanState.line.number)
	case jsExprState:
		emitToken(jsExprTokenType)
		reportError("in expression at end", scanState.line.number)
	case convState:
		reportError("in directive at end", scanState.line.number)
	case htmlState:
		emitToken(htmlTokenType)
		reportError("in html at end", scanState.line.number)
	}
	logIfLogging("emitting end of file token")
	emitToken(eofTokenType)
	linesDone <- 1
	logIfLogging("all lines read")
}

/* */ func getConvToken() {
	examineNextChar()
	//log.Printf("looking for conv token, examining from: %c at: %d type: %s", scanState.nextChar, scanState.charPos, scanState.nextCharType)
	switch scanState.nextCharType {
	case identCharType:
		for scanState.nextCharType == identCharType || scanState.nextCharType == digitCharType {
			charToToken()
		}
		emitToken(identTokenType)
	case digitCharType:
		for scanState.nextCharType == digitCharType {
			charToToken()
		}
		emitToken(numberTokenType)
	case spaceCharType:
		for scanState.nextCharType == spaceCharType {
			scanState.charPos++
			examineNextChar()
		}
	case specialCharType:
		charToToken()
		emitToken(delimTokenType)
	case otherCharType:
		for scanState.nextCharType == otherCharType {
			charToToken()
		}
		emitToken(delimTokenType)
	case stopCharType:
		log.Printf("stop char detected %c in %s...", scanState.nextChar, scanChars(5))
		scanState.charPos++
		examineNextChar()
	default:
		log.Printf("unknown detected %c", scanState.nextChar)
		scanState.charPos++
		examineNextChar()
	}
}

/* */ var scanState struct {
	lineLen      int
	charPos      int
	line         scanLine
	state        scanLineState
	tokenText    []byte
	nextChar     byte
	nextCharType charType
	more         bool
}

/* */ func scanChars(leng int) (chars string) {
	useLen := leng
	if leng+scanState.charPos > scanState.lineLen {
		useLen = scanState.lineLen - scanState.charPos
		//log.Printf("only returning %d chars", useLen)
	}
	return scanState.line.text[scanState.charPos : scanState.charPos+useLen]
}

/* */ func charToToken() {
	currentChar := scanState.line.text[scanState.charPos]
	scanState.tokenText = append(scanState.tokenText, currentChar)
	//	if scanState.charPos >= scanState.lineLen {
	//		if logging {if logging {log.Print("ran off end of line!")
	//	}
	scanState.charPos++
	examineNextChar()
}

/* */ func examineNextChar() {
	if scanState.charPos >= scanState.lineLen {
		//if logging {if logging {log.Print("at end of line")
		scanState.nextCharType = stopCharType
		scanState.more = false
		return
	}
	scanState.nextChar = scanState.line.text[scanState.charPos]
	scanState.nextCharType = charTypes[scanState.nextChar]
}

/* */ func setTokenState(newState scanLineState) {
	scanState.state = newState
	scanState.tokenText = make([]byte, 0)
}

/* */ func emitToken(theTokenType tokenType) {
	tokenText := string(scanState.tokenText)
	if tokenText == "" && theTokenType != eofTokenType {
		return
	}
	logfIfLogging("emitting token: %s", tokenText)
	theToken := token{theType: theTokenType, text: tokenText, lineNumber: scanState.line.number}
	scanState.tokenText = make([]byte, 0)
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
