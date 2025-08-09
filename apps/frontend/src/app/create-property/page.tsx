'use client';

import { ModernPropertyForm2025 } from '@/components/forms/modern-property-form-2025';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Building2, ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';
import Link from 'next/link';

/**
 * Create Property Page - Simplified and Focused
 * 
 * Features:
 * - No authentication (for development focus)
 * - No dashboard complexity
 * - Just the property creation form
 * - Clean and minimal design
 * - Direct testing of Property CRUD functionality
 */
export default function CreatePropertyPage() {
  const handleSuccess = () => {
    console.log('🎉 Property created successfully!');
    // You can add redirect or success message here
  };

  const handleCancel = () => {
    console.log('❌ Property creation cancelled');
    // You can add navigation logic here
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center gap-4 mb-4">
            <Link href="/dashboard">
              <Button variant="outline" size="sm">
                <ArrowLeft className="w-4 h-4 mr-2" />
                Volver
              </Button>
            </Link>
            <div className="flex items-center gap-2">
              <Building2 className="w-8 h-8 text-blue-600" />
              <h1 className="text-3xl font-bold text-gray-900">Crear Nueva Propiedad</h1>
            </div>
          </div>
          
          <div className="flex items-center gap-2 mb-4">
            <Badge variant="secondary" className="text-sm">
              React 19 + Server Actions
            </Badge>
            <Badge variant="outline" className="text-sm">
              63 Campos Completos
            </Badge>
            <Badge variant="outline" className="text-sm">
              Progressive Enhancement
            </Badge>
          </div>
          
          <p className="text-gray-600 max-w-2xl">
            Formulario modernizado con React 19 Server Actions, validación Zod y 
            progressive enhancement. Funciona con y sin JavaScript habilitado.
          </p>
        </div>

        {/* Form Card */}
        <Card className="shadow-lg">
          <CardHeader className="bg-blue-50 border-b">
            <CardTitle className="flex items-center gap-2 text-xl">
              🏠 Información de la Propiedad
              <Badge variant="secondary" className="text-xs ml-auto">
                Formulario Modernizado 2025
              </Badge>
            </CardTitle>
          </CardHeader>
          <CardContent className="p-6">
            <ModernPropertyForm2025
              onSuccess={handleSuccess}
              onCancel={handleCancel}
            />
          </CardContent>
        </Card>

        {/* Footer Info */}
        <div className="mt-8 text-center text-sm text-gray-500">
          <p>
            🚀 <strong>Property CRUD Core:</strong> Formulario conectado directamente con 
            backend Go (localhost:8080) → PostgreSQL (puerto 5433)
          </p>
          <p className="mt-2">
            📊 <strong>63 campos totales:</strong> Información básica, ubicación, características, 
            amenidades, precios adicionales y contacto
          </p>
        </div>
      </div>
    </div>
  );
}