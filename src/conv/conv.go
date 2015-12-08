// Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.

// conv.go (main program)
/* OldRope is a convertor for  text games. See README.mdown for a description.

 */
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

/* */ var logging, hashText bool

/* */ func main() {
	defer func() {
		rec := recover()
		if rec == nil {
			return
		}
		os.Stderr.WriteString(fmt.Sprintf("internal error: %v", rec))
	}()
	var inFileName, outFileName, baseDir, jsFileName, logFileName string
	var help bool
	flag.StringVar(&baseDir, "dir", ".", "directory for files")
	flag.StringVar(&inFileName, "in", "test.data", "input file name")
	flag.StringVar(&outFileName, "out", "testout.html", "output file name")
	flag.StringVar(&jsFileName, "jsout", "", "Javascript output file name (if not specified, Javascript will be embedded in the HTML)")
	flag.StringVar(&logFileName, "log", "", "log file name (for debugging)")
	flag.BoolVar(&hashText, "hash", false, "use hash escapes for text")
	flag.BoolVar(&help, "help", false, "display help")
	flag.BoolVar(&help, "h", false, "display help")
	flag.Parse()
	if help {
		flag.PrintDefaults()
		return
	}
	filePrefix := fmt.Sprintf("%s%c", baseDir, os.PathSeparator)
	initLog(filePrefix, logFileName)
	lineChan = make(chan scanLine)
	go getLines(filePrefix + inFileName)
	theParseChan := make(chan *pageSet)
	var theTokeniser tokeniser
	theTokeniser.init()
	go theTokeniser.getTokens()
	tokenChan = make(chan token)
	var theParser parser
	go theParser.parse(theParseChan)
	thePageSet := <-theParseChan
	if logging {
		log.Print("all lines scanned and parsed.")
	}
	dumpPages(thePageSet)
	var theGenerator generator
	//	theGenerator.init()
	theGenerator.makeTemplate()
	var theOutData outData
	theOutData.makeGenData(thePageSet)
	file, err := os.Create(filePrefix + outFileName)
	if err != nil {
		reportError(("cannot create file (" + filePrefix + outFileName + "): " + err.Error()), 0)
		log.Fatal(err)
	}
	var jsFile *os.File = file
	jsSeparateFile := jsFileName != ""
	if jsSeparateFile {
		var err error
		jsFile, err = os.Create(filePrefix + jsFileName)
		if err != nil {
			reportError(("cannot create file (" + filePrefix + jsFileName + "): " + err.Error()), 0)
			log.Fatal(err)
		}
	}
	theGenerator.genStart(file)
	if jsSeparateFile {
		file.WriteString("<script src='" + jsFileName + "'></script>")
	} else {
		file.WriteString("<script>")
	}
	theGenerator.genHeader(jsFile)
	theGenerator.genJsStart(jsFile, thePageSet.startPageName)
	theGenerator.expandTemplate(jsFile, theOutData)
	theGenerator.genJsEnd(jsFile)
	if !jsSeparateFile {
		file.WriteString("</script>")
	}
	theGenerator.genEnd(file)
	logIfLogging("file generated.")
}

/* */ func ifThenElse(p bool, st string, sf string) string {
	if p {
		return st
	}
	return sf
}

// initlog creates a log file for debugging.
/* */ func initLog(filePrefix, logFileName string) {
	if logFileName == "" {
		logging = false
		return
	}
	f, theError := os.OpenFile(filePrefix+logFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if theError != nil {
		reportError(("cannot create file (" + filePrefix + logFileName + "): " + theError.Error()), 0)
		log.Fatalf("error opening file: %v", theError)
	}
	log.SetOutput(f)
	logging = true
}
func logIfLogging(msg string) {
	if logging {
		log.Print(msg)
	}
}
func logfIfLogging(msg string, params ...interface{}) {
	if logging {
		log.Printf(msg, params...)
	}
}

/* */ func reportError(msg string, lineNumber int) {
	os.Stderr.WriteString(fmt.Sprintf("(%d): %s\n", lineNumber, msg))
}

// This file is part of OldRope. OldRope is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. OldRope is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with OldRope. If not, see <http://www.gnu.org/licenses/>.
