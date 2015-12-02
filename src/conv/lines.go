package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

/* */ type scanLine struct {
	text   string
	number int
	eof    bool
}

/* */ var lineChan chan scanLine

/* */ var linesDone chan int

/* */ func getLines(path string) {
	defer func() {
		rec := recover()
		if rec == nil {
			return
		}
		os.Stderr.WriteString(fmt.Sprintf("getLines: internal error: %v", rec))
	}()
	lineNumber := 0
	file, err := os.Open(path)
	if err != nil {
		reportError(fmt.Sprintf("get lines open err: %v", err), lineNumber)
		log.Printf("get lines open err: %v", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for err == nil {
		var line string
		line, err = reader.ReadString('\n')
		lineNumber++
		lineChan <- scanLine{text: line, number: lineNumber, eof: false}
	}
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		reportError(fmt.Sprintf("read error: %v", err), lineNumber)
		log.Printf("read error: %v", err)
	}
	lineChan <- scanLine{eof: true}
	if err != nil {
		reportError(fmt.Sprintf("get lines err: %v", err), lineNumber)
		log.Printf("get lines err: %v", err)

	}
	if logging {
		log.Printf("%d lines read", lineNumber)
	}
}
