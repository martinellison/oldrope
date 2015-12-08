// Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.
// conv_test.go
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestITE1(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("a", ifThenElse(true, "a", "b"))
	assert.Equal("b", ifThenElse(false, "a", "b"))
}

func TestLOL1(t *testing.T) {
	assert := assert.New(t)
	logging = true
	logIfLogging("this should appear in the test output")
	logging = false
	logIfLogging("this should NOT appear in the test output")
	assert.True(true) //dummy to keep Go happy
}

// This file is part of OldRope. OldRope is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. OldRope is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with OldRope. If not, see <http://www.gnu.org/licenses/>.
