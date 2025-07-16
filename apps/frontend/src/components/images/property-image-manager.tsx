'use client';

import { useState } from 'react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ImageUpload } from './image-upload';
import { ImageGallery } from './image-gallery';
import { 
  Upload, 
  Images, 
  Grid, 
  List, 
  Settings,
  Eye,
  BarChart3
} from 'lucide-react';

interface PropertyImageManagerProps {
  propertyId: string;
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

export function PropertyImageManager({ propertyId, className }: PropertyImageManagerProps) {
  const [activeTab, setActiveTab] = useState<'upload' | 'gallery'>('upload');
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [uploadedImages, setUploadedImages] = useState<UploadedImage[]>([]);

  const handleImageUploaded = (images: UploadedImage[]) => {
    setUploadedImages(prev => [...prev, ...images]);
    // Cambiar a galería después de subir la primera imagen
    if (uploadedImages.length === 0) {
      setActiveTab('gallery');
    }
  };

  return (
    <div className={className}>
      <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as 'upload' | 'gallery')}>
        <div className="flex items-center justify-between mb-6">
          <TabsList className="grid w-full max-w-md grid-cols-2">
            <TabsTrigger value="upload" className="flex items-center gap-2">
              <Upload className="h-4 w-4" />
              Subir Imágenes
            </TabsTrigger>
            <TabsTrigger value="gallery" className="flex items-center gap-2">
              <Images className="h-4 w-4" />
              Galería
              {uploadedImages.length > 0 && (
                <Badge variant="secondary" className="ml-1">
                  {uploadedImages.length}
                </Badge>
              )}
            </TabsTrigger>
          </TabsList>

          {activeTab === 'gallery' && (
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setViewMode(viewMode === 'grid' ? 'list' : 'grid')}
              >
                {viewMode === 'grid' ? (
                  <List className="h-4 w-4 mr-2" />
                ) : (
                  <Grid className="h-4 w-4 mr-2" />
                )}
                {viewMode === 'grid' ? 'Lista' : 'Grilla'}
              </Button>
            </div>
          )}
        </div>

        <TabsContent value="upload" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Upload className="h-5 w-5" />
                Subir Imágenes de la Propiedad
              </CardTitle>
            </CardHeader>
            <CardContent>
              <ImageUpload
                propertyId={propertyId}
                onImageUploaded={handleImageUploaded}
                maxImages={20}
                maxSizeInMB={10}
                acceptedTypes={['image/jpeg', 'image/png', 'image/webp']}
              />
            </CardContent>
          </Card>

          {/* Tips para upload */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-sm">
                <Settings className="h-4 w-4" />
                Consejos para mejores resultados
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                <div className="space-y-2">
                  <h4 className="font-medium text-gray-900">📸 Calidad de Imagen</h4>
                  <ul className="text-gray-600 space-y-1">
                    <li>• Resolución mínima: 1920x1080px</li>
                    <li>• Formato recomendado: JPEG o WebP</li>
                    <li>• Tamaño máximo: 10MB por imagen</li>
                  </ul>
                </div>
                <div className="space-y-2">
                  <h4 className="font-medium text-gray-900">🏠 Contenido Sugerido</h4>
                  <ul className="text-gray-600 space-y-1">
                    <li>• Sala principal, cocina, dormitorios</li>
                    <li>• Baños, exteriores, vista desde ventanas</li>
                    <li>• Detalles especiales y acabados</li>
                  </ul>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="gallery" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Images className="h-5 w-5" />
                Galería de Imágenes
              </CardTitle>
            </CardHeader>
            <CardContent>
              <ImageGallery
                propertyId={propertyId}
                viewMode={viewMode}
                editable={true}
                showActions={true}
              />
            </CardContent>
          </Card>

          {/* Image Statistics */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-sm">
                <BarChart3 className="h-4 w-4" />
                Estadísticas de Imágenes
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="text-center p-3 bg-blue-50 rounded-lg">
                  <div className="text-2xl font-bold text-blue-600">
                    {uploadedImages.length}
                  </div>
                  <div className="text-sm text-blue-600">Total</div>
                </div>
                <div className="text-center p-3 bg-green-50 rounded-lg">
                  <div className="text-2xl font-bold text-green-600">
                    {uploadedImages.filter(img => img.is_main).length}
                  </div>
                  <div className="text-sm text-green-600">Principal</div>
                </div>
                <div className="text-center p-3 bg-purple-50 rounded-lg">
                  <div className="text-2xl font-bold text-purple-600">
                    {Math.round(uploadedImages.reduce((acc, img) => acc + img.size, 0) / 1024 / 1024)}
                  </div>
                  <div className="text-sm text-purple-600">MB Total</div>
                </div>
                <div className="text-center p-3 bg-orange-50 rounded-lg">
                  <div className="text-2xl font-bold text-orange-600">
                    {uploadedImages.filter(img => img.type.includes('jpeg')).length}
                  </div>
                  <div className="text-sm text-orange-600">JPEG</div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}