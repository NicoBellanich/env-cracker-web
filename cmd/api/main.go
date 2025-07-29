package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/upload", uploadHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Could not read the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if filepath.Ext(header.Filename) != ".env" {
		http.Error(w, "Only .env files are allowed", http.StatusUnsupportedMediaType)
		return
	}

	inputData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read the file", http.StatusInternalServerError)
		return
	}

	extractedFiles := ParseCustomFormat(inputData)

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for _, f := range extractedFiles {
		filename := f.Filename
		if !strings.HasSuffix(filename, f.Ext) {
			filename += f.Ext
		}

		fw, err := zipWriter.Create(filename)
		if err != nil {
			http.Error(w, "Error creating zip entry", http.StatusInternalServerError)
			return
		}

		if _, err := fw.Write(f.Content); err != nil {
			http.Error(w, "Error writing file to zip", http.StatusInternalServerError)
			return
		}
	}

	if err := zipWriter.Close(); err != nil {
		http.Error(w, "Error finalizing the zip file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=extracted.zip")
	w.Write(buf.Bytes())
}
