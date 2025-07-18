'use client';

import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/store/auth';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { AlertCircle, Home, User, Eye, BarChart3 } from 'lucide-react';
import Link from 'next/link';

export default function UnauthorizedPage() {
  const router = useRouter();
  const { user, isAuthenticated } = useAuthStore();

  const handleGoBack = () => {
    router.back();
  };

  const handleGoHome = () => {
    router.push('/');
  };

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="mx-auto mb-4 w-16 h-16 bg-red-100 rounded-full flex items-center justify-center">
            <AlertCircle className="w-8 h-8 text-red-600" />
          </div>
          <CardTitle className="text-2xl font-bold text-gray-900">
            Página no encontrada
          </CardTitle>
        </CardHeader>
        
        <CardContent className="space-y-4">
          <div className="text-center">
            <p className="text-gray-600 mb-2">
              La página que buscas no existe o ha sido movida.
            </p>
            <p className="text-sm text-gray-500">
              Error 404: La URL solicitada no fue encontrada en el servidor.
            </p>
          </div>

          {isAuthenticated && user && (
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-3">
              <p className="text-sm text-blue-800">
                <strong>Usuario actual:</strong> {user.email} ({user.role})
              </p>
            </div>
          )}

          <div className="flex flex-col gap-2">
            <Button onClick={handleGoHome} className="w-full">
              <Home className="w-4 h-4 mr-2" />
              Ir al Dashboard
            </Button>
            <Button variant="outline" onClick={handleGoBack} className="w-full">
              Volver
            </Button>
          </div>

          <div className="border-t pt-4">
            <p className="text-sm font-medium text-gray-700 mb-2">
              Páginas populares:
            </p>
            <div className="space-y-2">
              <Link 
                href="/properties" 
                className="flex items-center text-sm text-blue-600 hover:text-blue-800"
              >
                <Eye className="w-4 h-4 mr-2" />
                → Ver propiedades
              </Link>
              {isAuthenticated && user?.role === 'admin' && (
                <>
                  <Link 
                    href="/users" 
                    className="flex items-center text-sm text-blue-600 hover:text-blue-800"
                  >
                    <User className="w-4 h-4 mr-2" />
                    → Gestionar usuarios
                  </Link>
                  <Link 
                    href="/analytics" 
                    className="flex items-center text-sm text-blue-600 hover:text-blue-800"
                  >
                    <BarChart3 className="w-4 h-4 mr-2" />
                    → Ver estadísticas
                  </Link>
                </>
              )}
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}