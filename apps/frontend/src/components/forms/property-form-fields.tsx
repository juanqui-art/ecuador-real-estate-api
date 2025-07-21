/**
 * Componentes de campos reutilizables para formularios de propiedades
 * Implementa el composition pattern de TanStack Form 2025
 */

import React from 'react';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { AlertCircle, CheckCircle, Loader2 } from 'lucide-react';
import { 
  ECUADORIAN_PROVINCES, 
  PROPERTY_TYPES, 
  PROPERTY_STATUS 
} from '@/lib/constants';
import { formatPrice } from '@/lib/utils';
// Import eliminado: usePropertyForm ya no se utiliza

// Tipos locales para los campos de formulario
type PropertyFormData = {
  title: string;
  description: string;
  price: number;
  type: string;
  status: string;
  province: string;
  city: string;
  address: string;
  bedrooms: number;
  bathrooms: number;
  area_m2: number;
  parking_spaces: number;
  year_built: number;
  contact_phone: string;
  contact_email: string;
  notes: string;
  garden: boolean;
  pool: boolean;
  elevator: boolean;
  balcony: boolean;
  terrace: boolean;
  garage: boolean;
  furnished: boolean;
  air_conditioning: boolean;
  security: boolean;
};

// Componente base para campos de formulario
interface BaseFieldProps {
  name: keyof PropertyFormData;
  label: string;
  description?: string;
  required?: boolean;
  className?: string;
}

interface FormFieldProps extends BaseFieldProps {
  children: React.ReactNode;
}

export function FormField({ name, label, description, required, className, children }: FormFieldProps) {
  return (
    <div className={`space-y-2 ${className || ''}`}>
      <Label htmlFor={name} className="text-sm font-medium">
        {label}
        {required && <span className="text-red-500 ml-1">*</span>}
      </Label>
      {children}
      {description && (
        <p className="text-xs text-gray-600">{description}</p>
      )}
    </div>
  );
}

// Campo de texto básico
interface TextFieldProps extends BaseFieldProps {
  placeholder?: string;
  type?: 'text' | 'email' | 'tel';
  maxLength?: number;
  showCharCount?: boolean;
}

interface TextFieldWithFormProps extends TextFieldProps {
  form: any;
}

export function TextField({ 
  name, 
  label, 
  description, 
  required, 
  placeholder, 
  type = 'text',
  maxLength,
  showCharCount,
  className,
  form
}: TextFieldWithFormProps) {
  
  return (
    <form.Field
      name={name}
      children={(field) => (
        <FormField 
          name={name} 
          label={label} 
          description={description} 
          required={required}
          className={className}
        >
          <Input
            id={name}
            type={type}
            value={field.state.value as string}
            onChange={(e) => field.handleChange(e.target.value)}
            onBlur={field.handleBlur}
            placeholder={placeholder}
            maxLength={maxLength}
            className={field.state.meta.errors.length > 0 ? 'border-red-500' : ''}
          />
          
          {/* Indicador de validación async */}
          {field.state.meta.isValidating && (
            <div className="flex items-center gap-2 text-sm text-blue-600">
              <Loader2 className="h-3 w-3 animate-spin" />
              Validando...
            </div>
          )}
          
          {/* Contador de caracteres */}
          {showCharCount && maxLength && (
            <div className="flex justify-between items-center">
              <span className="text-xs text-gray-500">
                {(field.state.value as string)?.length || 0}/{maxLength} caracteres
              </span>
            </div>
          )}
          
          {/* Errores de validación */}
          {field.state.meta.errors.map((error, index) => (
            <Alert key={index} variant="destructive" className="py-2">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription className="text-sm">{typeof error === 'string' ? error : error.message || 'Error de validación'}</AlertDescription>
            </Alert>
          ))}
        </FormField>
      )}
    />
  );
}

// Campo de área de texto
interface TextAreaFieldProps extends BaseFieldProps {
  placeholder?: string;
  rows?: number;
  maxLength?: number;
  showCharCount?: boolean;
}

interface TextAreaFieldWithFormProps extends TextAreaFieldProps {
  form: any;
}

export function TextAreaField({ 
  name, 
  label, 
  description, 
  required, 
  placeholder, 
  rows = 4,
  maxLength,
  showCharCount,
  className,
  form
}: TextAreaFieldWithFormProps) {
  
  return (
    <form.Field
      name={name}
      children={(field) => (
        <FormField 
          name={name} 
          label={label} 
          description={description} 
          required={required}
          className={className}
        >
          <Textarea
            id={name}
            value={field.state.value as string}
            onChange={(e) => field.handleChange(e.target.value)}
            onBlur={field.handleBlur}
            placeholder={placeholder}
            rows={rows}
            maxLength={maxLength}
            className={field.state.meta.errors.length > 0 ? 'border-red-500' : ''}
          />
          
          {/* Indicador de validación async */}
          {field.state.meta.isValidating && (
            <div className="flex items-center gap-2 text-sm text-blue-600">
              <Loader2 className="h-3 w-3 animate-spin" />
              Validando...
            </div>
          )}
          
          {/* Contador de caracteres */}
          {showCharCount && maxLength && (
            <div className="flex justify-between items-center">
              <span className="text-xs text-gray-500">
                {(field.state.value as string)?.length || 0}/{maxLength} caracteres
              </span>
            </div>
          )}
          
          {/* Errores de validación */}
          {field.state.meta.errors.map((error, index) => (
            <Alert key={index} variant="destructive" className="py-2">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription className="text-sm">{typeof error === 'string' ? error : error.message || 'Error de validación'}</AlertDescription>
            </Alert>
          ))}
        </FormField>
      )}
    />
  );
}

