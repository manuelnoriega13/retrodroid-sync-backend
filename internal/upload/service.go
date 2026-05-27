package upload

import (
	"errors"
	"fmt"
	"path/filepath"
	"retrodroid-sync-backend/model"
	"strings"
	"time"
)

// Service contiene la lógica de negocio para las subidas
type BackupService struct {
	storage *BackupStorage
}

func NewBackupService(storage *BackupStorage) *BackupService {
	return &BackupService{storage: storage}
}

// ProcessUpload sanitiza las entradas, inyecta la hora militar, valida el archivo y lo envía al storage
func (s *BackupService) ProcessUpload(uploadRequestDTO *model.UploadRequestDTO) (string, string, string, error) {
	// 1. Sanitizar el nombre del emulador (ej: "ppsspp")

	emulator := uploadRequestDTO.Emulator
	filename := uploadRequestDTO.Header.Filename

	sanitizedEmulator := filepath.Base(filepath.Clean(emulator))
	if sanitizedEmulator == "." || sanitizedEmulator == ".." || sanitizedEmulator == "/" {
		return "", "", "", errors.New("invalid or unsafe 'emulator' parameter value")
	}

	// 2. Validar extensión estricta
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".zip" {
		return "", "", "", errors.New("invalid file format. Only '.zip' archives are permitted")
	}

	// 3. Sanitizar el nombre base del archivo (ej: "ppsspp-24-12-26.zip")
	baseName := filepath.Base(filepath.Clean(filename))
	if baseName == "." || baseName == ".." || baseName == "/" {
		return "", "", "", errors.New("invalid or unsafe upload filename")
	}

	// 4. Obtener SOLO la hora en formato militar (HHMMSS)
	militaryTime := time.Now().Format("150405")

	// 5. Quitar el nombre del emulador del archivo original si ya lo trae para no duplicarlo
	// Si baseName es "ppsspp-24-12-26.zip", remainder se convierte en "24-12-26.zip"
	// Si baseName es "24-12-26.zip", se queda igual.
	remainder := strings.TrimPrefix(baseName, sanitizedEmulator+"-")

	// 6. Armar el nombre final: ppsspp-135659-24-12-26.zip
	finalFilename := fmt.Sprintf("%s-%s-%s", sanitizedEmulator, militaryTime, remainder)

	// Delegar el guardado físico a la capa de storage
	savedPath, err := s.storage.SaveFile(uploadRequestDTO)
	if err != nil {
		return "", "", "", err
	}

	return sanitizedEmulator, finalFilename, savedPath, nil
}
