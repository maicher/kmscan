package kmscan

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/maicher/kmscan/internal/ui"
)

type Uploader struct {
	DirPath string
	Logger *ui.Logger
}

type FileUploadRequest struct {
	Filename string `json:"filename"`
	Data     string `json:"data"`
}

func (u Uploader) Upload(name string) error {
	if err := u.upload(name); err != nil {
		u.Logger.Err(err.Error())

		return err
	}

	return nil
}

func (u Uploader) upload(name string) error {
	outputPath := filepath.Join(u.DirPath, name)
	t := time.Now()
	outFile, err := os.Open(outputPath)
	if err != nil {
		return err
	}

	// Read the file
	fileData, err := io.ReadAll(outFile)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Encode the file to base64
	encodedData := base64.StdEncoding.EncodeToString(fileData)

	// Create the request payload
	fileReq := FileUploadRequest{
		Filename: name,
		Data:     encodedData,
	}
	jsonData, err := json.Marshal(fileReq)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	// Send the request
	resp, err := http.Post("http://localhost:9500/api/upload", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	u.Logger.MsgWithDuration(time.Since(t), "%s %s", name, respBody)

	return nil
}
