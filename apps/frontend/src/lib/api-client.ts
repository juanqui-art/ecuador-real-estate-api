/**
 * Modern API Client using native fetch with interceptors
 * Replaces axios for Next.js 15 best practices
 */

interface ApiResponse<T = any> {
  data: T;
  status: number;
  headers: Headers;
}

interface ApiError {
  message: string;
  status: number;
  code?: string;
}

interface RequestConfig {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
  headers?: Record<string, string>;
  body?: any;
  signal?: AbortSignal;
}

class ApiClient {
  private baseURL: string;
  private defaultHeaders: Record<string, string>;

  constructor() {
    this.baseURL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    this.defaultHeaders = {
      'Content-Type': 'application/json',
      'User-Agent': 'RealtyCore-Dashboard/1.0 (Next.js)',
    };
  }

  /**
   * Get access token from localStorage
   */
  private getAccessToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('access_token');
  }

  /**
   * Get refresh token from localStorage
   */
  private getRefreshToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('refresh_token');
  }

  /**
   * Set tokens in localStorage and cookies
   */
  private setTokens(accessToken: string, refreshToken?: string): void {
    if (typeof window === 'undefined') return;
    
    localStorage.setItem('access_token', accessToken);
    if (refreshToken) {
      localStorage.setItem('refresh_token', refreshToken);
    }

    // Also set cookies for server-side middleware
    document.cookie = `access_token=${accessToken}; path=/; max-age=900`; // 15 min
    if (refreshToken) {
      document.cookie = `refresh_token=${refreshToken}; path=/; max-age=604800`; // 7 days
    }
  }

  /**
   * Clear tokens from localStorage and cookies
   */
  private clearTokens(): void {
    if (typeof window === 'undefined') return;
    
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    
    // Clear cookies
    document.cookie = 'access_token=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
    document.cookie = 'refresh_token=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
  }

  /**
   * Refresh access token using refresh token
   */
  private async refreshAccessToken(): Promise<string | null> {
    const refreshToken = this.getRefreshToken();
    if (!refreshToken) return null;

    try {
      const response = await fetch(`${this.baseURL}/api/auth/refresh`, {
        method: 'POST',
        headers: this.defaultHeaders,
        body: JSON.stringify({ refresh_token: refreshToken }),
      });

      if (!response.ok) {
        throw new Error('Token refresh failed');
      }

      const data = await response.json();
      const newAccessToken = data.access_token;
      
      this.setTokens(newAccessToken, data.refresh_token);
      return newAccessToken;
    } catch (error) {
      this.clearTokens();
      // Redirect to login
      if (typeof window !== 'undefined') {
        window.location.href = '/login';
      }
      return null;
    }
  }

  /**
   * Prepare request headers with authentication
   */
  private async prepareHeaders(customHeaders: Record<string, string> = {}): Promise<Record<string, string>> {
    const headers = { ...this.defaultHeaders, ...customHeaders };
    
    const accessToken = this.getAccessToken();
    if (accessToken) {
      headers.Authorization = `Bearer ${accessToken}`;
    }

    return headers;
  }

  /**
   * Handle response and potential token refresh
   */
  private async handleResponse<T>(response: Response, originalRequest: () => Promise<Response>): Promise<ApiResponse<T>> {
    // If request succeeded, return response
    if (response.ok) {
      const data = await response.json();
      return {
        data,
        status: response.status,
        headers: response.headers,
      };
    }

    // Special handling for logout requests - don't attempt token refresh
    if (response.url.includes('/auth/logout')) {
      if (response.status === 401) {
        console.log('🔓 Logout request received 401 - token already invalid, treating as success');
        // For logout, 401 means token is already invalid - this is actually success
        return {
          data: {} as T,
          status: 200, // Treat as success
          headers: response.headers,
        };
      }
    }

    // If 401 and not a logout request, try to refresh token
    if (response.status === 401 && !response.url.includes('/auth/logout')) {
      const newAccessToken = await this.refreshAccessToken();
      
      if (newAccessToken) {
        // Retry original request with new token
        const retryResponse = await originalRequest();
        if (retryResponse.ok) {
          const data = await retryResponse.json();
          return {
            data,
            status: retryResponse.status,
            headers: retryResponse.headers,
          };
        }
      }
    }

    // Handle error response
    let errorMessage = 'Request failed';
    try {
      const errorData = await response.json();
      errorMessage = errorData.message || errorData.error || errorMessage;
    } catch {
      errorMessage = response.statusText || errorMessage;
    }

    const error: ApiError = {
      message: errorMessage,
      status: response.status,
    };

    throw error;
  }

  /**
   * Make HTTP request
   */
  private async request<T>(endpoint: string, config: RequestConfig = {}): Promise<ApiResponse<T>> {
    const { method = 'GET', headers: customHeaders = {}, body, signal } = config;
    
    // Special handling for logout requests - capture token before it might be cleared
    let authHeaders = customHeaders;
    if (endpoint === '/auth/logout') {
      const accessToken = this.getAccessToken();
      if (accessToken) {
        authHeaders = {
          ...customHeaders,
          'Authorization': `Bearer ${accessToken}`
        };
        console.log('🔓 Logout request: Token captured and header set');
      } else {
        console.log('⚠️ Logout request: No access token found');
      }
    }
    
    const headers = await this.prepareHeaders(authHeaders);
    const url = `${this.baseURL}${endpoint}`;

    const requestOptions: RequestInit = {
      method,
      headers,
      signal,
    };

    // Add body for non-GET requests
    if (body && method !== 'GET') {
      if (body instanceof FormData) {
        // Remove Content-Type for FormData to let browser set it with boundary
        delete headers['Content-Type'];
        requestOptions.body = body;
      } else {
        requestOptions.body = JSON.stringify(body);
      }
    }

    const makeRequest = () => fetch(url, requestOptions);
    const response = await makeRequest();

    return this.handleResponse<T>(response, makeRequest);
  }

  /**
   * GET request
   */
  async get<T>(endpoint: string, signal?: AbortSignal): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { method: 'GET', signal });
  }

  /**
   * POST request
   */
  async post<T>(endpoint: string, body?: any, signal?: AbortSignal): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { method: 'POST', body, signal });
  }

  /**
   * PUT request
   */
  async put<T>(endpoint: string, body?: any, signal?: AbortSignal): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { method: 'PUT', body, signal });
  }

  /**
   * DELETE request
   */
  async delete<T>(endpoint: string, signal?: AbortSignal): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { method: 'DELETE', signal });
  }

  /**
   * PATCH request
   */
  async patch<T>(endpoint: string, body?: any, signal?: AbortSignal): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { method: 'PATCH', body, signal });
  }

  /**
   * Upload file (multipart/form-data)
   */
  async upload<T>(endpoint: string, formData: FormData, signal?: AbortSignal): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { 
      method: 'POST', 
      body: formData,
      signal 
    });
  }
}

// Export singleton instance
export const apiClient = new ApiClient();

// Export types for use in other files
export type { ApiResponse, ApiError, RequestConfig };