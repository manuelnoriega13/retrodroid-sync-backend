package upload

import (
	"encoding/json"
	"fmt"
	"net/http"
	"retrodroid-sync-backend/model"
	"retrodroid-sync-backend/util"
)

var UploadSize = util.MaxSizeMB(50)

type BackupHandler struct {
	backupService *BackupService
}

func NewBackupHandler(backupService *BackupService) *BackupHandler {
	return &BackupHandler{backupService: backupService}
}

func (h *BackupHandler) HandleBackupUpload(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed. Only POST requests are permitted")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, UploadSize)

	if err := r.ParseMultipartForm(UploadSize); err != nil {
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
