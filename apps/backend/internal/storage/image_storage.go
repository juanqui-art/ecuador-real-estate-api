package storage

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"realty-core/internal/domain"
)

// ImageStorage defines the interface for image storage operations
type ImageStorage interface {
	// Store saves image data and returns the storage path
	Store(data []byte, fileName string) (string, error)
	
	// Retrieve gets image data from storage
	Retrieve(filePath string) ([]byte, error)
	
	// Delete removes image from storage
	Delete(filePath string) error
	
	// Exists checks if image exists in storage
	Exists(filePath string) bool
	
	// GetURL returns the public URL for the image
	GetURL(filePath string) string
	
	// GetStorageInfo returns storage information
	GetStorageInfo() StorageInfo
}

// StorageInfo contains storage backend information
type StorageInfo struct {
	Type         string `json:"type"`
	BasePath     string `json:"base_path"`
	BaseURL      string `json:"base_url"`
	MaxFileSize  int64  `json:"max_file_size"`
	TotalSize    int64  `json:"total_size"`
	FileCount    int    `json:"file_count"`
	LastModified time.Time `json:"last_modified"`
}

// LocalImageStorage implements ImageStorage for local filesystem
type LocalImageStorage struct {
	basePath string
	baseURL  string
	maxSize  int64
}

// NewLocalImageStorage creates a new local image storage
func NewLocalImageStorage(basePath, baseURL string, maxSize int64) (*LocalImageStorage, error) {
	if basePath == "" {
		return nil, fmt.Errorf("base path cannot be empty")
	}
	
	if maxSize <= 0 {
		maxSize = domain.MaxUploadSize
	}
	
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}
	
	// Create subdirectories for organization
	subdirs := []string{"originals", "thumbnails", "variants", "temp"}
	for _, subdir := range subdirs {
		subPath := filepath.Join(basePath, subdir)
		if err := os.MkdirAll(subPath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create subdirectory %s: %w", subdir, err)
		}
	}
	
	return &LocalImageStorage{
		basePath: basePath,
		baseURL:  baseURL,
		maxSize:  maxSize,
	}, nil
}

// Store saves image data to local filesystem
func (ls *LocalImageStorage) Store(data []byte, fileName string) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("empty image data")
	}
	
	if int64(len(data)) > ls.maxSize {
		return "", fmt.Errorf("image too large: %d bytes, max: %d bytes", len(data), ls.maxSize)
	}
	
	if fileName == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}
	
	// Clean filename
	fileName = filepath.Clean(fileName)
	if filepath.IsAbs(fileName) {
		return "", fmt.Errorf("filename cannot be absolute path")
	}
	
	// Create full path
	fullPath := filepath.Join(ls.basePath, "originals", fileName)
	
	// Check if file already exists
	if _, err := os.Stat(fullPath); err == nil {
		return "", fmt.Errorf("file already exists: %s", fileName)
	}
	
	// Create directory if needed
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Write file
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	
	if _, err := file.Write(data); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	
	// Return relative path
	relPath := filepath.Join("originals", fileName)
	return relPath, nil
}

// Retrieve gets image data from local filesystem
func (ls *LocalImageStorage) Retrieve(filePath string) ([]byte, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}
	
	// Clean and validate path
	filePath = filepath.Clean(filePath)
	if filepath.IsAbs(filePath) {
		return nil, fmt.Errorf("file path cannot be absolute")
	}
	
	fullPath := filepath.Join(ls.basePath, filePath)
	
	// Security check - ensure path is within base directory
	if !ls.isPathWithinBase(fullPath) {
		return nil, fmt.Errorf("path outside base directory: %s", filePath)
	}
	
	// Read file
	data, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", filePath)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	return data, nil
}

// Delete removes image from local filesystem
func (ls *LocalImageStorage) Delete(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}
	
	// Clean and validate path
	filePath = filepath.Clean(filePath)
	if filepath.IsAbs(filePath) {
		return fmt.Errorf("file path cannot be absolute")
	}
	
	fullPath := filepath.Join(ls.basePath, filePath)
	
	// Security check
	if !ls.isPathWithinBase(fullPath) {
		return fmt.Errorf("path outside base directory: %s", filePath)
	}
	
	// Delete file
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // Already deleted
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}
	
	return nil
}

// Exists checks if image exists in local filesystem
func (ls *LocalImageStorage) Exists(filePath string) bool {
	if filePath == "" {
		return false
	}
	
	filePath = filepath.Clean(filePath)
	if filepath.IsAbs(filePath) {
		return false
	}
	
	fullPath := filepath.Join(ls.basePath, filePath)
	
	if !ls.isPathWithinBase(fullPath) {
		return false
	}
	
	_, err := os.Stat(fullPath)
	return err == nil
}

