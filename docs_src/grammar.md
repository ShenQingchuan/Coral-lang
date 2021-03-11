# 基础语法

Coral 程序可以由多个标记组成，可以是关键字，标识符，常量，字符串，符号。

如以下语句由 6 个标记（Token）组成：

```coral
println("My name is " + name)
```

6 个标记是(每行一个)：

```txt
println
(
"My name is "
+
name
)
```

## 文件概念

### 模块化

1. 模块化能够使得 Coral 组织分散的代码模块，
2. 利用 `as` 命别名，有助于解决导入时可能造成的 **"名字空间"** 冲突等问题。

如何建立起各个文件模块之间的关系呢？Coral 并不想让这件事情有太多学习成本，
所以我们学习了 JavaScript 采用 "文件路径字符串" 来定位一个模块。

**这里我们会有一些好用的 IDE 的支持...**

```
import "math";                   // 标准库会被编译器自动识别
import "../src/routes.cr";    // 相对路径
from "math" import pow;          // 提取部分导入
from "httplib" import {
    Request, Response
}
```

### 编译与运行

Coral 的目标是将 `.coral` 源代码编译到 `.cbytes` 字节码文件，然后通过虚拟机执行。

我们倾向于让这门语言看上去像脚本语言一样，所以无需 `main()` 函数。

```bash
# 运行源代码
coral run hello.cr
```

## 终结分隔符

在 Coral 程序中，每个语句正如 C 家族中的其它语言一样以分号 `;` 结尾。以下为两个语句：

```coral
var msg = "Hello World"
println(msg);
```

## 注释

注释不会被编译，甚至在词法分析阶段就被忽略了。

单行注释是最常见的注释形式，你可以在任何地方使用以 `//` 开头的单行注释。多行注释也叫块注释，均已以 `/*` 开头，并以 `*/` 结尾。如：

块注释编译时允许嵌套，但最多 5 层，具体可查看源代码 `lexer` 部分的实现。

```coral
// 单行注释
/*
 Author by 菜鸟教程
 我是多行注释
 */
```

## 标识符

标识符用来命名变量、类型等程序实体。一个标识符实际上就是
一个或是多个 Unicode 字符（**没错，支持中文变量名**）或下划线 `_` 组成的序列，
但是第一个字符必须是 Unicode 字符或下划线，**而不能是数字**。

以下是有效的标识符：

```txt
msesh     圆周率   abc   move_name   a_123
myname50  _temp   j     a23b9       retVal
```

以下是无效的标识符：

- `1ab` 以数字开头
- `case` Coral 语言的关键字
- `a+b` 运算符是不允许的

## 关键字

下面列举了 Coral 代码中可能会使用到的所有关键字和保留字：

```txt
 import    from         as         enum       break  
 continue  return       var        val        if           
 elif      else         switch     default    case       
 while     for          each       in         fn         
 class     interface    this       super      static   
 new       nil          true       false      try       
 catch     finally      throws
```

## 转义字符

Coral 语言支持以下一些特殊的转义字符序列：

|符号|字符含义|
|--- |------|
|`\a`	|响铃符 |
|`\b`	|退格 |
|`\t`	|制表符|
|`\v`	|竖直方向制表符|
|`\n`	|换行 |
|`\r`	|回车 |
|`\f`	|换页符|
|`\"`	|双引号（字符串中）|
|`\'`	|单引号（字符中）|
|`\\`	|反斜杠|
|`\uXXXX`	|4位 16进制 Unicode 字符|
|`\xXX`	|2位 16进制 Unicode 字符|