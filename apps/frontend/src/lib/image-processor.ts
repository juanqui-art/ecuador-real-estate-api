/**
 * Client-side image processing utilities
 * Handles compression, resizing, and optimization before upload
 */

export type ImageFormat = 'jpeg' | 'webp' | 'png' | 'avif' | 'auto';

export interface ImageProcessingOptions {
  maxWidth?: number;
  maxHeight?: number;
  quality?: number;
  format?: ImageFormat;
  maintainAspectRatio?: boolean;
  progressive?: boolean;
  optimize?: boolean;
  fallbackFormat?: 'jpeg' | 'webp' | 'png';
  variants?: {
    thumbnail?: { width: number; height: number; quality: number };
    medium?: { width: number; height: number; quality: number };
    large?: { width: number; height: number; quality: number };
  };
}

export interface ProcessedImage {
  file: File;
  originalSize: number;
  compressedSize: number;
  compressionRatio: number;
  dimensions: {
    width: number;
    height: number;
  };
  preview: string;
  format: string;
  variants?: {
    thumbnail?: { file: File; preview: string; dimensions: { width: number; height: number } };
    medium?: { file: File; preview: string; dimensions: { width: number; height: number } };
    large?: { file: File; preview: string; dimensions: { width: number; height: number } };
  };
  metadata?: {
    aspectRatio: number;
    orientation: number;
    hasTransparency: boolean;
    colorProfile?: string;
  };
}

export class ImageProcessor {
  private canvas: HTMLCanvasElement | null = null;
  private ctx: CanvasRenderingContext2D | null = null;
  private formatSupport: Map<string, boolean> = new Map();

  constructor() {
    // Don't initialize canvas during SSR
    if (typeof window !== 'undefined' && typeof document !== 'undefined') {
      this.initializeCanvas();
      this.detectFormatSupport();
    }
  }

  private initializeCanvas(): void {
    if (!this.canvas) {
      this.canvas = document.createElement('canvas');
      this.ctx = this.canvas.getContext('2d')!;
    }
  }

  private ensureCanvas(): void {
    if (typeof window === 'undefined' || typeof document === 'undefined') {
      throw new Error('ImageProcessor can only be used in browser environments');
    }
    if (!this.canvas) {
      this.initializeCanvas();
      this.detectFormatSupport();
    }
  }

  private detectFormatSupport(): void {
    if (!this.canvas) return;
    
    const testFormats = ['webp', 'avif', 'jpeg', 'png'];
    
    testFormats.forEach(format => {
      try {
        const dataURL = this.canvas!.toDataURL(`image/${format}`);
        this.formatSupport.set(format, dataURL.includes(`data:image/${format}`));
      } catch (error) {
        this.formatSupport.set(format, false);
      }
    });
  }

  private getBestFormat(requestedFormat: ImageFormat, hasTransparency: boolean): string {
    if (requestedFormat === 'auto') {
      // Smart format selection based on browser support and image characteristics
      if (hasTransparency) {
        if (this.formatSupport.get('avif')) return 'avif';
        if (this.formatSupport.get('webp')) return 'webp';
        return 'png';
      } else {
        if (this.formatSupport.get('avif')) return 'avif';
        if (this.formatSupport.get('webp')) return 'webp';
        return 'jpeg';
      }
    }
    
    if (this.formatSupport.get(requestedFormat)) {
      return requestedFormat;
    }
    
    // Fallback to supported format
    return hasTransparency ? 'png' : 'jpeg';
  }

  private analyzeImage(imageData: ImageData): { hasTransparency: boolean; complexity: number } {
    const data = imageData.data;
    let transparentPixels = 0;
    let colorVariations = new Set<number>();
    
    // Sample every 4th pixel for performance
    for (let i = 0; i < data.length; i += 16) {
      const r = data[i];
      const g = data[i + 1];
      const b = data[i + 2];
      const a = data[i + 3];
      
      if (a < 255) transparentPixels++;
      
      // Simple color hash for complexity analysis
      const colorHash = (r << 16) | (g << 8) | b;
      colorVariations.add(colorHash);
    }
    
    const hasTransparency = transparentPixels > 0;
    const complexity = colorVariations.size / (imageData.width * imageData.height / 4);
    
    return { hasTransparency, complexity };
  }

