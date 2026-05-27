package upload

import (
	"encoding/json"
	"fmt"
	"net/http"
	"retrodroid-sync-backend/model"
)

type BackupHandler struct {
	backupService *BackupService
	uploadSize    int64
}

func NewBackupHandler(backupService *BackupService, uploadSize int64) *BackupHandler {
	return &BackupHandler{
		backupService: backupService,
		uploadSize:    uploadSize,
	}
}

func (h *BackupHandler) HandleBackupUpload(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed. Only POST requests are permitted")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, h.uploadSize)

	if err := r.ParseMultipartForm(h.uploadSize); err != nil {
		h.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse multipart form: %v", err))
		return
	}

	UploadRequestDTO, err := model.ParseUploadRequest(r)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer UploadRequestDTO.File.Close()

	sanitizedEmulator, sanitizedFilename, destFilePath, err := h.backupService.ProcessUpload(UploadRequestDTO)

	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.UploadResponse{
		Status:   "success",
		Message:  "Save file synchronised successfully",
		Emulator: sanitizedEmulator,
		Filename: sanitizedFilename,
		Path:     destFilePath,
	})
}

func (h *BackupHandler) respondWithError(w http.ResponseWriter, code int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(model.ErrorResponse{
		Status: "error",
		Error:  errMsg,
	})
}
