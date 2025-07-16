'use client';

import { useState, useCallback } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  Star, 
  Trash2, 
  Edit, 
  Download, 
  Eye, 
  MoreHorizontal,
  ArrowUp,
  ArrowDown,
  X,
  ZoomIn,
  ZoomOut,
  RotateCcw,
  Share2
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { 
  Dialog, 
  DialogContent, 
  DialogHeader, 
  DialogTitle,
  DialogDescription,
  DialogFooter
} from '@/components/ui/dialog';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuSeparator,
} from '@/components/ui/dropdown-menu';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { apiClient } from '@/lib/api-client';
import { cn } from '@/lib/utils';
import { LoadingSpinner } from '@/components/ui/loading';

interface PropertyImage {
  id: string;
  url: string;
  thumbnail_url: string;
  filename: string;
  alt_text?: string;
  caption?: string;
  size: number;
  type: string;
  is_main: boolean;
  order: number;
  created_at: string;
  updated_at: string;
}

interface ImageGalleryProps {
  propertyId: string;
  className?: string;
  editable?: boolean;
  viewMode?: 'grid' | 'list';
  showActions?: boolean;
}

export function ImageGallery({
  propertyId,
  className,
  editable = true,
  viewMode = 'grid',
  showActions = true,
}: ImageGalleryProps) {
  const [selectedImage, setSelectedImage] = useState<PropertyImage | null>(null);
  const [isViewerOpen, setIsViewerOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [editingImage, setEditingImage] = useState<PropertyImage | null>(null);
  const [editForm, setEditForm] = useState({
    alt_text: '',
    caption: '',
  });

  const queryClient = useQueryClient();

  // Fetch images
  const { data: images, isLoading, error } = useQuery<PropertyImage[]>({
    queryKey: ['property-images', propertyId],
    queryFn: async () => {
      const response = await apiClient.get(`/properties/${propertyId}/images`);
      return response.data?.images || [];
    },
    enabled: !!propertyId,
  });

  // Set main image mutation
  const setMainImageMutation = useMutation({
    mutationFn: async (imageId: string) => {
      const response = await apiClient.post(`/properties/${propertyId}/images/main`, {
        image_id: imageId,
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['property-images', propertyId] });
    },
  });

  // Delete image mutation
  const deleteImageMutation = useMutation({
    mutationFn: async (imageId: string) => {
      const response = await apiClient.delete(`/images/${imageId}`);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['property-images', propertyId] });
      setIsDeleteDialogOpen(false);
      setSelectedImage(null);
    },
  });

  // Update image metadata mutation
  const updateImageMutation = useMutation({
    mutationFn: async ({ imageId, data }: { imageId: string; data: any }) => {
      const response = await apiClient.put(`/images/${imageId}/metadata`, data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['property-images', propertyId] });
      setIsEditDialogOpen(false);
      setEditingImage(null);
    },
  });

  // Reorder images mutation
  const reorderImagesMutation = useMutation({
    mutationFn: async (imageIds: string[]) => {
      const response = await apiClient.post(`/properties/${propertyId}/images/reorder`, {
        image_ids: imageIds,
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['property-images', propertyId] });
    },
  });

  const handleImageClick = useCallback((image: PropertyImage) => {
    setSelectedImage(image);
    setIsViewerOpen(true);
  }, []);

  const handleSetMainImage = useCallback((imageId: string) => {
    setMainImageMutation.mutate(imageId);
  }, [setMainImageMutation]);

  const handleDeleteImage = useCallback((image: PropertyImage) => {
    setSelectedImage(image);
    setIsDeleteDialogOpen(true);
  }, []);

  const handleEditImage = useCallback((image: PropertyImage) => {
    setEditingImage(image);
    setEditForm({
      alt_text: image.alt_text || '',
      caption: image.caption || '',
    });
    setIsEditDialogOpen(true);
  }, []);

  const handleMoveUp = useCallback((image: PropertyImage) => {
    if (!images) return;
    const currentIndex = images.findIndex(img => img.id === image.id);
    if (currentIndex <= 0) return;

    const newOrder = [...images];
    [newOrder[currentIndex], newOrder[currentIndex - 1]] = 
    [newOrder[currentIndex - 1], newOrder[currentIndex]];

    reorderImagesMutation.mutate(newOrder.map(img => img.id));
  }, [images, reorderImagesMutation]);

  const handleMoveDown = useCallback((image: PropertyImage) => {
    if (!images) return;
    const currentIndex = images.findIndex(img => img.id === image.id);
    if (currentIndex >= images.length - 1) return;

    const newOrder = [...images];
    [newOrder[currentIndex], newOrder[currentIndex + 1]] = 
    [newOrder[currentIndex + 1], newOrder[currentIndex]];

    reorderImagesMutation.mutate(newOrder.map(img => img.id));
  }, [images, reorderImagesMutation]);

  const handleDownload = useCallback((image: PropertyImage) => {
    const link = document.createElement('a');
    link.href = image.url;
    link.download = image.filename;
    link.click();
  }, []);

  const confirmDelete = useCallback(() => {
    if (!selectedImage) return;
    deleteImageMutation.mutate(selectedImage.id);
  }, [selectedImage, deleteImageMutation]);

  const handleEditSubmit = useCallback(() => {
    if (!editingImage) return;
    updateImageMutation.mutate({
      imageId: editingImage.id,
      data: editForm,
    });
  }, [editingImage, editForm, updateImageMutation]);

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  if (isLoading) {
    return (
      <div className={cn('grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4', className)}>
        {Array.from({ length: 4 }).map((_, i) => (
          <Card key={i} className="aspect-square">
            <CardContent className="p-0 h-full flex items-center justify-center">
              <LoadingSpinner size="sm" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  if (error || !images) {
    return (
      <Card className={className}>
        <CardContent className="p-8 text-center">
          <p className="text-red-600 mb-4">Error al cargar las imágenes</p>
          <Button onClick={() => window.location.reload()}>Reintentar</Button>
        </CardContent>
      </Card>
    );
  }

  if (images.length === 0) {
    return (
      <Card className={className}>
        <CardContent className="p-8 text-center">
          <div className="text-gray-400 mb-4">
            <div className="w-16 h-16 mx-auto mb-4 bg-gray-100 rounded-full flex items-center justify-center">
              <Eye className="w-8 h-8" />
            </div>
          </div>
          <p className="text-gray-500">No hay imágenes para esta propiedad</p>
          <p className="text-sm text-gray-400 mt-2">
            Sube imágenes para que los usuarios puedan ver la propiedad
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className={cn('space-y-4', className)}>
      {/* Image Grid */}
      <div className={cn(
        viewMode === 'grid' 
          ? 'grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4'
          : 'space-y-4'
      )}>
        {images.map((image, index) => (
          <Card 
            key={image.id} 
            className={cn(
              'group relative overflow-hidden cursor-pointer transition-all hover:shadow-lg',
              viewMode === 'grid' ? 'aspect-square' : 'aspect-video'
            )}
            onClick={() => handleImageClick(image)}
          >
            <div className="relative h-full">
              <img
                src={image.thumbnail_url || image.url}
                alt={image.alt_text || image.filename}
                className="w-full h-full object-cover"
              />
              
              {/* Overlay */}
              <div className="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-20 transition-all duration-200 flex items-center justify-center">
                <div className="opacity-0 group-hover:opacity-100 transition-opacity">
                  <ZoomIn className="w-6 h-6 text-white" />
                </div>
              </div>

              {/* Main Image Badge */}
              {image.is_main && (
                <Badge className="absolute top-2 left-2 bg-yellow-500 text-white">
                  <Star className="w-3 h-3 mr-1" />
                  Principal
                </Badge>
              )}

              {/* Order Badge */}
              <Badge 
                variant="secondary" 
                className="absolute top-2 right-2 bg-black bg-opacity-50 text-white"
              >
                {index + 1}
              </Badge>

              {/* Actions */}
              {editable && showActions && (
                <div className="absolute bottom-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="secondary" size="sm">
                        <MoreHorizontal className="w-4 h-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      {!image.is_main && (
                        <DropdownMenuItem onClick={(e) => {
                          e.stopPropagation();
                          handleSetMainImage(image.id);
                        }}>
                          <Star className="w-4 h-4 mr-2" />
                          Establecer como principal
                        </DropdownMenuItem>
                      )}
                      <DropdownMenuItem onClick={(e) => {
                        e.stopPropagation();
                        handleEditImage(image);
                      }}>
                        <Edit className="w-4 h-4 mr-2" />
                        Editar
                      </DropdownMenuItem>
                      <DropdownMenuItem onClick={(e) => {
                        e.stopPropagation();
                        handleDownload(image);
                      }}>
                        <Download className="w-4 h-4 mr-2" />
                        Descargar
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem onClick={(e) => {
                        e.stopPropagation();
                        handleMoveUp(image);
                      }} disabled={index === 0}>
                        <ArrowUp className="w-4 h-4 mr-2" />
                        Mover arriba
                      </DropdownMenuItem>
                      <DropdownMenuItem onClick={(e) => {
                        e.stopPropagation();
                        handleMoveDown(image);
                      }} disabled={index === images.length - 1}>
                        <ArrowDown className="w-4 h-4 mr-2" />
                        Mover abajo
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem 
                        onClick={(e) => {
                          e.stopPropagation();
                          handleDeleteImage(image);
                        }}
                        className="text-red-600"
                      >
                        <Trash2 className="w-4 h-4 mr-2" />
                        Eliminar
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              )}
            </div>

            {/* Image Info */}
            <div className="absolute bottom-0 left-0 right-0 bg-black bg-opacity-50 text-white p-2 opacity-0 group-hover:opacity-100 transition-opacity">
              <p className="text-xs truncate">{image.filename}</p>
              <p className="text-xs text-gray-300">{formatFileSize(image.size)}</p>
            </div>
          </Card>
        ))}
      </div>

      {/* Image Viewer Modal */}
      <Dialog open={isViewerOpen} onOpenChange={setIsViewerOpen}>
        <DialogContent className="max-w-4xl">
          <DialogHeader>
            <DialogTitle>
              {selectedImage?.filename}
            </DialogTitle>
          </DialogHeader>
          
          {selectedImage && (
            <div className="space-y-4">
              <div className="relative">
                <img
                  src={selectedImage.url}
                  alt={selectedImage.alt_text || selectedImage.filename}
                  className="w-full h-auto max-h-96 object-contain rounded-lg"
                />
              </div>
              
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <p className="font-medium">Tamaño:</p>
                  <p className="text-gray-600">{formatFileSize(selectedImage.size)}</p>
                </div>
                <div>
                  <p className="font-medium">Tipo:</p>
                  <p className="text-gray-600">{selectedImage.type}</p>
                </div>
                {selectedImage.caption && (
                  <div className="col-span-2">
                    <p className="font-medium">Descripción:</p>
                    <p className="text-gray-600">{selectedImage.caption}</p>
                  </div>
                )}
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Editar Imagen</DialogTitle>
            <DialogDescription>
              Actualiza la información de la imagen
            </DialogDescription>
          </DialogHeader>
          
          <div className="space-y-4">
            <div>
              <Label htmlFor="alt_text">Texto alternativo</Label>
              <Input
                id="alt_text"
                value={editForm.alt_text}
                onChange={(e) => setEditForm(prev => ({ ...prev, alt_text: e.target.value }))}
                placeholder="Descripción para accesibilidad"
              />
            </div>
            
            <div>
              <Label htmlFor="caption">Descripción</Label>
              <Textarea
                id="caption"
                value={editForm.caption}
                onChange={(e) => setEditForm(prev => ({ ...prev, caption: e.target.value }))}
                placeholder="Descripción de la imagen"
                rows={3}
              />
            </div>
          </div>
          
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsEditDialogOpen(false)}>
              Cancelar
            </Button>
            <Button onClick={handleEditSubmit} disabled={updateImageMutation.isPending}>
              {updateImageMutation.isPending ? 'Guardando...' : 'Guardar'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Eliminar Imagen</DialogTitle>
            <DialogDescription>
              ¿Estás seguro de que quieres eliminar esta imagen? Esta acción no se puede deshacer.
            </DialogDescription>
          </DialogHeader>
          
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsDeleteDialogOpen(false)}>
              Cancelar
            </Button>
            <Button 
              variant="destructive" 
              onClick={confirmDelete}
              disabled={deleteImageMutation.isPending}
            >
              {deleteImageMutation.isPending ? 'Eliminando...' : 'Eliminar'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}