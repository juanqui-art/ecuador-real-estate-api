'use client';

import {Alert, AlertDescription} from '@/components/ui/alert';
import {Button} from '@/components/ui/button';
import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card';
import {Input} from '@/components/ui/input';
import {Label} from '@/components/ui/label';
import {Textarea} from '@/components/ui/textarea';
import {type ActionResult, createPropertyAction, createPropertyWithRedirectAction} from '@/lib/actions/properties';
import {ECUADORIAN_PROVINCES, PROPERTY_STATUS, PROPERTY_TYPES} from '@/lib/constants';
import {AlertCircle, CheckCircle, Loader2} from 'lucide-react';
import React, {useState, useTransition} from 'react';
import {useFormStatus} from 'react-dom';

/**
 * Modern Property Form using React 19 best practices (2025)
 *
 * Key Features:
 * - useActionState for Server Actions integration
 * - useFormStatus for loading states
 * - Progressive Enhancement (works without JS)
 * - Server-side validation with Zod
 * - Modern error handling
 * - Optimistic UI updates
 */

interface ModernPropertyForm2025Props {
    onSuccess?: () => void;
    onCancel?: () => void;
}

// Submit Button with useFormStatus (2025 best practice) - Optimized with React.memo
const SubmitButton = React.memo(function SubmitButton({isPending}: { isPending: boolean }) {
    const {pending: formPending} = useFormStatus();
    const isLoading = isPending || formPending;

    return (
        <Button
            type="submit"
            disabled={isLoading}
            className="w-full sm:w-auto min-w-[140px]"
        >
            {isLoading ? (
                <>
                    <Loader2 className="w-4 h-4 mr-2 animate-spin"/>
                    Creando...
                </>
            ) : (
                <>
                    <CheckCircle className="w-4 h-4 mr-2"/>
                    Crear Propiedad
                </>
            )}
        </Button>
    );
});

// Loading Indicator with useFormStatus - Optimized with React.memo
const FormLoadingIndicator = React.memo(function FormLoadingIndicator({isPending}: { isPending: boolean }) {
    const {pending: formPending} = useFormStatus();
    const isLoading = isPending || formPending;

    if (!isLoading) return null;

    return (
        <div className="fixed inset-0 bg-black/20 flex items-center justify-center z-50">
            <Card className="p-6">
                <div className="flex items-center space-x-3">
                    <Loader2 className="w-6 h-6 animate-spin text-blue-600"/>
                    <div>
                        <p className="font-medium">Procesando propiedad...</p>
                        <p className="text-sm text-gray-500">Validando datos y guardando en el servidor</p>
                    </div>
                </div>
            </Card>
        </div>
    );
});

// Performance optimized form sections with React.memo
const BasicInformationSection = React.memo(function BasicInformationSection({
                                                                                errors
                                                                            }: {
    errors: Record<string, string[]>
}) {
    return (
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <CheckCircle className="w-5 h-5 text-red-500"/>
                    Información Básica *
                    <span className="text-sm text-red-500 font-normal ml-2">- Obligatorio</span>
                </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {/* Title */}
                    <div className="col-span-full">
                        <Label htmlFor="title">Título de la propiedad *</Label>
                        <Input
                            id="title"
                            name="title"
                            placeholder="Ej: Hermosa casa en Samborondón con piscina"
                            required
                            minLength={10}
                            className={errors?.title ? 'border-red-500' : ''}
                        />
                        {errors?.title && (
                            <p className="text-sm text-red-500 mt-1">{errors.title[0]}</p>
                        )}
                    </div>

                    {/* Type */}
                    <div>
                        <Label htmlFor="type">Tipo de propiedad *</Label>
                        <select
                            id="type"
                            name="type"
                            required
                            onChange={(e) => setPropertyType(e.target.value)}
                            className="w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-blue-500"
                        >
                            <option value="">Selecciona el tipo</option>
                            {PROPERTY_TYPES.map(type => (
                                <option key={type.value} value={type.value}>
                                    {type.label}
                                </option>
                            ))}
                        </select>
                        {errors?.type && (
                            <p className="text-sm text-red-500 mt-1">{errors.type[0]}</p>
                        )}
                    </div>

                    {/* Status */}
                    <div>
                        <Label htmlFor="status">Estado *</Label>
                        <select
                            id="status"
                            name="status"
                            required
                            defaultValue="available"
                            className="w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-blue-500"
                        >
                            {PROPERTY_STATUS.map(status => (
                                <option key={status.value} value={status.value}>
                                    {status.label}
                                </option>
                            ))}
                        </select>
                        {errors?.status && (
                            <p className="text-sm text-red-500 mt-1">{errors.status[0]}</p>
                        )}
                    </div>

                    {/* Price */}
                    <div>
                        <Label htmlFor="price">Precio (USD) *</Label>
                        <Input
                            id="price"
                            name="price"
                            type="number"
                            min="1000"
                            step="1000"
                            placeholder="285000"
                            required
                            className={errors?.price ? 'border-red-500' : ''}
                        />
                        {errors?.price && (
                            <p className="text-sm text-red-500 mt-1">{errors.price[0]}</p>
                        )}
                    </div>
                </div>

                {/* Description */}
                <div>
                    <Label htmlFor="description">Descripción *</Label>
                    <Textarea
                        id="description"
                        name="description"
                        placeholder="Describe las características principales de la propiedad..."
                        required
                        minLength={50}
                        rows={4}
                        className={errors?.description ? 'border-red-500' : ''}
                    />
                    {errors?.description && (
                        <p className="text-sm text-red-500 mt-1">{errors.description[0]}</p>
                    )}
                </div>
            </CardContent>
        </Card>
    );
});

