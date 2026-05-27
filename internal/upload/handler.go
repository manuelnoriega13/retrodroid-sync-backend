package upload

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const MaxUploadSize = 50 * 1024 * 1024 // 50 Megabytes

// Estructuras de respuesta JSON
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

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ServeHTTP procesa la petición multipart
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed. Only POST requests are permitted")
		return
	}

	// Limitar el tamaño del body
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		h.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse multipart form: %v", err))
		return
	}

	emulator := r.FormValue("emulator")
	if strings.TrimSpace(emulator) == "" {
		h.respondWithError(w, http.StatusBadRequest, "Missing parameter 'emulator'")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to retrieve 'file' field: %v", err))
		return
	}
	defer file.Close()

	// Enviar al servicio
	sanitizedEmulator, sanitizedFilename, destFilePath, err := h.service.ProcessUpload(emulator, header.Filename, file)
	if err != nil {
		// En un entorno de producción, aquí diferenciaríamos entre errores 400 y 500
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Responder éxito
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(UploadResponse{
		Status:   "success",
		Message:  "Save file synchronised successfully",
		Emulator: sanitizedEmulator,
		Filename: sanitizedFilename,
		Path:     destFilePath,
	})
}

func (h *Handler) respondWithError(w http.ResponseWriter, code int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Status: "error",
		Error:  errMsg,
	})
}
