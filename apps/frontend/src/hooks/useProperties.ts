import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';

// Types - will be imported from shared later
interface Property {
  id: string;
  title: string;
  description: string;
  price: number;
  province: string;
  city: string;
  type: string;
  status: string;
  bedrooms: number;
  bathrooms: number;
  area_m2: number;
  parking_spaces: number;
  main_image?: string;
  images: string[];
  created_at: string;
  updated_at: string;
}

interface PropertyFilters {
  search?: string;
  province?: string;
  city?: string;
  type?: string;
  status?: string;
  min_price?: number;
  max_price?: number;
  min_bedrooms?: number;
  max_bedrooms?: number;
  min_bathrooms?: number;
  max_bathrooms?: number;
  min_area?: number;
  max_area?: number;
  features?: string[];
  page?: number;
  limit?: number;
}

interface PropertyResponse {
  properties: Property[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

interface PropertyStats {
  total_properties: number;
  available_properties: number;
  sold_properties: number;
  rented_properties: number;
  average_price: number;
  total_value: number;
  properties_by_type: Record<string, number>;
  properties_by_province: Record<string, number>;
}

// API functions
const propertiesApi = {
  // Get all properties with filters
  getProperties: async (filters: PropertyFilters = {}): Promise<PropertyResponse> => {
    const params = new URLSearchParams();
    
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        params.append(key, value.toString());
      }
    });

    const response = await api.get(`/properties?${params.toString()}`);
    return response.data;
  },

  // Get single property
  getProperty: async (id: string): Promise<Property> => {
    const response = await api.get(`/properties/${id}`);
    return response.data;
  },

  // Get property by slug
  getPropertyBySlug: async (slug: string): Promise<Property> => {
    const response = await api.get(`/properties/slug/${slug}`);
    return response.data;
  },

  // Create property
  createProperty: async (property: Partial<Property>): Promise<Property> => {
    const response = await api.post('/properties', property);
    return response.data;
  },

  // Update property
  updateProperty: async (id: string, property: Partial<Property>): Promise<Property> => {
    const response = await api.put(`/properties/${id}`, property);
    return response.data;
  },

  // Delete property
  deleteProperty: async (id: string): Promise<void> => {
    await api.delete(`/properties/${id}`);
  },

  // Get property statistics
  getPropertyStats: async (): Promise<PropertyStats> => {
    const response = await api.get('/properties/statistics');
    return response.data;
  },

  // Search properties with ranking
  searchProperties: async (query: string, filters: PropertyFilters = {}): Promise<PropertyResponse> => {
    const params = new URLSearchParams({ q: query });
    
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        params.append(key, value.toString());
      }
    });

    const response = await api.get(`/properties/search/ranked?${params.toString()}`);
    return response.data;
  },

  // Get search suggestions
  getSearchSuggestions: async (query: string): Promise<string[]> => {
    const response = await api.get(`/properties/search/suggestions?q=${query}`);
    return response.data.suggestions || [];
  },
};

// React Query hooks
export const useProperties = (filters: PropertyFilters = {}) => {
  return useQuery({
    queryKey: ['properties', filters],
    queryFn: () => propertiesApi.getProperties(filters),
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes
  });
};

export const useProperty = (id: string) => {
  return useQuery({
    queryKey: ['property', id],
    queryFn: () => propertiesApi.getProperty(id),
    enabled: !!id,
    staleTime: 10 * 60 * 1000, // 10 minutes
  });
};

export const usePropertyBySlug = (slug: string) => {
  return useQuery({
    queryKey: ['property', 'slug', slug],
    queryFn: () => propertiesApi.getPropertyBySlug(slug),
    enabled: !!slug,
    staleTime: 10 * 60 * 1000, // 10 minutes
  });
};

export const usePropertyStats = () => {
  return useQuery({
    queryKey: ['property-stats'],
    queryFn: propertiesApi.getPropertyStats,
    staleTime: 15 * 60 * 1000, // 15 minutes
  });
};

export const useSearchProperties = (query: string, filters: PropertyFilters = {}) => {
  return useQuery({
    queryKey: ['search-properties', query, filters],
    queryFn: () => propertiesApi.searchProperties(query, filters),
    enabled: !!query.trim(),
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
};

export const useSearchSuggestions = (query: string) => {
  return useQuery({
    queryKey: ['search-suggestions', query],
    queryFn: () => propertiesApi.getSearchSuggestions(query),
    enabled: !!query.trim() && query.length > 2,
    staleTime: 1 * 60 * 1000, // 1 minute
  });
};

// Mutations
export const useCreateProperty = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: propertiesApi.createProperty,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['properties'] });
      queryClient.invalidateQueries({ queryKey: ['property-stats'] });
    },
  });
};

export const useUpdateProperty = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, property }: { id: string; property: Partial<Property> }) =>
      propertiesApi.updateProperty(id, property),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['properties'] });
      queryClient.invalidateQueries({ queryKey: ['property', data.id] });
      queryClient.invalidateQueries({ queryKey: ['property-stats'] });
    },
  });
};

export const useDeleteProperty = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: propertiesApi.deleteProperty,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['properties'] });
      queryClient.invalidateQueries({ queryKey: ['property-stats'] });
    },
  });
};

// Prefetch utilities
export const prefetchProperty = (queryClient: ReturnType<typeof useQueryClient>, id: string) => {
  queryClient.prefetchQuery({
    queryKey: ['property', id],
    queryFn: () => propertiesApi.getProperty(id),
    staleTime: 10 * 60 * 1000,
  });
};

export const prefetchProperties = (queryClient: ReturnType<typeof useQueryClient>, filters: PropertyFilters = {}) => {
  queryClient.prefetchQuery({
    queryKey: ['properties', filters],
    queryFn: () => propertiesApi.getProperties(filters),
    staleTime: 5 * 60 * 1000,
  });
};