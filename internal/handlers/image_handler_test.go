package handlers

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"realty-core/internal/cache"
	"realty-core/internal/domain"
)

// MockImageService es un mock del servicio de imágenes
type MockImageService struct {
	mock.Mock
}

func (m *MockImageService) Upload(propertyID string, file multipart.File, header *multipart.FileHeader, altText string) (*domain.ImageInfo, error) {
	args := m.Called(propertyID, file, header, altText)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ImageInfo), args.Error(1)
}

func (m *MockImageService) GetImage(id string) (*domain.ImageInfo, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ImageInfo), args.Error(1)
}

func (m *MockImageService) GetImagesByProperty(propertyID string) ([]domain.ImageInfo, error) {
	args := m.Called(propertyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ImageInfo), args.Error(1)
}

func (m *MockImageService) UpdateImageMetadata(id string, altText string, sortOrder int) error {
	args := m.Called(id, altText, sortOrder)
	return args.Error(0)
}

func (m *MockImageService) DeleteImage(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockImageService) ReorderImages(propertyID string, imageIDs []string) error {
	args := m.Called(propertyID, imageIDs)
	return args.Error(0)
}

func (m *MockImageService) SetMainImage(propertyID string, imageID string) error {
	args := m.Called(propertyID, imageID)
	return args.Error(0)
}

func (m *MockImageService) GetMainImage(propertyID string) (*domain.ImageInfo, error) {
	args := m.Called(propertyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ImageInfo), args.Error(1)
}

func (m *MockImageService) GetImageVariant(imageID string, width, height int, format string, quality int) ([]byte, error) {
	args := m.Called(imageID, width, height, format, quality)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockImageService) GenerateThumbnail(imageID string, size int) ([]byte, error) {
	args := m.Called(imageID, size)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockImageService) GetImageStats() (map[string]interface{}, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockImageService) ValidateUpload(header *multipart.FileHeader) error {
	args := m.Called(header)
	return args.Error(0)
}

func (m *MockImageService) CleanupTempFiles(olderThan time.Duration) error {
	args := m.Called(olderThan)
	return args.Error(0)
}

func (m *MockImageService) GetCacheStats() cache.ImageCacheStats {
	args := m.Called()
	return args.Get(0).(cache.ImageCacheStats)
}

func TestNewImageHandler(t *testing.T) {
	mockService := &MockImageService{}
	handler := NewImageHandler(mockService)
	
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.imageService)
}

