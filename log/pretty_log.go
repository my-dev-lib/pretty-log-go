package log

import (
	"fmt"
	"log"
	"strings"
)

func getHighlightLine(content string, width int) string {
	// 1. top line
	line := "┏"
	for i := 0; i < width-1; i++ {
		line += "━"
	}

	// 2. text
	line += "\n"
	line += content
	line += "\n"

	// 3. bottom line
	line += "┗"
	for i := 0; i < width-1; i++ {
		line += "━"
	}

	return line
}

// GetHighlightLine 获取高亮突出显示的一行
func GetHighlightLine(text string, width int) string {
	content := "┃ " + text
	return getHighlightLine(content, width)
}

// GetHighlightLines 获取高亮突出显示的若干行
func GetHighlightLines(texts []string, width int) string {
	content := ""
	for _, text := range texts {
		content += "┃ " + text + "\n"
	}

	content = strings.TrimSuffix(content, "\n")

	return getHighlightLine(content, width)
}

var infoColor = []int{3, 169, 244}
var debugColor = []int{139, 195, 74}
var warnColor = []int{255, 152, 0}
var errorColor = []int{244, 67, 54}

const (
	Version          = "1.0.0"
	FlagColorEnabled = 0x01

	LevelInfo  = 1
	LevelDebug = 2
	LevelWarn  = 3
	LevelError = 4
)

// Log 调试日志
type Log struct {
	tag  string
	flag int
}

// P log.Print 原始包装
func P(v ...any) {
	log.Print(v...)
}

// Pln log.Println 原始包装
func Pln(v ...any) {
	log.Println(v...)
}

// Pf log.Printf 原始包装
func Pf(format string, a ...any) {
	log.Printf(format, a...)
}

func (l *Log) SetFlag(flag int) {
	l.flag = flag
}

// I 以 Info 级别输出
func (l *Log) I(format string, a ...any) {
	l.println(LevelInfo, format, a...)
}

// D 以 Debug 级别输出
func (l *Log) D(format string, a ...any) {
	l.println(LevelDebug, format, a...)
}

// W 以 Warn 级别输出
func (l *Log) W(format string, a ...any) {
	l.println(LevelWarn, format, a...)
}

// E 以 Error 级别输出
func (l *Log) E(format string, a ...any) {
	l.println(LevelError, format, a...)
}

func getColorfulText(color []int, text string) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", color[0], color[1], color[2], text)
}

func (l *Log) fixLog(color []int, levelTag string, format string, a []any) (string, string) {
	text := fmt.Sprintf(format, a...)
	tag := levelTag

	if l.flag&FlagColorEnabled == FlagColorEnabled {
		tag = getColorfulText(color, levelTag)
		text = getColorfulText(color, text)
	}

	return tag, text
}

func (l *Log) println(level int, format string, a ...any) {
	tag := ""
	color := []int{255, 255, 255}

	switch level {
	case LevelInfo:
		tag = "INFO"
		color = infoColor
	case LevelDebug:
		tag = "DEBUG"
		color = debugColor
	case LevelWarn:
		tag = "WARN"
		color = warnColor
	case LevelError:
		tag = "ERROR"
		color = errorColor
	default:
		return
	}

	tag, text := l.fixLog(color, tag, format, a)

	// PT[DEBUG] This is a log.
	log.Printf("%s[%s] %s\n", l.tag, tag, text)
}

// NewLog 需要提供一个 tag
func NewLog(tag string) *Log {
	return &Log{tag: tag}
}

// GetPrettyTable 获得美观的表格
func GetPrettyTable(content [][]any) string {
	return GetPrettyTableWithName(content, "")
}

// GetPrettyTableWithName 获得美观的表格
func GetPrettyTableWithName(content [][]any, tableName string) string {
	if len(content) <= 1 {
		return ""
	}

	prettyLog := NewPrettyTable()
	if len(tableName) != 0 {
		prettyLog.SetTableName(tableName)
	}

	prettyLog.SetTitles(content[0]...)
	for i, line := range content {
		if i == 0 {
			continue
		}

		prettyLog.AddValues(line...)
	}

	return prettyLog.Get()
}

type PrettyTable interface {
	SetTableName(tableName string)
	SetTitles(titles ...any)
	AddValues(values ...any)
	Get() string
}

// HorizontalPrettyTable 创建水平美观表格
type HorizontalPrettyTable struct {
	nameWidths []int
	titles     []any
	content    [][]any
	tableName  string
}

// SetTableName 设置表格名称
func (hpt *HorizontalPrettyTable) SetTableName(tableName string) {
	hpt.tableName = tableName
}

