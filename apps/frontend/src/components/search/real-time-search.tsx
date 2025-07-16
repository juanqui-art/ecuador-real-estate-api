'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Search, X, MapPin, Home, DollarSign, Filter, Loader2 } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { apiClient } from '@/lib/api-client';
import { formatPrice } from '@/lib/utils';
import { useDebounce } from '@/hooks/useDebounce';

interface SearchFilters {
  query: string;
  type: string;
  province: string;
  minPrice: string;
  maxPrice: string;
  status: string;
}

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
  slug?: string;
  sector?: string;
}

interface RankedSearchResult {
  Property: SearchResult;
  Rank: number;
}

interface SearchResponse {
  results: SearchResult[];
  total: number;
  query: string;
  filters: Record<string, any>;
  suggestions: string[];
}

const provinces = [
  'Azuay', 'Bolívar', 'Cañar', 'Carchi', 'Chimborazo', 'Cotopaxi', 'El Oro', 
  'Esmeraldas', 'Galápagos', 'Guayas', 'Imbabura', 'Loja', 'Los Ríos', 'Manabí', 
  'Morona Santiago', 'Napo', 'Orellana', 'Pastaza', 'Pichincha', 'Santa Elena', 
  'Santo Domingo', 'Sucumbíos', 'Tungurahua', 'Zamora Chinchipe'
];

const propertyTypes = [
  { value: 'house', label: 'Casa' },
  { value: 'apartment', label: 'Apartamento' },
  { value: 'land', label: 'Terreno' },
  { value: 'commercial', label: 'Comercial' }
];

const propertyStatuses = [
  { value: 'available', label: 'Disponible' },
  { value: 'sold', label: 'Vendida' },
  { value: 'rented', label: 'Rentada' }
];

interface RealTimeSearchProps {
  onResultSelect?: (result: SearchResult) => void;
  placeholder?: string;
  showFilters?: boolean;
  className?: string;
}

