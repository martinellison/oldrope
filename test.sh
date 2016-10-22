#!/usr/bin/env bash
# Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.

# run automated tests

reset
echo "making test"
export BASE=$(git rev-parse --show-toplevel)

cd $BASE
$BASE/build.sh
BUILDRES=$?
if [[ $BUILDRES != 0 ]]
then
    echo "Build main result is" $BUILDRES
    exit 1
fi
go test -c -o $BASE/bin/runtest -cover conv
COMPILERESULT=$?
if [[ $COMPILERESULT != 0 ]]
then
    echo "build failed"
    exit 1
fi
echo "running test"
bin/runtest -test.coverprofile $BASE/test/profile.out
if [[ $? != 0 ]]
then
    exit 1
fi
go tool cover -func=$BASE/test/profile.out|grep -v '100.0%'
go tool cover -html=$BASE/test/profile.out -o $BASE/test/profile.html
echo "test done"
# This file is part of OldRope. OldRope is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. OldRope is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with OldRope. If not, see <http://www.gnu.org/licenses/>.