package ui

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type Logger struct {
	prefix string
}

func NewMainLogger() *Logger {
	return &Logger{prefix: color.BlueString("main  ")}
}

func NewScanLogger() *Logger {
	return &Logger{prefix: color.CyanString("scan%d ", 1)}
}

func NewProcLogger() *Logger {
	return &Logger{prefix: color.MagentaString("proc  ")}
}

func NewAPILogger() *Logger {
	return &Logger{prefix: color.HiGreenString("api   ")}
}

func NewUILogger() *Logger {
	return &Logger{prefix: color.HiYellowString("ui    ")}
}

func (l *Logger) Msg(subtext, format string, a ...any) {
	fmt.Println(
		l.prefix,
		fmt.Sprintf(format, a...),
		color.WhiteString(" %s", subtext),
	)
}

func (l *Logger) Err(format string, a ...any) {
	fmt.Println(
		l.prefix,
		color.RedString(format, a...),
	)
}

func (l *Logger) Warn(format string, a ...any) {
	fmt.Println(
		l.prefix,
		color.YellowString(format, a...),
	)
}

func (l *Logger) MsgWithDuration(t time.Duration, format string, a ...any) {
	fmt.Println(
		l.prefix,
		fmt.Sprintf(format, a...),
		color.WhiteString(" %dms", t.Milliseconds()),
	)
}
