package utils

import (
	"fmt"
	"os"
)

func StoreAvatar(avatarFile []byte, userID string) (string, error) {
	basePath := "./static/avatar"
	userPath := basePath + "/" + userID
	if len(avatarFile) == 0 {
		return "", fmt.Errorf("avatar file is empty (0 bytes)")
	}

	if err := os.MkdirAll(userPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directories: %w", err)
	}

	filePath := fmt.Sprintf("%s/avatar.png", userPath)
	err := os.WriteFile(filePath, avatarFile, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}
	if fileInfo.Size() == 0 {
		return "", fmt.Errorf("file was created but is empty (0 bytes)")
	}

	// 返回 URL
	avatarURL := fmt.Sprintf("http://127.0.0.1:8888/avatar/%s/avatar.png", userID)
	return avatarURL, nil
}