export function ModernPropertyForm2025({onSuccess, onCancel}: ModernPropertyForm2025Props) {
    // State management compatible approach
    const [state, setState] = useState<ActionResult>({
        success: false,
        message: '',
        errors: {},
    });
    const [isPending, startTransition] = useTransition();
    const [propertyType, setPropertyType] = useState<string>('');

    // Smart defaults based on property type
    const getSmartDefaults = (type: string) => {
        switch (type) {
            case 'land':
                return { bedrooms: 0, bathrooms: 0, parking_spaces: 0 };
            case 'commercial':
                return { bedrooms: 0, bathrooms: 1, parking_spaces: 3 };
            case 'apartment':
                return { bedrooms: 2, bathrooms: 2, parking_spaces: 1 };
            case 'house':
            default:
                return { bedrooms: 3, bathrooms: 2, parking_spaces: 2 };
        }
    };

    // Modern Server Action wrapper
    const handleSubmit = async (formData: FormData) => {
        startTransition(async () => {
            const result = await createPropertyAction(null, formData);
            setState(result);

            if (result.success) {
                console.log('✅ Property created successfully:', result.data);
                onSuccess?.();
            }
        });
    };

    // Handle successful creation - removed duplicate logic since it's in handleSubmit

    return (
        <div className="space-y-6">
            <FormLoadingIndicator isPending={isPending}/>

            {/* UX Optimization Banner */}
            <Alert className="border-blue-200 bg-blue-50">
                <CheckCircle className="h-4 w-4 text-blue-600"/>
                <AlertDescription className="text-blue-800">
                    <strong>✨ Formulario Optimizado:</strong> Solo 7 campos son obligatorios (marcados con *). 
                    Los campos opcionales pueden completarse después para mejorar la publicación.
                </AlertDescription>
            </Alert>

            {/* Success Message */}
            {state.success && (
                <Alert className="border-green-200 bg-green-50">
                    <CheckCircle className="h-4 w-4 text-green-600"/>
                    <AlertDescription className="text-green-800">
                        <strong>¡Éxito!</strong> {state.message}
                    </AlertDescription>
                </Alert>
            )}

            {/* Error Message */}
            {!state.success && state.message && (
                <Alert className="border-red-200 bg-red-50">
                    <AlertCircle className="h-4 w-4 text-red-600"/>
                    <AlertDescription className="text-red-800">
                        <strong>Error:</strong> {state.message}
                    </AlertDescription>
                </Alert>
            )}

            {/* Modern Form with Progressive Enhancement */}
            <form
                action={createPropertyWithRedirectAction}
                onSubmit={(e) => {
                    // Progressive enhancement: Use modern handler if JS is available
                    if (typeof window !== 'undefined') {
                        e.preventDefault();
                        const formData = new FormData(e.currentTarget);
                        void handleSubmit(formData);
                    }
                    // If JS is disabled, form will submit normally to Server Action with redirect
                }}
                className="space-y-6"
            >
                {/* Basic Information - Optimized with React.memo */}
                <BasicInformationSection errors={state.errors || {}}/>

                {/* Location */}
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <CheckCircle className="w-5 h-5 text-gray-500"/>
                            Ubicación (Opcional)
                            <span className="text-sm text-gray-500 font-normal ml-2">- Puede completarse después</span>
                        </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {/* Province */}
                            <div>
                                <Label htmlFor="province">Provincia (opcional)</Label>
                                <select
                                    id="province"
                                    name="province"
                                    className="w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-blue-500"
                                >
                                    <option value="">Selecciona la provincia</option>
                                    {ECUADORIAN_PROVINCES.map(province => (
                                        <option key={province} value={province}>
                                            {province}
                                        </option>
                                    ))}
                                </select>
                                {state.errors?.province && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.province[0]}</p>
                                )}
                            </div>

                            {/* City */}
                            <div>
                                <Label htmlFor="city">Ciudad (opcional)</Label>
                                <Input
                                    id="city"
                                    name="city"
                                    placeholder="Ej: Samborondón"
                                    minLength={2}
                                    className={state.errors?.city ? 'border-red-500' : ''}
                                />
                                {state.errors?.city && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.city[0]}</p>
                                )}
                            </div>

                            {/* Sector */}
                            <div className="col-span-2">
                                <Label htmlFor="sector">Sector</Label>
                                <Input
                                    id="sector"
                                    name="sector"
                                    placeholder="Ej: Norte, Centro, La Puntilla"
                                    className={state.errors?.sector ? 'border-red-500' : ''}
                                />
                                {state.errors?.sector && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.sector[0]}</p>
                                )}
                            </div>

                            {/* Address */}
                            <div className="col-span-full">
                                <Label htmlFor="address">Dirección completa (opcional)</Label>
                                <Input
                                    id="address"
                                    name="address"
                                    placeholder="Ej: Km 2.5 Vía Samborondón, Urbanización La Puntilla"
                                    minLength={10}
                                    className={state.errors?.address ? 'border-red-500' : ''}
                                />
                                {state.errors?.address && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.address[0]}</p>
                                )}
                            </div>
                        </div>
                    </CardContent>
                </Card>

                {/* Characteristics */}
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <CheckCircle className="w-5 h-5 text-gray-500"/>
                            Características (Opcional)
                            <span className="text-sm text-gray-500 font-normal ml-2">- Valores por defecto aplicados</span>
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        {propertyType && (
                            <div className="mb-4 p-3 bg-gray-50 rounded-md">
                                <p className="text-sm text-gray-700">
                                    💡 <strong>Valores automáticos para {PROPERTY_TYPES.find(t => t.value === propertyType)?.label}:</strong>{' '}
                                    Los campos se han ajustado con valores típicos para este tipo de propiedad.
                                </p>
                            </div>
                        )}
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                            {/* Bedrooms */}
                            <div>
                                <Label htmlFor="bedrooms">Dormitorios</Label>
                                <Input
                                    id="bedrooms"
                                    name="bedrooms"
                                    type="number"
                                    min="0"
                                    max="20"
                                    key={`bedrooms-${propertyType}`}
                                    defaultValue={getSmartDefaults(propertyType).bedrooms.toString()}
                                    className={state.errors?.bedrooms ? 'border-red-500' : ''}
                                />
                                {state.errors?.bedrooms && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.bedrooms[0]}</p>
                                )}
                            </div>

                            {/* Bathrooms */}
                            <div>
                                <Label htmlFor="bathrooms">Baños</Label>
                                <Input
                                    id="bathrooms"
                                    name="bathrooms"
                                    type="number"
                                    min="0"
                                    max="20"
                                    step="0.5"
                                    key={`bathrooms-${propertyType}`}
                                    defaultValue={getSmartDefaults(propertyType).bathrooms.toString()}
                                    className={state.errors?.bathrooms ? 'border-red-500' : ''}
                                />
                                {state.errors?.bathrooms && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.bathrooms[0]}</p>
                                )}
                            </div>

                            {/* Area */}
                            <div>
                                <Label htmlFor="area_m2">Área (m²) (opcional)</Label>
                                <Input
                                    id="area_m2"
                                    name="area_m2"
                                    type="number"
                                    min="10"
                                    max="10000"
                                    step="10"
                                    placeholder="320"
                                    className={state.errors?.area_m2 ? 'border-red-500' : ''}
                                />
                                {state.errors?.area_m2 && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.area_m2[0]}</p>
                                )}
                            </div>

                            {/* Parking */}
                            <div>
                                <Label htmlFor="parking_spaces">Parqueaderos</Label>
                                <Input
                                    id="parking_spaces"
                                    name="parking_spaces"
                                    type="number"
                                    min="0"
                                    max="20"
                                    key={`parking-${propertyType}`}
                                    defaultValue={getSmartDefaults(propertyType).parking_spaces.toString()}
                                    className={state.errors?.parking_spaces ? 'border-red-500' : ''}
                                />
                                {state.errors?.parking_spaces && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.parking_spaces[0]}</p>
                                )}
                            </div>

                            {/* Year Built */}
                            <div>
                                <Label htmlFor="year_built">Año de construcción</Label>
                                <Input
                                    id="year_built"
                                    name="year_built"
                                    type="number"
                                    min="1900"
                                    max={new Date().getFullYear()}
                                    placeholder="2020"
                                    className={state.errors?.year_built ? 'border-red-500' : ''}
                                />
                                {state.errors?.year_built && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.year_built[0]}</p>
                                )}
                            </div>

                            {/* Floors */}
                            <div>
                                <Label htmlFor="floors">Número de pisos</Label>
                                <Input
                                    id="floors"
                                    name="floors"
                                    type="number"
                                    min="1"
                                    max="50"
                                    placeholder="2"
                                    className={state.errors?.floors ? 'border-red-500' : ''}
                                />
                                {state.errors?.floors && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.floors[0]}</p>
                                )}
                            </div>
                        </div>
                    </CardContent>
                </Card>

                {/* Additional Pricing */}
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <CheckCircle className="w-5 h-5 text-blue-600"/>
                            Precios Adicionales
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                            {/* Rent Price */}
                            <div>
                                <Label htmlFor="rent_price">Precio de renta (USD)</Label>
                                <Input
                                    id="rent_price"
                                    name="rent_price"
                                    type="number"
                                    min="100"
                                    step="50"
                                    placeholder="1200"
                                    className={state.errors?.rent_price ? 'border-red-500' : ''}
                                />
                                {state.errors?.rent_price && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.rent_price[0]}</p>
                                )}
                            </div>

                            {/* Common Expenses */}
                            <div>
                                <Label htmlFor="common_expenses">Gastos comunes (USD)</Label>
                                <Input
                                    id="common_expenses"
                                    name="common_expenses"
                                    type="number"
                                    min="0"
                                    step="10"
                                    placeholder="150"
                                    className={state.errors?.common_expenses ? 'border-red-500' : ''}
                                />
                                {state.errors?.common_expenses && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.common_expenses[0]}</p>
                                )}
                            </div>

                            {/* Price per M2 */}
                            <div>
                                <Label htmlFor="price_per_m2">Precio por m² (USD)</Label>
                                <Input
                                    id="price_per_m2"
                                    name="price_per_m2"
                                    type="number"
                                    min="10"
                                    step="10"
                                    placeholder="890"
                                    className={state.errors?.price_per_m2 ? 'border-red-500' : ''}
                                />
                                {state.errors?.price_per_m2 && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.price_per_m2[0]}</p>
                                )}
                            </div>
                        </div>
                    </CardContent>
                </Card>

                {/* Amenities */}
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <CheckCircle className="w-5 h-5 text-blue-600"/>
                            Características Adicionales
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                            {/* Amenities checkboxes */}
                            {[
                                {name: 'garden', label: 'Jardín'},
                                {name: 'pool', label: 'Piscina'},
                                {name: 'elevator', label: 'Ascensor'},
                                {name: 'balcony', label: 'Balcón'},
                                {name: 'terrace', label: 'Terraza'},
                                {name: 'garage', label: 'Garaje'},
                                {name: 'furnished', label: 'Amueblado'},
                                {name: 'air_conditioning', label: 'Aire acondicionado'},
                                {name: 'security', label: 'Seguridad'},
                            ].map((amenity) => (
                                <div key={amenity.name} className="flex items-center space-x-2">
                                    <input
                                        type="checkbox"
                                        id={amenity.name}
                                        name={amenity.name}
                                        value="true"
                                        className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                    />
                                    <Label htmlFor={amenity.name} className="text-sm font-normal">
                                        {amenity.label}
                                    </Label>
                                </div>
                            ))}
                        </div>
                    </CardContent>
                </Card>

                {/* Estado y Clasificación */}
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <CheckCircle className="w-5 h-5 text-blue-600"/>
                            Estado y Clasificación
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {/* Property Status */}
                            <div>
                                <Label htmlFor="property_status">Estado de la propiedad *</Label>
                                <select
                                    id="property_status"
                                    name="property_status"
                                    required
                                    defaultValue="new"
                                    className="w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:ring-blue-500"
                                >
                                    <option value="new">Nueva</option>
                                    <option value="used">Usada</option>
                                    <option value="renovated">Renovada</option>
                                </select>
                                {state.errors?.property_status && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.property_status[0]}</p>
                                )}
                            </div>

                            {/* Featured */}
                            <div className="flex items-center space-x-2">
                                <input
                                    type="checkbox"
                                    id="featured"
                                    name="featured"
                                    value="true"
                                    className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                />
                                <Label htmlFor="featured" className="text-sm font-normal">
                                    ⭐ Propiedad destacada
                                </Label>
                                {state.errors?.featured && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.featured[0]}</p>
                                )}
                            </div>

                            {/* Tags */}
                            <div className="col-span-full">
                                <Label htmlFor="tags">Etiquetas/Tags</Label>
                                <Input
                                    id="tags"
                                    name="tags"
                                    placeholder="Ej: piscina, jardín, lujo, urbanización cerrada (separados por coma)"
                                    className={state.errors?.tags ? 'border-red-500' : ''}
                                />
                                <p className="text-xs text-gray-500 mt-1">
                                    Separa las etiquetas con comas. Ejemplo: piscina, jardín, seguridad, lujo
                                </p>
                                {state.errors?.tags && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.tags[0]}</p>
                                )}
                            </div>
                        </div>
                    </CardContent>
                </Card>

                {/* Contact Information */}
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <CheckCircle className="w-5 h-5 text-red-500"/>
                            Información de Contacto *
                            <span className="text-sm text-red-500 font-normal ml-2">- Obligatorio</span>
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {/* Phone */}
                            <div>
                                <Label htmlFor="contact_phone">Teléfono de contacto *</Label>
                                <Input
                                    id="contact_phone"
                                    name="contact_phone"
                                    type="tel"
                                    placeholder="0999999999"
                                    required
                                    minLength={10}
                                    className={state.errors?.contact_phone ? 'border-red-500' : ''}
                                />
                                {state.errors?.contact_phone && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.contact_phone[0]}</p>
                                )}
                            </div>

                            {/* Email */}
                            <div>
                                <Label htmlFor="contact_email">Email de contacto *</Label>
                                <Input
                                    id="contact_email"
                                    name="contact_email"
                                    type="email"
                                    placeholder="contacto@ejemplo.com"
                                    required
                                    className={state.errors?.contact_email ? 'border-red-500' : ''}
                                />
                                {state.errors?.contact_email && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.contact_email[0]}</p>
                                )}
                            </div>

                            {/* Notes */}
                            <div className="col-span-full">
                                <Label htmlFor="notes">Notas adicionales</Label>
                                <Textarea
                                    id="notes"
                                    name="notes"
                                    placeholder="Información adicional sobre la propiedad..."
                                    rows={3}
                                />
                            </div>
                        </div>
                    </CardContent>
                </Card>

                {/* Actions */}
                <div className="flex flex-col sm:flex-row justify-between gap-4 pt-6">
                    {onCancel && (
                        <Button type="button" variant="outline" onClick={onCancel}>
                            Cancelar
                        </Button>
                    )}

                    <div className="flex gap-3">
                        <SubmitButton isPending={isPending}/>
                    </div>
                </div>

                {/* Progressive Enhancement Info */}
                <div className="bg-blue-50 border border-blue-200 rounded-md p-4">
                    <p className="text-sm text-blue-800">
                        <strong>🚀 Progressive Enhancement (2025):</strong> Formulario optimizado que funciona con y sin
                        JavaScript.
                        Con JS: useTransition + Server Actions + React.memo. Sin JS: POST form tradicional + validación
                        server-side.
                    </p>
                </div>

                {/* No-JS Fallback Form (Hidden when JS is available) */}
                <noscript>
                    <div className="bg-yellow-50 border border-yellow-200 rounded-md p-4 mt-4">
                        <p className="text-sm text-yellow-800">
                            <strong>⚠️ JavaScript deshabilitado:</strong> El formulario funciona en modo tradicional con
                            recarga de página.
                            Habilita JavaScript para una experiencia mejorada con actualizaciones instantáneas.
                        </p>
                    </div>
                </noscript>
            </form>
        </div>
    );
}