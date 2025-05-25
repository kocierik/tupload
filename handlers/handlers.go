package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"path/filepath"
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
	var file io.Reader
	var filename string
	var filesize int64

	// Handle PUT request with raw body
	if c.Request.Method == "PUT" {
		file = c.Request.Body
		filename = filepath.Base(c.Request.URL.Path)
		if filename == "" || filename == "/" {
			filename = "uploaded_file"
		}
		filesize = c.Request.ContentLength
	} else {
		// Try to get the file from the form
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
			return
		}

		// Look for files in any field
		var foundFile bool
		for _, files := range form.File {
			if len(files) > 0 {
				formFile, err := files[0].Open()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
					return
				}
				file = formFile
				filename = files[0].Filename
				filesize = files[0].Size
				foundFile = true
				break
			}
		}

		if !foundFile {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
			return
		}
	}

	fileID := generateSimpleID()
	_, err := h.storage.SaveFileWithID(file, filename, fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	downloadURL := fmt.Sprintf("https://%s/download/%s", h.domain, fileID)
	output := fmt.Sprintf("\n=========================\n\nUploaded 1 file, %d bytes\n\nwget %s\n\n=========================",
		filesize, downloadURL)

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
