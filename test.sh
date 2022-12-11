#!/bin/sh


echo === FIB2 TEST ===
test=$(cat fibout)
res=$(go run . fib2.axe)
if [ "$test" == "$res" ]; then
    echo PASS
else
    echo FAIL
fi
