'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/store/auth';

export default function Home() {
  const { isAuthenticated } = useAuthStore();
  const router = useRouter();
  const [hasNavigated, setHasNavigated] = useState(false);

  useEffect(() => {
    if (!hasNavigated) {
      console.log('🏠 Home Page - Checking auth and navigating');
      setHasNavigated(true);
      
      if (isAuthenticated) {
        console.log('🏠 Home Page - User authenticated, going to dashboard');
        router.replace('/dashboard');
      } else {
        console.log('🏠 Home Page - User not authenticated, going to login');
        router.replace('/login');
      }
    }
  }, [isAuthenticated, router, hasNavigated]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center">
        <h1 className="text-2xl font-bold text-gray-900 mb-4">
          Sistema Inmobiliario Ecuador
        </h1>
        <p className="text-gray-600">Cargando...</p>
      </div>
    </div>
  );
}