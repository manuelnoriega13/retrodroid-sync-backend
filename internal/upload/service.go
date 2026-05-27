package upload

import (
	"errors"
	"io"
	"path/filepath"
	"strings"
)

// Service contiene la lógica de negocio para las subidas
type Service struct {
	storage *DiskStorage
}

func NewService(storage *DiskStorage) *Service {
	return &Service{storage: storage}
}

// ProcessUpload sanitiza las entradas, valida el archivo y lo envía al storage
func (s *Service) ProcessUpload(emulator, originalFilename string, file io.Reader) (string, string, string, error) {
	// Sanitizar el nombre del emulador
	sanitizedEmulator := filepath.Base(filepath.Clean(emulator))
	if sanitizedEmulator == "." || sanitizedEmulator == ".." || sanitizedEmulator == "/" {
		return "", "", "", errors.New("invalid or unsafe 'emulator' parameter value")
	}

	// Validar extensión estricta
	ext := strings.ToLower(filepath.Ext(originalFilename))
	if ext != ".zip" {
		return "", "", "", errors.New("invalid file format. Only '.zip' archives are permitted")
	}

	// Sanitizar el nombre del archivo
	sanitizedFilename := filepath.Base(filepath.Clean(originalFilename))
	if sanitizedFilename == "." || sanitizedFilename == ".." || sanitizedFilename == "/" {
		return "", "", "", errors.New("invalid or unsafe upload filename")
	}

	// Delegar el guardado físico a la capa de storage
	savedPath, err := s.storage.SaveFile(sanitizedEmulator, sanitizedFilename, file)
	if err != nil {
		return "", "", "", err
	}

	return sanitizedEmulator, sanitizedFilename, savedPath, nil
}
