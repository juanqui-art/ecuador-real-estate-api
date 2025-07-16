'use client';

import React from 'react';
import type { FieldApi } from '@tanstack/react-form';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import { FieldWrapper, useFieldState } from './base-form';
import { Eye, EyeOff } from 'lucide-react';

interface BaseFieldProps {
  label?: string;
  description?: string;
  required?: boolean;
  placeholder?: string;
}

/**
 * Text input field component
 */
interface TextFieldProps extends BaseFieldProps {
  field: FieldApi<any, any, any, any>;
  type?: 'text' | 'email' | 'tel' | 'url';
}

export function TextField({ 
  field, 
  label, 
  description, 
  required, 
  placeholder,
  type = 'text' 
}: TextFieldProps) {
  const { error, hasError } = useFieldState(field);

  return (
    <FieldWrapper 
      label={label} 
      description={description} 
      required={required} 
      error={hasError ? error : undefined}
    >
      <Input
        type={type}
        value={field.state.value || ''}
        onChange={(e) => field.handleChange(e.target.value)}
        onBlur={field.handleBlur}
        placeholder={placeholder}
        className={hasError ? 'border-red-500 focus:border-red-500' : ''}
      />
    </FieldWrapper>
  );
}

/**
 * Password input field component with toggle visibility
 */
interface PasswordFieldProps extends BaseFieldProps {
  field: FieldApi<any, any, any, any>;
}

export function PasswordField({ 
  field, 
  label, 
  description, 
  required, 
  placeholder 
}: PasswordFieldProps) {
  const [showPassword, setShowPassword] = React.useState(false);
  const { error, hasError } = useFieldState(field);

  return (
    <FieldWrapper 
      label={label} 
      description={description} 
      required={required} 
      error={hasError ? error : undefined}
    >
      <div className="relative">
        <Input
          type={showPassword ? 'text' : 'password'}
          value={field.state.value || ''}
          onChange={(e) => field.handleChange(e.target.value)}
          onBlur={field.handleBlur}
          placeholder={placeholder}
          className={`pr-10 ${hasError ? 'border-red-500 focus:border-red-500' : ''}`}
        />
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
          onClick={() => setShowPassword(!showPassword)}
        >
          {showPassword ? (
            <EyeOff className="h-4 w-4 text-gray-500" />
          ) : (
            <Eye className="h-4 w-4 text-gray-500" />
          )}
        </Button>
      </div>
    </FieldWrapper>
  );
}

/**
 * Number input field component
 */
interface NumberFieldProps extends BaseFieldProps {
  field: FieldApi<any, any, any, any>;
  min?: number;
  max?: number;
  step?: number;
}

export function NumberField({ 
  field, 
  label, 
  description, 
  required, 
  placeholder,
  min,
  max,
  step 
}: NumberFieldProps) {
  const { error, hasError } = useFieldState(field);

  return (
    <FieldWrapper 
      label={label} 
      description={description} 
      required={required} 
      error={hasError ? error : undefined}
    >
      <Input
        type="number"
        value={field.state.value || ''}
        onChange={(e) => {
          const value = e.target.value;
          field.handleChange(value === '' ? undefined : Number(value));
        }}
        onBlur={field.handleBlur}
        placeholder={placeholder}
        min={min}
        max={max}
        step={step}
        className={hasError ? 'border-red-500 focus:border-red-500' : ''}
      />
    </FieldWrapper>
  );
}

/**
 * Textarea field component
 */
interface TextareaFieldProps extends BaseFieldProps {
  field: FieldApi<any, any, any, any>;
  rows?: number;
}

export function TextareaField({ 
  field, 
  label, 
  description, 
  required, 
  placeholder,
  rows = 3 
}: TextareaFieldProps) {
  const { error, hasError } = useFieldState(field);

  return (
    <FieldWrapper 
      label={label} 
      description={description} 
      required={required} 
      error={hasError ? error : undefined}
    >
      <Textarea
        value={field.state.value || ''}
        onChange={(e) => field.handleChange(e.target.value)}
        onBlur={field.handleBlur}
        placeholder={placeholder}
        rows={rows}
        className={hasError ? 'border-red-500 focus:border-red-500' : ''}
      />
    </FieldWrapper>
  );
}

