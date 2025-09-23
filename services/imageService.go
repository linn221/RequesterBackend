package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/linn221/RequesterBackend/models"
	"gorm.io/gorm"
)

type ImageService struct {
	DB              *gorm.DB
	UploadDirectory string
	ValidTypes      []string
}

// UploadImages handles multiple file uploads and creates image records
func (s *ImageService) UploadImages(files []*multipart.FileHeader, referenceType string, referenceId int) ([]*models.Image, error) {
	// Validate reference type
	if !s.isValidReferenceType(referenceType) {
		return nil, fmt.Errorf("invalid reference type: %s", referenceType)
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(s.UploadDirectory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %v", err)
	}

	var images []*models.Image
	var uploadedFiles []string // Track uploaded files for cleanup on error

	// Process each file
	for _, file := range files {
		// Generate unique filename
		originalName := file.Filename
		ext := filepath.Ext(originalName)
		uniqueFilename := s.generateUniqueFilename(ext)

		// Save file to disk
		filePath := filepath.Join(s.UploadDirectory, uniqueFilename)
		if err := s.saveFile(file, filePath); err != nil {
			// Clean up any already uploaded files
			s.cleanupFiles(uploadedFiles)
			return nil, fmt.Errorf("failed to save file %s: %v", originalName, err)
		}
		uploadedFiles = append(uploadedFiles, filePath)

		// Create image record
		image := &models.Image{
			Filename:      uniqueFilename,
			OriginalName:  originalName,
			FilePath:      filePath,
			FileSize:      file.Size,
			MimeType:      file.Header.Get("Content-Type"),
			ReferenceType: referenceType,
			ReferenceID:   referenceId,
		}

		images = append(images, image)
	}

	// Save all images to database
	if err := s.DB.Create(&images).Error; err != nil {
		// Clean up uploaded files if database save fails
		s.cleanupFiles(uploadedFiles)
		return nil, fmt.Errorf("failed to save images to database: %v", err)
	}

	return images, nil
}

// DeleteImage deletes image record and file
func (s *ImageService) DeleteImage(imageId int) error {
	// Find image
	var image models.Image
	if err := s.DB.First(&image, imageId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("image not found")
		}
		return fmt.Errorf("failed to find image: %v", err)
	}

	// Delete file from disk
	if err := os.Remove(image.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	// Delete from database
	if err := s.DB.Delete(&image).Error; err != nil {
		return fmt.Errorf("failed to delete image from database: %v", err)
	}

	return nil
}

// GetImage retrieves image by Id
func (s *ImageService) GetImage(imageId int) (*models.Image, error) {
	var image models.Image
	if err := s.DB.First(&image, imageId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("image not found")
		}
		return nil, fmt.Errorf("failed to get image: %v", err)
	}
	return &image, nil
}

// GetImageByFilename retrieves image by filename
func (s *ImageService) GetImageByFilename(filename string) (*models.Image, error) {
	var image models.Image
	if err := s.DB.Where("filename = ?", filename).First(&image).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("image not found")
		}
		return nil, fmt.Errorf("failed to get image: %v", err)
	}
	return &image, nil
}

// ListImagesByReference retrieves images by reference type and Id
func (s *ImageService) ListImagesByReference(referenceType string, referenceId int) ([]models.Image, error) {
	var images []models.Image
	query := s.DB.Where("reference_type = ? AND reference_id = ?", referenceType, referenceId)

	if err := query.Find(&images).Error; err != nil {
		return nil, fmt.Errorf("failed to list images: %v", err)
	}

	return images, nil
}

// Helper methods

func (s *ImageService) isValidReferenceType(referenceType string) bool {
	for _, validType := range s.ValidTypes {
		if referenceType == validType {
			return true
		}
	}
	return false
}

func (s *ImageService) generateUniqueFilename(ext string) string {
	// Generate random bytes
	bytes := make([]byte, 16)
	rand.Read(bytes)
	randomStr := hex.EncodeToString(bytes)
	return randomStr + ext
}

func (s *ImageService) saveFile(file *multipart.FileHeader, filePath string) error {
	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, src)
	return err
}

func (s *ImageService) cleanupFiles(filePaths []string) {
	for _, filePath := range filePaths {
		os.Remove(filePath)
	}
}

// ServeImage serves the image file
func (s *ImageService) ServeImage(w http.ResponseWriter, r *http.Request, image *models.Image) error {
	// Check if file exists
	if _, err := os.Stat(image.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found")
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", image.MimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", image.OriginalName))

	// Serve file
	http.ServeFile(w, r, image.FilePath)
	return nil
}