func TestImageHandler_UploadImage(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		setupRequest   func() *http.Request
		mockSetup      func(*MockImageService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "successful upload",
			method: http.MethodPost,
			setupRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.WriteField("property_id", "test-property-id")
				writer.WriteField("alt_text", "Test image")
				
				fileWriter, _ := writer.CreateFormFile("image", "test.jpg")
				fileWriter.Write([]byte("fake-image-data"))
				writer.Close()
				
				req := httptest.NewRequest(http.MethodPost, "/api/images", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			mockSetup: func(m *MockImageService) {
				expectedImage := &domain.ImageInfo{
					ID:         "test-id",
					PropertyID: "test-property-id",
					FileName:   "test.jpg",
					AltText:    "Test image",
					CreatedAt:  time.Now(),
				}
				m.On("Upload", "test-property-id", mock.Anything, mock.Anything, "Test image").Return(expectedImage, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Image uploaded successfully",
		},
		{
			name:   "method not allowed",
			method: http.MethodGet,
			setupRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/api/images", nil)
			},
			mockSetup:      func(m *MockImageService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed",
		},
		{
			name:   "missing property ID",
			method: http.MethodPost,
			setupRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.WriteField("alt_text", "Test image")
				writer.Close()
				
				req := httptest.NewRequest(http.MethodPost, "/api/images", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			mockSetup:      func(m *MockImageService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Property ID is required",
		},
		{
			name:   "service error",
			method: http.MethodPost,
			setupRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.WriteField("property_id", "test-property-id")
				
				fileWriter, _ := writer.CreateFormFile("image", "test.jpg")
				fileWriter.Write([]byte("fake-image-data"))
				writer.Close()
				
				req := httptest.NewRequest(http.MethodPost, "/api/images", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			mockSetup: func(m *MockImageService) {
				m.On("Upload", "test-property-id", mock.Anything, mock.Anything, "").Return(nil, fmt.Errorf("upload failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Failed to upload image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockImageService{}
			handler := NewImageHandler(mockService)
			
			tt.mockSetup(mockService)
			
			req := tt.setupRequest()
			rr := httptest.NewRecorder()
			
			handler.UploadImage(rr, req)
			
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestImageHandler_GetImage(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		mockSetup      func(*MockImageService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful get image",
			path: "/api/images/test-id",
			mockSetup: func(m *MockImageService) {
				expectedImage := &domain.ImageInfo{
					ID:         "test-id",
					PropertyID: "test-property-id",
					FileName:   "test.jpg",
					AltText:    "Test image",
					CreatedAt:  time.Now(),
				}
				m.On("GetImage", "test-id").Return(expectedImage, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "test-id",
		},
		{
			name: "image not found",
			path: "/api/images/nonexistent-id",
			mockSetup: func(m *MockImageService) {
				m.On("GetImage", "nonexistent-id").Return(nil, fmt.Errorf("image not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Image not found",
		},
		{
			name: "invalid path",
			path: "/api/images/",
			mockSetup: func(m *MockImageService) {
				// No mock needed for invalid path
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Image ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockImageService{}
			handler := NewImageHandler(mockService)
			
			tt.mockSetup(mockService)
			
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rr := httptest.NewRecorder()
			
			handler.GetImage(rr, req)
			
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestImageHandler_GetImagesByProperty(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		mockSetup      func(*MockImageService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful get images by property",
			path: "/api/properties/test-property-id/images",
			mockSetup: func(m *MockImageService) {
				expectedImages := []domain.ImageInfo{
					{
						ID:         "img1",
						PropertyID: "test-property-id",
						FileName:   "test1.jpg",
						AltText:    "Test image 1",
						CreatedAt:  time.Now(),
					},
					{
						ID:         "img2",
						PropertyID: "test-property-id",
						FileName:   "test2.jpg",
						AltText:    "Test image 2",
						CreatedAt:  time.Now(),
					},
				}
				m.On("GetImagesByProperty", "test-property-id").Return(expectedImages, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "img1",
		},
		{
			name: "property not found",
			path: "/api/properties/nonexistent-id/images",
			mockSetup: func(m *MockImageService) {
				m.On("GetImagesByProperty", "nonexistent-id").Return(nil, fmt.Errorf("property not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to get images",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockImageService{}
			handler := NewImageHandler(mockService)
			
			tt.mockSetup(mockService)
			
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rr := httptest.NewRecorder()
			
			handler.GetImagesByProperty(rr, req)
			
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestImageHandler_UpdateImageMetadata(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		body           string
		mockSetup      func(*MockImageService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful update",
			path: "/api/images/test-id/metadata",
			body: `{"alt_text": "Updated alt text", "sort_order": 1}`,
			mockSetup: func(m *MockImageService) {
				m.On("UpdateImageMetadata", "test-id", "Updated alt text", 1).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Image metadata updated successfully",
		},
		{
			name: "image not found",
			path: "/api/images/nonexistent-id/metadata",
			body: `{"alt_text": "Updated alt text", "sort_order": 0}`,
			mockSetup: func(m *MockImageService) {
				m.On("UpdateImageMetadata", "nonexistent-id", "Updated alt text", 0).Return(fmt.Errorf("image not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Image not found",
		},
		{
			name: "invalid JSON",
			path: "/api/images/test-id/metadata",
			body: `{"invalid": json}`,
			mockSetup: func(m *MockImageService) {
				// No mock needed for invalid JSON
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockImageService{}
			handler := NewImageHandler(mockService)
			
			tt.mockSetup(mockService)
			
			req := httptest.NewRequest(http.MethodPut, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			
			handler.UpdateImageMetadata(rr, req)
			
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestImageHandler_DeleteImage(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		mockSetup      func(*MockImageService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful delete",
			path: "/api/images/test-id",
			mockSetup: func(m *MockImageService) {
				m.On("DeleteImage", "test-id").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Image deleted successfully",
		},
		{
			name: "image not found",
			path: "/api/images/nonexistent-id",
			mockSetup: func(m *MockImageService) {
				m.On("DeleteImage", "nonexistent-id").Return(fmt.Errorf("image not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Image not found",
		},
		{
			name: "invalid path",
			path: "/api/images/",
			mockSetup: func(m *MockImageService) {
				// No mock needed for invalid path
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Image ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockImageService{}
			handler := NewImageHandler(mockService)
			
			tt.mockSetup(mockService)
			
			req := httptest.NewRequest(http.MethodDelete, tt.path, nil)
			rr := httptest.NewRecorder()
			
			handler.DeleteImage(rr, req)
			
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestImageHandler_GetImageStats(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockImageService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful get stats",
			mockSetup: func(m *MockImageService) {
				expectedStats := map[string]interface{}{
					"total_images":     100,
					"total_size":       1024000,
					"average_size":     10240,
					"most_used_format": "jpg",
				}
				m.On("GetImageStats").Return(expectedStats, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "100",
		},
		{
			name: "service error",
			mockSetup: func(m *MockImageService) {
				m.On("GetImageStats").Return(nil, fmt.Errorf("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to get image stats",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockImageService{}
			handler := NewImageHandler(mockService)
			
			tt.mockSetup(mockService)
			
			req := httptest.NewRequest(http.MethodGet, "/api/images/stats", nil)
			rr := httptest.NewRecorder()
			
			handler.GetImageStats(rr, req)
			
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestImageHandler_ReorderImages(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		body           string
		mockSetup      func(*MockImageService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful reorder",
			path: "/api/properties/test-property-id/images/reorder",
			body: `{"image_ids": ["img1", "img2", "img3"]}`,
			mockSetup: func(m *MockImageService) {
				m.On("ReorderImages", "test-property-id", []string{"img1", "img2", "img3"}).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Images reordered successfully",
		},
		{
			name: "invalid JSON",
			path: "/api/properties/test-property-id/images/reorder",
			body: `{"invalid": json}`,
			mockSetup: func(m *MockImageService) {
				// No mock needed for invalid JSON
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body",
		},
		{
			name: "service error",
			path: "/api/properties/test-property-id/images/reorder",
			body: `{"image_ids": ["img1", "img2"]}`,
			mockSetup: func(m *MockImageService) {
				m.On("ReorderImages", "test-property-id", []string{"img1", "img2"}).Return(fmt.Errorf("reorder failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Failed to reorder images",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockImageService{}
			handler := NewImageHandler(mockService)
			
			tt.mockSetup(mockService)
			
			req := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			
			handler.ReorderImages(rr, req)
			
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestImageHandler_SetMainImage(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		body           string
		mockSetup      func(*MockImageService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful set main image",
			path: "/api/properties/test-property-id/images/main",
			body: `{"image_id": "img1"}`,
			mockSetup: func(m *MockImageService) {
				m.On("SetMainImage", "test-property-id", "img1").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Main image set successfully",
		},
		{
			name: "invalid JSON",
			path: "/api/properties/test-property-id/images/main",
			body: `{"invalid": json}`,
			mockSetup: func(m *MockImageService) {
				// No mock needed for invalid JSON
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body",
		},
		{
			name: "service error",
			path: "/api/properties/test-property-id/images/main",
			body: `{"image_id": "img1"}`,
			mockSetup: func(m *MockImageService) {
				m.On("SetMainImage", "test-property-id", "img1").Return(fmt.Errorf("set main image failed"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Failed to set main image",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockImageService{}
			handler := NewImageHandler(mockService)
			
			tt.mockSetup(mockService)
			
			req := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			
			handler.SetMainImage(rr, req)
			
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestImageHandler_GetMainImage(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		mockSetup      func(*MockImageService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful get main image",
			path: "/api/properties/test-property-id/images/main",
			mockSetup: func(m *MockImageService) {
				expectedImage := &domain.ImageInfo{
					ID:         "main-img",
					PropertyID: "test-property-id",
					FileName:   "main.jpg",
					AltText:    "Main image",
					CreatedAt:  time.Now(),
				}
				m.On("GetMainImage", "test-property-id").Return(expectedImage, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "main-img",
		},
		{
			name: "main image not found",
			path: "/api/properties/test-property-id/images/main",
			mockSetup: func(m *MockImageService) {
				m.On("GetMainImage", "test-property-id").Return(nil, fmt.Errorf("main image not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "No images found for property",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockImageService{}
			handler := NewImageHandler(mockService)
			
			tt.mockSetup(mockService)
			
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rr := httptest.NewRecorder()
			
			handler.GetMainImage(rr, req)
			
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestImageHandler_GetImageVariant(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		mockSetup      func(*MockImageService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful get variant",
			path: "/api/images/test-id/variant?w=200&h=200&f=jpg&q=80",
			mockSetup: func(m *MockImageService) {
				expectedData := []byte("fake-variant-data")
				m.On("GetImageVariant", "test-id", 200, 200, "jpg", 80).Return(expectedData, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "fake-variant-data",
		},
		{
			name: "variant not found",
			path: "/api/images/nonexistent-id/variant?w=200&h=200",
			mockSetup: func(m *MockImageService) {
				m.On("GetImageVariant", "nonexistent-id", 200, 200, "jpg", 85).Return(nil, fmt.Errorf("variant not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Image not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockImageService{}
			handler := NewImageHandler(mockService)
			
			tt.mockSetup(mockService)
			
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rr := httptest.NewRecorder()
			
			handler.GetImageVariant(rr, req)
			
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestImageHandler_GetThumbnail(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		mockSetup      func(*MockImageService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful get thumbnail",
			path: "/api/images/test-id/thumbnail?size=150",
			mockSetup: func(m *MockImageService) {
				expectedData := []byte("fake-thumbnail-data")
				m.On("GenerateThumbnail", "test-id", 150).Return(expectedData, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "fake-thumbnail-data",
		},
		{
			name: "thumbnail not found",
			path: "/api/images/nonexistent-id/thumbnail",
			mockSetup: func(m *MockImageService) {
				m.On("GenerateThumbnail", "nonexistent-id", 150).Return(nil, fmt.Errorf("thumbnail not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Image not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockImageService{}
			handler := NewImageHandler(mockService)
			
			tt.mockSetup(mockService)
			
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rr := httptest.NewRecorder()
			
			handler.GetThumbnail(rr, req)
			
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestImageHandler_CleanupTempFiles(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockImageService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful cleanup",
			mockSetup: func(m *MockImageService) {
				m.On("CleanupTempFiles", mock.AnythingOfType("time.Duration")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Temporary files cleaned up successfully",
		},
		{
			name: "service error",
			mockSetup: func(m *MockImageService) {
				m.On("CleanupTempFiles", mock.AnythingOfType("time.Duration")).Return(fmt.Errorf("cleanup failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to cleanup temp files",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockImageService{}
			handler := NewImageHandler(mockService)
			
			tt.mockSetup(mockService)
			
			req := httptest.NewRequest(http.MethodPost, "/api/images/cleanup", nil)
			rr := httptest.NewRecorder()
			
			handler.CleanupTempFiles(rr, req)
			
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
			
			mockService.AssertExpectations(t)
		})
	}
}

