# Pretty Log

简单日志，在一定程度上使打印的日志简洁清晰。

## 使用

```shell
go get github.com/my-dev-lib/pretty-log-go
```

```go
import (
    prettylog "github.com/my-dev-lib/pretty-log-go"
)
```

## 输出不同级别日志

```go
// 根据模块确定日志标签
const logTag = "log_test"

// 设置项目根目录名称，以便获得干净的栈信息路径
prettylog.Setup("pretty-log-go")

// 使用全局日志对象打印
prettylog.Iln(logTag, "This is an info level log.")
prettylog.Dln(logTag, "This is a debug level log.")
prettylog.Wln(logTag, "This is a warning level log.")
prettylog.Eln(logTag, "This is a error level log.")
// prettylog.Fatalln("log_test", "This is a fatal level log.")
// prettylog.Panicln(logTag, "This is a panic level log.")

// 使用局部日志对象打印
localLogger := prettylog.NewLogger()
localLogger.SetFlag(prettylog.FlagStackEnabled)
localLogger.Iln(logTag, "This is a custom info level log.")
```

```shell
2024/07/25 20:03:29 1308849 I/log_test: This is an info level log. < log_test.go:25
2024/07/25 20:03:29 1308849 D/log_test: This is a debug level log. < log_test.go:26
2024/07/25 20:03:29 1308849 W/log_test: This is a warning level log. < log_test.go:27
2024/07/25 20:03:29 1308849 E/log_test: This is a error level log. < log_test.go:28
2024/07/25 20:03:29 1308849 I/log_test: This is a custom info level log. < log_test.go:35
```

开启颜色（不能保证所有终端都支持）。

```go
log.SetFlag(prettylog.FlagColorEnabled)
```

![colorful_log.png](./arts/colorful_log.png)

## 输出醒目的信息

```go
fmt.Println(prettylog.GetHighlightLine("欢迎进入 V1.0 系统", 30))

lines := []string{"欢迎进入 V1.0 系统", "运行中…"}
fmt.Println(prettylog.GetHighlightLines(lines, 25))
```

```shell
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
┃ 欢迎进入 V1.0 系统
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

┏━━━━━━━━━━━━━━━━━━━━━━━━
┃ 欢迎进入 V1.0 系统
┃ 运行中…
┗━━━━━━━━━━━━━━━━━━━━━━━━
```

## 输出表格

不支持中文，因为无法确保对齐。

```go
// 直接获得表格
content := [][]interface{}{
    {"Name", "Age", "City", "High"},
    {"Alice", 25, "Beijing", "170cm"},
    {"Bob", 30, "San Francisco", "180cm"},
}

fmt.Println(prettylog.GetHorizontalPrettyTable(content))

// 带名称
fmt.Println(prettylog.GetHorizontalPrettyTableWithName(content, "Members"))
```

```shell
┌──────────────────────────────────┐
│ Name   Age  City           High  │
│ ─────  ───  ─────────────  ───── │
│ Alice  25   Beijing        170cm │
│ Bob    30   San Francisco  180cm │
└──────────────────────────────────┘
┌──────────────────────────────────┐
│ Members                          │
├──────────────────────────────────┤
│ Name   Age  City           High  │
│ ─────  ───  ─────────────  ───── │
│ Alice  25   Beijing        170cm │
│ Bob    30   San Francisco  180cm │
└──────────────────────────────────┘
```

以创建对象的方式输出水平表格。

```go
// 逐行记录表格，统一获得
prettyTable := prettylog.NewPrettyTable()
prettyTable.SetGravity(prettylog.GravityHorizontal)
prettyTable.SetTableName("Members")
prettyTable.SetTitles("Name", "Age", "City", "High")
prettyTable.AddValues("Alice", 25, "Beijing", "170cm")
prettyTable.AddValues("Bob", 30, "San Francisco", "180cm")
fmt.Println(prettyTable.Get())
```

```shell
┌──────────────────────────────────┐
│ Members                          │
├──────────────────────────────────┤
│ Name   Age  City           High  │
│ ─────  ───  ─────────────  ───── │
│ Alice  25   Beijing        170cm │
│ Bob    30   San Francisco  180cm │
└──────────────────────────────────┘
```

如果表格列过多，导致折行 ，可选择输出垂直表格。

```go
// 垂直表格
verticalTable := prettylog.NewPrettyTable()
verticalTable.SetGravity(prettylog.GravityVertical)
verticalTable.SetTableName("Members")
verticalTable.SetTitles("Name", "Age", "City", "High")
verticalTable.AddValues("Alice", 25, "Beijing", "170cm")
verticalTable.AddValues("Bob", 30, "San Francisco", "180cm")
fmt.Println(verticalTable.Get())
```

```shell
┌────────────────────╼
│       Members       
├────────[ 0 ]───────┈
│ Name: Alice
│  Age: 25
│ City: Beijing
│ High: 170cm
├────────[ 1 ]───────┈
│ Name: Bob
│  Age: 30
│ City: San Francisco
│ High: 180cm
└────────────────────╼
```