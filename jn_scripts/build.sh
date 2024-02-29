#!/usr/bin/sh

if [ -f command/jn/main.go ]; then
  MAIN_FILE="command/jn/main.go"
else
  MAIN_FILE="../command/jn/main.go"
fi

go build -o jane -v $MAIN_FILE

if [ $? -eq 0 ]; then
  echo "Compile is successful!"
else
  echo "erro compiler! check error1"
fi
