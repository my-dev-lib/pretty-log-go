package pretty_log

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
