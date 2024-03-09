@echo off

if exist jane.exe and exist jn.set (
  del jane.exe; del jn.set; del dist
)

if exist command/jn/main.go (
  go build -o jane.exe -v command/jn/main.go
) else (
  go build -o jane.exe -v ../command/jn/main.go
)