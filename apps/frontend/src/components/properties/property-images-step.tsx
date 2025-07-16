'use client';

import { useState } from 'react';
import { ArrowLeft, ArrowRight, Check, SkipForward } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { PropertyImageManager } from '@/components/images/property-image-manager';

interface PropertyImagesStepProps {
  propertyId: string;
  propertyTitle: string;
  onComplete: () => void;
  onSkip: () => void;
  onBack?: () => void;
}

export function PropertyImagesStep({
  propertyId,
  propertyTitle,
  onComplete,
  onSkip,
  onBack,
}: PropertyImagesStepProps) {
  const [uploadedCount, setUploadedCount] = useState(0);

  const handleImageUploaded = (images: any[]) => {
    setUploadedCount(prev => prev + images.length);
  };

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      {/* Header */}
      <div className="text-center space-y-2">
        <Badge variant="secondary" className="mb-4">
          Paso 2 de 2
        </Badge>
        <h1 className="text-3xl font-bold text-gray-900">
          Agregar Imágenes a tu Propiedad
        </h1>
        <p className="text-gray-600 max-w-2xl mx-auto">
          Las imágenes son cruciales para atraer compradores. Sube fotos de alta calidad 
          que muestren las mejores características de "{propertyTitle}".
        </p>
      </div>

      {/* Property Info */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            <span>🏠 {propertyTitle}</span>
            <Badge variant="outline">
              ID: {propertyId}
            </Badge>
          </CardTitle>
        </CardHeader>
      </Card>

      {/* Image Manager */}
      <PropertyImageManager 
        propertyId={propertyId}
        className="space-y-6"
      />

      {/* Progress & Actions */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              {onBack && (
                <Button variant="outline" onClick={onBack}>
                  <ArrowLeft className="h-4 w-4 mr-2" />
                  Volver
                </Button>
              )}
              
              <div className="flex items-center gap-2">
                {uploadedCount > 0 && (
                  <Badge variant="secondary">
                    {uploadedCount} imagen{uploadedCount !== 1 ? 's' : ''} subida{uploadedCount !== 1 ? 's' : ''}
                  </Badge>
                )}
              </div>
            </div>

            <div className="flex items-center gap-3">
              <Button variant="outline" onClick={onSkip}>
                <SkipForward className="h-4 w-4 mr-2" />
                Saltar por ahora
              </Button>
              
              <Button onClick={onComplete} className="min-w-[120px]">
                <Check className="h-4 w-4 mr-2" />
                {uploadedCount > 0 ? 'Finalizar' : 'Continuar'}
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Tips */}
      <Card className="bg-blue-50 border-blue-200">
        <CardContent className="pt-6">
          <div className="flex items-start gap-3">
            <div className="text-blue-600">💡</div>
            <div>
              <h3 className="font-medium text-blue-900 mb-2">
                Consejos para mejorar tus ventas
              </h3>
              <ul className="text-sm text-blue-800 space-y-1">
                <li>• Las propiedades con imágenes reciben 3x más consultas</li>
                <li>• Sube al menos 5-10 imágenes para mejores resultados</li>
                <li>• Incluye fotos de todas las habitaciones principales</li>
                <li>• Puedes agregar más imágenes después desde la página de edición</li>
              </ul>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}