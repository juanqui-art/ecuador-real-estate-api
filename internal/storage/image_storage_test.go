package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLocalImageStorage(t *testing.T) {
	tempDir := t.TempDir()
	
	tests := []struct {
		name     string
		basePath string
		baseURL  string
		maxSize  int64
		wantErr  bool
	}{
		{
			name:     "valid storage creation",
			basePath: tempDir,
			baseURL:  "http://localhost:8080/images",
			maxSize:  10 * 1024 * 1024,
			wantErr:  false,
		},
		{
			name:     "empty base path",
			basePath: "",
			baseURL:  "http://localhost:8080/images",
			maxSize:  10 * 1024 * 1024,
			wantErr:  true,
		},
		{
			name:     "zero max size should use default",
			basePath: tempDir,
			baseURL:  "http://localhost:8080/images",
			maxSize:  0,
			wantErr:  false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage, err := NewLocalImageStorage(tt.basePath, tt.baseURL, tt.maxSize)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, storage)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, storage)
				
				// Check that subdirectories were created
				if tt.basePath != "" {
					subdirs := []string{"originals", "thumbnails", "variants", "temp"}
					for _, subdir := range subdirs {
						path := filepath.Join(tt.basePath, subdir)
						assert.DirExists(t, path)
					}
				}
			}
		})
	}
}

func TestLocalImageStorage_Store(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewLocalImageStorage(tempDir, "http://localhost:8080/images", 1024*1024)
	require.NoError(t, err)
	
	tests := []struct {
		name     string
		data     []byte
		fileName string
		wantErr  bool
	}{
		{
			name:     "valid file storage",
			data:     []byte("test image data"),
			fileName: "test.jpg",
			wantErr:  false,
		},
		{
			name:     "empty data",
			data:     []byte{},
			fileName: "test.jpg",
			wantErr:  true,
		},
		{
			name:     "empty filename",
			data:     []byte("test image data"),
			fileName: "",
			wantErr:  true,
		},
		{
			name:     "data too large",
			data:     make([]byte, 2*1024*1024), // 2MB
			fileName: "large.jpg",
			wantErr:  true,
		},
		{
			name:     "absolute path filename",
			data:     []byte("test image data"),
			fileName: "/etc/passwd",
			wantErr:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := storage.Store(tt.data, tt.fileName)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, path)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, path)
				
				// Verify file exists
				fullPath := filepath.Join(tempDir, path)
				assert.FileExists(t, fullPath)
				
				// Verify file contents
				savedData, err := os.ReadFile(fullPath)
				assert.NoError(t, err)
				assert.Equal(t, tt.data, savedData)
			}
		})
	}
}

func TestLocalImageStorage_Retrieve(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewLocalImageStorage(tempDir, "http://localhost:8080/images", 1024*1024)
	require.NoError(t, err)
	
	// Store a test file first
	testData := []byte("test image data")
	path, err := storage.Store(testData, "test.jpg")
	require.NoError(t, err)
	
	tests := []struct {
		name     string
		filePath string
		wantData []byte
		wantErr  bool
	}{
		{
			name:     "valid file retrieval",
			filePath: path,
			wantData: testData,
			wantErr:  false,
		},
		{
			name:     "empty file path",
			filePath: "",
			wantErr:  true,
		},
		{
			name:     "non-existent file",
			filePath: "originals/nonexistent.jpg",
			wantErr:  true,
		},
		{
			name:     "absolute path",
			filePath: "/etc/passwd",
			wantErr:  true,
		},
		{
			name:     "path outside base directory",
			filePath: "../../../etc/passwd",
			wantErr:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := storage.Retrieve(tt.filePath)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, data)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantData, data)
			}
		})
	}
}

func TestLocalImageStorage_Delete(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewLocalImageStorage(tempDir, "http://localhost:8080/images", 1024*1024)
	require.NoError(t, err)
	
	// Store a test file first
	testData := []byte("test image data")
	path, err := storage.Store(testData, "test.jpg")
	require.NoError(t, err)
	
	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "valid file deletion",
			filePath: path,
			wantErr:  false,
		},
		{
			name:     "empty file path",
			filePath: "",
			wantErr:  true,
		},
		{
			name:     "non-existent file (should not error)",
			filePath: "originals/nonexistent.jpg",
			wantErr:  false,
		},
		{
			name:     "absolute path",
			filePath: "/etc/passwd",
			wantErr:  true,
		},
		{
			name:     "path outside base directory",
			filePath: "../../../etc/passwd",
			wantErr:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.Delete(tt.filePath)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				
				// Verify file is deleted if it was the valid path
				if tt.filePath == path {
					fullPath := filepath.Join(tempDir, tt.filePath)
					assert.NoFileExists(t, fullPath)
				}
			}
		})
	}
}

