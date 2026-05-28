# 🎮 Proyecto: Retrodroid Sync - Backend

## 1. Descripción General

**Retrodroid Sync** es un microservicio backend escrito en Go. Su única responsabilidad es recibir archivos `.zip` generados por una aplicación de Android (que ya contienen las partidas guardadas en su interior). El backend no debe analizar el contenido del `.zip`; simplemente debe aceptarlo y guardarlo en el disco local de forma organizada según la consola/emulador.

## 2. Estado Actual y Estructura del Proyecto

Actualmente, el proyecto tiene una estructura minimalista y un archivo base que debe ser modificado.

**Árbol de directorios:**

```text
.
├── CONTEXT.md
├── cmd
│   └── server
│       └── main.go
└── go.mod
```

**Estado base de `cmd/server/main.go`:**

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello World")
}
```

## 3. Contrato de la API y Lógica de Negocio

- **Endpoint:** `POST /upload`
- **Content-Type:** `multipart/form-data`
- **Campos del form esperados:**
  - `emulator` (string): Nombre del emulador o app (ej. `ppsspp`, `duckstation`).
  - `file` (file): El archivo `.zip` con el nombre ya generado por el cliente Android (ej. `ppsspp-24-12-26.zip`).
- **Lógica de Guardado:** El backend debe leer el campo `emulator`, crear dinámicamente la carpeta `./storage/{emulator}/` (si no existe), y guardar el archivo recibido exactamente con el nombre que trajo, resultando en: `./storage/ppsspp/ppsspp-24-12-26.zip`.

## 4. Reglas Estrictas de Código

- **Lenguaje:** Go (Golang) 1.21+.
- **Librerías:** Únicamente la librería estándar de Go (`net/http`, `os`, `io`, `path/filepath`, `fmt`). Cero frameworks externos.
- **Seguridad y Memoria:** Limitar el tamaño del body de la petición HTTP usando `http.MaxBytesReader` para evitar saturación de memoria.
- **Manejo de Errores:** Verificar cada error explícitamente (`if err != nil`) y devolver los HTTP Status Codes correctos (400, 500) en formato texto o JSON.

## 5. Pasos Exactos de Desarrollo para la IA

La IA debe ejecutar las siguientes modificaciones sobre `cmd/server/main.go` en este orden estricto:

1. **Reemplazo del código base:**
   Eliminar el `fmt.Println("Hello World")` dentro de la función `main()` para dar lugar al servidor HTTP.
2. **Configuración del Servidor y Rutas:**
   Registrar la ruta `/upload` apuntando a una función manejadora (ej. `uploadHandler`) y usar `http.ListenAndServe` para iniciar el servidor en el puerto `:8080`.
3. **Implementación del Handler (`uploadHandler`):**
   Validar que el método HTTP sea `POST`. Aplicar `http.MaxBytesReader` al cuerpo de la petición (limitar a 50MB). Parsear el formulario multipart mediante `r.ParseMultipartForm`.
4. **Extracción y Validación de Datos:**
   Obtener el valor del campo `emulator`. Si está vacío, devolver error `400 Bad Request`. Sanitizar el string de `emulator` usando `filepath.Base` para prevenir vulnerabilidades de ruta. Obtener el archivo del campo `file` (`r.FormFile`). Validar que la extensión del archivo sea `.zip`.
5. **I/O y Sistema de Archivos:**
   Construir la ruta del directorio de destino: `./storage/{emulator}`. Usar `os.MkdirAll` con permisos adecuados (0755) para crear el directorio si no existe. Crear el archivo de destino localmente con `os.Create` usando la ruta completa (directorio + nombre del archivo original). Transferir los bytes del archivo subido al disco usando `io.Copy` desde el lector del archivo al escritor.
6. **Finalización:**
   Asegurar el cierre de recursos usando `defer file.Close()` y `defer dst.Close()`. Devolver una respuesta `http.StatusOK` (200) en formato JSON o texto plano confirmando que el archivo se guardó correctamente en la ruta especificada.
