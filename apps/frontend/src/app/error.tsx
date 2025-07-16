'use client';

import { useEffect } from 'react';
import { AlertTriangle, RefreshCw, Home } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

interface ErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
}

export default function Error({ error, reset }: ErrorProps) {
  useEffect(() => {
    // Log the error to an error reporting service
    console.error('Application error:', error);
  }, [error]);

  const goHome = () => {
    window.location.href = '/dashboard';
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-gray-50">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="mx-auto mb-4 w-16 h-16 bg-red-100 rounded-full flex items-center justify-center">
            <AlertTriangle className="w-8 h-8 text-red-600" />
          </div>
          <CardTitle className="text-2xl">¡Ops! Algo salió mal</CardTitle>
          <CardDescription className="text-base">
            Se produjo un error inesperado en la aplicación. 
            Por favor, inténtalo de nuevo o regresa al inicio.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {process.env.NODE_ENV === 'development' && (
            <div className="bg-red-50 border border-red-200 rounded-md p-3">
              <p className="text-sm font-semibold text-red-800 mb-1">Error técnico:</p>
              <p className="text-xs font-mono text-red-600 break-all">
                {error.message}
              </p>
              {error.digest && (
                <p className="text-xs text-red-500 mt-1">
                  Digest: {error.digest}
                </p>
              )}
            </div>
          )}
          
          <div className="flex flex-col sm:flex-row gap-3">
            <Button 
              onClick={reset} 
              className="flex-1"
              variant="default"
            >
              <RefreshCw className="w-4 h-4 mr-2" />
              Reintentar
            </Button>
            <Button 
              onClick={goHome} 
              variant="outline"
              className="flex-1"
            >
              <Home className="w-4 h-4 mr-2" />
              Ir al Dashboard
            </Button>
          </div>
          
          <p className="text-xs text-center text-gray-500">
            Si el problema persiste, contacta al administrador del sistema.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}