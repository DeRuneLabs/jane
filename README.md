![jane](.github/jane.png)

Jane is compiled programming language, static type, fast, modern and simple. the
flow of jane source compiled, its to translate to C++ code and compiled it from
C++ code. Transpole to C++ only instead of compiling is also an Optional. the
mission to be advance, readable and good choic for system programming

| name      | description                 |
| --------- | --------------------------- |
| `ast`     | abstact syntax tree builder |
| `command` | main and compile files      |
| `parser`  | interpreter                 |
| `package` | uitlity package jane        |

example jane code

```
function main() int32 {
    return 0;
}
```

will be transpile to cpp output which is

```cpp
#include <iostream>
#include <locale.h>

template <typename any>

int main() {
    setlocale(0x0, "");
    return 0;
}
```
