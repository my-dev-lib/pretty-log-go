package pretty_log

import (
	"fmt"
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

// GetHorizontalPrettyTable 获得横向美观的表格
func GetHorizontalPrettyTable(content [][]any) string {
	return GetHorizontalPrettyTableWithName(content, "")
}

// GetVerticalPrettyTable 获得纵向美观的表格
func GetVerticalPrettyTable(content [][]any) string {
	return GetVerticalPrettyTableWithName(content, "")
}

func GetHorizontalPrettyTableWithName(content [][]any, tableName string) string {
	return getPrettyTableWithOption(content, tableName, GravityHorizontal)
}

func GetVerticalPrettyTableWithName(content [][]any, tableName string) string {
	return getPrettyTableWithOption(content, tableName, GravityVertical)
}

// GetPrettyTableWithOption 获得美观的表格
func getPrettyTableWithOption(content [][]any, tableName string, gravity Gravity) string {
	if len(content) <= 1 {
		return ""
	}

	prettyLog := NewPrettyTable()
	prettyLog.SetGravity(gravity)
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

func NewPrettyTable() PrettyTable {
	return &prettyTableImpl{}
}

type Gravity int

const (
	GravityVertical   Gravity = 0
	GravityHorizontal         = 1
)

type PrettyTable interface {
	SetTableName(tableName string)
	SetTitles(titles ...any)
	AddValues(values ...any)
	SetGravity(gravity Gravity)
	Get() string
}

type prettyTableImpl struct {
	nameWidths []int
	titles     []any
	content    [][]any
	tableName  string
	gravity    Gravity
}

func (pti *prettyTableImpl) SetTableName(tableName string) {
	pti.tableName = tableName
}

func (pti *prettyTableImpl) SetTitles(titles ...any) {
	if pti.isVertical() {
		pti.SetTitlesVertical(titles)
	} else {
		pti.SetTitlesHorizontal(titles)
	}
}

func (pti *prettyTableImpl) SetTitlesVertical(titles []any) {
	clear(pti.titles)
	pti.titles = pti.titles[:0]
	pti.titles = append(pti.titles, titles...)
}

func (pti *prettyTableImpl) SetTitlesHorizontal(titles []any) {
	clear(pti.titles)
	pti.titles = pti.titles[:0]
	pti.titles = append(pti.titles, titles...)
	pti.updateNameWidthsHorizontal(titles)
}

func (pti *prettyTableImpl) AddValues(values ...any) {
	if pti.isVertical() {
		pti.AddValuesVertical(values...)
	} else {
		pti.AddValuesHorizontal(values...)
	}
}

func (pti *prettyTableImpl) AddValuesVertical(values ...any) {
	pti.content = append(pti.content, values)
	pti.updateNameWidthsVertical(values)
}

func (pti *prettyTableImpl) AddValuesHorizontal(values ...any) {
	pti.content = append(pti.content, values)
	pti.updateNameWidthsHorizontal(values)
}

func (pti *prettyTableImpl) SetGravity(gravity Gravity) {
	pti.gravity = gravity
}

func (pti *prettyTableImpl) Get() string {
	if pti.isVertical() {
		return pti.GetVertical()
	} else {
		return pti.GetHorizontal()
	}
}

func (pti *prettyTableImpl) isVertical() bool {
	return pti.gravity == GravityVertical
}

func (pti *prettyTableImpl) updateNameWidthsHorizontal(arr []any) {
	if len(pti.nameWidths) == 0 {
		for i := 0; i < len(arr); i++ {
			pti.nameWidths = append(pti.nameWidths, 0)
		}
	}

	if len(pti.nameWidths) < len(arr) {
		arr = arr[:len(pti.nameWidths)]
	}

	for i, v := range arr {
		vl := len(fmt.Sprint(v))
		if vl > pti.nameWidths[i] {
			pti.nameWidths[i] = vl
		}
	}
}

func (pti *prettyTableImpl) getMaxWidthHorizontal() int {
	maxWidth := 0
	for _, v := range pti.nameWidths {
		maxWidth += v
		maxWidth += 2
	}

	maxWidth += 2
	return maxWidth
}

func (pti *prettyTableImpl) getPrettyTableHorizontal(tableName string, nameWidths []int, titles []any, content [][]any) string {
	if len(content) == 0 {
		return ""
	}

	lastIndex := len(nameWidths) - 1

	maxWidth := pti.getMaxWidthHorizontal()

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

// GetHorizontal 获得横向表格内容
func (pti *prettyTableImpl) GetHorizontal() string {
	return pti.getPrettyTableHorizontal(pti.tableName, pti.nameWidths, pti.titles, pti.content)
}

func (pti *prettyTableImpl) updateNameWidthsVertical(arr []any) {
	if len(pti.nameWidths) == 0 {
		// 固定 2 个
		pti.nameWidths = append(pti.nameWidths, 0, 0)
	}

	if len(pti.titles) == 0 {
		for range arr {
			pti.titles = append(pti.titles, "")
		}
	}

	if len(pti.titles) < len(arr) {
		arr = arr[:len(pti.titles)]
	}

	for i, t := range pti.titles {
		var v string
		if i >= len(arr) {
			v = ""
		} else {
			v = fmt.Sprint(arr[i])
		}

		tl := len(fmt.Sprint(t))
		vl := len(v)

		if tl > pti.nameWidths[0] {
			pti.nameWidths[0] = tl
		}

		if vl > pti.nameWidths[1] {
			pti.nameWidths[1] = vl
		}
	}
}

func (pti *prettyTableImpl) getMaxWidthVertical() int {
	maxWidth := 0
	for _, v := range pti.nameWidths {
		maxWidth += v
		maxWidth += 2
	}
	maxWidth += 2
	return maxWidth
}

func (pti *prettyTableImpl) getPrettyTableVertical(tableName string, nameWidths []int, titles []any, content [][]any) string {
	if len(content) == 0 {
		return ""
	}

	maxWidth := pti.getMaxWidthVertical()

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

			for i := 0; i < pti.nameWidths[0]-len(fmt.Sprint(t)); i++ {
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

func (pti *prettyTableImpl) GetVertical() string {
	return pti.getPrettyTableVertical(pti.tableName, pti.nameWidths, pti.titles, pti.content)
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