// Campo numérico
interface NumberFieldProps extends BaseFieldProps {
  placeholder?: string;
  min?: number;
  max?: number;
  step?: number;
  showFormatted?: boolean;
  formatFunction?: (value: number) => string;
}

interface NumberFieldWithFormProps extends NumberFieldProps {
  form: any;
}

export function NumberField({ 
  name, 
  label, 
  description, 
  required, 
  placeholder, 
  min, 
  max, 
  step = 1,
  showFormatted = false,
  formatFunction,
  className,
  form
}: NumberFieldWithFormProps) {
  
  return (
    <form.Field
      name={name}
      children={(field) => (
        <FormField 
          name={name} 
          label={label} 
          description={description} 
          required={required}
          className={className}
        >
          <Input
            id={name}
            type="number"
            value={field.state.value as number}
            onChange={(e) => field.handleChange(Number(e.target.value))}
            onBlur={field.handleBlur}
            placeholder={placeholder}
            min={min}
            max={max}
            step={step}
            className={field.state.meta.errors.length > 0 ? 'border-red-500' : ''}
          />
          
          {/* Valor formateado */}
          {showFormatted && field.state.value && formatFunction && (
            <div className="flex items-center gap-2 text-sm text-gray-600">
              <CheckCircle className="h-3 w-3 text-green-500" />
              {formatFunction(field.state.value as number)}
            </div>
          )}
          
          {/* Indicador de validación async */}
          {field.state.meta.isValidating && (
            <div className="flex items-center gap-2 text-sm text-blue-600">
              <Loader2 className="h-3 w-3 animate-spin" />
              Validando...
            </div>
          )}
          
          {/* Errores de validación */}
          {field.state.meta.errors.map((error, index) => (
            <Alert key={index} variant="destructive" className="py-2">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription className="text-sm">{typeof error === 'string' ? error : error.message || 'Error de validación'}</AlertDescription>
            </Alert>
          ))}
        </FormField>
      )}
    />
  );
}

// Campo de selección
interface SelectFieldProps extends BaseFieldProps {
  options: { value: string; label: string }[];
  placeholder?: string;
}

interface SelectFieldWithFormProps extends SelectFieldProps {
  form: any;
}

