# Realty Image Management

Manage property images for: $ARGUMENTS

## Context - Current Image System
Complete image management with:
- **Upload:** Multi-part form upload with validation
- **Processing:** Resize, compress, thumbnails, watermarks
- **Storage:** Local filesystem with organized structure
- **Cache:** LRU cache for thumbnails and variants
- **Metadata:** PostgreSQL storage with property relations

## Image Operations:
1. **Upload and validation:**
   - Support formats: JPEG, PNG, WebP
   - Size limits: 10MB max upload
   - Dimensions: 3000x2000 max
   - Metadata extraction: EXIF, size, format

2. **Processing pipeline:**
   - Resize to standard sizes
   - Compress with quality settings
   - Generate thumbnails (150px, 300px, 600px)
   - Apply watermarks for public display
   - Convert formats as needed

3. **Storage organization:**
   ```
   storage/images/
   ├── originals/
   ├── thumbnails/
   ├── variants/
   └── temp/
   ```

4. **Cache integration:**
   - Cache thumbnails and variants
   - Invalidate on image updates
   - Preload popular images
   - Monitor cache hit rates

## Image processing examples:
```go
// Thumbnail generation
func (p *ImageProcessor) GenerateThumbnail(imageData []byte, size int) ([]byte, error) {
    img, err := p.DecodeImage(imageData)
    if err != nil {
        return nil, err
    }
    
    thumbnail := p.ResizeImage(img, size, size)
    return p.EncodeJPEG(thumbnail, 85)
}

// Variant creation
func (p *ImageProcessor) CreateVariant(imageData []byte, width, height, quality int, format string) ([]byte, error) {
    // Process and cache variant
}
```

## API endpoints:
- `POST /api/images` - Upload image
- `GET /api/images/{id}/thumbnail?size=300` - Get thumbnail
- `GET /api/images/{id}/variant?w=800&h=600&q=85` - Get variant
- `PUT /api/images/{id}/metadata` - Update metadata
- `DELETE /api/images/{id}` - Delete image

## Common use cases:
- "add watermark to property images"
- "optimize image compression for web"
- "generate responsive image variants"
- "implement image lazy loading"
- "create image gallery component"
- "add image metadata editing"

## Performance considerations:
- Async processing for large images
- Cache frequently accessed variants
- Optimize storage paths
- Use image CDN for production
- Monitor processing times

## Output format:
- Image processing functions
- Storage optimization
- Cache integration
- API endpoints
- Performance optimizations