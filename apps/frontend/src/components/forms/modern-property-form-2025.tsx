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
                    <CheckCircle className="w-5 h-5 text-blue-600"/>
                    Información Básica
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
                            <CheckCircle className="w-5 h-5 text-blue-600"/>
                            Ubicación
                        </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {/* Province */}
                            <div>
                                <Label htmlFor="province">Provincia *</Label>
                                <select
                                    id="province"
                                    name="province"
                                    required
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
                                <Label htmlFor="city">Ciudad *</Label>
                                <Input
                                    id="city"
                                    name="city"
                                    placeholder="Ej: Samborondón"
                                    required
                                    minLength={2}
                                    className={state.errors?.city ? 'border-red-500' : ''}
                                />
                                {state.errors?.city && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.city[0]}</p>
                                )}
                            </div>

                            {/* Address */}
                            <div className="col-span-full">
                                <Label htmlFor="address">Dirección completa *</Label>
                                <Input
                                    id="address"
                                    name="address"
                                    placeholder="Ej: Km 2.5 Vía Samborondón, Urbanización La Puntilla"
                                    required
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
                            <CheckCircle className="w-5 h-5 text-blue-600"/>
                            Características
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                            {/* Bedrooms */}
                            <div>
                                <Label htmlFor="bedrooms">Dormitorios *</Label>
                                <Input
                                    id="bedrooms"
                                    name="bedrooms"
                                    type="number"
                                    min="0"
                                    max="20"
                                    defaultValue="3"
                                    required
                                    className={state.errors?.bedrooms ? 'border-red-500' : ''}
                                />
                                {state.errors?.bedrooms && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.bedrooms[0]}</p>
                                )}
                            </div>

                            {/* Bathrooms */}
                            <div>
                                <Label htmlFor="bathrooms">Baños *</Label>
                                <Input
                                    id="bathrooms"
                                    name="bathrooms"
                                    type="number"
                                    min="0"
                                    max="20"
                                    step="0.5"
                                    defaultValue="2"
                                    required
                                    className={state.errors?.bathrooms ? 'border-red-500' : ''}
                                />
                                {state.errors?.bathrooms && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.bathrooms[0]}</p>
                                )}
                            </div>

                            {/* Area */}
                            <div>
                                <Label htmlFor="area_m2">Área (m²) *</Label>
                                <Input
                                    id="area_m2"
                                    name="area_m2"
                                    type="number"
                                    min="10"
                                    max="10000"
                                    step="10"
                                    placeholder="320"
                                    required
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
                                    defaultValue="2"
                                    className={state.errors?.parking_spaces ? 'border-red-500' : ''}
                                />
                                {state.errors?.parking_spaces && (
                                    <p className="text-sm text-red-500 mt-1">{state.errors.parking_spaces[0]}</p>
                                )}
                            </div>

                            {/* Year Built */}
                            <div className="col-span-2">
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

                {/* Contact Information */}
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <CheckCircle className="w-5 h-5 text-blue-600"/>
                            Información de Contacto
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