export function SelectField({ 
  name, 
  label, 
  description, 
  required, 
  options,
  placeholder,
  className,
  form
}: SelectFieldWithFormProps) {
  
  return (
    <form.Field
      name={name}
      children={(field) => (
        <FormField 
          name={name} 
          label={label} 
          description={description} 
          required={required}
          className={className}
        >
          <Select 
            value={field.state.value as string} 
            onValueChange={field.handleChange}
          >
            <SelectTrigger className={field.state.meta.errors.length > 0 ? 'border-red-500' : ''}>
              <SelectValue placeholder={placeholder} />
            </SelectTrigger>
            <SelectContent>
              {options.map(option => (
                <SelectItem key={option.value} value={option.value}>
                  {option.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          
          {/* Errores de validación */}
          {field.state.meta.errors.map((error, index) => (
            <Alert key={index} variant="destructive" className="py-2">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription className="text-sm">{typeof error === 'string' ? error : error.message || 'Error de validación'}</AlertDescription>
            </Alert>
          ))}
        </FormField>
      )}
    />
  );
}

// Campo de checkbox
interface CheckboxFieldProps extends BaseFieldProps {
  description?: string;
}

interface CheckboxFieldWithFormProps extends CheckboxFieldProps {
  form: any;
}

export function CheckboxField({ 
  name, 
  label, 
  description, 
  className,
  form
}: CheckboxFieldWithFormProps) {
  
  return (
    <form.Field
      name={name}
      children={(field) => (
        <div className={`flex items-center space-x-2 ${className || ''}`}>
          <Checkbox
            id={name}
            checked={field.state.value as boolean}
            onCheckedChange={field.handleChange}
          />
          <Label htmlFor={name} className="text-sm font-normal">
            {label}
          </Label>
          {description && (
            <p className="text-xs text-gray-600 ml-2">{description}</p>
          )}
        </div>
      )}
    />
  );
}

// Campos específicos para propiedades
export function PropertyTitleField({ form }: { form: any }) {
  return (
    <TextField
      form={form}
      name="title"
      label="Título de la propiedad"
      placeholder="Ej: Hermosa casa en Samborondón con piscina"
      required
      maxLength={255}
      showCharCount
      description="Un título atractivo ayuda a captar la atención de los compradores"
    />
  );
}

export function PropertyDescriptionField({ form }: { form: any }) {
  return (
    <TextAreaField
      form={form}
      name="description"
      label="Descripción"
      placeholder="Describe las características principales de la propiedad..."
      required
      rows={4}
      maxLength={5000}
      showCharCount
      description="Describe la propiedad de manera detallada y atractiva"
    />
  );
}

export function PropertyPriceField({ form }: { form: any }) {
  return (
    <NumberField
      form={form}
      name="price"
      label="Precio (USD)"
      placeholder="285000"
      required
      min={1000}
      max={999999999}
      showFormatted
      formatFunction={formatPrice}
      description="Precio de venta en dólares americanos"
    />
  );
}

export function PropertyTypeField({ form }: { form: any }) {
  return (
    <SelectField
      form={form}
      name="type"
      label="Tipo de propiedad"
      options={PROPERTY_TYPES}
      placeholder="Selecciona el tipo"
      required
    />
  );
}

export function PropertyStatusField({ form }: { form: any }) {
  return (
    <SelectField
      form={form}
      name="status"
      label="Estado"
      options={PROPERTY_STATUS}
      placeholder="Selecciona el estado"
      required
    />
  );
}

export function PropertyProvinceField({ form }: { form: any }) {
  return (
    <SelectField
      form={form}
      name="province"
      label="Provincia"
      options={ECUADORIAN_PROVINCES.map(province => ({ value: province, label: province }))}
      placeholder="Selecciona la provincia"
      required
    />
  );
}

export function PropertyCityField({ form }: { form: any }) {
  return (
    <TextField
      form={form}
      name="city"
      label="Ciudad"
      placeholder="Ej: Samborondón"
      required
      maxLength={100}
    />
  );
}

export function PropertyAddressField({ form }: { form: any }) {
  return (
    <TextField
      form={form}
      name="address"
      label="Dirección completa"
      placeholder="Ej: Km 2.5 Vía Samborondón, Urbanización La Puntilla"
      required
      maxLength={500}
    />
  );
}

export function PropertyBedroomsField({ form }: { form: any }) {
  return (
    <NumberField
      form={form}
      name="bedrooms"
      label="Dormitorios"
      min={0}
      max={20}
      required
    />
  );
}

export function PropertyBathroomsField({ form }: { form: any }) {
  return (
    <NumberField
      form={form}
      name="bathrooms"
      label="Baños"
      min={0}
      max={20}
      step={0.5}
      required
    />
  );
}

export function PropertyAreaField({ form }: { form: any }) {
  return (
    <NumberField
      form={form}
      name="area_m2"
      label="Área (m²)"
      min={10}
      max={10000}
      required
      showFormatted
      formatFunction={(value) => `${value} metros cuadrados`}
    />
  );
}

export function PropertyParkingField({ form }: { form: any }) {
  return (
    <NumberField
      form={form}
      name="parking_spaces"
      label="Parqueaderos"
      min={0}
      max={20}
    />
  );
}

export function PropertyYearBuiltField({ form }: { form: any }) {
  return (
    <NumberField
      form={form}
      name="year_built"
      label="Año de construcción"
      min={1900}
      max={new Date().getFullYear()}
    />
  );
}

export function PropertyContactPhoneField({ form }: { form: any }) {
  return (
    <TextField
      form={form}
      name="contact_phone"
      label="Teléfono de contacto"
      type="tel"
      placeholder="0999999999"
      required
      maxLength={20}
    />
  );
}

export function PropertyContactEmailField({ form }: { form: any }) {
  return (
    <TextField
      form={form}
      name="contact_email"
      label="Email de contacto"
      type="email"
      placeholder="contacto@ejemplo.com"
      required
      maxLength={255}
    />
  );
}

export function PropertyNotesField({ form }: { form: any }) {
  return (
    <TextAreaField
      form={form}
      name="notes"
      label="Notas adicionales"
      placeholder="Información adicional sobre la propiedad..."
      rows={3}
      maxLength={1000}
      showCharCount
    />
  );
}

// Componente para las características adicionales
export function PropertyFeaturesFields({ form }: { form: any }) {
  const features = [
    { name: 'garden', label: 'Jardín' },
    { name: 'pool', label: 'Piscina' },
    { name: 'elevator', label: 'Ascensor' },
    { name: 'balcony', label: 'Balcón' },
    { name: 'terrace', label: 'Terraza' },
    { name: 'garage', label: 'Garaje' },
    { name: 'furnished', label: 'Amueblado' },
    { name: 'air_conditioning', label: 'Aire acondicionado' },
    { name: 'security', label: 'Seguridad' },
  ];

  return (
    <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
      {features.map((feature) => (
        <CheckboxField
          key={feature.name}
          form={form}
          name={feature.name as keyof PropertyFormData}
          label={feature.label}
        />
      ))}
    </div>
  );
}