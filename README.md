# 说明

简单日志，在一定程度上使打印的日志简洁清晰。

## 输出不同级别日志

```go
log.P("log.print(...)\n")
log.Pf("log.printf %s\n", "(...)")
log.Pln("log.println(...)")

log_ := log.NewLog("[Test]")
log_.I("This is an info level log.")
log_.D("This is a debug level log.")
log_.W("This is a warn level log.")
log_.E("This is an error level log.")
```

```shell
2024/01/14 19:19:29 log.print(...)
2024/01/14 19:19:29 log.printf (...)
2024/01/14 19:19:29 log.println(...)
2024/01/14 19:19:29 [Test][INFO] This is an info level log.
2024/01/14 19:19:29 [Test][DEBUG] This is a debug level log.
2024/01/14 19:19:29 [Test][WARN] This is a warn level log.
2024/01/14 19:19:29 [Test][ERROR] This is an error level log.
```

开启颜色（不能保证所有终端都支持）。

```go
log_.SetFlag(log.FlagColorEnabled)
```

![colorful_log.png](./arts/colorful_log.png)

## 输出醒目的信息

```go
fmt.Println(log.GetHighlightLine("欢迎进入 V1.0 系统", 30))

lines := []string{"欢迎进入 V1.0 系统", "运行中…"}
fmt.Println(log.GetHighlightLines(lines, 25))
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
fmt.Println(log.GetPrettyTable(content))

// 带名称
fmt.Println(log.GetPrettyTableWithName(content, "Members"))
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
prettyTable := log.NewPrettyTable()
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
verticalTable := log.NewVerticalPrettyTable()
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