#!/bin/bash

go build -o ./csrf-v3 cmd/main.go

set +m

./csrf-v3 &
fiber=$!
sleep 0.2


function csrf() {
  awk '{
    if (($1 == "localhost") && ($6 == "csrf_")) {
      print $7
    }
  }' cookies.txt
}

function clean() {
  kill -9 $fiber
  rm -f cookies.txt output.txt code.txt
  exit
}

trap clean SIGINT

while true ; do
  curl -v -c cookies.txt http://localhost:3000/login > output.txt 2>&1
  curl -v -s -b cookies.txt -c cookies.txt http://localhost:3000/login >> output.txt 2>&1
  curl -v -b cookies.txt -c cookies.txt -w '%output{code.txt}%{http_code}\n' -d '_csrf='"$(csrf)"'&username=user' http://localhost:3000/login >> output.txt 2>&1
  fail_code="$(cat code.txt)"
  curl -v -b cookies.txt -c cookies.txt -w '%output{code.txt}%{http_code}\n' -d '_csrf='"$(csrf)"'&username=user&password=user' http://localhost:3000/login >> output.txt 2>&1
  ok_code="$(cat code.txt)"

  if [ $fail_code -ne 401 ] || [ $ok_code -ne 200 ] ; then
    cat output.txt
    clean
  fi
  rm -f cookies.txt output.txt code.txt
done

clean