// GetURL returns the public URL for the image
func (ls *LocalImageStorage) GetURL(filePath string) string {
	if filePath == "" {
		return ""
	}
	
	// Clean path
	filePath = filepath.Clean(filePath)
	if filepath.IsAbs(filePath) {
		return ""
	}
	
	// Convert to URL path (use forward slashes)
	urlPath := filepath.ToSlash(filePath)
	
	if ls.baseURL == "" {
		return urlPath
	}
	
	return ls.baseURL + "/" + urlPath
}

// GetStorageInfo returns storage information
func (ls *LocalImageStorage) GetStorageInfo() StorageInfo {
	info := StorageInfo{
		Type:     "local",
		BasePath: ls.basePath,
		BaseURL:  ls.baseURL,
		MaxFileSize: ls.maxSize,
	}
	
	// Calculate storage statistics
	totalSize, fileCount, lastModified := ls.calculateStorageStats()
	info.TotalSize = totalSize
	info.FileCount = fileCount
	info.LastModified = lastModified
	
	return info
}

// StoreVariant stores an image variant (thumbnail, medium, etc.)
func (ls *LocalImageStorage) StoreVariant(data []byte, fileName string, variant string) (string, error) {
	if variant == "" {
		return "", fmt.Errorf("variant cannot be empty")
	}
	
	// Validate variant
	validVariants := []string{"thumbnails", "variants", "temp"}
	isValid := false
	for _, v := range validVariants {
		if v == variant {
			isValid = true
			break
		}
	}
	
	if !isValid {
		return "", fmt.Errorf("invalid variant: %s", variant)
	}
	
	if len(data) == 0 {
		return "", fmt.Errorf("empty image data")
	}
	
	if fileName == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}
	
	// Clean filename
	fileName = filepath.Clean(fileName)
	if filepath.IsAbs(fileName) {
		return "", fmt.Errorf("filename cannot be absolute path")
	}
	
	// Create full path
	fullPath := filepath.Join(ls.basePath, variant, fileName)
	
	// Create directory if needed
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Write file
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	
	if _, err := file.Write(data); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	
	// Return relative path
	relPath := filepath.Join(variant, fileName)
	return relPath, nil
}

// CleanupTempFiles removes temporary files older than specified duration
func (ls *LocalImageStorage) CleanupTempFiles(olderThan time.Duration) error {
	tempDir := filepath.Join(ls.basePath, "temp")
	cutoff := time.Now().Add(-olderThan)
	
	return filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			
			if info.ModTime().Before(cutoff) {
				if err := os.Remove(path); err != nil {
					return fmt.Errorf("failed to remove temp file %s: %w", path, err)
				}
			}
		}
		
		return nil
	})
}

// isPathWithinBase checks if path is within base directory (security check)
func (ls *LocalImageStorage) isPathWithinBase(path string) bool {
	absBase, err := filepath.Abs(ls.basePath)
	if err != nil {
		return false
	}
	
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	
	rel, err := filepath.Rel(absBase, absPath)
	if err != nil {
		return false
	}
	
	return !filepath.IsAbs(rel) && !filepath.HasPrefix(rel, "..")
}

// calculateStorageStats calculates storage statistics
func (ls *LocalImageStorage) calculateStorageStats() (int64, int, time.Time) {
	var totalSize int64
	var fileCount int
	var lastModified time.Time
	
	filepath.WalkDir(ls.basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			
			totalSize += info.Size()
			fileCount++
			
			if info.ModTime().After(lastModified) {
				lastModified = info.ModTime()
			}
		}
		
		return nil
	})
	
	return totalSize, fileCount, lastModified
}

// CopyFile copies a file from source to destination within storage
func (ls *LocalImageStorage) CopyFile(srcPath, dstPath string) error {
	if srcPath == "" || dstPath == "" {
		return fmt.Errorf("source and destination paths cannot be empty")
	}
	
	srcPath = filepath.Clean(srcPath)
	dstPath = filepath.Clean(dstPath)
	
	if filepath.IsAbs(srcPath) || filepath.IsAbs(dstPath) {
		return fmt.Errorf("paths cannot be absolute")
	}
	
	srcFullPath := filepath.Join(ls.basePath, srcPath)
	dstFullPath := filepath.Join(ls.basePath, dstPath)
	
	if !ls.isPathWithinBase(srcFullPath) || !ls.isPathWithinBase(dstFullPath) {
		return fmt.Errorf("paths outside base directory")
	}
	
	// Open source file
	srcFile, err := os.Open(srcFullPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()
	
	// Create destination directory if needed
	if err := os.MkdirAll(filepath.Dir(dstFullPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	
	// Create destination file
	dstFile, err := os.Create(dstFullPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()
	
	// Copy data
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	
	return nil
}