/**
 * Select field component
 */
interface SelectFieldProps extends BaseFieldProps {
  field: FieldApi<any, any, any, any>;
  options: { value: string; label: string }[];
}

export function SelectField({ 
  field, 
  label, 
  description, 
  required, 
  placeholder,
  options 
}: SelectFieldProps) {
  const { error, hasError } = useFieldState(field);

  return (
    <FieldWrapper 
      label={label} 
      description={description} 
      required={required} 
      error={hasError ? error : undefined}
    >
      <Select
        value={field.state.value || ''}
        onValueChange={(value) => field.handleChange(value)}
      >
        <SelectTrigger className={hasError ? 'border-red-500 focus:border-red-500' : ''}>
          <SelectValue placeholder={placeholder} />
        </SelectTrigger>
        <SelectContent>
          {options.map((option) => (
            <SelectItem key={option.value} value={option.value}>
              {option.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </FieldWrapper>
  );
}

/**
 * Checkbox field component
 */
interface CheckboxFieldProps extends BaseFieldProps {
  field: FieldApi<any, any, any, any>;
}

export function CheckboxField({ 
  field, 
  label, 
  description, 
  required 
}: CheckboxFieldProps) {
  const { error, hasError } = useFieldState(field);

  return (
    <FieldWrapper 
      description={description} 
      error={hasError ? error : undefined}
    >
      <div className="flex items-center space-x-2">
        <Checkbox
          id={`checkbox-${field.name}`}
          checked={field.state.value || false}
          onCheckedChange={(checked) => field.handleChange(!!checked)}
        />
        <Label 
          htmlFor={`checkbox-${field.name}`}
          className="text-sm font-medium"
        >
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </Label>
      </div>
    </FieldWrapper>
  );
}

/**
 * File input field component
 */
interface FileFieldProps extends BaseFieldProps {
  field: FieldApi<any, any, any, any>;
  accept?: string;
  multiple?: boolean;
}

export function FileField({ 
  field, 
  label, 
  description, 
  required, 
  accept,
  multiple = false 
}: FileFieldProps) {
  const { error, hasError } = useFieldState(field);

  return (
    <FieldWrapper 
      label={label} 
      description={description} 
      required={required} 
      error={hasError ? error : undefined}
    >
      <Input
        type="file"
        accept={accept}
        multiple={multiple}
        onChange={(e) => {
          const files = e.target.files;
          if (files) {
            field.handleChange(multiple ? Array.from(files) : files[0]);
          }
        }}
        onBlur={field.handleBlur}
        className={hasError ? 'border-red-500 focus:border-red-500' : ''}
      />
    </FieldWrapper>
  );
}

/**
 * Currency input field component (for Ecuador - USD)
 */
interface CurrencyFieldProps extends BaseFieldProps {
  field: FieldApi<any, any, any, any>;
  min?: number;
  max?: number;
}

export function CurrencyField({ 
  field, 
  label, 
  description, 
  required, 
  placeholder,
  min = 0,
  max 
}: CurrencyFieldProps) {
  const { error, hasError } = useFieldState(field);

  return (
    <FieldWrapper 
      label={label} 
      description={description} 
      required={required} 
      error={hasError ? error : undefined}
    >
      <div className="relative">
        <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500">
          $
        </span>
        <Input
          type="number"
          value={field.state.value || ''}
          onChange={(e) => {
            const value = e.target.value;
            field.handleChange(value === '' ? undefined : Number(value));
          }}
          onBlur={field.handleBlur}
          placeholder={placeholder}
          min={min}
          max={max}
          step="0.01"
          className={`pl-8 ${hasError ? 'border-red-500 focus:border-red-500' : ''}`}
        />
      </div>
    </FieldWrapper>
  );
}