// tokens.go
package main

import "log"

type token struct {
	theType    tokenType
	text       string
	lineNumber int
}
type tokenType int

const (
	eofTokenType tokenType = iota
	jsCodeTokenType
	jsExprTokenType
	textTokenType
	numberTokenType
	identTokenType
	delimTokenType
)

var tokenChan chan token

type scanLineState int

const (
	textState scanLineState = iota
	commentState
	jsCodeState
	jsExprState
	convState
)

func getTokens() {
	setTokenState(textState)
	for {
		scanState.line = <-lineChan
		if scanState.line.eof {
			tokenChan <- token{theType: eofTokenType}
			linesDone <- 1
			break
		}
		log.Printf("%d: %s (%d)", scanState.line.number, scanState.line.text, len(scanState.line.text))
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
				case "$<":
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
				default:
					charToToken()
				}
			case commentState:
				if scanChars(2) == "*/" {
					log.Print("end of comment")
					setTokenState(textState)
					scanState.charPos += 2
				} else {
					scanState.charPos++
				}
			case jsCodeState:
				if scanChars(2) == ">$" {
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
				if scanChars(2) == "]$" {
					//emitToken(convTokenType)
					scanState.state = textState
					scanState.charPos += 2
				} else {
					getConvToken()
				}
			default:
				log.Fatalf("unknown scan state: %d", scanState.state)
			}
		}
	}
	switch scanState.state {
	case textState:
		emitToken(textTokenType)
	case commentState:
		log.Print("in comment at end!")
	case jsCodeState:
		emitToken(jsCodeTokenType)
	case jsExprState:
		emitToken(jsExprTokenType)
	case convState:
		//emitToken(convTokenType)
	}
	log.Print("all lines read")
}
func getConvToken() {
	examineNextChar()
	log.Printf("looking for conv token, examining from: %b at: %d type: %d", scanState.nextChar, scanState.charPos, scanState.nextCharType)
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
		for scanState.nextCharType == digitCharType {
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
		log.Printf("stop char detected %b", scanState.nextChar)
	default:
		log.Printf("unknown detected %b", scanState.nextChar)
	}
}

var scanState struct {
	lineLen      int
	charPos      int
	line         scanLine
	state        scanLineState
	tokenText    []byte
	nextChar     byte
	nextCharType charType
}

func scanChars(leng int) (chars string) {
	useLen := leng
	if leng+scanState.charPos > scanState.lineLen {
		useLen = scanState.lineLen - scanState.charPos
		log.Printf("only returning %d chars", useLen)
	}
	return scanState.line.text[scanState.charPos : scanState.charPos+useLen]
}
func charToToken() {
	currentChar := scanState.line.text[scanState.charPos]
	scanState.tokenText = append(scanState.tokenText, currentChar)
	scanState.charPos++
	//	if scanState.charPos >= scanState.lineLen {
	//		log.Print("ran off end of line!")
	//	}
	examineNextChar()
}
func examineNextChar() {
	if scanState.charPos >= scanState.lineLen {
		log.Print("atend of line")
		scanState.nextCharType = stopCharType
		return
	}
	scanState.nextChar = scanState.line.text[scanState.charPos]
	scanState.nextCharType = charTypes[scanState.nextChar]
}
func setTokenState(newState scanLineState) {
	scanState.state = newState
	scanState.tokenText = make([]byte, 0)
}
func emitToken(theTokenType tokenType) {
	tokenText := string(scanState.tokenText)
	log.Printf("emitting token: %s", tokenText)
	theToken := token{theType: theTokenType, text: tokenText, lineNumber: scanState.line.number}
	scanState.tokenText = make([]byte, 0)
	tokenChan <- theToken
}

type charType int

const (
	identCharType charType = iota
	digitCharType
	spaceCharType
	specialCharType
	stopCharType
	otherCharType
)

var charTypes []charType

func init() {
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
