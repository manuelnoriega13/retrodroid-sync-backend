package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"retrodroid-sync-backend/internal/upload"
	"retrodroid-sync-backend/util"
)

func main() {
	dirSavePath := flag.String("dir", "", "Directorio base donde se guardarán los saves de los emuladores (REQUERIDO)")
	uploadSizeMB := flag.Int64("uploadSize", 50, "Tamaño máximo del archivo en Megabytes (opcional)")

	flag.Parse()

	if *dirSavePath == "" {
		log.Println("Error: El argumento '-dir' es obligatorio.")
		log.Println("Uso correcto: ./retrodroid-sync -dir=/ruta/absoluta/al/directorio -uploadSize=50")
		log.Println("\nOpciones disponibles:")
		flag.PrintDefaults() // Muestra la ayuda del comando automáticamente
		os.Exit(1)           // Código 1 indica que el programa terminó por un error
	}

	uploadSizeBytes := util.MaxSizeMB(*uploadSizeMB)

	storage := upload.NewBackupStorage(*dirSavePath)
	service := upload.NewBackupService(storage)
	backup := upload.NewBackupHandler(service, uploadSizeBytes)

	http.HandleFunc("/upload", backup.HandleBackupUpload)

	port := ":8080"
	log.Println("🚀 Retrodroid Sync Backend iniciando...")
	log.Printf("📁 Guardando los archivos en el directorio: %s", *dirSavePath)
	log.Printf("📦 Tamaño máximo de subida: %d MB", *uploadSizeMB)
	log.Printf("📡 Escuchando en http://localhost%s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Printf("Error al iniciar el servidor: %v", err)
		os.Exit(1)
	}
}
