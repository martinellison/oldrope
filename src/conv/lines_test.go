// lines_test.go
package main

// Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLines1(t *testing.T) {
	assert := assert.New(t)
	oldLogging := logging
	logging = true
	lineChan = make(chan scanLine)
	go getLines("test1.oldrope")
	ln := 0
	var theScanLine scanLine
	for !getLine(&theScanLine, assert) {
		ln++
		assert.Equal(ln, theScanLine.number)
	}
	logging = oldLogging
}
func getLine(theScanLine *scanLine, assert *assert.Assertions) (eof bool) {
	*theScanLine = <-lineChan
	return theScanLine.eof
}

// This file is part of OldRope. OldRope is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. OldRope is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with OldRope. If not, see <http://www.gnu.org/licenses/>.
