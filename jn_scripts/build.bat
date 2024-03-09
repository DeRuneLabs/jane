@REM MIT License

@REM Copyright (c) 2024 arfy slowy - DeRuneLabs

@echo off

if exist .\jane.exe (del /f jane.exe)

if exist command\jane\main.go (
    go build -o jane.exe -v command\jane\main.go
) else (
    go build -o jane.exe -v ..\command\jane\main.go
)

if exist .\jane.exe (
    echo jane created successfully
) else (
    echo something wrong when build jane
)