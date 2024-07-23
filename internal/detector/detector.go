package detector

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/maicher/kmscan/internal/monitor"
)

type Detector struct {
	monitor *monitor.Monitor
	cmdName string
}

func New(m *monitor.Monitor, cmdName string) *Detector {
	return &Detector{
		monitor: m,
		cmdName: cmdName,
	}
}
func (d *Detector) ReadOrDetect(cacheDir string, forceDetect bool) (devices []string, err error) {
	devices, err = d.ReadDevices(cacheDir)
	if err != nil {
		d.monitor.Msg(err.Error(), "unable to read stored devices")
	}

	// Detect devices.
	if err != nil || len(devices) == 0 || forceDetect {
		devices, err = d.DetectScanimageDevices()
		if err != nil {
			return devices, fmt.Errorf("unable to detect scanimage device: %s", err)
		}

		if err := d.StoreDevices(devices, cacheDir); err != nil {
			return devices, fmt.Errorf("unable to save scanimage devices: %s", err)
		}
	}

	return devices, nil
}

func (d *Detector) DetectScanimageDevices() (devices []string, err error) {
	d.monitor.Msg("", "detecting scanimage devices")

	var buf bytes.Buffer
	var r *strings.Reader
	cmd := exec.Command(d.cmdName, "-L")
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.Error); ok {
			return devices, fmt.Errorf("unable to detect scanner: command %s was not found in the system", d.cmdName)
		} else {
			return devices, fmt.Errorf("unable to detect scanner: %s", d.cmdName)
		}
	}

	s := bufio.NewScanner(&buf)
	s.Split(bufio.ScanLines)
	var device string
	var blank string
	replacements := strings.NewReplacer("'", "", "`", "")
	for s.Scan() {

		r = strings.NewReader(s.Text())
		if _, err := fmt.Fscanf(r, "%s %s", &blank, &device); err == nil {
			device = replacements.Replace(device)
			if strings.HasPrefix(device, "pixma:") {
				d.monitor.Msg("detected", s.Text())
				devices = append(devices, device)
			} else {
				d.monitor.Msg("", s.Text())
			}
		}
	}
	if len(devices) == 0 {
		return devices, errors.New("no scanimage devices detected")
	}

	return devices, nil
}

func (d *Detector) ReadDevices(cacheDir string) (devices []string, err error) {
	filePath := path.Join(cacheDir, "devices")
	d.monitor.Msg("", "reading devices from %s", filePath)

	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return devices, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read lines from the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		d.monitor.Msg("", scanner.Text())
		devices = append(devices, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return devices, err
	}

	return devices, nil
}

func (d *Detector) StoreDevices(devices []string, cacheDir string) error {
	filePath := path.Join(cacheDir, "devices")
	d.monitor.Msg("", "storing detected devices in %s", filePath)

	// Ensure the directory exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Open the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write lines to the file
	writer := bufio.NewWriter(file)
	for _, device := range devices {
		fmt.Fprintln(writer, device)
	}

	// Ensure all buffered operations are applied
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %v", err)
	}

	return nil
}
