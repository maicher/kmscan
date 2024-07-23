package kmscan

import (
	"log"

	"github.com/eiannone/keyboard"
	"github.com/maicher/kmscan/internal/monitor"
)

type KeyboardController struct {
	Monitor *monitor.Monitor
}

func (k *KeyboardController) Listen() error {
	k.Monitor.Msg("press 'q' to exit, 'h' to print help", "listening for keyboard events")
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
			k.Monitor.Msg("", "q")
			break
		}

		if char == 'i' || char == 'I' {
			k.Monitor.Msg("", "i")
			k.Monitor.Msg("", `
This is an example message.`)
		}

		if char == 'h' || char == 'H' {
			k.Monitor.Msg("", "h")
			k.Monitor.Msg("", `
----Help----
i   print info message
q   quit
h   print help
------------`)
		}
	}

	return nil
}
