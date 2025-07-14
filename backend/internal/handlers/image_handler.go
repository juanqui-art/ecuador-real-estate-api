package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"realty-core/internal/domain"
	"realty-core/internal/service"
)

// ImageHandler handles HTTP requests for image operations
type ImageHandler struct {
	imageService service.ImageServiceInterface
}

// NewImageHandler creates a new image handler
func NewImageHandler(imageService service.ImageServiceInterface) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
	}
}

// UploadImage handles image upload requests
func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		h.sendErrorResponse(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get property ID from form
	propertyID := r.FormValue("property_id")
	if propertyID == "" {
		h.sendErrorResponse(w, "Property ID is required", http.StatusBadRequest)
		return
	}

	// Get alt text (optional)
	altText := r.FormValue("alt_text")

	// Get uploaded file
	file, handler, err := r.FormFile("image")
	if err != nil {
		h.sendErrorResponse(w, "Failed to get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Upload and process image
	imageInfo, err := h.imageService.Upload(propertyID, file, handler, altText)
	if err != nil {
		h.sendErrorResponse(w, fmt.Sprintf("Failed to upload image: %v", err), http.StatusBadRequest)
		return
	}

	h.sendSuccessResponse(w, "Image uploaded successfully", imageInfo)
}

// GetImage handles requests to get image metadata
func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract image ID from URL path
	imageID := h.extractIDFromPath(r.URL.Path, "/api/images/")
	if imageID == "" {
		h.sendErrorResponse(w, "Image ID is required", http.StatusBadRequest)
		return
	}

	// Get image
	image, err := h.imageService.GetImage(imageID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.sendErrorResponse(w, "Image not found", http.StatusNotFound)
		} else {
			h.sendErrorResponse(w, fmt.Sprintf("Failed to get image: %v", err), http.StatusInternalServerError)
		}
		return
	}

	h.sendSuccessResponse(w, "Image retrieved successfully", image)
}

// GetImagesByProperty handles requests to get all images for a property
func (h *ImageHandler) GetImagesByProperty(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract property ID from URL path
	propertyID := h.extractIDFromPath(r.URL.Path, "/api/properties/")
	if propertyID == "" {
		h.sendErrorResponse(w, "Property ID is required", http.StatusBadRequest)
		return
	}

	// Get images for property
	images, err := h.imageService.GetImagesByProperty(propertyID)
	if err != nil {
		h.sendErrorResponse(w, fmt.Sprintf("Failed to get images: %v", err), http.StatusInternalServerError)
		return
	}

	h.sendSuccessResponse(w, "Images retrieved successfully", images)
}

// UpdateImageMetadata handles requests to update image metadata
func (h *ImageHandler) UpdateImageMetadata(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract image ID from URL path
	imageID := h.extractIDFromPath(r.URL.Path, "/api/images/")
	if imageID == "" {
		h.sendErrorResponse(w, "Image ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req struct {
		AltText   string `json:"alt_text"`
		SortOrder int    `json:"sort_order"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update metadata
	err := h.imageService.UpdateImageMetadata(imageID, req.AltText, req.SortOrder)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.sendErrorResponse(w, "Image not found", http.StatusNotFound)
		} else {
			h.sendErrorResponse(w, fmt.Sprintf("Failed to update image: %v", err), http.StatusInternalServerError)
		}
		return
	}

	h.sendSuccessResponse(w, "Image metadata updated successfully", nil)
}

// DeleteImage handles requests to delete an image
func (h *ImageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract image ID from URL path
	imageID := h.extractIDFromPath(r.URL.Path, "/api/images/")
	if imageID == "" {
		h.sendErrorResponse(w, "Image ID is required", http.StatusBadRequest)
		return
	}

	// Delete image
	err := h.imageService.DeleteImage(imageID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.sendErrorResponse(w, "Image not found", http.StatusNotFound)
		} else {
			h.sendErrorResponse(w, fmt.Sprintf("Failed to delete image: %v", err), http.StatusInternalServerError)
		}
		return
	}

	h.sendSuccessResponse(w, "Image deleted successfully", nil)
}

// ReorderImages handles requests to reorder images for a property
func (h *ImageHandler) ReorderImages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract property ID from URL path
	propertyID := h.extractIDFromPath(r.URL.Path, "/api/properties/")
	if propertyID == "" {
		h.sendErrorResponse(w, "Property ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req domain.ImageReorderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if len(req.ImageIDs) == 0 {
		h.sendErrorResponse(w, "Image IDs are required", http.StatusBadRequest)
		return
	}

	// Reorder images
	err := h.imageService.ReorderImages(propertyID, req.ImageIDs)
	if err != nil {
		h.sendErrorResponse(w, fmt.Sprintf("Failed to reorder images: %v", err), http.StatusBadRequest)
		return
	}

	h.sendSuccessResponse(w, "Images reordered successfully", nil)
}

// SetMainImage handles requests to set an image as the main image
func (h *ImageHandler) SetMainImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract property ID from URL path
	propertyID := h.extractIDFromPath(r.URL.Path, "/api/properties/")
	if propertyID == "" {
		h.sendErrorResponse(w, "Property ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req struct {
		ImageID string `json:"image_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ImageID == "" {
		h.sendErrorResponse(w, "Image ID is required", http.StatusBadRequest)
		return
	}

	// Set main image
	err := h.imageService.SetMainImage(propertyID, req.ImageID)
	if err != nil {
		h.sendErrorResponse(w, fmt.Sprintf("Failed to set main image: %v", err), http.StatusBadRequest)
		return
	}

	h.sendSuccessResponse(w, "Main image set successfully", nil)
}

// GetMainImage handles requests to get the main image for a property
func (h *ImageHandler) GetMainImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract property ID from URL path
	propertyID := h.extractIDFromPath(r.URL.Path, "/api/properties/")
	if propertyID == "" {
		h.sendErrorResponse(w, "Property ID is required", http.StatusBadRequest)
		return
	}

	// Get main image
	image, err := h.imageService.GetMainImage(propertyID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.sendErrorResponse(w, "No images found for property", http.StatusNotFound)
		} else {
			h.sendErrorResponse(w, fmt.Sprintf("Failed to get main image: %v", err), http.StatusInternalServerError)
		}
		return
	}

	h.sendSuccessResponse(w, "Main image retrieved successfully", image)
}

