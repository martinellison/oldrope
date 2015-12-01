// conv.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var logging bool

func main() {
	var inFileName, outFileName, baseDir, jsFileName, logFileName string
	var help bool
	flag.StringVar(&baseDir, "dir", ".", "directory for files")
	flag.StringVar(&inFileName, "in", "test.data", "input file name")
	flag.StringVar(&outFileName, "out", "testout.html", "output file name")
	flag.StringVar(&jsFileName, "jsout", "", "Javascript output file name (if not specified, Javascript will be embedded in the HTML)")
	flag.StringVar(&logFileName, "log", "", "log file name (for debugging)")
	flag.BoolVar(&help, "help", false, "display help")
	flag.BoolVar(&help, "h", false, "display help")
	flag.Parse()
	if help {
		flag.PrintDefaults()
		return
	}
	filePrefix := fmt.Sprintf("%s%c", baseDir, os.PathSeparator)
	initLog(filePrefix, logFileName)
	lineChan = make(chan scanLine, 1)
	go getLines(filePrefix + inFileName)
	linesDone = make(chan int)
	go getTokens()
	tokenChan = make(chan token)
	go parse()
	<-linesDone
	if logging {
		log.Print("all lines scanned and parsed.")
	}
	dumpPages()
	makeTemplate()
	makeGenData()
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
	genStart(file)
	if jsSeparateFile {
		file.WriteString("<script src='" + jsFileName + "'></script>")
	} else {
		file.WriteString("<script>")
	}
	genHeader(jsFile)
	genJsStart(jsFile)
	expandTemplate(jsFile)
	genJsEnd(jsFile)
	if !jsSeparateFile {
		file.WriteString("</script>")
	}
	genEnd(file)
	if logging {
		log.Print("file generated.")
	}
}
func ifThenElse(p bool, st string, sf string) string {
	if p {
		return st
	}
	return sf
}

// initlog creates a log file for debugging.
func initLog(filePrefix, logFileName string) {
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
func reportError(msg string, lineNumber int) {
	os.Stderr.WriteString(fmt.Sprintf("(%d): %s", lineNumber, msg))
}
