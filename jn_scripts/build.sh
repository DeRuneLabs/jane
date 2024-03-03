#!/bin/bash

command=$1
FILE="command/jn/main.go"
FILE2="../command/jn/main.go"

build_jane() {
    if [ "$command" = "--clear-project" ]; then
        sudo rm -rf  dist/* jane jn.set
        echo "project are cleaned"
        exit 0
    elif [ "$command" = "--build" ]; then
        if [ -f command/jn/main.go ];then
            go build -o jane -v $FILE
        else
            go build -o jane -v $FILE2
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

install_debian_family() {
    if [ "$( dpkg -l | awk '/golang-go/ {print }'|wc -l)" -ge 1 | "$(grep -ic go /usr/local/)" -eq 0 ]; then
        build_jane
    else
        wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz && export PATH=$PATH:/usr/local/go/bin
    fi
}

install_arch_family() {
    if pacman -Qs go > /dev/null; then
        build_jane
    else
        sudo pacman -S go
    fi
}

check_os() {
     . /etc/os-release
    case $ID in
        ubuntu)
        install_debian_family
        ;;
        arch)
        install_arch_family
        ;;
        darwin)
        install_debian_family
        ;;
    esac
}

check_os