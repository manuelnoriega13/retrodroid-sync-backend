package model

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
)

type UploadResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Emulator string `json:"emulator"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
}

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type UploadRequestDTO struct {
	Emulator string
	File     multipart.File
	Header   *multipart.FileHeader
}

func ParseUploadRequest(r *http.Request) (*UploadRequestDTO, error) {
	emulator := strings.TrimSpace(r.FormValue("emulator"))
	if emulator == "" {
		return nil, errors.New("missing parameter 'emulator'")
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve 'file' field: %w", err)
	}

	return &UploadRequestDTO{
		Emulator: emulator,
		File:     file,
		Header:   header,
	}, nil
}
