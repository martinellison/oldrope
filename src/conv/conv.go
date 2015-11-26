// conv.go
package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	log.Print("start")
	var inFileName, outFileName string
	var help bool
	flag.StringVar(&outFileName, "out", "test/testout.html", "output file name")
	flag.StringVar(&inFileName, "in", "test.data", "input file name")
	flag.BoolVar(&help, "help", false, "display help")
	flag.BoolVar(&help, "h", false, "display help")
	flag.Parse()
	if help {
		flag.PrintDefaults()
		return
	}
	lineChan = make(chan scanLine, 1)
	go getLines(inFileName)
	linesDone = make(chan int)
	go getTokens()
	tokenChan = make(chan token)
	go parse()
	<-linesDone
	log.Print("all done.")
	dumpPages()
	makeTemplate()
	makeGenData()
	file, err := os.Create(outFileName)
	if err != nil {
		log.Fatal(err)
	}
	genStart(file)
	expandTemplate(file)
	genEnd(file)
}

//func dumpLines() {
//	for {
//		var line scanLine
//		line = <-lineChan
//		if line.eof {
//			break
//		}
//		log.Printf("%d: %s", line.number, line.text)
//	}
//	log.Print("all lines read")
//	linesDone <- 1
//}
