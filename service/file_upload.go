package service

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Field struct {
	FieldName  string
	FieldValue string
}

type FileUploadConfig struct {
	FilePath      string
	FieldValues   []Field
	UploadURL     string
	Authorization string
}

func PostFile(config FileUploadConfig) error {
	// Open the file
	file, err := os.Open(config.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	// Create a buffer for the POST request body
	body := &bytes.Buffer{}

	// Create a multipart writer
	writer := multipart.NewWriter(body)

	// Create a form file writer for the given file field
	formFile, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return fmt.Errorf("failed to create form file writer: %w", err)
	}

	// Copy the file into the form file writer
	if _, err = io.Copy(formFile, file); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Iterate over FieldValues and add them to the form
	for _, f := range config.FieldValues {
		if f.FieldName == "" || f.FieldValue == "" {
			continue
		}
		err := writer.WriteField(f.FieldName, f.FieldValue)
		if err != nil {
			return err
		}
	}

	err = writer.WriteField("private", "true")
	if err != nil {
		return err
	}

	// Close the multipart writer
	if err = writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	// Create a new POST request

	log.Printf("Uploading %s...\n", config.FilePath)

	req, err := http.NewRequest("POST", config.UploadURL, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", config.Authorization)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}

	body = &bytes.Buffer{}

	// Copy the response body to the buffer
	if _, err = io.Copy(body, response.Body); err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("%s uploaded succesfully\n", config.FilePath)
	return nil
}