export function RealTimeSearch({ 
  onResultSelect, 
  placeholder = "Buscar propiedades...",
  showFilters = true,
  className = ""
}: RealTimeSearchProps) {
  const [filters, setFilters] = useState<SearchFilters>({
    query: '',
    type: '',
    province: '',
    minPrice: '',
    maxPrice: '',
    status: ''
  });
  
  const [isOpen, setIsOpen] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [recentSearches, setRecentSearches] = useState<string[]>([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  
  const inputRef = useRef<HTMLInputElement>(null);
  const resultsRef = useRef<HTMLDivElement>(null);
  
  const debouncedQuery = useDebounce(filters.query, 300);
  const debouncedFilters = useDebounce(filters, 500);

  // Load recent searches from localStorage
  useEffect(() => {
    const stored = localStorage.getItem('recentSearches');
    if (stored) {
      setRecentSearches(JSON.parse(stored));
    }
  }, []);

  // Save recent searches to localStorage
  const saveRecentSearch = useCallback((query: string) => {
    if (query.trim().length < 2) return;
    
    const updated = [query, ...recentSearches.filter(s => s !== query)].slice(0, 5);
    setRecentSearches(updated);
    localStorage.setItem('recentSearches', JSON.stringify(updated));
  }, [recentSearches]);

  // Search API call
  const { data: searchResults, isLoading, error } = useQuery<SearchResponse>({
    queryKey: ['real-time-search', debouncedQuery, debouncedFilters],
    queryFn: async () => {
      const trimmedQuery = debouncedQuery.trim();
      
      // Don't search if query is too short (backend requires minimum 2 characters)
      if (trimmedQuery && trimmedQuery.length < 2) {
        return { results: [], total: 0, query: debouncedQuery, filters: debouncedFilters, suggestions: [] };
      }
      
      if (!trimmedQuery && !Object.values(debouncedFilters).some(v => v && v !== debouncedQuery)) {
        return { results: [], total: 0, query: '', filters: {}, suggestions: [] };
      }

      const searchParams = new URLSearchParams();
      
      // Use 'q' parameter instead of 'query' to match backend
      if (debouncedQuery.trim()) {
        searchParams.append('q', debouncedQuery.trim());
      }
      if (debouncedFilters.type) {
        searchParams.append('type', debouncedFilters.type);
      }
      if (debouncedFilters.province) {
        searchParams.append('province', debouncedFilters.province);
      }
      if (debouncedFilters.minPrice) {
        searchParams.append('min_price', debouncedFilters.minPrice);
      }
      if (debouncedFilters.maxPrice) {
        searchParams.append('max_price', debouncedFilters.maxPrice);
      }
      if (debouncedFilters.status) {
        searchParams.append('status', debouncedFilters.status);
      }

      try {
        // Try ranked search first (works with text queries)
        if (debouncedQuery.trim()) {
          try {
            const response = await apiClient.get(`/properties/search/ranked?${searchParams.toString()}`);
            
            if (response.data.success) {
              const rankedResults = response.data.data as RankedSearchResult[];
              const results = rankedResults.map(item => ({
                ...item.Property,
                // Map backend fields to frontend interface
                created_at: item.Property.created_at || new Date().toISOString(),
                images_count: 0, // Backend doesn't return this
                main_image_url: item.Property.main_image || undefined,
                status: item.Property.status || 'available'
              }));
              
              return {
                results,
                total: results.length,
                query: debouncedQuery,
                filters: debouncedFilters,
                suggestions: []
              };
            }
          } catch (rankedError) {
            console.warn('Ranked search failed, falling back to filter search:', rankedError);
          }
        }
        
        // Fallback to filter endpoint for other filters
        try {
          const response = await apiClient.get(`/properties/filter?${searchParams.toString()}`);
          
          if (response.data.success) {
            const results = (response.data.data || []).map((item: any) => ({
              ...item,
              created_at: item.created_at || new Date().toISOString(),
              images_count: 0,
              main_image_url: item.main_image || undefined,
              status: item.status || 'available'
            }));
            
            return {
              results,
              total: results.length,
              query: debouncedQuery,
              filters: debouncedFilters,
              suggestions: []
            };
          }
        } catch (filterError) {
          console.warn('Filter search failed:', filterError);
        }
        
        return { results: [], total: 0, query: debouncedQuery, filters: debouncedFilters, suggestions: [] };
      } catch (error) {
        console.error('Search error:', error);
        
        // Enhanced error logging to help debug
        if (error instanceof Error) {
          console.error('Error message:', error.message);
          console.error('Error stack:', error.stack);
        }
        
        return { results: [], total: 0, query: debouncedQuery, filters: debouncedFilters, suggestions: [] };
      }
    },
    enabled: debouncedQuery.trim().length >= 2 || Object.values(debouncedFilters).some(v => v && v !== debouncedQuery),
    refetchOnWindowFocus: false,
  });

  // Handle input focus
  const handleFocus = useCallback(() => {
    setIsOpen(true);
    setShowSuggestions(true);
  }, []);

  // Handle input blur
  const handleBlur = useCallback(() => {
    // Delay to allow clicks on results
    setTimeout(() => {
      setIsOpen(false);
      setShowSuggestions(false);
      setSelectedIndex(-1);
    }, 200);
  }, []);

  // Handle input change
  const handleInputChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setFilters(prev => ({ ...prev, query: value }));
    setSelectedIndex(-1);
    setIsOpen(true);
    setShowSuggestions(true);
  }, []);

  // Handle filter change
  const handleFilterChange = useCallback((key: keyof SearchFilters, value: string) => {
    setFilters(prev => ({ ...prev, [key]: value }));
  }, []);

  // Handle clear search
  const handleClear = useCallback(() => {
    setFilters(prev => ({ ...prev, query: '' }));
    setSelectedIndex(-1);
    setIsOpen(false);
    inputRef.current?.focus();
  }, []);

  // Handle result selection
  const handleResultSelect = useCallback((result: SearchResult) => {
    saveRecentSearch(result.title);
    setIsOpen(false);
    setShowSuggestions(false);
    onResultSelect?.(result);
  }, [onResultSelect, saveRecentSearch]);

  // Handle recent search selection
  const handleRecentSearchSelect = useCallback((query: string) => {
    setFilters(prev => ({ ...prev, query }));
    setShowSuggestions(false);
  }, []);

  // Handle keyboard navigation
  const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
    if (!isOpen) return;

    const results = searchResults?.results || [];
    const totalItems = results.length + (showSuggestions ? recentSearches.length : 0);

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setSelectedIndex(prev => (prev + 1) % totalItems);
        break;
      case 'ArrowUp':
        e.preventDefault();
        setSelectedIndex(prev => (prev - 1 + totalItems) % totalItems);
        break;
      case 'Enter':
        e.preventDefault();
        if (selectedIndex >= 0) {
          if (showSuggestions && selectedIndex < recentSearches.length) {
            handleRecentSearchSelect(recentSearches[selectedIndex]);
          } else {
            const resultIndex = showSuggestions ? selectedIndex - recentSearches.length : selectedIndex;
            if (results[resultIndex]) {
              handleResultSelect(results[resultIndex]);
            }
          }
        }
        break;
      case 'Escape':
        setIsOpen(false);
        setShowSuggestions(false);
        setSelectedIndex(-1);
        inputRef.current?.blur();
        break;
    }
  }, [isOpen, selectedIndex, searchResults, recentSearches, showSuggestions, handleResultSelect, handleRecentSearchSelect]);

  // Get status color
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'available': return 'bg-green-100 text-green-800';
      case 'sold': return 'bg-red-100 text-red-800';
      case 'rented': return 'bg-blue-100 text-blue-800';
      default: return 'bg-gray-100 text-gray-800';
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

  const getTypeLabel = (type: string) => {
    const found = propertyTypes.find(t => t.value === type);
    return found ? found.label : type;
  };

  return (
    <div className={`relative ${className}`}>
      {/* Search Input */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
        <Input
          ref={inputRef}
          type="text"
          placeholder={placeholder}
          value={filters.query}
          onChange={handleInputChange}
          onFocus={handleFocus}
          onBlur={handleBlur}
          onKeyDown={handleKeyDown}
          className="pl-10 pr-10"
        />
        {filters.query && (
          <Button
            variant="ghost"
            size="sm"
            onClick={handleClear}
            className="absolute right-1 top-1/2 transform -translate-y-1/2 h-8 w-8 p-0"
          >
            <X className="h-4 w-4" />
          </Button>
        )}
      </div>

      {/* Filters */}
      {showFilters && (
        <div className="grid grid-cols-2 md:grid-cols-6 gap-2 mt-3">
          <Select value={filters.type} onValueChange={(value) => handleFilterChange('type', value)}>
            <SelectTrigger className="h-9">
              <SelectValue placeholder="Tipo" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">Todos</SelectItem>
              {propertyTypes.map(type => (
                <SelectItem key={type.value} value={type.value}>
                  {type.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select value={filters.province} onValueChange={(value) => handleFilterChange('province', value)}>
            <SelectTrigger className="h-9">
              <SelectValue placeholder="Provincia" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">Todas</SelectItem>
              {provinces.map(province => (
                <SelectItem key={province} value={province}>
                  {province}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Input
            placeholder="Precio mín."
            value={filters.minPrice}
            onChange={(e) => handleFilterChange('minPrice', e.target.value)}
            className="h-9"
          />

          <Input
            placeholder="Precio máx."
            value={filters.maxPrice}
            onChange={(e) => handleFilterChange('maxPrice', e.target.value)}
            className="h-9"
          />

          <Select value={filters.status} onValueChange={(value) => handleFilterChange('status', value)}>
            <SelectTrigger className="h-9">
              <SelectValue placeholder="Estado" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">Todos</SelectItem>
              {propertyStatuses.map(status => (
                <SelectItem key={status.value} value={status.value}>
                  {status.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Button
            variant="outline"
            size="sm"
            onClick={() => setFilters({
              query: '',
              type: '',
              province: '',
              minPrice: '',
              maxPrice: '',
              status: ''
            })}
            className="h-9"
          >
            <Filter className="h-4 w-4 mr-2" />
            Limpiar
          </Button>
        </div>
      )}

      {/* Search Results */}
      {isOpen && (
        <Card className="absolute top-full mt-1 w-full z-50 max-h-96 overflow-hidden">
          <CardContent className="p-0">
            <div ref={resultsRef} className="max-h-96 overflow-y-auto">
              {/* Loading State */}
              {isLoading && (
                <div className="p-4 text-center">
                  <Loader2 className="h-5 w-5 animate-spin mx-auto mb-2" />
                  <p className="text-sm text-gray-600">Buscando...</p>
                </div>
              )}

              {/* Error State */}
              {error && (
                <div className="p-4 text-center text-red-600">
                  <p className="text-sm">Error al realizar la búsqueda</p>
                </div>
              )}

              {/* Recent Searches */}
              {showSuggestions && recentSearches.length > 0 && !filters.query && (
                <div className="p-3">
                  <h3 className="text-sm font-medium text-gray-700 mb-2">Búsquedas recientes</h3>
                  {recentSearches.map((search, index) => (
                    <div
                      key={index}
                      className={`flex items-center px-3 py-2 cursor-pointer rounded-md ${
                        selectedIndex === index ? 'bg-blue-50' : 'hover:bg-gray-50'
                      }`}
                      onClick={() => handleRecentSearchSelect(search)}
                    >
                      <Search className="h-4 w-4 text-gray-400 mr-3" />
                      <span className="text-sm">{search}</span>
                    </div>
                  ))}
                </div>
              )}

              {/* Search Results */}
              {searchResults && searchResults.results.length > 0 && (
                <div className="p-3">
                  <div className="flex items-center justify-between mb-2">
                    <h3 className="text-sm font-medium text-gray-700">
                      Resultados ({searchResults.total})
                    </h3>
                    {searchResults.query && (
                      <Badge variant="secondary" className="text-xs">
                        "{searchResults.query}"
                      </Badge>
                    )}
                  </div>
                  
                  <div className="space-y-2">
                    {searchResults.results.map((result, index) => {
                      const adjustedIndex = showSuggestions ? index + recentSearches.length : index;
                      return (
                        <div
                          key={result.id}
                          className={`flex items-start p-3 cursor-pointer rounded-md border ${
                            selectedIndex === adjustedIndex ? 'bg-blue-50 border-blue-200' : 'hover:bg-gray-50'
                          }`}
                          onClick={() => handleResultSelect(result)}
                        >
                          <div className="flex-shrink-0 w-12 h-12 bg-gray-200 rounded-md mr-3 flex items-center justify-center">
                            {result.main_image_url ? (
                              <img 
                                src={result.main_image_url} 
                                alt={result.title}
                                className="w-full h-full object-cover rounded-md"
                              />
                            ) : (
                              <Home className="h-6 w-6 text-gray-400" />
                            )}
                          </div>
                          
                          <div className="flex-1 min-w-0">
                            <div className="flex items-start justify-between">
                              <h4 className="text-sm font-medium text-gray-900 truncate">
                                {result.title}
                              </h4>
                              <div className="flex items-center gap-1 ml-2">
                                <Badge className={`text-xs ${getStatusColor(result.status)}`}>
                                  {getStatusLabel(result.status)}
                                </Badge>
                              </div>
                            </div>
                            
                            <div className="flex items-center gap-4 mt-1">
                              <div className="flex items-center text-xs text-gray-600">
                                <DollarSign className="h-3 w-3 mr-1" />
                                {formatPrice(result.price)}
                              </div>
                              <div className="flex items-center text-xs text-gray-600">
                                <MapPin className="h-3 w-3 mr-1" />
                                {result.city}, {result.province}
                              </div>
                            </div>
                            
                            <div className="flex items-center gap-3 mt-1">
                              <span className="text-xs text-gray-500">
                                {getTypeLabel(result.type)}
                              </span>
                              {result.bedrooms > 0 && (
                                <span className="text-xs text-gray-500">
                                  {result.bedrooms} dorm.
                                </span>
                              )}
                              {result.bathrooms > 0 && (
                                <span className="text-xs text-gray-500">
                                  {result.bathrooms} baños
                                </span>
                              )}
                              {result.area_m2 > 0 && (
                                <span className="text-xs text-gray-500">
                                  {result.area_m2} m²
                                </span>
                              )}
                            </div>
                          </div>
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}

              {/* Short Query Info */}
              {filters.query && filters.query.trim().length === 1 && !isLoading && (
                <div className="p-4 text-center">
                  <p className="text-sm text-gray-500">Escribe al menos 2 caracteres para buscar</p>
                </div>
              )}

              {/* No Results */}
              {searchResults && searchResults.results.length === 0 && filters.query && filters.query.trim().length >= 2 && !isLoading && (
                <div className="p-4 text-center">
                  <p className="text-sm text-gray-600">No se encontraron resultados para "{filters.query}"</p>
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}