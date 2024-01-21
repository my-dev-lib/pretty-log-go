package pretty_log_test

import (
	"fmt"
	prettylog "github.com/my-dev-lib/pretty-log-go"
	"testing"
)

func TestHighlightLine(t *testing.T) {
	fmt.Println(prettylog.GetHighlightLine("欢迎进入 V1.0 系统", 30))

	lines := []string{"欢迎进入 V1.0 系统", "运行中…"}
	fmt.Println(prettylog.GetHighlightLines(lines, 25))
}

func TestLog(t *testing.T) {
	prettylog.P("prettylog.print(...)\n")
	prettylog.Pf("prettylog.printf %s\n", "(...)")
	prettylog.Pln("prettylog.println(...)")

	log := prettylog.NewLog("[Test]")
	log.SetFlag(prettylog.FlagColorEnabled)
	log.I("This is an info level log.")
	log.D("This is a debug level log.")
	log.W("This is a warn level log.")
	log.E("This is an error level log.")
}

func TestPrettyTable(t *testing.T) {
	// 直接获得表格
	content := [][]interface{}{
		{"Name", "Age", "City", "High"},
		{"Alice", 25, "Beijing", "170cm"},
		{"Bob", 30, "San Francisco", "180cm"},
	}
	fmt.Println(prettylog.GetPrettyTable(content))

	// 带名称
	fmt.Println(prettylog.GetPrettyTableWithName(content, "Members"))

	// 逐行记录表格，统一获得
	prettyTable := prettylog.NewPrettyTable()
	prettyTable.SetTableName("Members")
	prettyTable.SetTitles("Name", "Age", "City", "High")
	prettyTable.AddValues("Alice", 25, "Beijing", "170cm")
	prettyTable.AddValues("Bob", 30, "San Francisco", "180cm")
	fmt.Println(prettyTable.Get())

	// 垂直表格
	verticalTable := prettylog.NewVerticalPrettyTable()
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
