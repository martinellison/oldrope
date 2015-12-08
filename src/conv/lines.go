// Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

/* a scanLine is a line as read in ready for tokenising*/
type scanLine struct {
	text   string
	number int
	eof    bool
}

/* lineChan is the current scanLine*/
var lineChan chan scanLine

/* linesDone is a channel that receives the scanLines*/
var linesDone chan int

/* getLines reads the input from a file and coverts it to scanLines */
func getLines(path string) {
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
	logfIfLogging("%d lines read", lineNumber)
}

// This file is part of OldRope. OldRope is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. OldRope is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with OldRope. If not, see <http://www.gnu.org/licenses/>.