func (hpt *HorizontalPrettyTable) updateNameWidths(arr []any) {
	if len(hpt.nameWidths) == 0 {
		for i := 0; i < len(arr); i++ {
			hpt.nameWidths = append(hpt.nameWidths, 0)
		}
	}

	if len(hpt.nameWidths) < len(arr) {
		arr = arr[:len(hpt.nameWidths)]
	}

	for i, v := range arr {
		vl := len(fmt.Sprint(v))
		if vl > hpt.nameWidths[i] {
			hpt.nameWidths[i] = vl
		}
	}
}

// SetTitles 设置表格标题
func (hpt *HorizontalPrettyTable) SetTitles(titles ...any) {
	clear(hpt.titles)
	hpt.titles = append(hpt.titles, titles...)
	hpt.updateNameWidths(titles)
}

// AddValues 增加表格数值
func (hpt *HorizontalPrettyTable) AddValues(values ...any) {
	hpt.content = append(hpt.content, values)
	hpt.updateNameWidths(values)
}

func (hpt *HorizontalPrettyTable) getMaxWidth() int {
	maxWidth := 0
	for _, v := range hpt.nameWidths {
		maxWidth += v
		maxWidth += 2
	}

	maxWidth += 2
	return maxWidth
}

func (hpt *HorizontalPrettyTable) getPrettyTable(tableName string, nameWidths []int, titles []any, content [][]any) string {
	if len(content) == 0 {
		return ""
	}

	lastIndex := len(nameWidths) - 1

	maxWidth := hpt.getMaxWidth()

	var pretty strings.Builder
	// 1. print top line
	topLine := getLines(nameWidths, lastIndex, "┌─", "─", "──", "─┐")
	pretty.WriteString(topLine)
	pretty.WriteString("\n")

	nameLength := len(tableName)
	if nameLength != 0 {
		tableName = fixLongString(tableName, maxWidth-6)
		nameLength = len(tableName)

		// print table name line
		pretty.WriteString("│ ")
		totalSize := 0
		for _, nw := range nameWidths {
			totalSize += nw
			totalSize += 2
		}

		pretty.WriteString(tableName)
		for i := 0; i < totalSize-nameLength-1; i++ {
			pretty.WriteString(" ")
		}

		pretty.WriteString("│\n")

		// print split line
		splitLine := getLines(nameWidths, lastIndex, "├─", "─", "──", "─┤")
		pretty.WriteString(splitLine)
		pretty.WriteString("\n")
	}

	// 2. print titles
	if len(titles) != 0 {
		pretty.WriteString("│ ")
		for ii, name := range titles {
			nameStr := fmt.Sprint(name)
			nameWidth := nameWidths[ii]
			pretty.WriteString(nameStr)
			for j := 0; j < nameWidth-len(nameStr); j++ {
				pretty.WriteString(" ")
			}

			if ii != lastIndex {
				pretty.WriteString("  ")
			}
		}
		pretty.WriteString(" │\n")

		// 3. print middle line
		middleLine := getLines(nameWidths, lastIndex, "│ ", "─", "  ", " │")
		pretty.WriteString(middleLine)
		pretty.WriteString("\n")
	}

	// 4. print content
	for _, rowData := range content {
		pretty.WriteString("│ ")
		columnLen := len(nameWidths)
		for j := 0; j < columnLen; j++ {
			nameWidth := nameWidths[j]
			value_ := ""
			if j < len(rowData) {
				value_ = fmt.Sprint(rowData[j])
			}

			pretty.WriteString(value_)
			for k := 0; k < nameWidth-len(value_); k++ {
				pretty.WriteString(" ")
			}

			if j != lastIndex {
				pretty.WriteString("  ")
			}
		}

		pretty.WriteString(" │\n")
	}

	// 5. print bottom line
	bottomLine := getLines(nameWidths, lastIndex, "└─", "─", "──", "─┘")
	pretty.WriteString(bottomLine)
	return pretty.String()
}

// Get 获得表格内容
func (hpt *HorizontalPrettyTable) Get() string {
	return hpt.getPrettyTable(hpt.tableName, hpt.nameWidths, hpt.titles, hpt.content)
}

// NewPrettyTable 创建表格对象
func NewPrettyTable() PrettyTable {
	return &HorizontalPrettyTable{}
}

// VerticalPrettyTable 创建垂直美观表格
type VerticalPrettyTable struct {
	nameWidths []int
	titles     []any
	content    [][]any
	tableName  string
}

