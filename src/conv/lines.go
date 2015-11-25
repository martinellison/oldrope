package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

type scanLine struct {
	text   string
	number int
	eof    bool
}

var lineChan chan scanLine
var linesDone chan int

func getLines(path string) {
	lineNumber := 0
	file, err := os.Open(path)
	if err != nil {
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
		log.Printf("read error: %v", err)
	}
	lineChan <- scanLine{eof: true}
	if err != nil {
		log.Printf("get lines err: %v", err)
	}
	log.Printf("%d lines sent", lineNumber)
}
