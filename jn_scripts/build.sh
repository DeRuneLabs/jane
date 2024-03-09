#!/bin/bash

# MIT License

# Copyright (c) 2024 arfy slowy - DeRuneLabs

if [ -f command/jn/main.go ]; then
    JANE_MAIN_FILE="command/jn/main.go"
else
    JANE_MAIN_FILE="../command/jn/main.go"
fi


build_jane() {
    if [ $? = "--clear-project" ]; then
        sudo rm -rf  dist/* jane jn.set
        echo "project are cleaned"
        exit 0
    elif [ "$command" = "--build" ]; then
        if [ -f command/jn/main.go ];then
            go build -o jane -v $JANE_MAIN_FILE
        else
            go build -o jane -v $JANE_MAIN_FILE
        fi
        if [ $? -eq 0 ]; then
            echo "Build succesfull"
            exit 0
        else
            echo "Build failed, check error"
        fi
    else
        echo "invalid argument"
    fi
}

go build -o jane -v $JANE_MAIN_FILE

if [ $? -eq 0 ]; then
    echo "binary create success"
else
    echo "something wrong when build jane"
fi