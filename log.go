package pretty_log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type (
	Flag  int // 日志标记
	Level int // 日志级别
)

const (
	FlagClear        Flag = 0x00
	FlagColorEnabled Flag = 0x01 // 启用颜色
	FlagStackEnabled Flag = 0x02 // 启用栈信息

	baseStackCount  = 4 // 栈数量
	baseStackOffset = 3 // 栈偏移

	LevelInfo  Level = 1
	LevelDebug Level = 2
	LevelWarn  Level = 3
	LevelError Level = 4
	LevelFatal Level = 5
	LevelPanic Level = 6

	// D/Hello: This is a debug level log.
	logTextFormat = "%s/%s: %s"
	// 2024/07/22 19:11:22 11243 D/Hello: This is a debug level log. > hello.go:33
	logFormat = "%d %s%s"
)

var codeRootPath = "."

func Setup(projectDir string) Logger {
	codeRootPath = detectProjectRoot(projectDir)
	return globalLogger
}

var levelInfoMap = map[Level]*levelInfo{
	LevelInfo:  {color: []int{3, 169, 244}, label: "I", printer: log.Print}, // 蓝色
	LevelDebug: {color: []int{76, 175, 80}, label: "D", printer: log.Print}, // 绿色
	LevelWarn:  {color: []int{255, 152, 0}, label: "W", printer: log.Print}, // 橙色
	LevelError: {color: []int{244, 67, 54}, label: "E", printer: log.Print}, // 红色
	LevelFatal: {color: []int{121, 85, 72}, label: "F", printer: log.Fatal}, // 褐色
	LevelPanic: {color: []int{121, 85, 72}, label: "P", printer: log.Panic}, // 褐色
}

type Logger interface {
	SetStackCount(int)
	SetStackOffset(int)
	SetFlag(Flag)
	AddFlag(Flag)
	Iln(tag string, a ...any)
	If(tag string, format string, a ...any)
	Dln(tag string, a ...any)
	Df(tag string, format string, a ...any)
	Wln(tag string, a ...any)
	Wf(tag string, format string, a ...any)
	Eln(tag string, a ...any)
	Ef(tag string, format string, a ...any)
	Fatalln(tag string, a ...any)
	Fatalf(tag string, format string, a ...any)
	Panicln(tag string, a ...any)
	Panicf(tag string, format string, a ...any)
}

type levelInfo struct {
	color   []int
	label   string
	printer func(...any)
}

type loggerImpl struct {
	flag        Flag
	stackCount  int
	stackOffset int
}

func NewLogger() Logger {
	return &loggerImpl{
		stackCount:  0,
		stackOffset: 0,
		flag:        FlagColorEnabled | FlagStackEnabled,
	}
}

func (l *loggerImpl) SetStackCount(stackCount int) {
	if stackCount <= 1 {
		stackCount = 0
	} else {
		stackCount -= 1
	}

	l.stackCount = stackCount
}

func (l *loggerImpl) SetStackOffset(stackOffset int) {
	l.stackOffset = stackOffset
}

func (l *loggerImpl) SetFlag(flag Flag) {
	l.flag = flag
}

func (l *loggerImpl) AddFlag(flag Flag) {
	l.flag &= flag
}

func (l *loggerImpl) Iln(tag string, a ...any) {
	l.println(LevelInfo, tag, fmt.Sprint(a...))
}

func (l *loggerImpl) If(tag string, format string, a ...any) {
	l.println(LevelInfo, tag, format, a...)
}

func (l *loggerImpl) Dln(tag string, a ...any) {
	l.println(LevelDebug, tag, fmt.Sprint(a...))
}

func (l *loggerImpl) Df(tag string, format string, a ...any) {
	l.println(LevelDebug, tag, format, a...)
}

func (l *loggerImpl) Wln(tag string, a ...any) {
	l.println(LevelWarn, tag, fmt.Sprint(a...))
}

func (l *loggerImpl) Wf(tag string, format string, a ...any) {
	l.println(LevelDebug, tag, format, a...)
}

func (l *loggerImpl) Eln(tag string, a ...any) {
	l.println(LevelError, tag, fmt.Sprint(a...))
}

func (l *loggerImpl) Ef(tag string, format string, a ...any) {
	l.println(LevelError, tag, format, a...)
}

