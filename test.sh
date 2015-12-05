reset
echo "making test"
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
./build.sh
BUILDRES=$?
if [[ $BUILDRES != 0 ]]
then
	echo "Build main result is" $BUILDRES
	exit 1
fi
go test -c -o runtest -cover conv
COMPILERESULT=$?
if [[ $COMPILERESULT != 0 ]]
then
	echo "build failed"
	exit 1
fi
echo "running test"
./runtest -test.coverprofile test/profile.out
if [[ $? != 0 ]]
then
	exit 1
fi
GOPATH=~/git/twine go tool cover -func=test/profile.out|grep -v '100.0%'
GOPATH=~/git/twine go tool cover -html=test/profile.out -o test/profile.html
echo "test done"
