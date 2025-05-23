package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) (*FileStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &FileStorage{basePath: basePath}, nil
}

func (fs *FileStorage) SaveFile(file io.Reader, filename string) (string, error) {
	// Generate a random ID for the file
	id := make([]byte, 16)
	if _, err := rand.Read(id); err != nil {
		return "", fmt.Errorf("failed to generate file ID: %w", err)
	}
	fileID := hex.EncodeToString(id)

	// Create the file path
	filePath := filepath.Join(fs.basePath, fileID)

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy the file content
	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Schedule file deletion after 24 hours
	go func() {
		time.Sleep(24 * time.Hour)
		os.Remove(filePath)
	}()

	return fileID, nil
}

func (fs *FileStorage) SaveFileWithID(file io.Reader, filename string, fileID string) (string, error) {
	// Create the file path
	filePath := filepath.Join(fs.basePath, fileID)

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy the file content
	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Schedule file deletion after 24 hours
	go func() {
		time.Sleep(24 * time.Hour)
		os.Remove(filePath)
	}()

	return filePath, nil
}

func (fs *FileStorage) GetFilePath(fileID string) (string, error) {
	filePath := filepath.Join(fs.basePath, fileID)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	}
	return filePath, nil
}
