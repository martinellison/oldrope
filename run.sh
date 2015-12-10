#!/usr/bin/env bash
# Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.

# run the convertor with some test data. Make the Javascript readable for checking.

reset
export BASE=$(git rev-parse --show-toplevel)
cd $BASE
$BASE/build.sh 
if [[ $? != 0 ]]
then
	exit 1
fi
$BASE/bin/conv -dir $BASE/test -in ../test1.oldrope -jsout testout.js
# for js-beautify see https://github.com/beautify-web/js-beautify
js-beautify $BASE/test/testout.js > temp
mv temp $BASE/test/testout.js
# This file is part of OldRope. OldRope is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version. OldRope is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details. You should have received a copy of the GNU General Public License along with OldRope. If not, see <http://www.gnu.org/licenses/>.