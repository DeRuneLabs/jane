#!/usr/bin/sh
#
if [ -f command/jn/main.go ]; then
  MAIN_FILE="command/jn/main.go"
else
  MAIN_FILE="../command/jn/main.go"
fi

go build -o jn.out -v $MAIN_FILE

if [ $? -eq 0 ]; then
  echo "Compile is successful!"
else
  echo "-----------------------------------------------------------------------"
  echo "An unexpected error occurred while compiling X. Check errors above."
fi
