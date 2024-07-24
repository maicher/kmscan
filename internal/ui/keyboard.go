package ui

import (
	"log"

	"github.com/eiannone/keyboard"
)

type Keyboard struct {
	Logger *Logger
}

func (k *Keyboard) Listen() error {
	k.Logger.Msg("press 'q' to exit, 'h' to print help", "listening for keyboard events")

	if err := keyboard.Open(); err != nil {
		return err
	}
	defer keyboard.Close()

	for {
		char, _, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if char == 'q' || char == 'Q' {
			k.Logger.Msg("", `quit`)
			break
		}

		if char == 'i' || char == 'I' {
			k.Logger.Msg("", `info
This is an example message.`)
		}

		if char == 'h' || char == 'H' {
			k.Logger.Msg("", `help
Commands:
  i - print info message
  q - quit
  h - print help`)
		}
	}

	return nil
}
