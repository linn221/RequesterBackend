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
	"strings"

	"github.com/linn221/RequesterBackend/models"
	"gorm.io/gorm"
)

type AttachmentService struct {
	DB              *gorm.DB
	UploadDirectory string
	MaxFileSize     int64
	AllowedTypes    []string
	ValidTypes      []string
}

// UploadAttachment handles file upload and creates attachment record
func (s *AttachmentService) UploadAttachment(file *multipart.FileHeader, referenceType string, referenceId int) (*models.Attachment, error) {
	// File size validation removed - no limit on attachment file upload size

	// Validate file type
	if !s.isAllowedType(file.Header.Get("Content-Type")) {
		return nil, fmt.Errorf("file type %s is not allowed", file.Header.Get("Content-Type"))
	}

	// Validate reference type
	if !s.isValidReferenceType(referenceType) {
		return nil, fmt.Errorf("invalid reference type: %s", referenceType)
	}

	// Generate unique filename
	originalName := file.Filename
	ext := filepath.Ext(originalName)
	uniqueFilename := s.generateUniqueFilename(ext)

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(s.UploadDirectory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Save file to disk
	filePath := filepath.Join(s.UploadDirectory, uniqueFilename)
	if err := s.saveFile(file, filePath); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// Create attachment record
	attachment := &models.Attachment{
		Filename:      uniqueFilename,
		OriginalName:  originalName,
		FilePath:      filePath,
		FileSize:      file.Size,
		MimeType:      file.Header.Get("Content-Type"),
		ReferenceType: referenceType,
		ReferenceID:   referenceId,
	}

	// Save to database
	if err := s.DB.Create(attachment).Error; err != nil {
		// Clean up file if database save fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save attachment to database: %v", err)
	}

	return attachment, nil
}

// DeleteAttachment deletes attachment record and file
func (s *AttachmentService) DeleteAttachment(attachmentId int) error {
	// Find attachment
	var attachment models.Attachment
	if err := s.DB.First(&attachment, attachmentId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("attachment not found")
		}
		return fmt.Errorf("failed to find attachment: %v", err)
	}

	// Delete file from disk
	if err := os.Remove(attachment.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	// Delete from database
	if err := s.DB.Delete(&attachment).Error; err != nil {
		return fmt.Errorf("failed to delete attachment from database: %v", err)
	}

	return nil
}

// GetAttachment retrieves attachment by Id
func (s *AttachmentService) GetAttachment(attachmentId int) (*models.Attachment, error) {
	var attachment models.Attachment
	if err := s.DB.First(&attachment, attachmentId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("attachment not found")
		}
		return nil, fmt.Errorf("failed to get attachment: %v", err)
	}
	return &attachment, nil
}

// ListAttachmentsByReference retrieves attachments by reference type and Id
func (s *AttachmentService) ListAttachmentsByReference(referenceType string, referenceId int) ([]models.Attachment, error) {
	var attachments []models.Attachment
	query := s.DB.Where("reference_type = ? AND reference_id = ?", referenceType, referenceId)

	if err := query.Find(&attachments).Error; err != nil {
		return nil, fmt.Errorf("failed to list attachments: %v", err)
	}

	return attachments, nil
}

// Helper methods

func (s *AttachmentService) isAllowedType(mimeType string) bool {
	for _, allowedType := range s.AllowedTypes {
		if strings.HasPrefix(mimeType, allowedType) {
			return true
		}
	}
	return false
}

func (s *AttachmentService) isValidReferenceType(referenceType string) bool {
	for _, validType := range s.ValidTypes {
		if referenceType == validType {
			return true
		}
	}
	return false
}

func (s *AttachmentService) generateUniqueFilename(ext string) string {
	// Generate random bytes
	bytes := make([]byte, 16)
	rand.Read(bytes)
	randomStr := hex.EncodeToString(bytes)
	return randomStr + ext
}

func (s *AttachmentService) saveFile(file *multipart.FileHeader, filePath string) error {
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

// ServeFile serves the attachment file
func (s *AttachmentService) ServeFile(w http.ResponseWriter, r *http.Request, attachment *models.Attachment) error {
	// Check if file exists
	if _, err := os.Stat(attachment.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found")
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", attachment.MimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", attachment.OriginalName))

	// Serve file
	http.ServeFile(w, r, attachment.FilePath)
	return nil
}
