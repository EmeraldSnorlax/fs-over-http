package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func ReadFileUnsafe(file string, removeNewline bool) string {
	content, err := ReadFile(file)

	if err != nil {
		log.Printf("- Failed to read '%s'", file)
		return content
	}

	lastByte, cb := LastByte(content)

	if removeNewline && lastByte == 10 {
		cb = cb[:len(cb)-1]
		return string(cb)
	}

	return content
}

func ReadFile(file string) (string, error) {
	dat, err := ioutil.ReadFile(file)
	return string(dat), err
}

func WriteToFile(file string, content string) error {
	data := []byte(content)
	err := ioutil.WriteFile(file, data, 0600)

	if err != nil {
		log.Printf("- Failed to read '%s'", file)
	}

	return err
}

func IsDirectory(path string) (bool, error) {
	fi, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	if fi.Mode().IsDir() {
		return true, nil
	} else {
		return false, nil
	}
}

func GetFileContentTypeExt(out *os.File, file string) (string, error) {
	ext := filepath.Ext(file)

	switch ext {
	case ".txt", ".text":
		return "text/plain", nil
	case ".htm", ".html":
		return "text/html", nil
	case ".css":
		return "text/css", nil
	case ".js", ".mjs":
		return "application/javascript", nil
	case ".mov":
		return "video/quicktime", nil
	}

	return GetFileContentType(out)
}

// GetFileContentType detects the content type
// and returns a valid MIME type
func GetFileContentType(out *os.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DetectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

// ReadLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// ReadNonEmptyLines reads non-empty lines
// using ReadLines and a filter.
func ReadNonEmptyLines(file string) (ret []string) {
	lines, err := ReadLines(file)

	if err != nil {
		return []string{}
	}

	for _, s := range lines {
		lastByte, _ := LastByte(s)

		if len(s) > 1 || lastByte != 10 {
			ret = append(ret, s)
		}
	}

	return ret
}
