'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { CheckCircle, XCircle, Info } from 'lucide-react';
import { imageProcessor } from '@/lib/image-processor';

interface ProcessingStats {
  formatsSupported: string[];
  canvasInitialized: boolean;
  browserSupport: {
    webp: boolean;
    avif: boolean;
    progressive: boolean;
  };
}

export function ImageProcessorStats() {
  const [stats, setStats] = useState<ProcessingStats | null>(null);

  useEffect(() => {
    // Get stats after component mounts (client-side only)
    const processingStats = imageProcessor.getProcessingStats();
    setStats(processingStats);
  }, []);

  if (!stats) {
    return null;
  }

  const formatBadges = stats.formatsSupported.map(format => {
    const colors = {
      webp: 'bg-green-100 text-green-800',
      avif: 'bg-blue-100 text-blue-800',
      jpeg: 'bg-yellow-100 text-yellow-800',
      png: 'bg-purple-100 text-purple-800',
    };

    return (
      <Badge 
        key={format} 
        variant="secondary" 
        className={colors[format as keyof typeof colors] || 'bg-gray-100 text-gray-800'}
      >
        {format.toUpperCase()}
      </Badge>
    );
  });

  return (
    <Card className="border-blue-200 bg-blue-50">
      <CardHeader className="pb-3">
        <CardTitle className="text-sm font-medium text-blue-900 flex items-center gap-2">
          <Info className="h-4 w-4" />
          Procesamiento de Imágenes - Estado del Sistema
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <div>
          <p className="text-xs text-blue-700 font-medium mb-2">Formatos Soportados:</p>
          <div className="flex flex-wrap gap-2">
            {formatBadges}
          </div>
        </div>
        
        <div className="grid grid-cols-2 gap-4 text-xs">
          <div className="flex items-center gap-2">
            {stats.browserSupport.webp ? (
              <CheckCircle className="h-3 w-3 text-green-600" />
            ) : (
              <XCircle className="h-3 w-3 text-red-600" />
            )}
            <span className="text-blue-800">WebP</span>
          </div>
          
          <div className="flex items-center gap-2">
            {stats.browserSupport.avif ? (
              <CheckCircle className="h-3 w-3 text-green-600" />
            ) : (
              <XCircle className="h-3 w-3 text-red-600" />
            )}
            <span className="text-blue-800">AVIF</span>
          </div>
          
          <div className="flex items-center gap-2">
            {stats.canvasInitialized ? (
              <CheckCircle className="h-3 w-3 text-green-600" />
            ) : (
              <XCircle className="h-3 w-3 text-red-600" />
            )}
            <span className="text-blue-800">Canvas</span>
          </div>
          
          <div className="flex items-center gap-2">
            {stats.browserSupport.progressive ? (
              <CheckCircle className="h-3 w-3 text-green-600" />
            ) : (
              <XCircle className="h-3 w-3 text-red-600" />
            )}
            <span className="text-blue-800">Progressive</span>
          </div>
        </div>
        
        <div className="pt-2 border-t border-blue-200">
          <p className="text-xs text-blue-600">
            ✨ Optimización automática activa - Formato inteligente según navegador
          </p>
        </div>
      </CardContent>
    </Card>
  );
}