// GetImageVariant handles requests to get image variants (thumbnails, resized images)
func (h *ImageHandler) GetImageVariant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract image ID from URL path
	imageID := h.extractIDFromPath(r.URL.Path, "/api/images/")
	if imageID == "" {
		h.sendErrorResponse(w, "Image ID is required", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	query := r.URL.Query()
	width := h.parseIntParam(query.Get("w"), 0)
	height := h.parseIntParam(query.Get("h"), 0)
	quality := h.parseIntParam(query.Get("q"), domain.DefaultQuality)
	format := query.Get("f")

	// Default format
	if format == "" {
		format = "jpg"
	}

	// Validate format
	if !domain.IsValidImageFormat(format) {
		h.sendErrorResponse(w, fmt.Sprintf("Invalid format: %s", format), http.StatusBadRequest)
		return
	}

	// Get image variant
	imageData, err := h.imageService.GetImageVariant(imageID, width, height, format, quality)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.sendErrorResponse(w, "Image not found", http.StatusNotFound)
		} else {
			h.sendErrorResponse(w, fmt.Sprintf("Failed to get image variant: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Set appropriate headers
	contentType := fmt.Sprintf("image/%s", format)
	if format == "jpg" {
		contentType = "image/jpeg"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	w.Header().Set("Content-Length", strconv.Itoa(len(imageData)))

	// Write image data
	w.Write(imageData)
}

// GetThumbnail handles requests to get image thumbnails
func (h *ImageHandler) GetThumbnail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract image ID from URL path
	imageID := h.extractIDFromPath(r.URL.Path, "/api/images/")
	if imageID == "" {
		h.sendErrorResponse(w, "Image ID is required", http.StatusBadRequest)
		return
	}

	// Parse size parameter
	sizeParam := r.URL.Query().Get("size")
	size := domain.ThumbnailSize
	if sizeParam != "" {
		if parsedSize, err := strconv.Atoi(sizeParam); err == nil && parsedSize > 0 {
			size = parsedSize
		}
	}

	// Get thumbnail
	thumbnailData, err := h.imageService.GenerateThumbnail(imageID, size)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.sendErrorResponse(w, "Image not found", http.StatusNotFound)
		} else {
			h.sendErrorResponse(w, fmt.Sprintf("Failed to generate thumbnail: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	w.Header().Set("Content-Length", strconv.Itoa(len(thumbnailData)))

	// Write thumbnail data
	w.Write(thumbnailData)
}

// GetImageStats handles requests to get image statistics
func (h *ImageHandler) GetImageStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get image statistics
	stats, err := h.imageService.GetImageStats()
	if err != nil {
		h.sendErrorResponse(w, fmt.Sprintf("Failed to get image stats: %v", err), http.StatusInternalServerError)
		return
	}

	h.sendSuccessResponse(w, "Image statistics retrieved successfully", stats)
}

// CleanupTempFiles handles requests to cleanup temporary files
func (h *ImageHandler) CleanupTempFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse duration parameter (default: 24 hours)
	hoursParam := r.URL.Query().Get("hours")
	hours := 24
	if hoursParam != "" {
		if parsedHours, err := strconv.Atoi(hoursParam); err == nil && parsedHours > 0 {
			hours = parsedHours
		}
	}

	duration := time.Duration(hours) * time.Hour

	// Cleanup temporary files
	err := h.imageService.CleanupTempFiles(duration)
	if err != nil {
		h.sendErrorResponse(w, fmt.Sprintf("Failed to cleanup temp files: %v", err), http.StatusInternalServerError)
		return
	}

	h.sendSuccessResponse(w, "Temporary files cleaned up successfully", nil)
}

// GetCacheStats handles requests to get cache statistics
func (h *ImageHandler) GetCacheStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get cache statistics
	stats := h.imageService.GetCacheStats()

	h.sendSuccessResponse(w, "Cache statistics retrieved successfully", stats)
}

// Helper methods

// extractIDFromPath extracts ID from URL path
func (h *ImageHandler) extractIDFromPath(path, prefix string) string {
	if !strings.HasPrefix(path, prefix) {
		return ""
	}

	remaining := strings.TrimPrefix(path, prefix)
	parts := strings.Split(remaining, "/")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}

	return ""
}

// parseIntParam parses integer parameter with default value
func (h *ImageHandler) parseIntParam(param string, defaultValue int) int {
	if param == "" {
		return defaultValue
	}

	if value, err := strconv.Atoi(param); err == nil {
		return value
	}

	return defaultValue
}

// sendSuccessResponse sends a success response
func (h *ImageHandler) sendSuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// sendErrorResponse sends an error response
func (h *ImageHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Success: false,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

