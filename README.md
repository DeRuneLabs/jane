![jane](.github/jane.png)


![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/DeRuneLabs/jane/workflow_go_linux.yml?style=flat-square&logo=github&label=Build%20Linux)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/DeRuneLabs/jane/workflow_go_macos.yml?style=flat-square&logo=github&label=Build%20MacOS)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/DeRuneLabs/jane/workflow_go_windows.yml?style=flat-square&logo=github&label=Build%20Windows)


Jane is early experimental compiled programming language, static type, fast, modern and simple. the
flow of jane source compiled, its to translate to C++ code and compiled it from
C++ code. Transpile to C++ only instead of compiling is also an Optional. the
mission to be advance, readable and good choice for system programming

```py
main() {
    println("welcome")
}
```

| name      | description                   |
| --------- | ----------------------------- |
| `ast`     | abstact syntax tree builder   |
| `command` | main and compile files        |
| `parser`  | interpreter                   |
| `package` | uitlity package jane          |
| `janelib` | builtin jane standard library |

operator

| operator | description | support type                |
| -------- | ----------- | --------------------------- |
| `+`      | sum         | integer, float, string      |
| `-`      | difference  | integer, float              |
| `*`      | product     | integer, float              |
| `/`      | quotient    | integer, float              |
| `%`      | remainder   | integer                     |
| `~`      | bitwise NOT | integer                     |
| `&`      | bitwise AND | integer                     |
| `^`      | bitwise XOR | integer                     |
| `!`      | logical NOT | bool                        |
| `&&`     | logical AND | bool                        |
| `!=`     | LOGICAL XOR | bool                        |
| `<<`     | left shift  | integer << unsigned integer |
| `>>`     | right shift | integer >> unsigned integer |

precedence

```
| precedence | operator               |
| ---------- | ---------------------- |
| 5          | `*  /  %  <<  >>  &`   |
| 4          | `+ - || ^`             |
| 3          | `==  !=  <  <=  >  >=` |
| 2          | `&&`                   |
| 1          | `||`                   |
```

## Build

**powershell**
```psh
go build -o jane.exe -v command/jn/main.go
```

**bash**
```sh
go build -o jane -v command/jn/main.go
```

