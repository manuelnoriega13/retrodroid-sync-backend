package upload

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// DiskStorage maneja las operaciones del sistema de archivos local
type DiskStorage struct {
	baseDir string
}

func NewDiskStorage(baseDir string) *DiskStorage {
	return &DiskStorage{baseDir: baseDir}
}

// SaveFile crea los directorios necesarios y guarda el archivo en el disco
func (s *DiskStorage) SaveFile(emulator, filename string, file io.Reader) (string, error) {
	storageDir := filepath.Join(s.baseDir, emulator)

	// Crear el directorio si no existe
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage directories: %w", err)
	}

	destFilePath := filepath.Join(storageDir, filename)

	// Crear el archivo de destino
	dst, err := os.Create(destFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copiar los datos
	if _, err = io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save uploaded file content: %w", err)
	}

	return destFilePath, nil
}