func (vpt *VerticalPrettyTable) updateNameWidths(arr []any) {
	if len(vpt.nameWidths) == 0 {
		// 固定 2 个
		vpt.nameWidths = append(vpt.nameWidths, 0, 0)
	}

	if len(vpt.titles) == 0 {
		for range arr {
			vpt.titles = append(vpt.titles, "")
		}
	}

	if len(vpt.titles) < len(arr) {
		arr = arr[:len(vpt.titles)]
	}

	for i, t := range vpt.titles {
		var v string
		if i >= len(arr) {
			v = ""
		} else {
			v = fmt.Sprint(arr[i])
		}

		tl := len(fmt.Sprint(t))
		vl := len(v)

		if tl > vpt.nameWidths[0] {
			vpt.nameWidths[0] = tl
		}

		if vl > vpt.nameWidths[1] {
			vpt.nameWidths[1] = vl
		}
	}
}

func (vpt *VerticalPrettyTable) SetTableName(tableName string) {
	vpt.tableName = tableName
}

func (vpt *VerticalPrettyTable) SetTitles(titles ...any) {
	clear(vpt.titles)
	vpt.titles = append(vpt.titles, titles...)
}

func (vpt *VerticalPrettyTable) AddValues(values ...any) {
	vpt.content = append(vpt.content, values)
	vpt.updateNameWidths(values)
}

func (vpt *VerticalPrettyTable) getMaxWidth() int {
	maxWidth := 0
	for _, v := range vpt.nameWidths {
		maxWidth += v
		maxWidth += 2
	}
	maxWidth += 2
	return maxWidth
}

func (vpt *VerticalPrettyTable) getPrettyTable(tableName string, nameWidths []int, titles []any, content [][]any) string {
	if len(content) == 0 {
		return ""
	}

	maxWidth := vpt.getMaxWidth()

	lastIndex := len(nameWidths) - 1

	var pretty strings.Builder
	// 1. print top line
	topLine := getLines(nameWidths, lastIndex, "┌─", "─", "──", "╼")
	pretty.WriteString(topLine)
	pretty.WriteString("\n")

	tableNameLength := len(tableName)
	if tableNameLength != 0 {
		tableName = fixLongString(tableName, maxWidth-6)
		tableNameLength := len(tableName)

		// 2. print table name
		tableNameStart := (maxWidth-tableNameLength)/2 - 1
		tableNameLine := getLinesReplace(maxWidth, "│", " ", tableNameStart, tableName)
		pretty.WriteString(tableNameLine)
		pretty.WriteString("\n")
	}

	numStart := (maxWidth-len("[ 0 ]"))/2 - 1
	for i, v := range content {
		// print number line
		num := fmt.Sprintf("[ %d ]", i)
		numLine := getLinesReplace(maxWidth, "├", "─", numStart, num)
		numLine = numLine[:len(numLine)-3] + "┈"
		pretty.WriteString(numLine)
		pretty.WriteString("\n")

		for j, t := range titles {
			pretty.WriteString("│ ")

			for i := 0; i < vpt.nameWidths[0]-len(fmt.Sprint(t)); i++ {
				pretty.WriteString(" ")
			}

			pretty.WriteString(fmt.Sprintf("%v: %v", t, v[j]))
			pretty.WriteString("\n")
		}
	}

	bottomLine := getLines(nameWidths, lastIndex, "└─", "─", "──", "╼")
	pretty.WriteString(bottomLine)
	pretty.WriteString("\n")
	return pretty.String()
}

func (vpt *VerticalPrettyTable) Get() string {
	return vpt.getPrettyTable(vpt.tableName, vpt.nameWidths, vpt.titles, vpt.content)
}

// NewVerticalPrettyTable 创建纵向表格
func NewVerticalPrettyTable() PrettyTable {
	return &VerticalPrettyTable{}
}

// utils:

// Hello World -> Hello...
func fixLongString(str string, maxLength int) string {
	nameLength := len(str)
	if nameLength > maxLength {
		str = str[:maxLength-1] + "..."
	}

	return str
}

func getLines(nameWidths []int, lastIndex int, left string, m1 string, m2 string, right string) string {
	var pretty strings.Builder
	pretty.WriteString(left)
	for ii, width := range nameWidths {
		for i := 0; i < width; i++ {
			pretty.WriteString(m1)
		}

		if ii != lastIndex {
			pretty.WriteString(m2)
		}
	}

	pretty.WriteString(right)
	return pretty.String()
}

func getLinesReplace(maxWidth int, left string, m string, start int, replace string) string {
	var pretty strings.Builder
	pretty.WriteString(left)
	for i := 0; i < maxWidth-1; i++ {
		if i == start {
			pretty.WriteString(replace)
			i += len(replace)
			continue
		}

		pretty.WriteString(m)
	}

	return pretty.String()
}