func (l *loggerImpl) Fatalln(tag string, a ...any) {
	l.println(LevelFatal, tag, fmt.Sprint(a...))
}

func (l *loggerImpl) Fatalf(tag string, format string, a ...any) {
	l.println(LevelFatal, tag, format, a...)
}

func (l *loggerImpl) Panicln(tag string, a ...any) {
	l.println(LevelPanic, tag, fmt.Sprint(a...))
}

func (l *loggerImpl) Panicf(tag string, format string, a ...any) {
	l.println(LevelPanic, tag, format, a...)
}

func (l *loggerImpl) println(level Level, tag, format string, a ...any) {
	li, _ := levelInfoMap[level]
	text := l.buildLog(li, tag, format, a...)
	li.printer(text)
}

func (l *loggerImpl) buildLog(li *levelInfo, tag, format string, a ...any) string {
	logText := fmt.Sprintf(logTextFormat, li.label, tag, fmt.Sprintf(format, a...))
	logText = l.getColorfulText(li.color, logText)

	stackInfo := ""
	if l.flag&FlagStackEnabled == FlagStackEnabled {
		stackInfo = getStackInfo(baseStackCount+l.stackCount, baseStackOffset+l.stackOffset)
	}

	text := fmt.Sprintf(logFormat, os.Getpid(), logText, stackInfo)
	return text
}

func (l *loggerImpl) getColorfulText(color []int, text string) string {
	colorfulText := text

	if l.flag&FlagColorEnabled == FlagColorEnabled {
		colorfulText = getColorfulText(color, colorfulText)
	}

	return colorfulText
}

var globalLogger Logger

func init() {
	globalLogger = NewLogger()
	globalLogger.SetStackCount(2)
	globalLogger.SetStackOffset(1)
}

func GetGlobalLogger() Logger {
	return globalLogger
}

func Iln(tag string, a ...any) {
	globalLogger.Iln(tag, a...)
}

func If(tag string, format string, a ...any) {
	globalLogger.If(tag, format, a...)
}

func Dln(tag string, a ...any) {
	globalLogger.Dln(tag, a...)
}

func Df(tag string, format string, a ...any) {
	globalLogger.Df(tag, format, a...)
}

func Wln(tag string, a ...any) {
	globalLogger.Wln(tag, a...)
}

func Wf(tag string, format string, a ...any) {
	globalLogger.Wf(tag, format, a...)
}

func Eln(tag string, a ...any) {
	globalLogger.Eln(tag, a...)
}

func Ef(tag string, format string, a ...any) {
	globalLogger.Ef(tag, format, a...)
}

func Fatalln(tag string, a ...any) {
	globalLogger.Fatalln(tag, a...)
}

func Fatalf(tag string, format string, a ...any) {
	globalLogger.Fatalf(tag, format, a...)
}

func Panicln(tag string, a ...any) {
	globalLogger.Panicln(tag, a...)
}
func Panicf(tag string, format string, a ...any) {
	globalLogger.Panicf(tag, format, a...)
}

// util

func getColorfulText(color []int, text string) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", color[0], color[1], color[2], text)
}

func getStackInfo(stackCount, stackOffset int) string {
	if stackCount < 1 {
		stackCount = 1
	}

	pc := make([]uintptr, stackCount+stackOffset)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	var sb strings.Builder
	for i := 0; i < stackCount; i++ {
		frame, more := frames.Next()
		if i < stackOffset {
			if more {
				continue
			}

			break
		}

		// 计算相对路径
		codeFile, err := filepath.Rel(codeRootPath, frame.File)
		if err != nil {
			codeFile = frame.File
		}

		sb.WriteString(fmt.Sprintf("%s:%d\n", codeFile, frame.Line))
		if !more {
			break
		}
	}

	return " < " + sb.String()
}

// 确定项目根目录
func detectProjectRoot(projectDir string) string {
	var buf [1]uintptr
	runtime.Callers(3, buf[:])

	frames := runtime.CallersFrames(buf[:])
	frame, _ := frames.Next()

	file := frame.File

	for {
		if strings.HasSuffix(file, projectDir) ||
			strings.HasSuffix(file, projectDir+"/") {
			break
		}

		file = filepath.Dir(file)
	}

	return file
}
