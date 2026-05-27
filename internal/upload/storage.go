package upload

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"retrodroid-sync-backend/model"
)

// DiskStorage maneja las operaciones del sistema de archivos local
type BackupStorage struct {
	baseDir string
}

func NewBackupStorage(baseDir string) *BackupStorage {
	return &BackupStorage{baseDir: baseDir}
}

// SaveFile crea los directorios necesarios y guarda el archivo en el disco
func (s *BackupStorage) SaveFile(uploadRequestDTO *model.UploadRequestDTO) (string, error) {

	emulator := uploadRequestDTO.Emulator
	file := uploadRequestDTO.File
	filename := uploadRequestDTO.Header.Filename

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
