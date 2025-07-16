'use client';

import React from 'react';
import { useForm, type FieldApi } from '@tanstack/react-form';
import { zodValidator } from '@tanstack/zod-form-adapter';
import type { z } from 'zod';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { LoadingSpinner } from '@/components/ui/loading';
import { AlertTriangle, CheckCircle } from 'lucide-react';

interface BaseFormProps<TFormData> {
  title?: string;
  description?: string;
  schema: z.ZodSchema<TFormData>;
  defaultValues?: Partial<TFormData>;
  onSubmit: (values: TFormData) => Promise<any> | any;
  submitText?: string;
  cancelText?: string;
  onCancel?: () => void;
  isLoading?: boolean;
  children: (form: any) => React.ReactNode;
  className?: string;
}

/**
 * Base form component using TanStack Form with Zod validation
 */
export function BaseForm<TFormData extends Record<string, any>>({
  title,
  description,
  schema,
  defaultValues,
  onSubmit,
  submitText = 'Guardar',
  cancelText = 'Cancelar',
  onCancel,
  isLoading = false,
  children,
  className = '',
}: BaseFormProps<TFormData>) {
  const [submitStatus, setSubmitStatus] = React.useState<{
    type: 'success' | 'error' | null;
    message: string;
  }>({ type: null, message: '' });

  const form = useForm({
    defaultValues: defaultValues || {} as Partial<TFormData>,
    validatorAdapter: zodValidator(),
    validators: {
      onChange: schema,
    },
    onSubmit: async ({ value }) => {
      try {
        setSubmitStatus({ type: null, message: '' });
        const result = await onSubmit(value as TFormData);
        
        if (result?.success === false) {
          setSubmitStatus({
            type: 'error',
            message: result.message || 'Error al procesar el formulario',
          });
        } else {
          setSubmitStatus({
            type: 'success',
            message: result?.message || 'Formulario enviado exitosamente',
          });
        }
        
        return result;
      } catch (error) {
        const message = error instanceof Error ? error.message : 'Error inesperado';
        setSubmitStatus({
          type: 'error',
          message,
        });
        throw error;
      }
    },
  });

  return (
    <Card className={className}>
      {(title || description) && (
        <CardHeader>
          {title && <CardTitle>{title}</CardTitle>}
          {description && <p className="text-sm text-muted-foreground">{description}</p>}
        </CardHeader>
      )}
      
      <CardContent>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
          className="space-y-6"
        >
          {children(form)}
          
          {/* Submit Status Message */}
          {submitStatus.type && (
            <div className={`p-3 rounded-md border ${
              submitStatus.type === 'success' 
                ? 'bg-green-50 border-green-200 text-green-800'
                : 'bg-red-50 border-red-200 text-red-800'
            }`}>
              <div className="flex items-center">
                {submitStatus.type === 'success' ? (
                  <CheckCircle className="w-4 h-4 mr-2" />
                ) : (
                  <AlertTriangle className="w-4 h-4 mr-2" />
                )}
                <span className="text-sm">{submitStatus.message}</span>
              </div>
            </div>
          )}

          {/* Form Actions */}
          <div className="flex flex-col sm:flex-row gap-3 pt-4">
            <Button
              type="submit"
              disabled={isLoading || !form.state.canSubmit}
              className="flex-1 sm:flex-none"
            >
              {isLoading ? (
                <LoadingSpinner size="sm" />
              ) : (
                submitText
              )}
            </Button>
            
            {onCancel && (
              <Button
                type="button"
                variant="outline"
                onClick={onCancel}
                disabled={isLoading}
                className="flex-1 sm:flex-none"
              >
                {cancelText}
              </Button>
            )}
          </div>
        </form>
      </CardContent>
    </Card>
  );
}

/**
 * Field wrapper component for consistent styling
 */
interface FieldWrapperProps {
  label?: string;
  description?: string;
  required?: boolean;
  error?: string;
  children: React.ReactNode;
}

export function FieldWrapper({ 
  label, 
  description, 
  required, 
  error, 
  children 
}: FieldWrapperProps) {
  return (
    <div className="space-y-2">
      {label && (
        <label className="text-sm font-medium text-gray-700">
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </label>
      )}
      
      {children}
      
      {description && (
        <p className="text-xs text-gray-500">{description}</p>
      )}
      
      {error && (
        <p className="text-xs text-red-600 flex items-center">
          <AlertTriangle className="w-3 h-3 mr-1" />
          {error}
        </p>
      )}
    </div>
  );
}

/**
 * Helper function to get field error message
 */
export function getFieldError(field: FieldApi<any, any, any, any>) {
  return field.state.meta.errors?.[0];
}

/**
 * Utility hook for form field state
 */
export function useFieldState(field: FieldApi<any, any, any, any>) {
  const error = getFieldError(field);
  const isDirty = field.state.meta.isDirty;
  const isTouched = field.state.meta.isTouched;
  const isValidating = field.state.meta.isValidating;
  
  return {
    error,
    isDirty,
    isTouched,
    isValidating,
    hasError: !!(error && isTouched),
    isValid: !error && isTouched && isDirty,
  };
}