func TestLocalImageStorage_Exists(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewLocalImageStorage(tempDir, "http://localhost:8080/images", 1024*1024)
	require.NoError(t, err)
	
	// Store a test file first
	testData := []byte("test image data")
	path, err := storage.Store(testData, "test.jpg")
	require.NoError(t, err)
	
	tests := []struct {
		name     string
		filePath string
		want     bool
	}{
		{
			name:     "existing file",
			filePath: path,
			want:     true,
		},
		{
			name:     "non-existent file",
			filePath: "originals/nonexistent.jpg",
			want:     false,
		},
		{
			name:     "empty file path",
			filePath: "",
			want:     false,
		},
		{
			name:     "absolute path",
			filePath: "/etc/passwd",
			want:     false,
		},
		{
			name:     "path outside base directory",
			filePath: "../../../etc/passwd",
			want:     false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists := storage.Exists(tt.filePath)
			assert.Equal(t, tt.want, exists)
		})
	}
}

func TestLocalImageStorage_GetURL(t *testing.T) {
	tempDir := t.TempDir()
	baseURL := "http://localhost:8080/images"
	storage, err := NewLocalImageStorage(tempDir, baseURL, 1024*1024)
	require.NoError(t, err)
	
	tests := []struct {
		name     string
		filePath string
		want     string
	}{
		{
			name:     "valid file path",
			filePath: "originals/test.jpg",
			want:     baseURL + "/originals/test.jpg",
		},
		{
			name:     "empty file path",
			filePath: "",
			want:     "",
		},
		{
			name:     "absolute path",
			filePath: "/etc/passwd",
			want:     "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := storage.GetURL(tt.filePath)
			assert.Equal(t, tt.want, url)
		})
	}
}

func TestLocalImageStorage_GetURL_NoBaseURL(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewLocalImageStorage(tempDir, "", 1024*1024)
	require.NoError(t, err)
	
	filePath := "originals/test.jpg"
	url := storage.GetURL(filePath)
	
	assert.Equal(t, filePath, url)
}

func TestLocalImageStorage_StoreVariant(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewLocalImageStorage(tempDir, "http://localhost:8080/images", 1024*1024)
	require.NoError(t, err)
	
	tests := []struct {
		name     string
		data     []byte
		fileName string
		variant  string
		wantErr  bool
	}{
		{
			name:     "valid thumbnail storage",
			data:     []byte("thumbnail data"),
			fileName: "thumb.jpg",
			variant:  "thumbnails",
			wantErr:  false,
		},
		{
			name:     "valid variant storage",
			data:     []byte("variant data"),
			fileName: "variant.jpg",
			variant:  "variants",
			wantErr:  false,
		},
		{
			name:     "valid temp storage",
			data:     []byte("temp data"),
			fileName: "temp.jpg",
			variant:  "temp",
			wantErr:  false,
		},
		{
			name:     "empty variant",
			data:     []byte("data"),
			fileName: "test.jpg",
			variant:  "",
			wantErr:  true,
		},
		{
			name:     "invalid variant",
			data:     []byte("data"),
			fileName: "test.jpg",
			variant:  "invalid",
			wantErr:  true,
		},
		{
			name:     "empty data",
			data:     []byte{},
			fileName: "test.jpg",
			variant:  "thumbnails",
			wantErr:  true,
		},
		{
			name:     "empty filename",
			data:     []byte("data"),
			fileName: "",
			variant:  "thumbnails",
			wantErr:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := storage.StoreVariant(tt.data, tt.fileName, tt.variant)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, path)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, path)
				
				// Verify file exists
				fullPath := filepath.Join(tempDir, path)
				assert.FileExists(t, fullPath)
				
				// Verify file contents
				savedData, err := os.ReadFile(fullPath)
				assert.NoError(t, err)
				assert.Equal(t, tt.data, savedData)
			}
		})
	}
}

