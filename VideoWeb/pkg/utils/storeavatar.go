package utils

import (
	"fmt"
	"os"
)

func StoreAvatar(avatarFile []byte, userID string) (string, error) {
	basePath := "./avatars"
	userPath := basePath + "/" + userID
	if err := os.MkdirAll(userPath, os.ModePerm); err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("%s/avatar.png", userPath)
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	avatarURL := fmt.Sprintf("http://127.0.0.1:8080/avatars/%s", userID)
	return avatarURL, nil
}
