reset
export BASE=$(git rev-parse --show-toplevel)
cd $BASE
export GOPATH=$BASE 
go fmt conv
FMTRES=$?
if [[ $FMTRES != 0 ]]
then
	echo "Format result is" $FMTRES
	exit 1
fi

go build conv	
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

