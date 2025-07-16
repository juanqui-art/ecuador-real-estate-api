'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Building } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { LoginForm } from '@/components/forms/login-form';
import { useAuthStore } from '@/store/auth';

export default function LoginPage() {
  const [hasRedirected, setHasRedirected] = useState(false);
  const router = useRouter();
  const { isAuthenticated } = useAuthStore();

  // Redirect if already authenticated - but only once
  useEffect(() => {
    console.log('📄 Login Page - Checking auth state:', { isAuthenticated, hasRedirected });
    
    if (isAuthenticated && !hasRedirected) {
      console.log('📄 Login Page - User is authenticated, preparing redirect');
      setHasRedirected(true);
      
      // Get redirect parameter
      const urlParams = new URLSearchParams(window.location.search);
      const redirectTo = urlParams.get('redirect') || '/dashboard';
      
      console.log('📄 Login Page - Redirecting to:', redirectTo);
      
      // Use Next.js router for navigation
      router.replace(redirectTo);
    }
  }, [isAuthenticated, hasRedirected, router]);

  const handleSuccess = () => {
    console.log('📄 Login Page - Login successful');
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="w-full max-w-md space-y-6">
        {/* Header */}
        <Card>
          <CardHeader className="space-y-1">
            <div className="flex items-center justify-center mb-4">
              <Building className="h-10 w-10 text-blue-600" />
            </div>
            <CardTitle className="text-2xl text-center">InmoEcuador</CardTitle>
            <CardDescription className="text-center">
              Sistema de gestión inmobiliaria para Ecuador
            </CardDescription>
          </CardHeader>
        </Card>

        {/* Login Form - TanStack Form */}
        <LoginForm
          onSuccess={handleSuccess}
          redirectTo="/dashboard"
        />

        {/* Test Credentials */}
        <Card>
          <CardContent className="pt-6">
            <div className="text-center text-sm text-gray-600">
              <p className="font-medium mb-2">Credenciales de prueba:</p>
              <p>Email: test@example.com</p>
              <p>Password: test123</p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}