func TestLocalImageStorage_CleanupTempFiles(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewLocalImageStorage(tempDir, "http://localhost:8080/images", 1024*1024)
	require.NoError(t, err)
	
	// Create some temp files
	oldFile := "old_temp.jpg"
	newFile := "new_temp.jpg"
	
	// Store old file
	_, err = storage.StoreVariant([]byte("old data"), oldFile, "temp")
	require.NoError(t, err)
	
	// Modify the file time to be old
	oldPath := filepath.Join(tempDir, "temp", oldFile)
	oldTime := time.Now().Add(-2 * time.Hour)
	err = os.Chtimes(oldPath, oldTime, oldTime)
	require.NoError(t, err)
	
	// Store new file
	_, err = storage.StoreVariant([]byte("new data"), newFile, "temp")
	require.NoError(t, err)
	
	// Cleanup files older than 1 hour
	err = storage.CleanupTempFiles(1 * time.Hour)
	assert.NoError(t, err)
	
	// Verify old file was deleted
	assert.NoFileExists(t, oldPath)
	
	// Verify new file still exists
	newPath := filepath.Join(tempDir, "temp", newFile)
	assert.FileExists(t, newPath)
}

func TestLocalImageStorage_GetStorageInfo(t *testing.T) {
	tempDir := t.TempDir()
	baseURL := "http://localhost:8080/images"
	maxSize := int64(1024 * 1024)
	storage, err := NewLocalImageStorage(tempDir, baseURL, maxSize)
	require.NoError(t, err)
	
	// Store some test files
	testData := []byte("test data")
	storage.Store(testData, "test1.jpg")
	storage.Store(testData, "test2.jpg")
	
	info := storage.GetStorageInfo()
	
	assert.Equal(t, "local", info.Type)
	assert.Equal(t, tempDir, info.BasePath)
	assert.Equal(t, baseURL, info.BaseURL)
	assert.Equal(t, maxSize, info.MaxFileSize)
	assert.Greater(t, info.TotalSize, int64(0))
	assert.Greater(t, info.FileCount, 0)
	assert.False(t, info.LastModified.IsZero())
}

func TestLocalImageStorage_CopyFile(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewLocalImageStorage(tempDir, "http://localhost:8080/images", 1024*1024)
	require.NoError(t, err)
	
	// Store a test file first
	testData := []byte("test image data")
	srcPath, err := storage.Store(testData, "source.jpg")
	require.NoError(t, err)
	
	tests := []struct {
		name    string
		srcPath string
		dstPath string
		wantErr bool
	}{
		{
			name:    "valid file copy",
			srcPath: srcPath,
			dstPath: "originals/copy.jpg",
			wantErr: false,
		},
		{
			name:    "empty source path",
			srcPath: "",
			dstPath: "originals/copy.jpg",
			wantErr: true,
		},
		{
			name:    "empty destination path",
			srcPath: srcPath,
			dstPath: "",
			wantErr: true,
		},
		{
			name:    "non-existent source file",
			srcPath: "originals/nonexistent.jpg",
			dstPath: "originals/copy.jpg",
			wantErr: true,
		},
		{
			name:    "absolute source path",
			srcPath: "/etc/passwd",
			dstPath: "originals/copy.jpg",
			wantErr: true,
		},
		{
			name:    "absolute destination path",
			srcPath: srcPath,
			dstPath: "/tmp/copy.jpg",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.CopyFile(tt.srcPath, tt.dstPath)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				
				// Verify destination file exists
				dstFullPath := filepath.Join(tempDir, tt.dstPath)
				assert.FileExists(t, dstFullPath)
				
				// Verify file contents match
				copiedData, err := os.ReadFile(dstFullPath)
				assert.NoError(t, err)
				assert.Equal(t, testData, copiedData)
			}
		})
	}
}

func TestLocalImageStorage_isPathWithinBase(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewLocalImageStorage(tempDir, "http://localhost:8080/images", 1024*1024)
	require.NoError(t, err)
	
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "path within base",
			path: filepath.Join(tempDir, "originals", "test.jpg"),
			want: true,
		},
		{
			name: "path outside base",
			path: "/etc/passwd",
			want: false,
		},
		{
			name: "path traversal attempt",
			path: filepath.Join(tempDir, "..", "..", "etc", "passwd"),
			want: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := storage.isPathWithinBase(tt.path)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestLocalImageStorage_calculateStorageStats(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewLocalImageStorage(tempDir, "http://localhost:8080/images", 1024*1024)
	require.NoError(t, err)
	
	// Store some test files
	testData1 := []byte("test data 1")
	testData2 := []byte("test data 2 - longer")
	
	storage.Store(testData1, "test1.jpg")
	storage.Store(testData2, "test2.jpg")
	
	totalSize, fileCount, lastModified := storage.calculateStorageStats()
	
	assert.Greater(t, totalSize, int64(0))
	assert.Greater(t, fileCount, 0)
	assert.False(t, lastModified.IsZero())
	
	// Should account for all files
	expectedMinSize := int64(len(testData1) + len(testData2))
	assert.GreaterOrEqual(t, totalSize, expectedMinSize)
}