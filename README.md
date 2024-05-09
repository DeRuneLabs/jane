![jane](.github/jane.png)

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/DeRuneLabs/jane/workflow_go_linux.yml?style=flat-square&logo=github&label=Build%20Linux)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/DeRuneLabs/jane/workflow_go_macos.yml?style=flat-square&logo=github&label=Build%20MacOS)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/DeRuneLabs/jane/workflow_go_windows.yml?style=flat-square&logo=github&label=Build%20Windows)

Jane is early experimental compiled programming language, static type, fast, modern and simple. jane design for maintainable and reliable software development. ensure memory safety and commits not to contain undefined behavior, contains a reference compiler withfeatures that help developers to design safe applications.

## CPP interopability

jane is meant to work with cpp, a cpp header file depedence can be addedto the jane code, allowing its functions to be linked. writting cpp code that is compatible with the jane code created by compiler is rather straightforward. jane stores all of cpp code it uses for jane in the `api` directory. this API make it feasible and easy to develop cpp programming that can be completely integrated with jane

![summary_hpp_image](.github/codesnap/summary_hpp.png)
![summary_jn_image](.github/codesnap/summary_jn.png)

## feature of jn

- simple
- fast and scaleable development
- performance-critical software
- memory safety
- fun

## information

the project structure, including its lexical and syntactic structure, has now revealed. however, if there reference compiler is rewritten in jane, it is expected that AST, Lexer, and certain packages will be included in the standard library. this a modification that need official compiler project structure to be rebuilt. reference compiler is likely to make extensive use of standard library. this will also allow dev create language specific utilities using jn std library.

## Build

**powershell**

```psh
go build -o jane.exe -v command/jn/main.go
```

**bash**

```sh
go build -o jane -v command/jn/main.go
```

using makefile

```sh
# to build all
make all
# to clean binary files after make
make clean
```
