'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
// import { useAuthStore } from '@/store/auth';

export default function Home() {
  // NO AUTH MODE: Direct access to dashboard
  // const { isAuthenticated } = useAuthStore();
  const router = useRouter();
  const [hasNavigated, setHasNavigated] = useState(false);

  useEffect(() => {
    if (!hasNavigated) {
      console.log('🏠 Home Page - NO AUTH MODE - Going directly to dashboard');
      setHasNavigated(true);
      
      // Redirect directly to dashboard (no auth check)
      router.replace('/dashboard');
    }
  }, [router, hasNavigated]);

  // Show temporary page while redirecting
  return (
    <div className="flex items-center justify-center min-h-screen bg-gradient-to-br from-blue-50 to-gray-100">
      <div className="text-center space-y-6 p-8 bg-white rounded-lg shadow-lg">
        <div className="space-y-4">
          <h1 className="text-3xl font-bold text-gray-900">
            🏠 Sistema Inmobiliario Ecuador
          </h1>
          <p className="text-gray-600">
            Plataforma de gestión de propiedades inmobiliarias
          </p>
          <div className="bg-yellow-50 border border-yellow-200 rounded-md p-3">
            <p className="text-sm text-yellow-800">
              <strong>Modo Desarrollo:</strong> Autenticación desactivada
            </p>
          </div>
        </div>
        
        <div className="space-y-3">
          <p className="text-gray-500">Redirigiendo al dashboard...</p>
          <div className="flex justify-center space-x-3">
            <Button 
              onClick={() => router.push('/dashboard')}
              className="bg-blue-600 hover:bg-blue-700"
            >
              🏠 Dashboard
            </Button>
            <Button 
              onClick={() => router.push('/properties')}
              variant="outline"
            >
              📋 Propiedades
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}