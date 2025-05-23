package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/erik/tupload/storage"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage *storage.FileStorage
	domain  string
}

func NewHandler(storage *storage.FileStorage, domain string) *Handler {
	return &Handler{
		storage: storage,
		domain:  domain,
	}
}

func generateSimpleID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 5
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (h *Handler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	fileID := generateSimpleID()
	_, err = h.storage.SaveFileWithID(file, header.Filename, fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	downloadURL := fmt.Sprintf("https://%s/download/%s", h.domain, fileID)
	output := fmt.Sprintf("\n=========================\n\nUploaded 1 file, %d bytes\n\nwget %s\n\n=========================",
		header.Size, downloadURL)

	c.String(http.StatusOK, output)
}

func (h *Handler) DownloadFile(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file ID provided"})
		return
	}

	filePath, err := h.storage.GetFilePath(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(filePath)
}
