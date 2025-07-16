'use client';

import { useState, useCallback, useRef } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Upload, X, Image as ImageIcon, AlertCircle, Check, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Progress } from '@/components/ui/progress';
import { apiClient } from '@/lib/api-client';
import { cn } from '@/lib/utils';
import { imageProcessor, formatFileSize } from '@/lib/image-processor';

interface ImageUploadProps {
  propertyId?: string;
  onImageUploaded?: (images: UploadedImage[]) => void;
  maxImages?: number;
  maxSizeInMB?: number;
  acceptedTypes?: string[];
  className?: string;
}

interface UploadedImage {
  id: string;
  url: string;
  filename: string;
  size: number;
  type: string;
  is_main: boolean;
  order: number;
}

interface UploadProgress {
  file: File;
  originalFile: File;
  progress: number;
  status: 'processing' | 'uploading' | 'success' | 'error';
  error?: string;
  result?: UploadedImage;
  compressionRatio?: number;
  preview?: string;
}

export function ImageUpload({
  propertyId,
  onImageUploaded,
  maxImages = 10,
  maxSizeInMB = 5,
  acceptedTypes = ['image/jpeg', 'image/png', 'image/webp'],
  className
}: ImageUploadProps) {
  const [isDragging, setIsDragging] = useState(false);
  const [uploads, setUploads] = useState<UploadProgress[]>([]);
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const queryClient = useQueryClient();

  const uploadMutation = useMutation({
    mutationFn: async (file: File) => {
      const formData = new FormData();
      formData.append('image', file);
      if (propertyId) {
        formData.append('property_id', propertyId);
      }

      // Simular progreso de upload
      const response = await apiClient.post('/images', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });

      return response.data;
    },
    onSuccess: (data, file) => {
      setUploads(prev => prev.map(upload => 
        upload.file === file 
          ? { ...upload, status: 'success', progress: 100, result: data }
          : upload
      ));
      
      if (onImageUploaded && data) {
        onImageUploaded([data]);
      }
      
      // Invalidar cache de imágenes
      queryClient.invalidateQueries({ queryKey: ['property-images', propertyId] });
    },
    onError: (error, file) => {
      setUploads(prev => prev.map(upload => 
        upload.file === file 
          ? { ...upload, status: 'error', error: error.message }
          : upload
      ));
    },
  });

  const validateFile = (file: File): string | null => {
    if (!acceptedTypes.includes(file.type)) {
      return `Tipo de archivo no soportado. Usa: ${acceptedTypes.join(', ')}`;
    }
    
    if (file.size > maxSizeInMB * 1024 * 1024) {
      return `Archivo muy grande. Máximo ${maxSizeInMB}MB`;
    }
    
    return null;
  };

  const handleFiles = useCallback(async (files: File[]) => {
    setError(null);
    const validFiles: File[] = [];
    
    for (const file of files) {
      const validation = validateFile(file);
      if (validation) {
        setError(validation);
        continue;
      }
      
      if (uploads.length + validFiles.length >= maxImages) {
        setError(`Máximo ${maxImages} imágenes permitidas`);
        break;
      }
      
      validFiles.push(file);
    }
    
    if (validFiles.length === 0) return;
    
    // Agregar archivos a la cola de procesamiento
    const newUploads: UploadProgress[] = validFiles.map(file => ({
      file,
      originalFile: file,
      progress: 0,
      status: 'processing' as const,
    }));
    
    setUploads(prev => [...prev, ...newUploads]);
    
    // Procesar e iniciar uploads
    for (const file of validFiles) {
      try {
        // Procesar imagen (solo en el cliente)
        if (typeof window === 'undefined') {
          throw new Error('Image processing only available in browser');
        }
        
        const processed = await imageProcessor.processImage(file, {
          maxWidth: 1920,
          maxHeight: 1080,
          quality: 0.85,
          format: 'jpeg',
        });
        
        // Actualizar estado con archivo procesado
        setUploads(prev => prev.map(upload => 
          upload.originalFile === file 
            ? { 
                ...upload, 
                file: processed.file,
                status: 'uploading' as const,
                compressionRatio: processed.compressionRatio,
                preview: processed.preview
              }
            : upload
        ));
        
        // Iniciar upload
        uploadMutation.mutate(processed.file);
      } catch (error) {
        setUploads(prev => prev.map(upload => 
          upload.originalFile === file 
            ? { ...upload, status: 'error', error: 'Error al procesar imagen' }
            : upload
        ));
      }
    }
  }, [uploads.length, maxImages, uploadMutation]);

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  }, []);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    
    const files = Array.from(e.dataTransfer.files);
    handleFiles(files);
  }, [handleFiles]);

  const handleFileInput = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || []);
    handleFiles(files);
    
    // Limpiar input
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  }, [handleFiles]);

  const removeUpload = useCallback((index: number) => {
    setUploads(prev => prev.filter((_, i) => i !== index));
  }, []);

  const clearError = useCallback(() => {
    setError(null);
  }, []);


  return (
    <div className={cn('space-y-4', className)}>
      {/* Error Alert */}
      {error && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription className="flex items-center justify-between">
            {error}
            <Button variant="ghost" size="sm" onClick={clearError}>
              <X className="h-4 w-4" />
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {/* Drop Zone */}
      <Card className={cn(
        'border-2 border-dashed transition-colors cursor-pointer',
        isDragging 
          ? 'border-primary bg-primary/5' 
          : 'border-gray-300 hover:border-gray-400'
      )}>
        <CardContent 
          className="p-8 text-center"
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          onClick={() => fileInputRef.current?.click()}
        >
          <div className="space-y-4">
            <div className="flex justify-center">
              <div className={cn(
                'p-4 rounded-full transition-colors',
                isDragging 
                  ? 'bg-primary text-primary-foreground' 
                  : 'bg-gray-100 text-gray-600'
              )}>
                <Upload className="h-8 w-8" />
              </div>
            </div>
            
            <div>
              <p className="text-lg font-medium text-gray-900">
                {isDragging ? 'Suelta las imágenes aquí' : 'Arrastra imágenes aquí'}
              </p>
              <p className="text-sm text-gray-500 mt-1">
                o <span className="text-primary font-medium">haz clic para seleccionar</span>
              </p>
            </div>
            
            <div className="flex justify-center gap-4 text-xs text-gray-500">
              <span>Máximo {maxImages} imágenes</span>
              <span>•</span>
              <span>Hasta {maxSizeInMB}MB cada una</span>
              <span>•</span>
              <span>JPG, PNG, WebP</span>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* File Input */}
      <input
        ref={fileInputRef}
        type="file"
        multiple
        accept={acceptedTypes.join(',')}
        onChange={handleFileInput}
        className="hidden"
      />

      {/* Upload Progress */}
      {uploads.length > 0 && (
        <Card>
          <CardContent className="p-4">
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <h3 className="font-medium">Subiendo imágenes</h3>
                <Badge variant="secondary">
                  {uploads.filter(u => u.status === 'success').length} / {uploads.length}
                </Badge>
              </div>
              
              <div className="space-y-2">
                {uploads.map((upload, index) => (
                  <div key={index} className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                    {/* Preview thumbnail */}
                    {upload.preview && (
                      <div className="flex-shrink-0">
                        <img 
                          src={upload.preview}
                          alt="Preview"
                          className="w-10 h-10 object-cover rounded"
                        />
                      </div>
                    )}
                    
                    <div className="flex-shrink-0">
                      {upload.status === 'processing' && (
                        <Loader2 className="h-4 w-4 animate-spin text-purple-600" />
                      )}
                      {upload.status === 'uploading' && (
                        <Loader2 className="h-4 w-4 animate-spin text-blue-600" />
                      )}
                      {upload.status === 'success' && (
                        <Check className="h-4 w-4 text-green-600" />
                      )}
                      {upload.status === 'error' && (
                        <AlertCircle className="h-4 w-4 text-red-600" />
                      )}
                    </div>
                    
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center justify-between mb-1">
                        <p className="text-sm font-medium truncate">
                          {upload.originalFile.name}
                        </p>
                        <div className="flex items-center gap-2">
                          {upload.compressionRatio && upload.compressionRatio > 0 && (
                            <Badge variant="secondary" className="text-xs">
                              -{upload.compressionRatio}%
                            </Badge>
                          )}
                          <span className="text-xs text-gray-500">
                            {formatFileSize(upload.file.size)}
                          </span>
                        </div>
                      </div>
                      
                      {upload.status === 'processing' && (
                        <div>
                          <Progress value={50} className="h-1 mb-1" />
                          <p className="text-xs text-purple-600">Procesando imagen...</p>
                        </div>
                      )}
                      
                      {upload.status === 'uploading' && (
                        <div>
                          <Progress value={upload.progress} className="h-1 mb-1" />
                          <p className="text-xs text-blue-600">Subiendo...</p>
                        </div>
                      )}
                      
                      {upload.status === 'error' && (
                        <p className="text-xs text-red-600">{upload.error}</p>
                      )}
                      
                      {upload.status === 'success' && (
                        <p className="text-xs text-green-600">Subida completada</p>
                      )}
                    </div>
                    
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => removeUpload(index)}
                      className="flex-shrink-0"
                    >
                      <X className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}