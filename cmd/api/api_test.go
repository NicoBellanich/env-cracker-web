package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindTXTContentIndex(t *testing.T) {
	chunk := append([]byte("some headers\r\n_SIG/D.C."), []byte("   Text content here")...)
	index := FindTXTContentIndex(chunk)

	if index == -1 {
		t.Fatalf("Expected to find text content start after signature, but got -1")
	}

	if !bytes.HasPrefix(chunk[index:], []byte("Text")) {
		t.Errorf("Expected content to start with 'Text', got: %s", chunk[index:index+10])
	}
}

func TestExtractMetadataAndData(t *testing.T) {
	sample := `FILENAME/testfile
EXT/.txt
_SIG/D.C. Text begins here`

	chunk := []byte(sample)
	result := extractMetadataAndData(chunk)

	if result == nil {
		t.Fatalf("Expected non-nil result")
	}

	if result.Filename != "testfile" {
		t.Errorf("Expected filename 'testfile', got '%s'", result.Filename)
	}
	if result.Ext != ".txt" {
		t.Errorf("Expected extension '.txt', got '%s'", result.Ext)
	}
	if !bytes.Contains(result.Content, []byte("Text begins here")) {
		t.Errorf("Expected content to include 'Text begins here'")
	}
}

func TestParseCustomFormat(t *testing.T) {
	input := []byte(`ignore this**%%DOCUFILENAME/file1
EXT/.txt
_SIG/D.C. Hello from file1**%%DOCUFILENAME/file2
EXT/.txt
_SIG/D.C. Hello from file2`)

	files := ParseCustomFormat(input)

	if len(files) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(files))
	}

	if files[0].Filename != "file1" || !bytes.Contains(files[0].Content, []byte("Hello from file1")) {
		t.Errorf("Unexpected content in first file")
	}

	if files[1].Filename != "file2" || !bytes.Contains(files[1].Content, []byte("Hello from file2")) {
		t.Errorf("Unexpected content in second file")
	}
}

func TestUploadHandler_WithValidEnvFile(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "sample.env")
	if err != nil {
		t.Fatal(err)
	}

	content := `**%%DOCUFILENAME/test
EXT/.txt
_SIG/D.C. Hello test file`
	part.Write([]byte(content))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	uploadHandler(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %d", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType != "application/zip" {
		t.Errorf("Expected content-type application/zip, got %s", contentType)
	}

	data, _ := io.ReadAll(res.Body)
	if len(data) == 0 {
		t.Error("Expected zip content in response")
	}
}

func TestUploadHandler_RejectsNonEnvFile(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("This is not a .env file"))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	uploadHandler(rr, req)

	if rr.Code != http.StatusUnsupportedMediaType {
		t.Errorf("Expected status 415 UnsupportedMediaType, got %d", rr.Code)
	}
}
