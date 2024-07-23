package kmscan

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/maicher/kmscan/internal/monitor"
)

type Scanner struct {
	Resolution int
	Monitor    *monitor.Monitor
}

func (s *Scanner) Scan(ctx context.Context, device string, deviceNo, i int) (string, error) {
	var calibrate string

	// Open a temporary file in RAM (in the /dev/shm directory)
	file, err := os.CreateTemp("/dev/shm", "scan_*.png")
	if err != nil {
		s.Monitor.Err("failed to create temporary file: %s", err)
	}

	if i%10 == 0 {
		s.Monitor.Msg("press SEND button on the scanner", "ready for scan %d with calibrate", i+1)
		calibrate = "always"
	} else {
		s.Monitor.Msg("press SEND button on the scanner", "ready for scan %d", i+1)
		calibrate = "never"
	}

	cmd := exec.CommandContext(ctx, "scanimage",
		"-d", device,
		"--resolution", fmt.Sprintf("%ddpi", s.Resolution),
		"--mode", "Color",
		"--format", "png",
		"--calibrate", calibrate,
		"--button-controlled=yes",
	)

	cmd.Stdout = file
	if err = cmd.Run(); err != nil {
		if ctx.Err() == context.Canceled {
			s.Monitor.Msg("cancellation", "scanning interrupted")
			return "", errors.New("cencelled")
		} else {
			s.Monitor.Err("error scanning: %s", err)
			return "", err
		}
	}
	s.Monitor.Msg("", "scanning %d complete", i+1)

	return file.Name(), nil
}
