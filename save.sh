TEXT=$1
git commit -m "$TEXT" -a 
find | grep "~$" | xargs rm
