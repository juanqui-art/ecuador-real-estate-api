/**
 * Client-side image processing utilities
 * Handles compression, resizing, and optimization before upload
 */

export interface ImageProcessingOptions {
  maxWidth?: number;
  maxHeight?: number;
  quality?: number;
  format?: 'jpeg' | 'webp' | 'png';
  maintainAspectRatio?: boolean;
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
}

export class ImageProcessor {
  private canvas: HTMLCanvasElement | null = null;
  private ctx: CanvasRenderingContext2D | null = null;

  constructor() {
    // Don't initialize canvas during SSR
    if (typeof window !== 'undefined' && typeof document !== 'undefined') {
      this.initializeCanvas();
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
    }
  }

  /**
   * Process a single image file
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
      format = 'jpeg',
      maintainAspectRatio = true,
    } = options;

    return new Promise((resolve, reject) => {
      const img = new Image();
      const reader = new FileReader();

      reader.onload = (e) => {
        img.src = e.target?.result as string;
      };

      img.onload = () => {
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

          // Apply image smoothing
          this.ctx!.imageSmoothingEnabled = true;
          this.ctx!.imageSmoothingQuality = 'high';

          // Draw resized image
          this.ctx!.drawImage(img, 0, 0, width, height);

          // Convert to blob
          this.canvas!.toBlob(
            (blob) => {
              if (!blob) {
                reject(new Error('Failed to process image'));
                return;
              }

              // Create new file
              const processedFile = new File(
                [blob],
                this.generateFilename(file.name, format),
                {
                  type: `image/${format}`,
                  lastModified: Date.now(),
                }
              );

              // Create preview URL
              const preview = URL.createObjectURL(blob);

              const result: ProcessedImage = {
                file: processedFile,
                originalSize: file.size,
                compressedSize: blob.size,
                compressionRatio: Math.round(((file.size - blob.size) / file.size) * 100),
                dimensions: { width, height },
                preview,
              };

              resolve(result);
            },
            `image/${format}`,
            quality
          );
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
   * Process multiple images
   */
  async processImages(
    files: File[],
    options: ImageProcessingOptions = {}
  ): Promise<ProcessedImage[]> {
    const results: ProcessedImage[] = [];
    
    for (const file of files) {
      try {
        const processed = await this.processImage(file, options);
        results.push(processed);
      } catch (error) {
        console.error(`Failed to process ${file.name}:`, error);
        // Continue with other files
      }
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
   * Generate filename with proper extension
   */
  private generateFilename(originalName: string, format: string): string {
    const nameWithoutExt = originalName.replace(/\.[^/.]+$/, '');
    const extension = format === 'jpeg' ? 'jpg' : format;
    return `${nameWithoutExt}.${extension}`;
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
   * Validate image file
   */
  validateImage(file: File): { isValid: boolean; error?: string } {
    const validTypes = ['image/jpeg', 'image/png', 'image/webp'];
    const maxSize = 10 * 1024 * 1024; // 10MB

    if (!validTypes.includes(file.type)) {
      return {
        isValid: false,
        error: 'Tipo de archivo no soportado. Usa JPEG, PNG o WebP.',
      };
    }

    if (file.size > maxSize) {
      return {
        isValid: false,
        error: 'Archivo muy grande. Máximo 10MB.',
      };
    }

    return { isValid: true };
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
  getImageMetadata: (file: File) => 
    getImageProcessor().getImageMetadata(file),
  createThumbnail: (file: File, size?: number) => 
    getImageProcessor().createThumbnail(file, size),
  validateImage: (file: File) => 
    getImageProcessor().validateImage(file),
  cleanup: (urls: string[]) => 
    getImageProcessor().cleanup(urls),
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