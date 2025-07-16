/*
 * SERVER ACTIONS - ARCHIVED
 * 
 * These Server Actions are kept for reference but are not currently used.
 * The project uses Client-side approach (TanStack Form + fetch API) for better UX.
 * 
 * Server Actions might be useful for:
 * - Public landing page forms (contact, newsletter)
 * - SEO-critical authentication flows
 * - Cases where JavaScript is disabled
 */

'use server';

import { redirect } from 'next/navigation';
import { cookies } from 'next/headers';
import { apiClient } from '@/lib/api-client';
import { loginSchema, changePasswordSchema } from '@/lib/validations/auth';
import type { LoginFormData, ChangePasswordFormData } from '@/lib/validations/auth';
import type { LoginResponse, User } from '@shared/types/auth';

/**
 * Server action for user login
 * This runs on the server and handles authentication
 */
export async function loginAction(prevState: any, formData: FormData) {
  try {
    // Extract and validate form data
    const rawData = {
      email: formData.get('email') as string,
      password: formData.get('password') as string,
    };

    // Validate with Zod
    const validatedData = loginSchema.parse(rawData);

    // Make API call to backend
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'User-Agent': 'RealtyCore-Dashboard/1.0 (Next.js)',
      },
      body: JSON.stringify(validatedData),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      return {
        success: false,
        message: errorData.message || 'Credenciales inválidas',
        errors: {},
      };
    }

    const loginData: LoginResponse = await response.json();

    // Set secure cookies for authentication
    const cookieStore = await cookies();
    
    // Set access token cookie (15 minutes)
    cookieStore.set('access_token', loginData.tokens.access_token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: 15 * 60, // 15 minutes
    });

    // Set refresh token cookie (7 days)
    cookieStore.set('refresh_token', loginData.tokens.refresh_token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: 7 * 24 * 60 * 60, // 7 days
    });

    // Set user data cookie for client-side access
    cookieStore.set('user_data', JSON.stringify(loginData.user), {
      httpOnly: false, // Accessible to client-side
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: 7 * 24 * 60 * 60, // 7 days
    });

    return {
      success: true,
      message: 'Login exitoso',
      user: loginData.user,
      redirect: '/dashboard',
    };

  } catch (error) {
    if (error instanceof Error && error.name === 'ZodError') {
      return {
        success: false,
        message: 'Datos inválidos',
        errors: (error as any).flatten().fieldErrors,
      };
    }

    console.error('Login action error:', error);
    return {
      success: false,
      message: 'Error interno del servidor',
      errors: {},
    };
  }
}

/**
 * Server action for user logout
 */
export async function logoutAction() {
  try {
    const cookieStore = await cookies();
    const accessToken = cookieStore.get('access_token')?.value;

    // Try to logout from backend (best effort)
    if (accessToken) {
      try {
        await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/auth/logout`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${accessToken}`,
            'Content-Type': 'application/json',
            'User-Agent': 'RealtyCore-Dashboard/1.0 (Next.js)',
          },
        });
      } catch (error) {
        // Backend logout failed, but we still want to clear client state
        console.warn('Backend logout failed:', error);
      }
    }

    // Clear all auth cookies
    cookieStore.delete('access_token');
    cookieStore.delete('refresh_token');
    cookieStore.delete('user_data');

    return { success: true };

  } catch (error) {
    console.error('Logout action error:', error);
    // Even if logout fails, clear cookies
    const cookieStore = await cookies();
    cookieStore.delete('access_token');
    cookieStore.delete('refresh_token');
    cookieStore.delete('user_data');
    
    return { success: true }; // Always succeed logout from user perspective
  }
}

/**
 * Server action for changing password
 */
export async function changePasswordAction(prevState: any, formData: FormData) {
  try {
    const cookieStore = await cookies();
    const accessToken = cookieStore.get('access_token')?.value;

    if (!accessToken) {
      return {
        success: false,
        message: 'No autenticado',
        errors: {},
      };
    }

    // Extract and validate form data
    const rawData = {
      current_password: formData.get('current_password') as string,
      new_password: formData.get('new_password') as string,
      confirm_password: formData.get('confirm_password') as string,
    };

    // Validate with Zod
    const validatedData = changePasswordSchema.parse(rawData);

    // Make API call to backend
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/auth/change-password`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        'User-Agent': 'RealtyCore-Dashboard/1.0 (Next.js)',
      },
      body: JSON.stringify({
        current_password: validatedData.current_password,
        new_password: validatedData.new_password,
      }),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      return {
        success: false,
        message: errorData.message || 'Error al cambiar contraseña',
        errors: {},
      };
    }

    return {
      success: true,
      message: 'Contraseña actualizada exitosamente',
    };

  } catch (error) {
    if (error instanceof Error && error.name === 'ZodError') {
      return {
        success: false,
        message: 'Datos inválidos',
        errors: (error as any).flatten().fieldErrors,
      };
    }

    console.error('Change password action error:', error);
    return {
      success: false,
      message: 'Error interno del servidor',
      errors: {},
    };
  }
}

/**
 * Server action to refresh token
 */
export async function refreshTokenAction() {
  try {
    const cookieStore = await cookies();
    const refreshToken = cookieStore.get('refresh_token')?.value;

    if (!refreshToken) {
      // No refresh token, redirect to login
      redirect('/login');
    }

    // Make API call to refresh
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/auth/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'User-Agent': 'RealtyCore-Dashboard/1.0 (Next.js)',
      },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!response.ok) {
      // Refresh failed, clear cookies and redirect
      cookieStore.delete('access_token');
      cookieStore.delete('refresh_token');
      cookieStore.delete('user_data');
      redirect('/login');
    }

    const tokenData = await response.json();

    // Update access token cookie
    cookieStore.set('access_token', tokenData.access_token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: 15 * 60, // 15 minutes
    });

    // Update refresh token if provided
    if (tokenData.refresh_token) {
      cookieStore.set('refresh_token', tokenData.refresh_token, {
        httpOnly: true,
        secure: process.env.NODE_ENV === 'production',
        sameSite: 'lax',
        maxAge: 7 * 24 * 60 * 60, // 7 days
      });
    }

    return {
      success: true,
      access_token: tokenData.access_token,
    };

  } catch (error) {
    console.error('Token refresh action error:', error);
    
    // Clear cookies and redirect on error
    const cookieStore = await cookies();
    cookieStore.delete('access_token');
    cookieStore.delete('refresh_token');
    cookieStore.delete('user_data');
    redirect('/login');
  }
}

/**
 * Get current user from server-side cookies
 */
export async function getCurrentUser(): Promise<User | null> {
  try {
    const cookieStore = await cookies();
    const userData = cookieStore.get('user_data')?.value;
    
    if (!userData) {
      return null;
    }

    return JSON.parse(userData) as User;
  } catch (error) {
    console.error('Get current user error:', error);
    return null;
  }
}

/**
 * Check if user is authenticated (server-side)
 */
export async function isAuthenticated(): Promise<boolean> {
  try {
    const cookieStore = await cookies();
    const accessToken = cookieStore.get('access_token')?.value;
    const userData = cookieStore.get('user_data')?.value;
    
    return !!(accessToken && userData);
  } catch (error) {
    console.error('Auth check error:', error);
    return false;
  }
}