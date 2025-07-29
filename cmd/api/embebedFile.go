package main

import (
	"bufio"
	"bytes"
	"strings"
)

type EmbeddedFile struct {
	Filename string
	Ext      string
	Content  []byte
}

func ParseCustomFormat(data []byte) []EmbeddedFile {
	chunks := bytes.Split(data, []byte("**%%DOCU"))
	var result []EmbeddedFile

	for i, chunk := range chunks {
		if i == 0 {
			continue
		}

		metaAndData := extractMetadataAndData(chunk)
		if metaAndData != nil {
			result = append(result, *metaAndData)
		}
	}

	return result
}

func extractMetadataAndData(chunk []byte) *EmbeddedFile {
	reader := bufio.NewReader(bytes.NewReader(chunk))
	var filename string
	var ext string
	var content []byte

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		lineStr := string(line)

		if strings.HasPrefix(lineStr, "FILENAME/") {
			filename = strings.TrimPrefix(lineStr, "FILENAME/")
		}
		if strings.HasPrefix(lineStr, "EXT/") {
			ext = strings.TrimPrefix(lineStr, "EXT/")
		}

		if strings.HasPrefix(lineStr, "<?xml") || strings.Contains(lineStr, "RIFF") || strings.Contains(lineStr, "\xFF\xD8\xFF") {
			break
		}
	}

	if filename == "" {
		return nil
	}

	contentIndex := bytes.Index(chunk, []byte("<?xml"))
	if contentIndex == -1 {
		contentIndex = bytes.Index(chunk, []byte("RIFF"))
	}
	if contentIndex == -1 {
		contentIndex = bytes.Index(chunk, []byte("\xFF\xD8\xFF"))
	}
	if contentIndex == -1 && ext == ".txt" {
		contentIndex = FindTXTContentIndex(chunk)
	}

	if contentIndex == -1 {
		return nil
	}

	content = chunk[contentIndex:]

	return &EmbeddedFile{
		Filename: filename,
		Ext:      ext,
		Content:  content,
	}
}

func FindTXTContentIndex(chunk []byte) int {
	sigMarker := []byte("_SIG/D.C.")
	sigIndex := bytes.Index(chunk, sigMarker)
	if sigIndex != -1 {
		start := sigIndex + len(sigMarker)
		for start < len(chunk) {
			c := chunk[start]
			if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
				return start
			}
			start++
		}
	}
	return -1
}
