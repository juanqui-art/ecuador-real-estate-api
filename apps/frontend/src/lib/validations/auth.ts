import { z } from 'zod';

/**
 * Login form validation schema
 */
export const loginSchema = z.object({
  email: z
    .string()
    .min(1, 'El email es requerido')
    .email('Formato de email inválido'),
  password: z
    .string()
    .min(1, 'La contraseña es requerida')
    .min(6, 'La contraseña debe tener al menos 6 caracteres'),
});

/**
 * Change password form validation schema
 */
export const changePasswordSchema = z.object({
  current_password: z
    .string()
    .min(1, 'La contraseña actual es requerida'),
  new_password: z
    .string()
    .min(6, 'La nueva contraseña debe tener al menos 6 caracteres')
    .regex(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/, 
      'La contraseña debe contener al menos: una minúscula, una mayúscula y un número'),
  confirm_password: z
    .string()
    .min(1, 'La confirmación de contraseña es requerida'),
}).refine((data) => data.new_password === data.confirm_password, {
  message: 'Las contraseñas no coinciden',
  path: ['confirm_password'],
});

/**
 * Registration form validation schema (for admin creating users)
 */
export const registerSchema = z.object({
  email: z
    .string()
    .min(1, 'El email es requerido')
    .email('Formato de email inválido'),
  password: z
    .string()
    .min(6, 'La contraseña debe tener al menos 6 caracteres')
    .regex(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/, 
      'La contraseña debe contener al menos: una minúscula, una mayúscula y un número'),
  confirm_password: z
    .string()
    .min(1, 'La confirmación de contraseña es requerida'),
  first_name: z
    .string()
    .min(1, 'El nombre es requerido')
    .min(2, 'El nombre debe tener al menos 2 caracteres')
    .max(50, 'El nombre no puede tener más de 50 caracteres'),
  last_name: z
    .string()
    .min(1, 'El apellido es requerido')
    .min(2, 'El apellido debe tener al menos 2 caracteres')
    .max(50, 'El apellido no puede tener más de 50 caracteres'),
  role: z.enum(['admin', 'agency', 'agent', 'owner', 'buyer'], {
    required_error: 'El rol es requerido',
    invalid_type_error: 'Rol inválido',
  }),
  phone: z
    .string()
    .optional()
    .refine((val) => {
      if (!val) return true; // Optional field
      // Ecuador phone number format: +593XXXXXXXXX or 09XXXXXXXX
      return /^(\+593|0)[1-9]\d{8}$/.test(val);
    }, 'Formato de teléfono ecuatoriano inválido'),
}).refine((data) => data.password === data.confirm_password, {
  message: 'Las contraseñas no coinciden',
  path: ['confirm_password'],
});

/**
 * Profile update validation schema
 */
export const profileUpdateSchema = z.object({
  first_name: z
    .string()
    .min(1, 'El nombre es requerido')
    .min(2, 'El nombre debe tener al menos 2 caracteres')
    .max(50, 'El nombre no puede tener más de 50 caracteres'),
  last_name: z
    .string()
    .min(1, 'El apellido es requerido')
    .min(2, 'El apellido debe tener al menos 2 caracteres')
    .max(50, 'El apellido no puede tener más de 50 caracteres'),
  phone: z
    .string()
    .optional()
    .refine((val) => {
      if (!val) return true; // Optional field
      // Ecuador phone number format: +593XXXXXXXXX or 09XXXXXXXX
      return /^(\+593|0)[1-9]\d{8}$/.test(val);
    }, 'Formato de teléfono ecuatoriano inválido'),
  avatar_url: z
    .string()
    .url('URL de avatar inválida')
    .optional()
    .or(z.literal('')),
  bio: z
    .string()
    .max(500, 'La biografía no puede tener más de 500 caracteres')
    .optional(),
});

// TypeScript types derived from schemas
export type LoginFormData = z.infer<typeof loginSchema>;
export type ChangePasswordFormData = z.infer<typeof changePasswordSchema>;
export type RegisterFormData = z.infer<typeof registerSchema>;
export type ProfileUpdateFormData = z.infer<typeof profileUpdateSchema>;