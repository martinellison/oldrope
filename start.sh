#!/usr/bin/env bash
# Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.

# set GOPATH (only on my environment)

reset
export BASE=$(git rev-parse --show-toplevel)
MACHINE=`uname -n` 
if [[ $MACHINE == 'edward' ]]
then
	export GOX=/home/martin/gocode
else
	export GOX=/work/golang
fi

export GOPATH="$GOX:$BASE"
echo "GOPATH is now" $GOPATH