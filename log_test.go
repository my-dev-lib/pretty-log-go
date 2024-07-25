package pretty_log_test

import (
	"fmt"
	prettylog "github.com/my-dev-lib/pretty-log-go"
	"testing"
)

const logTag = "log_test"

func TestHighlightLine(t *testing.T) {
	fmt.Println("--== TestHighlightLine ==--")
	fmt.Println(prettylog.GetHighlightLine("欢迎进入 V1.0 系统", 30))

	lines := []string{"欢迎进入 V1.0 系统", "运行中…"}
	fmt.Println(prettylog.GetHighlightLines(lines, 25))
}

func TestLog(t *testing.T) {
	fmt.Println("--== TestLog ==--")
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
}

func TestPrettyTable(t *testing.T) {
	fmt.Println("--== TestPrettyTable ==--")
	// 直接获得表格
	content := [][]interface{}{
		{"Name", "Age", "City", "High"},
		{"Alice", 25, "Beijing", "170cm"},
		{"Bob", 30, "San Francisco", "180cm"},
	}

	fmt.Println(prettylog.GetHorizontalPrettyTable(content))

	// 带名称
	fmt.Println(prettylog.GetHorizontalPrettyTableWithName(content, "Members"))

	// 逐行记录表格，统一获得
	prettyTable := prettylog.NewPrettyTable()
	prettyTable.SetGravity(prettylog.GravityHorizontal)
	prettyTable.SetTableName("Members")
	prettyTable.SetTitles("Name", "Age", "City", "High")
	prettyTable.AddValues("Alice", 25, "Beijing", "170cm")
	prettyTable.AddValues("Bob", 30, "San Francisco", "180cm")
	fmt.Println(prettyTable.Get())

	// 垂直表格
	verticalTable := prettylog.NewPrettyTable()
	verticalTable.SetGravity(prettylog.GravityVertical)
	verticalTable.SetTableName("Members")
	verticalTable.SetTitles("Name", "Age", "City", "High")
	verticalTable.AddValues("Alice", 25, "Beijing", "170cm")
	verticalTable.AddValues("Bob", 30, "San Francisco", "180cm")
	fmt.Println(verticalTable.Get())
}

func TestAll(t *testing.T) {
	TestHighlightLine(t)
	TestLog(t)
	TestPrettyTable(t)
}
