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

type Uploader interface {
	Upload(string) error
}

func NewUploader(dirPath, apiURL, apiKey string, logger *ui.Logger) Uploader {
	if apiURL == "" {
		logger.Warn("no api-url, images will not be uploaded")
		return &NullUploader{}
	}

	return &APIUploader{
		DirPath: dirPath,
		APIURL:  apiURL,
		APIKey:  apiKey,
		Logger:  logger,
	}
}

type NullUploader struct{}

func (u NullUploader) Upload(_ string) error {
	return nil
}

type APIUploader struct {
	DirPath string
	APIURL  string
	APIKey  string
	Logger  *ui.Logger

	client http.Client
}

type FileUploadRequest struct {
	Filename string `json:"filename"`
	Data     string `json:"data"`
}

func (u APIUploader) Upload(name string) error {
	if err := u.upload(name); err != nil {
		u.Logger.Err(err.Error())

		return err
	}

	return nil
}

func (u APIUploader) upload(name string) error {
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
	req, err := http.NewRequest("POST", u.APIURL, bytes.NewBuffer(jsonData))
	req.Header.Set("api-key", u.APIKey)
	if err != nil {
		return fmt.Errorf("error building request: %w", err)
	}
	resp, err := u.client.Do(req)
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
