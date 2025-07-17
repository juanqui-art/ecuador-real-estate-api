'use client';

import { useState, useCallback, useRef, useEffect } from 'react';
import { Upload, X, Image as ImageIcon, Plus, AlertCircle, GripVertical, ZoomIn, Star } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { imageProcessor } from '@/lib/image-processor';
import { ImageProcessorStats } from './image-processor-stats';

export interface TemporaryImageFile {
  id: string;
  file: File;
  preview: string;
  processed?: {
    file: File;
    preview: string;
  };
  isProcessing: boolean;
  error?: string;
  isMain: boolean;
}

interface TemporaryImageUploadProps {
  images: TemporaryImageFile[];
  onImagesChange: (images: TemporaryImageFile[]) => void;
  maxImages?: number;
  className?: string;
}

export function TemporaryImageUpload({ 
  images, 
  onImagesChange, 
  maxImages = 10,
  className = ""
}: TemporaryImageUploadProps) {
  const [isDragging, setIsDragging] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [draggedIndex, setDraggedIndex] = useState<number | null>(null);
  const [draggedOver, setDraggedOver] = useState<number | null>(null);
  const [previewImage, setPreviewImage] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const processImage = useCallback(async (file: File, tempImage: TemporaryImageFile) => {
    try {
      // Update processing state
      onImagesChange(images.map(img => 
        img.id === tempImage.id 
          ? { ...img, isProcessing: true, error: undefined }
          : img
      ));

      // Process the image with optimized settings for real estate
      const processed = await imageProcessor.processPropertyImage(file);

      // Update with processed result
      onImagesChange(images.map(img => 
        img.id === tempImage.id 
          ? { 
              ...img, 
              processed: {
                file: processed.file,
                preview: processed.preview,
              },
              isProcessing: false,
              error: undefined
            }
          : img
      ));
    } catch (error) {
      console.error('Error processing image:', error);
      onImagesChange(images.map(img => 
        img.id === tempImage.id 
          ? { 
              ...img, 
              isProcessing: false,
              error: error instanceof Error ? error.message : 'Error al procesar la imagen'
            }
          : img
      ));
    }
  }, [images, onImagesChange]);

  const handleFileSelect = useCallback(async (files: FileList | null) => {
    if (!files || files.length === 0) return;

    const remainingSlots = maxImages - images.length;
    const filesToProcess = Array.from(files).slice(0, remainingSlots);

    if (filesToProcess.length === 0) {
      alert(`Máximo ${maxImages} imágenes permitidas`);
      return;
    }

    setUploadProgress(0);
    
    const newImages: TemporaryImageFile[] = [];
    
    for (let i = 0; i < filesToProcess.length; i++) {
      const file = filesToProcess[i];
      
      // Validate file
      if (!file.type.startsWith('image/')) {
        continue;
      }
      
      if (file.size > 10 * 1024 * 1024) { // 10MB limit
        alert(`La imagen ${file.name} es demasiado grande. Máximo 10MB.`);
        continue;
      }

      // Create preview
      const preview = URL.createObjectURL(file);
      
      const tempImage: TemporaryImageFile = {
        id: `temp-${Date.now()}-${i}`,
        file,
        preview,
        isProcessing: false,
        isMain: images.length === 0 && i === 0, // First image is main
      };

      newImages.push(tempImage);
      
      // Update progress
      setUploadProgress(((i + 1) / filesToProcess.length) * 100);
    }

    // Add new images
    const updatedImages = [...images, ...newImages];
    onImagesChange(updatedImages);

    // Process images asynchronously
    newImages.forEach(tempImage => {
      processImage(tempImage.file, tempImage);
    });

    // Reset progress after a delay
    setTimeout(() => setUploadProgress(0), 2000);
  }, [images, maxImages, onImagesChange, processImage]);

  const handleFileDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  }, []);

  const handleFileDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  }, []);

  const handleFileDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    handleFileSelect(e.dataTransfer.files);
  }, [handleFileSelect]);

  const handleFileInputChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    handleFileSelect(e.target.files);
  }, [handleFileSelect]);

  const removeImage = useCallback((imageId: string) => {
    const imageToRemove = images.find(img => img.id === imageId);
    if (imageToRemove) {
      URL.revokeObjectURL(imageToRemove.preview);
      if (imageToRemove.processed) {
        URL.revokeObjectURL(imageToRemove.processed.preview);
      }
    }
    
    const updatedImages = images.filter(img => img.id !== imageId);
    
    // If removed image was main, make first image main
    if (imageToRemove?.isMain && updatedImages.length > 0) {
      updatedImages[0].isMain = true;
    }
    
    onImagesChange(updatedImages);
  }, [images, onImagesChange]);

  const setMainImage = useCallback((imageId: string) => {
    onImagesChange(images.map(img => ({
      ...img,
      isMain: img.id === imageId
    })));
  }, [images, onImagesChange]);

  const reorderImages = useCallback((fromIndex: number, toIndex: number) => {
    const reorderedImages = [...images];
    const [movedImage] = reorderedImages.splice(fromIndex, 1);
    reorderedImages.splice(toIndex, 0, movedImage);
    onImagesChange(reorderedImages);
  }, [images, onImagesChange]);

  // Drag and drop handlers for image reordering
  const handleDragStart = useCallback((e: React.DragEvent, index: number) => {
    setDraggedIndex(index);
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/html', '');
  }, []);

  const handleImageDragOver = useCallback((e: React.DragEvent, index: number) => {
    e.preventDefault();
    setDraggedOver(index);
  }, []);

  const handleImageDragLeave = useCallback(() => {
    setDraggedOver(null);
  }, []);

  const handleImageDrop = useCallback((e: React.DragEvent, dropIndex: number) => {
    e.preventDefault();
    if (draggedIndex !== null && draggedIndex !== dropIndex) {
      reorderImages(draggedIndex, dropIndex);
    }
    setDraggedIndex(null);
    setDraggedOver(null);
  }, [draggedIndex, reorderImages]);

  const handleImagePreview = useCallback((imageUrl: string) => {
    setPreviewImage(imageUrl);
  }, []);

  const closePreview = useCallback(() => {
    setPreviewImage(null);
  }, []);

  // Close preview with ESC key
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && previewImage) {
        closePreview();
      }
    };
    
    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [previewImage, closePreview]);

  const processingCount = images.filter(img => img.isProcessing).length;
  const errorCount = images.filter(img => img.error).length;

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Upload Area */}
      <Card>
        <CardContent className="pt-6">
          <div
            className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
              isDragging
                ? 'border-blue-500 bg-blue-50'
                : 'border-gray-300 hover:border-gray-400'
            }`}
            onDragOver={handleFileDragOver}
            onDragLeave={handleFileDragLeave}
            onDrop={handleFileDrop}
          >
            <div className="flex flex-col items-center space-y-2">
              <Upload className="h-12 w-12 text-gray-400" />
              <div className="text-lg font-medium">
                Sube imágenes de la propiedad
              </div>
              <div className="text-sm text-gray-500">
                Arrastra y suelta imágenes aquí, o haz clic para seleccionar
              </div>
              <div className="text-xs text-gray-400">
                Máximo {maxImages} imágenes • PNG, JPG, WEBP • Máximo 10MB cada una
              </div>
              <Button
                type="button"
                variant="outline"
                onClick={() => fileInputRef.current?.click()}
                disabled={images.length >= maxImages}
                className="mt-4"
              >
                <Plus className="h-4 w-4 mr-2" />
                Seleccionar Imágenes
              </Button>
            </div>
          </div>

          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            multiple
            onChange={handleFileInputChange}
            className="hidden"
          />
        </CardContent>
      </Card>

      {/* Upload Progress */}
      {uploadProgress > 0 && (
        <Card>
          <CardContent className="pt-6">
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span>Subiendo imágenes...</span>
                <span>{Math.round(uploadProgress)}%</span>
              </div>
              <Progress value={uploadProgress} className="h-2" />
            </div>
          </CardContent>
        </Card>
      )}

      {/* Status Messages */}
      {(processingCount > 0 || errorCount > 0) && (
        <div className="space-y-2">
          {processingCount > 0 && (
            <Alert>
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                Procesando {processingCount} imagen{processingCount > 1 ? 'es' : ''}...
              </AlertDescription>
            </Alert>
          )}
          {errorCount > 0 && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>
                {errorCount} imagen{errorCount > 1 ? 'es' : ''} con errores de procesamiento
              </AlertDescription>
            </Alert>
          )}
        </div>
      )}

      {/* Image Grid */}
      {images.length > 0 && (
        <Card>
          <CardContent className="pt-6">
            <div className="space-y-4">
              <div className="flex justify-between items-center">
                <h4 className="font-medium">Imágenes seleccionadas ({images.length}/{maxImages})</h4>
                <div className="flex items-center gap-2">
                  <Badge variant="outline">
                    {images.filter(img => img.isMain).length > 0 ? 'Imagen principal seleccionada' : 'Selecciona imagen principal'}
                  </Badge>
                  {images.length > 1 && (
                    <Badge variant="secondary" className="text-xs">
                      <GripVertical className="h-3 w-3 mr-1" />
                      Arrastra para reordenar
                    </Badge>
                  )}
                </div>
              </div>
              
              <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
                {images.map((image, index) => (
                  <div
                    key={image.id}
                    draggable={!image.isProcessing && !image.error}
                    onDragStart={(e) => handleDragStart(e, index)}
                    onDragOver={(e) => handleImageDragOver(e, index)}
                    onDragLeave={handleImageDragLeave}
                    onDrop={(e) => handleImageDrop(e, index)}
                    className={`relative group border-2 rounded-lg overflow-hidden cursor-move transition-all duration-200 ${
                      image.isMain ? 'border-blue-500 ring-2 ring-blue-200' : 'border-gray-200'
                    } ${
                      draggedIndex === index ? 'opacity-50 scale-95' : ''
                    } ${
                      draggedOver === index && draggedIndex !== index ? 'border-blue-400 bg-blue-50' : ''
                    }`}
                  >
                    <div className="aspect-square bg-gray-100 flex items-center justify-center relative">
                      {image.isProcessing ? (
                        <div className="flex flex-col items-center space-y-2">
                          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
                          <span className="text-xs text-gray-500">Procesando...</span>
                        </div>
                      ) : image.error ? (
                        <div className="flex flex-col items-center space-y-2 text-red-500">
                          <AlertCircle className="h-8 w-8" />
                          <span className="text-xs text-center px-2">{image.error}</span>
                        </div>
                      ) : (
                        <img
                          src={image.processed?.preview || image.preview}
                          alt={`Imagen ${index + 1}`}
                          className="w-full h-full object-cover"
                        />
                      )}
                      
                      {/* Drag Handle */}
                      {!image.isProcessing && !image.error && (
                        <div className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
                          <div className="bg-black bg-opacity-50 rounded p-1">
                            <GripVertical className="h-4 w-4 text-white" />
                          </div>
                        </div>
                      )}
                    </div>

                    {/* Overlay Controls */}
                    <div className="absolute inset-0 bg-black bg-opacity-50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                      <div className="flex space-x-2">
                        {!image.isProcessing && !image.error && (
                          <Button
                            type="button"
                            size="sm"
                            variant="secondary"
                            onClick={() => handleImagePreview(image.processed?.preview || image.preview)}
                            className="text-xs"
                          >
                            <ZoomIn className="h-3 w-3" />
                          </Button>
                        )}
                        {!image.isMain && !image.isProcessing && !image.error && (
                          <Button
                            type="button"
                            size="sm"
                            variant="secondary"
                            onClick={() => setMainImage(image.id)}
                            className="text-xs"
                          >
                            <Star className="h-3 w-3" />
                          </Button>
                        )}
                        <Button
                          type="button"
                          size="sm"
                          variant="destructive"
                          onClick={() => removeImage(image.id)}
                        >
                          <X className="h-3 w-3" />
                        </Button>
                      </div>
                    </div>

                    {/* Main Image Badge */}
                    {image.isMain && (
                      <div className="absolute top-2 left-2">
                        <Badge variant="default" className="text-xs flex items-center gap-1">
                          <Star className="h-3 w-3" />
                          Principal
                        </Badge>
                      </div>
                    )}

                    {/* Position Badge */}
                    <div className="absolute bottom-2 left-2">
                      <Badge variant="secondary" className="text-xs">
                        #{index + 1}
                      </Badge>
                    </div>

                    {/* File Size Badge */}
                    <div className="absolute bottom-2 right-2">
                      <Badge variant="secondary" className="text-xs">
                        {(image.file.size / 1024 / 1024).toFixed(1)}MB
                      </Badge>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Help Text */}
      <div className="text-sm text-gray-500 space-y-1">
        <div className="flex items-center gap-2">
          <span className="w-2 h-2 bg-blue-500 rounded-full"></span>
          <p>La primera imagen será la imagen principal por defecto</p>
        </div>
        <div className="flex items-center gap-2">
          <span className="w-2 h-2 bg-green-500 rounded-full"></span>
          <p>Las imágenes se procesarán automáticamente para optimizar el tamaño</p>
        </div>
        <div className="flex items-center gap-2">
          <span className="w-2 h-2 bg-yellow-500 rounded-full"></span>
          <p>Arrastra y suelta para reordenar las imágenes</p>
        </div>
        <div className="flex items-center gap-2">
          <span className="w-2 h-2 bg-purple-500 rounded-full"></span>
          <p>Formatos soportados: JPEG, PNG, WebP, AVIF</p>
        </div>
      </div>

      {/* Processing Stats */}
      <ImageProcessorStats />

      {/* Image Preview Modal */}
      {previewImage && (
        <div className="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center z-50">
          <div className="relative max-w-4xl max-h-4xl p-4">
            <img
              src={previewImage}
              alt="Preview"
              className="max-w-full max-h-full object-contain"
            />
            <Button
              type="button"
              variant="secondary"
              size="sm"
              onClick={closePreview}
              className="absolute top-2 right-2"
            >
              <X className="h-4 w-4" />
            </Button>
            <div className="absolute bottom-2 left-2 text-white text-sm bg-black bg-opacity-50 px-2 py-1 rounded">
              Presiona ESC para cerrar
            </div>
          </div>
        </div>
      )}
    </div>
  );
}