  /**
   * Process a single image file with advanced optimization
   */
  async processImage(
    file: File,
    options: ImageProcessingOptions = {}
  ): Promise<ProcessedImage> {
    this.ensureCanvas();
    
    const {
      maxWidth = 1920,
      maxHeight = 1080,
      quality = 0.8,
      format = 'auto',
      maintainAspectRatio = true,
      progressive = true,
      optimize = true,
      fallbackFormat = 'jpeg',
      variants
    } = options;

    return new Promise((resolve, reject) => {
      const img = new Image();
      const reader = new FileReader();

      reader.onload = (e) => {
        img.src = e.target?.result as string;
      };

      img.onload = async () => {
        try {
          const { width, height } = this.calculateDimensions(
            img.width,
            img.height,
            maxWidth,
            maxHeight,
            maintainAspectRatio
          );

          this.canvas!.width = width;
          this.canvas!.height = height;

          // Clear canvas
          this.ctx!.clearRect(0, 0, width, height);

          // Apply advanced image smoothing
          this.ctx!.imageSmoothingEnabled = true;
          this.ctx!.imageSmoothingQuality = 'high';

          // Draw resized image
          this.ctx!.drawImage(img, 0, 0, width, height);

          // Analyze image for optimal format selection
          const imageData = this.ctx!.getImageData(0, 0, width, height);
          const analysis = this.analyzeImage(imageData);
          
          // Get optimal format
          const optimalFormat = this.getBestFormat(format, analysis.hasTransparency);
          
          // Adjust quality based on image complexity
          let adjustedQuality = quality;
          if (optimize) {
            adjustedQuality = this.optimizeQuality(quality, analysis.complexity, optimalFormat);
          }

          // Generate main image
          const mainBlob = await this.canvasToBlob(this.canvas!, optimalFormat, adjustedQuality);
          if (!mainBlob) {
            reject(new Error('Failed to process image'));
            return;
          }

          // Create main processed file
          const processedFile = new File(
            [mainBlob],
            this.generateFilename(file.name, optimalFormat),
            {
              type: `image/${optimalFormat}`,
              lastModified: Date.now(),
            }
          );

          // Create preview URL
          const preview = URL.createObjectURL(mainBlob);

          // Generate variants if requested
          let processedVariants: ProcessedImage['variants'] = undefined;
          if (variants) {
            processedVariants = await this.generateVariants(
              img,
              variants,
              optimalFormat,
              file.name,
              analysis.hasTransparency
            );
          }

          // Create metadata
          const metadata = {
            aspectRatio: img.width / img.height,
            orientation: this.getImageOrientation(img),
            hasTransparency: analysis.hasTransparency,
            colorProfile: this.detectColorProfile(imageData),
          };

          const result: ProcessedImage = {
            file: processedFile,
            originalSize: file.size,
            compressedSize: mainBlob.size,
            compressionRatio: Math.round(((file.size - mainBlob.size) / file.size) * 100),
            dimensions: { width, height },
            preview,
            format: optimalFormat,
            variants: processedVariants,
            metadata,
          };

          resolve(result);
        } catch (error) {
          reject(error);
        }
      };

      img.onerror = () => {
        reject(new Error('Failed to load image'));
      };

      reader.onerror = () => {
        reject(new Error('Failed to read file'));
      };

      reader.readAsDataURL(file);
    });
  }

  /**
   * Helper method to convert canvas to blob with format support
   */
  private canvasToBlob(canvas: HTMLCanvasElement, format: string, quality: number): Promise<Blob | null> {
    return new Promise((resolve) => {
      canvas.toBlob(resolve, `image/${format}`, quality);
    });
  }

  /**
   * Optimize quality based on image complexity and format
   */
  private optimizeQuality(baseQuality: number, complexity: number, format: string): number {
    let adjustedQuality = baseQuality;
    
    // Adjust quality based on image complexity
    if (complexity > 0.8) {
      // High complexity images (lots of detail) - slightly higher quality
      adjustedQuality = Math.min(1.0, baseQuality + 0.1);
    } else if (complexity < 0.3) {
      // Low complexity images (simple graphics) - can use lower quality
      adjustedQuality = Math.max(0.5, baseQuality - 0.1);
    }
    
    // Format-specific adjustments
    if (format === 'webp' || format === 'avif') {
      // These formats are more efficient, can use slightly lower quality
      adjustedQuality = Math.max(0.6, adjustedQuality - 0.1);
    }
    
    return adjustedQuality;
  }

  /**
   * Generate image variants (thumbnails, medium, large)
   */
  private async generateVariants(
    img: HTMLImageElement,
    variants: NonNullable<ImageProcessingOptions['variants']>,
    format: string,
    originalName: string,
    hasTransparency: boolean
  ): Promise<ProcessedImage['variants']> {
    const result: ProcessedImage['variants'] = {};
    
    for (const [variantName, config] of Object.entries(variants)) {
      if (!config) continue;
      
      const { width, height } = this.calculateDimensions(
        img.width,
        img.height,
        config.width,
        config.height,
        true
      );
      
      this.canvas!.width = width;
      this.canvas!.height = height;
      
      this.ctx!.clearRect(0, 0, width, height);
      this.ctx!.imageSmoothingEnabled = true;
      this.ctx!.imageSmoothingQuality = 'high';
      this.ctx!.drawImage(img, 0, 0, width, height);
      
      const blob = await this.canvasToBlob(this.canvas!, format, config.quality);
      if (!blob) continue;
      
      const variantFile = new File(
        [blob],
        this.generateFilename(originalName, format, variantName),
        {
          type: `image/${format}`,
          lastModified: Date.now(),
        }
      );
      
      result[variantName as keyof ProcessedImage['variants']] = {
        file: variantFile,
        preview: URL.createObjectURL(blob),
        dimensions: { width, height },
      };
    }
    
    return result;
  }

