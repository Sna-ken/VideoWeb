package utils

import (
	"path/filepath"
	"strings"
)

var allowedFileTypes = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

func GetFileExtension(filename string, data []byte) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if allowedFileTypes[ext] {
		return ext
	}

	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return ".jpg"
	}
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return ".jpg"
	}
	if len(data) >= 8 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return ".png"
	}

	return ".jpg"
}
