'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { Search, X, MapPin, Home, DollarSign, Loader2 } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { formatPrice } from '@/lib/utils';
import { useDebounce } from '@/hooks/useDebounce';

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
  main_image?: string;
  slug?: string;
  sector?: string;
}

interface RankedSearchResult {
  Property: SearchResult;
  Rank: number;
}

interface PublicSearchProps {
  onResultSelect?: (result: SearchResult) => void;
  placeholder?: string;
  className?: string;
}

export function PublicSearch({ 
  onResultSelect, 
  placeholder = "Buscar propiedades...",
  className = ""
}: PublicSearchProps) {
  const [query, setQuery] = useState('');
  const [isOpen, setIsOpen] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [results, setResults] = useState<SearchResult[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const inputRef = useRef<HTMLInputElement>(null);
  const resultsRef = useRef<HTMLDivElement>(null);
  
  const debouncedQuery = useDebounce(query, 300);

  // Search function using direct fetch
  const searchProperties = useCallback(async (searchQuery: string) => {
    const trimmedQuery = searchQuery.trim();
    
    // Clear results for empty queries
    if (!trimmedQuery) {
      setResults([]);
      setError(null);
      return;
    }

    // Don't search if query is too short (backend requires minimum 2 characters)
    if (trimmedQuery.length < 2) {
      setResults([]);
      setError(null);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
      const searchParams = new URLSearchParams();
      searchParams.append('q', trimmedQuery);

      // Try ranked search first
      let response = await fetch(`${baseUrl}/api/properties/search/ranked?${searchParams.toString()}`);
      
      if (!response.ok) {
        // Handle specific HTTP errors
        if (response.status === 400) {
          const errorData = await response.json().catch(() => ({ message: 'Bad request' }));
          throw new Error(errorData.message || 'Consulta inválida');
        }
        throw new Error(`Error HTTP ${response.status}`);
      }

      const data = await response.json();
      
      if (data.success && data.data) {
        const rankedResults = data.data as RankedSearchResult[];
        const mappedResults = rankedResults.map(item => ({
          ...item.Property,
          created_at: item.Property.created_at || new Date().toISOString(),
          main_image: item.Property.main_image || undefined,
          status: item.Property.status || 'available'
        }));
        
        setResults(mappedResults);
      } else {
        // Fallback to filter endpoint
        response = await fetch(`${baseUrl}/api/properties/filter?${searchParams.toString()}`);
        
        if (!response.ok) {
          // Handle specific HTTP errors for filter endpoint
          if (response.status === 400) {
            const errorData = await response.json().catch(() => ({ message: 'Bad request' }));
            throw new Error(errorData.message || 'Consulta inválida');
          }
          throw new Error(`Error HTTP ${response.status}`);
        }

        const filterData = await response.json();
        
        if (filterData.success && filterData.data) {
          const mappedResults = filterData.data.map((item: any) => ({
            ...item,
            created_at: item.created_at || new Date().toISOString(),
            main_image: item.main_image || undefined,
            status: item.status || 'available'
          }));
          
          setResults(mappedResults);
        } else {
          setResults([]);
        }
      }
    } catch (err) {
      console.error('Search error:', err);
      setError(err instanceof Error ? err.message : 'Error en la búsqueda');
      setResults([]);
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Effect to trigger search when debounced query changes
  useEffect(() => {
    searchProperties(debouncedQuery);
  }, [debouncedQuery, searchProperties]);

  // Handle input focus
  const handleFocus = useCallback(() => {
    setIsOpen(true);
  }, []);

  // Handle input blur
  const handleBlur = useCallback(() => {
    // Delay to allow clicks on results
    setTimeout(() => {
      setIsOpen(false);
      setSelectedIndex(-1);
    }, 200);
  }, []);

  // Handle input change
  const handleInputChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setQuery(value);
    setSelectedIndex(-1);
    setIsOpen(true);
  }, []);

  // Handle clear search
  const handleClear = useCallback(() => {
    setQuery('');
    setResults([]);
    setSelectedIndex(-1);
    setIsOpen(false);
    inputRef.current?.focus();
  }, []);

  // Handle result selection
  const handleResultSelect = useCallback((result: SearchResult) => {
    setIsOpen(false);
    onResultSelect?.(result);
  }, [onResultSelect]);

  // Handle keyboard navigation
  const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
    if (!isOpen) return;

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setSelectedIndex(prev => (prev + 1) % results.length);
        break;
      case 'ArrowUp':
        e.preventDefault();
        setSelectedIndex(prev => (prev - 1 + results.length) % results.length);
        break;
      case 'Enter':
        e.preventDefault();
        if (selectedIndex >= 0 && results[selectedIndex]) {
          handleResultSelect(results[selectedIndex]);
        }
        break;
      case 'Escape':
        setIsOpen(false);
        setSelectedIndex(-1);
        inputRef.current?.blur();
        break;
    }
  }, [isOpen, selectedIndex, results, handleResultSelect]);

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
    switch (type) {
      case 'house': return 'Casa';
      case 'apartment': return 'Apartamento';
      case 'land': return 'Terreno';
      case 'commercial': return 'Comercial';
      default: return type;
    }
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
          value={query}
          onChange={handleInputChange}
          onFocus={handleFocus}
          onBlur={handleBlur}
          onKeyDown={handleKeyDown}
          className="pl-10 pr-10"
        />
        {query && (
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
                  <p className="text-sm">Error: {error}</p>
                </div>
              )}

              {/* Search Results */}
              {!isLoading && !error && results.length > 0 && (
                <div className="p-3">
                  <div className="flex items-center justify-between mb-2">
                    <h3 className="text-sm font-medium text-gray-700">
                      Resultados ({results.length})
                    </h3>
                    {query && (
                      <Badge variant="secondary" className="text-xs">
                        "{query}"
                      </Badge>
                    )}
                  </div>
                  
                  <div className="space-y-2">
                    {results.map((result, index) => (
                      <div
                        key={result.id}
                        className={`flex items-start p-3 cursor-pointer rounded-md border ${
                          selectedIndex === index ? 'bg-blue-50 border-blue-200' : 'hover:bg-gray-50'
                        }`}
                        onClick={() => handleResultSelect(result)}
                      >
                        <div className="flex-shrink-0 w-12 h-12 bg-gray-200 rounded-md mr-3 flex items-center justify-center">
                          {result.main_image ? (
                            <img 
                              src={result.main_image} 
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
                    ))}
                  </div>
                </div>
              )}

              {/* Short Query Info */}
              {!isLoading && !error && query && query.trim().length === 1 && (
                <div className="p-4 text-center">
                  <p className="text-sm text-gray-500">Escribe al menos 2 caracteres para buscar</p>
                </div>
              )}

              {/* No Results */}
              {!isLoading && !error && results.length === 0 && query && query.trim().length >= 2 && (
                <div className="p-4 text-center">
                  <p className="text-sm text-gray-600">No se encontraron resultados para "{query}"</p>
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}