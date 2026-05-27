package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"retrodroid-sync-backend/internal/upload"
)

func main() {
	// 1. Definir el flag con un string vacío como valor por defecto
	dirSavePath := flag.String("dir", "", "Directorio base donde se guardarán los saves de los emuladores (REQUERIDO)")

	// Parsear los argumentos pasados por consola
	flag.Parse()

	// Validación estricta: Si el usuario no pasó el flag, abortamos
	if *dirSavePath == "" {
		fmt.Println("Error: El argumento '-dir' es obligatorio.")
		fmt.Println("Uso correcto: ./retrodroid-sync -dir=/ruta/absoluta/al/directorio")
		fmt.Println("\nOpciones disponibles:")
		flag.PrintDefaults() // Muestra la ayuda del comando automáticamente
		os.Exit(1)           // Código 1 indica que el programa terminó por un error
	}

	// 2. Inicializar dependencias
	storage := upload.NewBackupStorage(*dirSavePath)
	service := upload.NewBackupService(storage)
	backup := upload.NewBackupHandler(service)

	// 3. Registrar rutas
	http.HandleFunc("/upload", backup.HandleBackupUpload)

	// 4. Iniciar el servidor
	port := ":8080"
	fmt.Printf("🚀 Retrodroid Sync Backend iniciando...\n")
	fmt.Printf("📁 Guardando los archivos en el directorio: %s\n", *dirSavePath)
	fmt.Printf("📡 Escuchando en http://localhost%s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
		os.Exit(1)
	}
}
