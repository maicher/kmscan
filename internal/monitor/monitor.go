package monitor

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type scope int

const (
	initializer scope = iota
	scanner
	processor
)

type log struct {
	scope   scope
	text    string
	subtext string
	isError bool
}

type Monitor struct {
	msgCh chan log
}

func New() *Monitor {
	m := Monitor{
		msgCh: make(chan log),
	}
	go m.loop()

	return &m
}

func (m *Monitor) Init(format string, a ...any) {
	m.msgCh <- log{
		scope: initializer,
		text:  fmt.Sprintf(format, a...),
	}
}

func (m *Monitor) InitError(format string, a ...any) {
	m.msgCh <- log{
		scope:   initializer,
		text:    fmt.Sprintf(format, a...),
		isError: true,
	}
}

func (m *Monitor) Processor(t time.Duration, format string, a ...any) {
	m.msgCh <- log{
		scope:   processor,
		text:    fmt.Sprintf(format, a...),
		subtext: fmt.Sprintf("%dms", t.Milliseconds()),
	}
}

func (m *Monitor) ProcessorRatio(ratio float64, format string, a ...any) {
	m.msgCh <- log{
		scope:   processor,
		text:    fmt.Sprintf(format, a...),
		subtext: fmt.Sprintf("ratio %.2f", ratio),
	}
}
func (m *Monitor) ProcessorError(format string, a ...any) {
	m.msgCh <- log{
		scope:   processor,
		text:    fmt.Sprintf(format, a...),
		isError: true,
	}
}

func (m *Monitor) Scanner(subtext, format string, a ...any) {
	m.msgCh <- log{
		scope:   scanner,
		text:    fmt.Sprintf(format, a...),
		subtext: subtext,
	}
}

func (m *Monitor) loop() {
	for msg := range m.msgCh {
		switch msg.scope {
		case initializer:
			color.Set(color.FgBlue, color.Bold)
			fmt.Printf("init      ")
			if msg.isError {
				color.Set(color.FgRed, color.Bold)
				fmt.Printf("%s\n", msg.text)
			} else {
				color.Unset()
				fmt.Printf("%s\n", msg.text)
			}
			color.Unset()
		case scanner:
			color.Set(color.FgMagenta, color.Bold)
			fmt.Printf("scanner   ")
			color.Unset()
			fmt.Printf("%s", msg.text)
			if msg.subtext == "" {
				fmt.Printf("\n")
			} else {
				color.Set(color.FgWhite)
				fmt.Printf(" (%s)\n", msg.subtext)
			}
			color.Unset()
		case processor:
			color.Set(color.FgBlue)
			fmt.Printf("processor ")
			if msg.isError {
				color.Set(color.FgRed, color.Bold)
				fmt.Printf("%s\n", msg.text)
			} else {
				color.Unset()
				fmt.Printf(msg.text)
				color.Set(color.FgWhite)
				if msg.subtext == "" {
					fmt.Printf("\n")
				} else {
					color.Set(color.FgWhite)
					fmt.Printf(" (%s)\n", msg.subtext)
				}
			}
			color.Unset()
		}
	}
}
