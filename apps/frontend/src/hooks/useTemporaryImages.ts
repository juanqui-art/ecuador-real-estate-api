import { useState, useCallback } from 'react';
import { useMutation } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import type { TemporaryImageFile } from '@/components/images/temporary-image-upload';

interface UploadImageResponse {
  id: string;
  filename: string;
  url: string;
  thumbnail_url: string;
  metadata: {
    width: number;
    height: number;
    size: number;
    format: string;
  };
}

interface UseTemporaryImagesProps {
  onUploadComplete?: (uploadedImages: UploadImageResponse[]) => void;
  onUploadError?: (error: Error) => void;
}

export function useTemporaryImages({ 
  onUploadComplete, 
  onUploadError 
}: UseTemporaryImagesProps = {}) {
  const [temporaryImages, setTemporaryImages] = useState<TemporaryImageFile[]>([]);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [isUploading, setIsUploading] = useState(false);

  // Mutation for uploading images to server
  const uploadImagesMutation = useMutation({
    mutationFn: async ({ propertyId, images }: { propertyId: string; images: TemporaryImageFile[] }) => {
      const uploadedImages: UploadImageResponse[] = [];
      setIsUploading(true);
      setUploadProgress(0);

      try {
        // Upload images sequentially to avoid overwhelming the server
        for (let i = 0; i < images.length; i++) {
          const image = images[i];
          
          // Skip if image has errors or is still processing
          if (image.error || image.isProcessing) {
            continue;
          }

          const formData = new FormData();
          
          // Use processed image if available, otherwise use original
          const fileToUpload = image.processed?.file || image.file;
          formData.append('image', fileToUpload);
          formData.append('property_id', propertyId);
          formData.append('is_main', image.isMain.toString());
          formData.append('order', i.toString());

          try {
            const response = await apiClient.upload('/images', formData);
            
            if (response.data.success) {
              uploadedImages.push(response.data.data);
            }
          } catch (error) {
            console.error(`Error uploading image ${i + 1}:`, error);
            // Continue with other images even if one fails
          }

          // Update progress
          setUploadProgress(((i + 1) / images.length) * 100);
        }

        // Set main image if there's one marked as main
        const mainImage = images.find(img => img.isMain);
        if (mainImage && uploadedImages.length > 0) {
          const mainUploadedImage = uploadedImages.find(uploaded => 
            uploaded.filename === (mainImage.processed?.file.name || mainImage.file.name)
          );
          
          if (mainUploadedImage) {
            try {
              await apiClient.post(`/properties/${propertyId}/images/main`, {
                image_id: mainUploadedImage.id
              });
            } catch (error) {
              console.error('Error setting main image:', error);
            }
          }
        }

        return uploadedImages;
      } finally {
        setIsUploading(false);
        setUploadProgress(0);
      }
    },
    onSuccess: (uploadedImages) => {
      // Clear temporary images after successful upload
      clearTemporaryImages();
      onUploadComplete?.(uploadedImages);
    },
    onError: (error: Error) => {
      console.error('Error uploading images:', error);
      onUploadError?.(error);
    },
  });

  const clearTemporaryImages = useCallback(() => {
    // Clean up object URLs to prevent memory leaks
    temporaryImages.forEach(image => {
      URL.revokeObjectURL(image.preview);
      if (image.processed) {
        URL.revokeObjectURL(image.processed.preview);
      }
    });
    setTemporaryImages([]);
  }, [temporaryImages]);

  const addTemporaryImages = useCallback((newImages: TemporaryImageFile[]) => {
    setTemporaryImages(prev => [...prev, ...newImages]);
  }, []);

  const removeTemporaryImage = useCallback((imageId: string) => {
    setTemporaryImages(prev => {
      const imageToRemove = prev.find(img => img.id === imageId);
      if (imageToRemove) {
        URL.revokeObjectURL(imageToRemove.preview);
        if (imageToRemove.processed) {
          URL.revokeObjectURL(imageToRemove.processed.preview);
        }
      }
      
      const updatedImages = prev.filter(img => img.id !== imageId);
      
      // If removed image was main, make first image main
      if (imageToRemove?.isMain && updatedImages.length > 0) {
        updatedImages[0].isMain = true;
      }
      
      return updatedImages;
    });
  }, []);

  const setMainTemporaryImage = useCallback((imageId: string) => {
    setTemporaryImages(prev => prev.map(img => ({
      ...img,
      isMain: img.id === imageId
    })));
  }, []);

  const reorderTemporaryImages = useCallback((fromIndex: number, toIndex: number) => {
    setTemporaryImages(prev => {
      const reorderedImages = [...prev];
      const [movedImage] = reorderedImages.splice(fromIndex, 1);
      reorderedImages.splice(toIndex, 0, movedImage);
      return reorderedImages;
    });
  }, []);

  const uploadImages = useCallback(async (propertyId: string) => {
    if (temporaryImages.length === 0) {
      return [];
    }

    // Filter out images with errors or still processing
    const validImages = temporaryImages.filter(img => !img.error && !img.isProcessing);
    
    if (validImages.length === 0) {
      throw new Error('No hay imágenes válidas para subir');
    }

    return uploadImagesMutation.mutateAsync({ propertyId, images: validImages });
  }, [temporaryImages, uploadImagesMutation]);

  const hasValidImages = temporaryImages.some(img => !img.error && !img.isProcessing);
  const processingCount = temporaryImages.filter(img => img.isProcessing).length;
  const errorCount = temporaryImages.filter(img => img.error).length;

  return {
    temporaryImages,
    setTemporaryImages,
    addTemporaryImages,
    removeTemporaryImage,
    setMainTemporaryImage,
    reorderTemporaryImages,
    clearTemporaryImages,
    uploadImages,
    uploadProgress,
    isUploading,
    hasValidImages,
    processingCount,
    errorCount,
    isUploadPending: uploadImagesMutation.isPending,
  };
}