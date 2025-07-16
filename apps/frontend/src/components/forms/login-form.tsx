'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import { BaseForm } from './base-form';
import { TextField, PasswordField } from './form-fields';
import { loginSchema } from '@/lib/validations/auth';
import { useLogin } from '@/hooks/useAuth';
import type { LoginFormData } from '@/lib/validations/auth';

interface LoginFormProps {
  onSuccess?: () => void;
  redirectTo?: string;
}

export function LoginForm({ onSuccess, redirectTo = '/dashboard' }: LoginFormProps) {
  const router = useRouter();
  const loginMutation = useLogin();

  const handleSubmit = async (values: LoginFormData) => {
    try {
      const result = await loginMutation.mutateAsync(values);
      
      if (result) {
        // Call success callback
        onSuccess?.();
        
        // Redirect to target page
        router.push(redirectTo);
        
        return { success: true, message: 'Login exitoso' };
      }
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Error de autenticación';
      return { success: false, message };
    }
  };

  return (
    <BaseForm
      title="Iniciar Sesión"
      description="Ingresa tus credenciales para acceder al sistema"
      schema={loginSchema}
      defaultValues={{
        email: 'test@example.com',
        password: 'test123',
      }}
      onSubmit={handleSubmit}
      submitText="Iniciar Sesión"
      isLoading={loginMutation.isPending}
      className="w-full max-w-md"
    >
      {(form) => (
        <>
          <form.Field name="email">
            {(field) => (
              <TextField
                field={field}
                label="Correo Electrónico"
                type="email"
                placeholder="correo@ejemplo.com"
                required
              />
            )}
          </form.Field>

          <form.Field name="password">
            {(field) => (
              <PasswordField
                field={field}
                label="Contraseña"
                placeholder="Ingresa tu contraseña"
                required
              />
            )}
          </form.Field>
        </>
      )}
    </BaseForm>
  );
}