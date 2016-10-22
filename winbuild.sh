#!/usr/bin/env bash
# Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.reset
export BASE=$(git rev-parse --show-toplevel)
cd $BASE
if [[ ! -d $BASE/bin ]]
then
    mkdir $BASE/bin
fi
GOOS=windows GOARCH=amd64 go build -o bin/conv.exe conv
# This file is part of OldRope. OldRope is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. OldRope is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with OldRope. If not, see <http://www.gnu.org/licenses/>.