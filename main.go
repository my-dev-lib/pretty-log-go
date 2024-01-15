package main

import (
	"fmt"
	"main/log"
)

func testHighlightLine() {
	fmt.Println(log.GetHighlightLine("欢迎进入 V1.0 系统", 30))

	lines := []string{"欢迎进入 V1.0 系统", "运行中…"}
	fmt.Println(log.GetHighlightLines(lines, 25))
}

func testLog() {
	log.P("log.print(...)\n")
	log.Pf("log.printf %s\n", "(...)")
	log.Pln("log.println(...)")

	log_ := log.NewLog("[Test]")
	log_.SetFlag(log.FlagColorEnabled)
	log_.I("This is an info level log.")
	log_.D("This is a debug level log.")
	log_.W("This is a warn level log.")
	log_.E("This is an error level log.")
}

func testPrettyTable() {
	// 直接获得表格
	content := [][]interface{}{
		{"Name", "Age", "City", "High"},
		{"Alice", 25, "Beijing", "170cm"},
		{"Bob", 30, "San Francisco", "180cm"},
	}
	fmt.Println(log.GetPrettyTable(content))

	// 带名称
	fmt.Println(log.GetPrettyTableWithName(content, "Members"))

	// 逐行记录表格，统一获得
	prettyTable := log.NewPrettyTable()
	prettyTable.SetTableName("Members")
	prettyTable.SetTitles("Name", "Age", "City", "High")
	prettyTable.AddValues("Alice", 25, "Beijing", "170cm")
	prettyTable.AddValues("Bob", 30, "San Francisco", "180cm")
	fmt.Println(prettyTable.Get())

	// 垂直表格
	verticalTable := log.NewVerticalPrettyTable()
	verticalTable.SetTableName("Members")
	verticalTable.SetTitles("Name", "Age", "City", "High")
	verticalTable.AddValues("Alice", 25, "Beijing", "170cm")
	verticalTable.AddValues("Bob", 30, "San Francisco", "180cm")
	fmt.Println(verticalTable.Get())
}

func main() {
	testHighlightLine()
	testLog()
	testPrettyTable()
}
