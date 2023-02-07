package utils

import (
	"os"
	"path/filepath"
)

func GetFilesInDirectory(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, info.Name())
		}
		return nil
	})
	return files, err
}

func GetFileExtension(file string) string {
	return filepath.Ext(file)
}

func IsDirectoryExists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}
