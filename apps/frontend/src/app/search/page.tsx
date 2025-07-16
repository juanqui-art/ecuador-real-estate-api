'use client';

import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { 
  Search, 
  Filter, 
  Grid, 
  List, 
  MapPin, 
  TrendingUp,
  Clock,
  Bookmark,
  Share2,
  Heart
} from 'lucide-react';
import { DashboardLayout } from '@/components/layout/dashboard-layout';
import { PublicSearch } from '@/components/search/public-search';
import { formatPrice } from '@/lib/utils';

interface SearchResult {
  id: string;
  title: string;
  description: string;
  price: number;
  type: string;
  province: string;
  city: string;
  status: string;
  bedrooms: number;
  bathrooms: number;
  area_m2: number;
  created_at: string;
  images_count: number;
  main_image_url?: string;
}

export default function SearchPage() {
  const [selectedResults, setSelectedResults] = useState<SearchResult[]>([]);
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [savedSearches, setSavedSearches] = useState<string[]>([]);

  const handleResultSelect = (result: SearchResult) => {
    // Add to recent results
    setSelectedResults(prev => {
      const filtered = prev.filter(r => r.id !== result.id);
      return [result, ...filtered].slice(0, 10); // Keep last 10 results
    });
    
    // Navigate to property detail page
    window.location.href = `/properties/${result.id}`;
  };

  const handleSaveSearch = (query: string) => {
    if (query.trim() && !savedSearches.includes(query)) {
      setSavedSearches(prev => [query, ...prev].slice(0, 5));
    }
  };

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'house': return 'bg-blue-100 text-blue-800';
      case 'apartment': return 'bg-green-100 text-green-800';
      case 'land': return 'bg-yellow-100 text-yellow-800';
      case 'commercial': return 'bg-purple-100 text-purple-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'available': return 'bg-green-100 text-green-800';
      case 'sold': return 'bg-red-100 text-red-800';
      case 'rented': return 'bg-blue-100 text-blue-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  const getTypeLabel = (type: string) => {
    switch (type) {
      case 'house': return 'Casa';
      case 'apartment': return 'Apartamento';
      case 'land': return 'Terreno';
      case 'commercial': return 'Comercial';
      default: return type;
    }
  };

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'available': return 'Disponible';
      case 'sold': return 'Vendida';
      case 'rented': return 'Rentada';
      default: return status;
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('es-EC', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Búsqueda de Propiedades</h1>
            <p className="text-gray-600 mt-1">
              Encuentra la propiedad perfecta con nuestra búsqueda en tiempo real
            </p>
          </div>
          
          <div className="flex items-center gap-2">
            <Button
              variant={viewMode === 'grid' ? 'default' : 'outline'}
              size="sm"
              onClick={() => setViewMode('grid')}
            >
              <Grid className="h-4 w-4" />
            </Button>
            <Button
              variant={viewMode === 'list' ? 'default' : 'outline'}
              size="sm"
              onClick={() => setViewMode('list')}
            >
              <List className="h-4 w-4" />
            </Button>
          </div>
        </div>

        {/* Advanced Search */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Search className="h-5 w-5" />
              Búsqueda Avanzada
            </CardTitle>
          </CardHeader>
          <CardContent>
            <PublicSearch
              placeholder="Buscar por título, ubicación, características..."
              onResultSelect={handleResultSelect}
              className="w-full"
            />
          </CardContent>
        </Card>

        {/* Search Results and Saved Searches */}
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Main Results */}
          <div className="lg:col-span-3">
            <Tabs defaultValue="results" className="w-full">
              <TabsList className="grid w-full grid-cols-3">
                <TabsTrigger value="results">
                  Resultados ({selectedResults.length})
                </TabsTrigger>
                <TabsTrigger value="trending">
                  Tendencias
                </TabsTrigger>
                <TabsTrigger value="recent">
                  Recientes
                </TabsTrigger>
              </TabsList>

              <TabsContent value="results" className="space-y-4">
                {selectedResults.length === 0 ? (
                  <Card>
                    <CardContent className="p-8 text-center">
                      <Search className="h-12 w-12 mx-auto text-gray-400 mb-4" />
                      <h3 className="text-lg font-medium text-gray-900 mb-2">
                        Comienza tu búsqueda
                      </h3>
                      <p className="text-gray-600">
                        Usa el buscador arriba para encontrar propiedades
                      </p>
                    </CardContent>
                  </Card>
                ) : (
                  <div className={`grid gap-4 ${
                    viewMode === 'grid' 
                      ? 'grid-cols-1 md:grid-cols-2' 
                      : 'grid-cols-1'
                  }`}>
                    {selectedResults.map((result) => (
                      <Card key={result.id} className="hover:shadow-lg transition-shadow">
                        <CardContent className="p-4">
                          <div className={`flex gap-4 ${
                            viewMode === 'grid' ? 'flex-col' : 'flex-row'
                          }`}>
                            {/* Image */}
                            <div className={`${
                              viewMode === 'grid' ? 'w-full h-48' : 'w-24 h-24'
                            } bg-gray-200 rounded-lg flex items-center justify-center flex-shrink-0`}>
                              {result.main_image_url ? (
                                <img 
                                  src={result.main_image_url} 
                                  alt={result.title}
                                  className="w-full h-full object-cover rounded-lg"
                                />
                              ) : (
                                <MapPin className="h-8 w-8 text-gray-400" />
                              )}
                            </div>

                            {/* Content */}
                            <div className="flex-1 min-w-0">
                              <div className="flex items-start justify-between mb-2">
                                <h3 className="font-semibold text-gray-900 truncate">
                                  {result.title}
                                </h3>
                                <div className="flex gap-1 ml-2">
                                  <Button variant="ghost" size="sm">
                                    <Heart className="h-4 w-4" />
                                  </Button>
                                  <Button variant="ghost" size="sm">
                                    <Share2 className="h-4 w-4" />
                                  </Button>
                                </div>
                              </div>

                              <div className="flex items-center gap-2 mb-2">
                                <Badge className={getTypeColor(result.type)}>
                                  {getTypeLabel(result.type)}
                                </Badge>
                                <Badge className={getStatusColor(result.status)}>
                                  {getStatusLabel(result.status)}
                                </Badge>
                              </div>

                              <div className="flex items-center gap-4 text-sm text-gray-600 mb-2">
                                <div className="flex items-center gap-1">
                                  <MapPin className="h-3 w-3" />
                                  {result.city}, {result.province}
                                </div>
                                <div className="flex items-center gap-1">
                                  <Clock className="h-3 w-3" />
                                  {formatDate(result.created_at)}
                                </div>
                              </div>

                              <div className="flex items-center justify-between">
                                <div className="text-lg font-bold text-green-600">
                                  {formatPrice(result.price)}
                                </div>
                                <div className="flex items-center gap-3 text-sm text-gray-600">
                                  {result.bedrooms > 0 && (
                                    <span>{result.bedrooms} dorm.</span>
                                  )}
                                  {result.bathrooms > 0 && (
                                    <span>{result.bathrooms} baños</span>
                                  )}
                                  {result.area_m2 > 0 && (
                                    <span>{result.area_m2} m²</span>
                                  )}
                                </div>
                              </div>

                              <div className="mt-3">
                                <Button className="w-full" size="sm">
                                  Ver Detalles
                                </Button>
                              </div>
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                  </div>
                )}
              </TabsContent>

              <TabsContent value="trending" className="space-y-4">
                <Card>
                  <CardContent className="p-8 text-center">
                    <TrendingUp className="h-12 w-12 mx-auto text-gray-400 mb-4" />
                    <h3 className="text-lg font-medium text-gray-900 mb-2">
                      Tendencias del Mercado
                    </h3>
                    <p className="text-gray-600">
                      Análisis de tendencias en desarrollo
                    </p>
                  </CardContent>
                </Card>
              </TabsContent>

              <TabsContent value="recent" className="space-y-4">
                <Card>
                  <CardContent className="p-8 text-center">
                    <Clock className="h-12 w-12 mx-auto text-gray-400 mb-4" />
                    <h3 className="text-lg font-medium text-gray-900 mb-2">
                      Búsquedas Recientes
                    </h3>
                    <p className="text-gray-600">
                      Historial de búsquedas en desarrollo
                    </p>
                  </CardContent>
                </Card>
              </TabsContent>
            </Tabs>
          </div>

          {/* Sidebar */}
          <div className="space-y-4">
            {/* Saved Searches */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg flex items-center gap-2">
                  <Bookmark className="h-5 w-5" />
                  Búsquedas Guardadas
                </CardTitle>
              </CardHeader>
              <CardContent>
                {savedSearches.length === 0 ? (
                  <p className="text-sm text-gray-600">
                    No tienes búsquedas guardadas
                  </p>
                ) : (
                  <div className="space-y-2">
                    {savedSearches.map((search, index) => (
                      <div key={index} className="flex items-center justify-between p-2 bg-gray-50 rounded">
                        <span className="text-sm truncate">{search}</span>
                        <Button variant="ghost" size="sm">
                          <Search className="h-3 w-3" />
                        </Button>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Quick Filters */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg flex items-center gap-2">
                  <Filter className="h-5 w-5" />
                  Filtros Rápidos
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div>
                    <h4 className="text-sm font-medium mb-2">Por Precio</h4>
                    <div className="flex flex-wrap gap-1">
                      {['< $50k', '$50k-$100k', '$100k-$200k', '> $200k'].map(range => (
                        <Badge key={range} variant="secondary" className="text-xs cursor-pointer">
                          {range}
                        </Badge>
                      ))}
                    </div>
                  </div>
                  
                  <div>
                    <h4 className="text-sm font-medium mb-2">Por Tipo</h4>
                    <div className="flex flex-wrap gap-1">
                      {['Casa', 'Apartamento', 'Terreno', 'Comercial'].map(type => (
                        <Badge key={type} variant="secondary" className="text-xs cursor-pointer">
                          {type}
                        </Badge>
                      ))}
                    </div>
                  </div>
                  
                  <div>
                    <h4 className="text-sm font-medium mb-2">Por Ubicación</h4>
                    <div className="flex flex-wrap gap-1">
                      {['Quito', 'Guayaquil', 'Cuenca', 'Ambato'].map(city => (
                        <Badge key={city} variant="secondary" className="text-xs cursor-pointer">
                          {city}
                        </Badge>
                      ))}
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Search Tips */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Consejos de Búsqueda</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2 text-sm text-gray-600">
                  <p>• Usa comillas para buscar frases exactas</p>
                  <p>• Combina filtros para mejores resultados</p>
                  <p>• Guarda tus búsquedas favoritas</p>
                  <p>• Configura alertas para nuevas propiedades</p>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
}