#!/usr/bin/env bash
# Copyright 2015 Martin Ellison. For GPL3 licence notice, see the end of this file.

# set GOPATH (only on my environment)

reset
export BASE=$(git rev-parse --show-toplevel)
MACHINE=`uname -n`
case $MACHINE in
    edward)	export GOX=/home/martin/gocode;;
    raspberrypi) export GOX=/work/golang;;
    *) export GOX=/work/golang;;
esac

export GOPATH="$GOX:$BASE"
echo "GOPATH is now" $GOPATH
