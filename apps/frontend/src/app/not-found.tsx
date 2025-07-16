import Link from 'next/link';
import { FileQuestion, Home, ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

export default function NotFound() {
  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-gray-50">
      <Card className="w-full max-w-md text-center">
        <CardHeader>
          <div className="mx-auto mb-4 w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center">
            <FileQuestion className="w-8 h-8 text-blue-600" />
          </div>
          <CardTitle className="text-2xl">Página no encontrada</CardTitle>
          <CardDescription className="text-base">
            La página que buscas no existe o ha sido movida.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="bg-blue-50 border border-blue-200 rounded-md p-3">
            <p className="text-sm text-blue-800">
              <strong>Error 404:</strong> La URL solicitada no fue encontrada en el servidor.
            </p>
          </div>
          
          <div className="flex flex-col sm:flex-row gap-3">
            <Button asChild className="flex-1">
              <Link href="/dashboard">
                <Home className="w-4 h-4 mr-2" />
                Ir al Dashboard
              </Link>
            </Button>
            <Button asChild variant="outline" className="flex-1">
              <Link href="javascript:history.back()">
                <ArrowLeft className="w-4 h-4 mr-2" />
                Volver
              </Link>
            </Button>
          </div>
          
          <div className="pt-4 border-t">
            <p className="text-sm text-gray-600 mb-2">
              Páginas populares:
            </p>
            <div className="space-y-1">
              <Link 
                href="/properties" 
                className="block text-sm text-blue-600 hover:underline"
              >
                → Ver propiedades
              </Link>
              <Link 
                href="/users" 
                className="block text-sm text-blue-600 hover:underline"
              >
                → Gestionar usuarios
              </Link>
              <Link 
                href="/statistics" 
                className="block text-sm text-blue-600 hover:underline"
              >
                → Ver estadísticas
              </Link>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}