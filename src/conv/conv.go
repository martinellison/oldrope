// conv.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	log.Print("start")
	var inFileName, outFileName, baseDir, jsFileName string
	var help bool
	flag.StringVar(&baseDir, "dir", ".", "directory for files")
	flag.StringVar(&inFileName, "in", "test.data", "input file name")
	flag.StringVar(&outFileName, "out", "testout.html", "output file name")
	flag.StringVar(&jsFileName, "jsout", "", "Javascript output file name (if not specified, Javascript will be embedded in the HTML)")
	flag.BoolVar(&help, "help", false, "display help")
	flag.BoolVar(&help, "h", false, "display help")
	flag.Parse()
	if help {
		flag.PrintDefaults()
		return
	}
	filePrefix := fmt.Sprintf("%s%c", baseDir, os.PathSeparator)
	lineChan = make(chan scanLine, 1)
	go getLines(filePrefix + inFileName)
	linesDone = make(chan int)
	go getTokens()
	tokenChan = make(chan token)
	go parse()
	<-linesDone
	log.Print("all done.")
	dumpPages()
	makeTemplate()
	makeGenData()
	file, err := os.Create(filePrefix + outFileName)
	if err != nil {
		log.Fatal(err)
	}
	var jsFile *os.File = file
	jsSeparateFile := jsFileName != ""
	if jsSeparateFile {
		var err error
		jsFile, err = os.Create(filePrefix + jsFileName)
		if err != nil {
			log.Fatal(err)
		}
	}
	genStart(file)
	if jsSeparateFile {
		file.WriteString("<script src='" + jsFileName + "'></script>")
	} else {
		file.WriteString("<script>")
	}
	genJsStart(jsFile)
	expandTemplate(jsFile)
	genJsEnd(jsFile)
	if !jsSeparateFile {
		file.WriteString("</script>")
	}
	genEnd(file)
}
