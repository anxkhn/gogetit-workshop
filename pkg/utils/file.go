package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func IsValidPath(path string) bool {
	return !strings.Contains(path, "..")
}

func IsPathTraversal(path string) bool {
	cleaned := filepath.Clean(path)
	return strings.Contains(cleaned, "..")
}

func SanitizePath(baseDir, userPath string) string {
	cleaned := filepath.Clean(userPath)
	fullPath := filepath.Join(baseDir, cleaned)
	return fullPath
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func CreateFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return nil, err
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func SanitizeFilename(name string) string {
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, ":", "_")
	name = strings.ReplaceAll(name, "*", "_")
	name = strings.ReplaceAll(name, "?", "_")
	name = strings.ReplaceAll(name, "\"", "_")
	name = strings.ReplaceAll(name, "<", "_")
	name = strings.ReplaceAll(name, ">", "_")
	name = strings.ReplaceAll(name, "|", "_")

	if len(name) > 200 {
		name = name[:200]
	}

	return name
}

func WriteFile(path string, data []byte) error {
	file, err := CreateFile(path)
	if err != nil {
		return err
	}

	_, writeErr := file.Write(data)
	closeErr := file.Close()

	// Surface the write error first when both fail, it's the
	// actionable cause. A close error after a successful write still
	// matters: buffered data may not have flushed and the file on
	// disk would be silently truncated.
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}
