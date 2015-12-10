#!/usr/bin/env bash
# Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.

# build the convertor

reset
export BASE=$(git rev-parse --show-toplevel)

if [[ "$GOPATH" == "" ]]
then
	echo "set GOPATH to include " $BASE
	exit 1
fi
cd $BASE
if [[ ! -d $BASE/bin ]]
then
	mkdir $BASE/bin
fi
gofmt -s -w -e src/conv
FMTRES=$?
if [[ $FMTRES != 0 ]]
then
	echo "Format result is" $FMTRES
	exit 1
fi

go build -o bin/conv conv
BUILDRES=$?
	if [[ $BUILDRES != 0 ]]
	then
		echo  "build result is" $BUILDRES
		exit 1
	fi

go vet conv
	VETRES=$?
	if [[ $VETRES != 0 ]]
	then
		echo "vet for $PACKAGE failed with status $VETRES"
		exit 1
	fi

# This file is part of OldRope. OldRope is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. OldRope is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with OldRope. If not, see <http://www.gnu.org/licenses/>.