  /**
   * Detect image orientation (simplified)
   */
  private getImageOrientation(img: HTMLImageElement): number {
    // This is a simplified version - in a real implementation
    // you'd parse EXIF data for accurate orientation
    return img.width > img.height ? 1 : 6; // 1 = landscape, 6 = portrait
  }

  /**
   * Detect color profile (simplified)
   */
  private detectColorProfile(imageData: ImageData): string {
    // This is a simplified implementation
    // In a real scenario, you'd analyze the color histogram
    return 'sRGB';
  }

  /**
   * Process multiple images with batch optimization
   */
  async processImages(
    files: File[],
    options: ImageProcessingOptions = {}
  ): Promise<ProcessedImage[]> {
    const results: ProcessedImage[] = [];
    
    // Process images in batches to avoid memory issues
    const batchSize = 3;
    for (let i = 0; i < files.length; i += batchSize) {
      const batch = files.slice(i, i + batchSize);
      
      const batchPromises = batch.map(async (file) => {
        try {
          return await this.processImage(file, options);
        } catch (error) {
          console.error(`Failed to process ${file.name}:`, error);
          return null;
        }
      });
      
      const batchResults = await Promise.all(batchPromises);
      results.push(...batchResults.filter(result => result !== null) as ProcessedImage[]);
    }
    
    return results;
  }

  /**
   * Calculate optimal dimensions while maintaining aspect ratio
   */
  private calculateDimensions(
    originalWidth: number,
    originalHeight: number,
    maxWidth: number,
    maxHeight: number,
    maintainAspectRatio: boolean
  ): { width: number; height: number } {
    if (!maintainAspectRatio) {
      return { width: maxWidth, height: maxHeight };
    }

    const aspectRatio = originalWidth / originalHeight;

    let width = originalWidth;
    let height = originalHeight;

    // Scale down if too large
    if (width > maxWidth) {
      width = maxWidth;
      height = width / aspectRatio;
    }

    if (height > maxHeight) {
      height = maxHeight;
      width = height * aspectRatio;
    }

    return {
      width: Math.round(width),
      height: Math.round(height),
    };
  }

  /**
   * Generate filename with proper extension and variant suffix
   */
  private generateFilename(originalName: string, format: string, variant?: string): string {
    const nameWithoutExt = originalName.replace(/\.[^/.]+$/, '');
    const extension = format === 'jpeg' ? 'jpg' : format;
    const variantSuffix = variant ? `_${variant}` : '';
    return `${nameWithoutExt}${variantSuffix}.${extension}`;
  }

  /**
   * Get image metadata
   */
  async getImageMetadata(file: File): Promise<{
    width: number;
    height: number;
    size: number;
    type: string;
    aspectRatio: number;
  }> {
    return new Promise((resolve, reject) => {
      const img = new Image();
      const reader = new FileReader();

      reader.onload = (e) => {
        img.src = e.target?.result as string;
      };

      img.onload = () => {
        resolve({
          width: img.width,
          height: img.height,
          size: file.size,
          type: file.type,
          aspectRatio: img.width / img.height,
        });
      };

      img.onerror = () => reject(new Error('Failed to load image'));
      reader.onerror = () => reject(new Error('Failed to read file'));
      reader.readAsDataURL(file);
    });
  }

  /**
   * Create thumbnail
   */
  async createThumbnail(
    file: File,
    size: number = 200
  ): Promise<{ file: File; preview: string }> {
    const processed = await this.processImage(file, {
      maxWidth: size,
      maxHeight: size,
      quality: 0.7,
      format: 'jpeg',
      maintainAspectRatio: true,
    });

    return {
      file: processed.file,
      preview: processed.preview,
    };
  }

