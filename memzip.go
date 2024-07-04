package memzip

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
)

type ZipArchive struct {
	buffer *bytes.Buffer
	writer *zip.Writer
}

// Creates a new in-memory ZIP archive
func NewZipArchive() *ZipArchive {
	buffer := new(bytes.Buffer)
	writer := zip.NewWriter(buffer)
	return &ZipArchive{buffer: buffer, writer: writer}
}

// Adds an existing file or directory to the archive with an optional destination path
func (za *ZipArchive) AddPath(srcPath string, destPath ...string) error {
	info, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	basePath, err := filepath.Abs(srcPath)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		basePath = filepath.Dir(basePath)
	}

	// Determine the base destination path
	var baseDestPath string
	if len(destPath) > 0 {
		baseDestPath = destPath[0]
	} else {
		baseDestPath = ""
	}

	return filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}

		destRelativePath := filepath.Join(baseDestPath, relativePath)

		if info.IsDir() {
			_, err = za.writer.Create(destRelativePath + "/")
			if err != nil {
				return err
			}
		} else {
			return za.addFileToArchive(path, destRelativePath)
		}
		return nil
	})
}

// Adds a file to the archive
func (za *ZipArchive) addFileToArchive(filePath, relativePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	w, err := za.writer.Create(relativePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, file)
	return err
}

// Creates a new file in the archive
func (za *ZipArchive) CreateFile(fileName, content string) error {
	w, err := za.writer.Create(fileName)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(content))
	return err
}

// Closes the zip archive and returns the contents
func (za *ZipArchive) Close() ([]byte, error) {
	if err := za.writer.Close(); err != nil {
		return nil, err
	}
	return za.buffer.Bytes(), nil
}
