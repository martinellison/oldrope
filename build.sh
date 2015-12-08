reset
export BASE=$(git rev-parse --show-toplevel)
MACHINE=`uname -n`
if [[ $MACHINE == 'edward' ]]
then
	export GOX=/home/martin/gocode
else
	export GOX=/work/golang
fi

cd $BASE
export GOPATH=$GOX:$BASE
go fmt conv
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

