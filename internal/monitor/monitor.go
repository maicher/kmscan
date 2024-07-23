package monitor

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type Monitor struct {
	prefix string
}

func NewMain() *Monitor {
	return &Monitor{prefix: color.BlueString("main     ")}
}

func NewScanner() *Monitor {
	return &Monitor{prefix: color.CyanString("scanner%d ", 1)}
}

func NewProcessor() *Monitor {
	return &Monitor{prefix: color.MagentaString("procesor ")}
}

func NewAPI() *Monitor {
	return &Monitor{prefix: color.YellowString("api      ")}
}

func NewUI() *Monitor {
	return &Monitor{prefix: color.HiYellowString("ui       ")}
}

func (m *Monitor) Msg(subtext, format string, a ...any) {
	fmt.Println(
		m.prefix,
		fmt.Sprintf(format, a...),
		color.WhiteString(" %s", subtext),
	)
}

func (m *Monitor) Err(format string, a ...any) {
	fmt.Println(
		m.prefix,
		color.RedString(format, a...),
	)
}

func (m *Monitor) MsgWithDuration(t time.Duration, format string, a ...any) {
	fmt.Println(
		m.prefix,
		fmt.Sprintf(format, a...),
		color.WhiteString(" %ds", t.Milliseconds()),
	)
}
