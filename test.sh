echo "making test"
export BASE=$(git rev-parse --show-toplevel)
cd $BASE
export GOPATH="/work/go:$BASE" 
go test -c -o runtest conv
echo "running test"
./runtest
echo "test done"