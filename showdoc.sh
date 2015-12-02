export BASE=$(git rev-parse --show-toplevel)
cd $BASE/src/conv
go doc -cmd -u