  /**
   * Validate image file with enhanced checks
   */
  validateImage(file: File): { isValid: boolean; error?: string; warnings?: string[] } {
    const validTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/avif', 'image/svg+xml'];
    const maxSize = 15 * 1024 * 1024; // 15MB
    const minSize = 1024; // 1KB
    const warnings: string[] = [];

    // Check file type
    if (!validTypes.includes(file.type)) {
      return {
        isValid: false,
        error: 'Tipo de archivo no soportado. Usa JPEG, PNG, WebP, AVIF o SVG.',
      };
    }

    // Check file size
    if (file.size > maxSize) {
      return {
        isValid: false,
        error: 'Archivo muy grande. Máximo 15MB.',
      };
    }

    if (file.size < minSize) {
      return {
        isValid: false,
        error: 'Archivo muy pequeño. Mínimo 1KB.',
      };
    }

    // Size warnings
    if (file.size > 5 * 1024 * 1024) {
      warnings.push('Archivo grande (>5MB). Considera comprimir la imagen.');
    }

    // Format warnings
    if (file.type === 'image/png' && file.size > 2 * 1024 * 1024) {
      warnings.push('PNG grande detectado. JPEG o WebP serían más eficientes.');
    }

    if (file.type === 'image/svg+xml') {
      warnings.push('SVG detectado. Se convertirá a formato raster.');
    }

    return { 
      isValid: true, 
      warnings: warnings.length > 0 ? warnings : undefined 
    };
  }

  /**
   * Process image with multiple variants for real estate use case
   */
  async processPropertyImage(file: File): Promise<ProcessedImage> {
    const variants = {
      thumbnail: { width: 300, height: 200, quality: 0.7 },
      medium: { width: 800, height: 600, quality: 0.8 },
      large: { width: 1920, height: 1080, quality: 0.85 },
    };

    return this.processImage(file, {
      maxWidth: 1920,
      maxHeight: 1080,
      quality: 0.85,
      format: 'auto',
      optimize: true,
      variants,
    });
  }

  /**
   * Get format support information
   */
  getFormatSupport(): Record<string, boolean> {
    const support: Record<string, boolean> = {};
    for (const [format, isSupported] of this.formatSupport.entries()) {
      support[format] = isSupported;
    }
    return support;
  }

  /**
   * Get processing statistics
   */
  getProcessingStats(): {
    formatsSupported: string[];
    canvasInitialized: boolean;
    browserSupport: {
      webp: boolean;
      avif: boolean;
      progressive: boolean;
    };
  } {
    return {
      formatsSupported: Array.from(this.formatSupport.keys()).filter(
        format => this.formatSupport.get(format)
      ),
      canvasInitialized: this.canvas !== null,
      browserSupport: {
        webp: this.formatSupport.get('webp') || false,
        avif: this.formatSupport.get('avif') || false,
        progressive: true, // Canvas API supports progressive JPEG
      },
    };
  }

  /**
   * Cleanup blob URLs to prevent memory leaks
   */
  cleanup(urls: string[]): void {
    urls.forEach(url => {
      if (url.startsWith('blob:')) {
        URL.revokeObjectURL(url);
      }
    });
  }

  /**
   * Clean up processed image variants
   */
  cleanupProcessedImage(processedImage: ProcessedImage): void {
    const urlsToClean = [processedImage.preview];
    
    if (processedImage.variants) {
      Object.values(processedImage.variants).forEach(variant => {
        if (variant) {
          urlsToClean.push(variant.preview);
        }
      });
    }
    
    this.cleanup(urlsToClean);
  }
}

// Lazy singleton instance - only initialize on first use
let _imageProcessor: ImageProcessor | null = null;

export function getImageProcessor(): ImageProcessor {
  if (!_imageProcessor) {
    _imageProcessor = new ImageProcessor();
  }
  return _imageProcessor;
}

// Export singleton instance for backward compatibility
export const imageProcessor = {
  processImage: (file: File, options?: ImageProcessingOptions) => 
    getImageProcessor().processImage(file, options),
  processImages: (files: File[], options?: ImageProcessingOptions) => 
    getImageProcessor().processImages(files, options),
  processPropertyImage: (file: File) => 
    getImageProcessor().processPropertyImage(file),
  getImageMetadata: (file: File) => 
    getImageProcessor().getImageMetadata(file),
  createThumbnail: (file: File, size?: number) => 
    getImageProcessor().createThumbnail(file, size),
  validateImage: (file: File) => 
    getImageProcessor().validateImage(file),
  getFormatSupport: () => 
    getImageProcessor().getFormatSupport(),
  getProcessingStats: () => 
    getImageProcessor().getProcessingStats(),
  cleanup: (urls: string[]) => 
    getImageProcessor().cleanup(urls),
  cleanupProcessedImage: (processedImage: ProcessedImage) => 
    getImageProcessor().cleanupProcessedImage(processedImage),
};

// Helper functions
export const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

export const getImageDimensions = (file: File): Promise<{ width: number; height: number }> => {
  return new Promise((resolve, reject) => {
    const img = new Image();
    const reader = new FileReader();

    reader.onload = (e) => {
      img.src = e.target?.result as string;
    };

    img.onload = () => {
      resolve({ width: img.width, height: img.height });
    };

    img.onerror = reject;
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
};