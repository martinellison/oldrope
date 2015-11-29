cd ~/git/twine/src
../build.sh 
../conv -dir ../test -jsout testout.js
js-beautify ../test/testout.js > temp
mv temp ../test/testout.js
