#!/bin/sh

set -e

runtest() {
    if [ "$test" = "$res" ]; then
        echo PASS
    else
        echo FAIL
    fi
}

cd test

echo === FIB2 TEST ===
test=$(cat "fib.out")
res=$(go run .. "fib2.axe")
runtest

echo === READ TEST ===
test=$(cat "fib.axe")
res=$(go run .. "read.axe")
runtest
