# Modern Frontend Stack Guide

## 🚀 Stack Tecnológico Modernizado

### **Core Framework**
- **Next.js 15** con React 19
- **App Router** (no Pages Router)
- **React Server Components** optimizados
- **Turbopack** para desarrollo rápido

### **State Management**
- **Zustand** para estado global
- **TanStack Query** para estado del servidor
- **React Hook Form** reemplazado por TanStack Form
- **Zod** para validación de esquemas

### **API & Networking**
- **Fetch API nativo** (no axios)
- **Interceptors personalizados** para autenticación
- **Automatic token refresh** con refresh tokens
- **Error handling** centralizado

### **Forms & Validation**
- **TanStack Form** con Zod validation
- **Real-time validation** y error feedback
- **Type-safe forms** con TypeScript
- **Accessible form components**

### **UI/UX**
- **shadcn/ui** como base component system
- **Tailwind CSS** para styling
- **Framer Motion** para animaciones
- **Responsive design** mobile-first

## 🔧 Configuración y Uso

### **API Client**
```typescript
// /lib/api-client.ts
import { apiClient } from '@/lib/api-client';

// Uso básico
const response = await apiClient.get('/users');
const user = await apiClient.post('/users', userData);

// Con tipos
const response = await apiClient.get<User[]>('/users');
```

### **Forms con TanStack Form**
```typescript
// /components/forms/user-form.tsx
import { useForm } from '@tanstack/react-form';
import { zodValidator } from '@tanstack/zod-form-adapter';
import { z } from 'zod';

const userSchema = z.object({
  email: z.string().email(),
  name: z.string().min(2),
});

export function UserForm() {
  const form = useForm({
    defaultValues: {
      email: '',
      name: '',
    },
    onSubmit: async ({ value }) => {
      const result = await apiClient.post('/users', value);
      // Handle success
    },
    validatorAdapter: zodValidator,
  });

  return (
    <form onSubmit={(e) => {
      e.preventDefault();
      form.handleSubmit();
    }}>
      <form.Field
        name="email"
        validators={{
          onChange: userSchema.shape.email,
        }}
        children={(field) => (
          <Input
            value={field.state.value}
            onChange={(e) => field.handleChange(e.target.value)}
            placeholder="Email"
          />
        )}
      />
    </form>
  );
}
```

### **State Management con Zustand**
```typescript
// /store/auth.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  login: (user: User, tokens: Tokens) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,
      login: (user, tokens) => {
        set({ user, isAuthenticated: true });
        // Store tokens securely
      },
      logout: () => {
        set({ user: null, isAuthenticated: false });
      },
    }),
    {
      name: 'auth-storage',
    }
  )
);
```

### **Server State con TanStack Query**
```typescript
// /hooks/useUsers.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';

export function useUsers() {
  return useQuery({
    queryKey: ['users'],
    queryFn: () => apiClient.get<User[]>('/users'),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

export function useCreateUser() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (userData: CreateUserData) => 
      apiClient.post('/users', userData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });
}
```

## 🎨 Componentes y Patrones

### **Base Form Component**
```typescript
// /components/forms/base-form.tsx
import { useForm } from '@tanstack/react-form';
import { zodValidator } from '@tanstack/zod-form-adapter';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

interface BaseFormProps<T> {
  title: string;
  schema: z.ZodSchema<T>;
  defaultValues: T;
  onSubmit: (values: T) => Promise<void>;
  children: (form: FormApi<T>) => React.ReactNode;
}

export function BaseForm<T>({ 
  title, 
  schema, 
  defaultValues, 
  onSubmit, 
  children 
}: BaseFormProps<T>) {
  const form = useForm({
    defaultValues,
    onSubmit: async ({ value }) => {
      await onSubmit(value);
    },
    validatorAdapter: zodValidator,
  });

  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={(e) => {
          e.preventDefault();
          form.handleSubmit();
        }}>
          {children(form)}
          <Button type="submit" disabled={form.state.isSubmitting}>
            {form.state.isSubmitting ? 'Guardando...' : 'Guardar'}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
```

### **Error Boundary**
```typescript
// /components/ui/error-boundary.tsx
'use client';

import { Component, ErrorInfo, ReactNode } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false,
  };

  public static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Error caught by boundary:', error, errorInfo);
  }

  public render() {
    if (this.state.hasError) {
      return this.props.fallback || (
        <Card className="w-full max-w-md mx-auto">
          <CardHeader>
            <CardTitle>¡Oops! Algo salió mal</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-gray-600">
              Ha ocurrido un error inesperado. Por favor, intenta nuevamente.
            </p>
            <Button 
              onClick={() => this.setState({ hasError: false })}
              variant="outline"
            >
              Intentar de nuevo
            </Button>
          </CardContent>
        </Card>
      );
    }

    return this.props.children;
  }
}
```

## 🔒 Autenticación y Seguridad

### **Auth Hook**
```typescript
// /hooks/useAuth.ts
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api-client';
import { useAuthStore } from '@/store/auth';

export const useLogin = () => {
  const { login } = useAuthStore();
  const router = useRouter();
  
  return useMutation({
    mutationFn: (credentials: LoginData) => 
      apiClient.post('/auth/login', credentials),
    onSuccess: (data) => {
      login(data.user, data.tokens);
      router.push('/dashboard');
    },
  });
};

export const useLogout = () => {
  const { logout } = useAuthStore();
  const queryClient = useQueryClient();
  const router = useRouter();
  
  return useMutation({
    mutationFn: () => apiClient.post('/auth/logout'),
    onSuccess: () => {
      logout();
      queryClient.clear();
      router.push('/login');
    },
    onError: (error) => {
      // Handle logout errors gracefully
      if (error.status === 401) {
        // Token already invalid, proceed with logout
        logout();
        queryClient.clear();
        router.push('/login');
      }
    },
  });
};
```

### **Protected Route Component**
```typescript
// /components/auth/protected-route.tsx
'use client';

import { useAuthStore } from '@/store/auth';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';
import { Loading } from '@/components/ui/loading';

interface ProtectedRouteProps {
  children: React.ReactNode;
  requiredRole?: string[];
}

export function ProtectedRoute({ children, requiredRole }: ProtectedRouteProps) {
  const { isAuthenticated, user } = useAuthStore();
  const router = useRouter();

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login');
      return;
    }

    if (requiredRole && !requiredRole.includes(user?.role || '')) {
      router.push('/unauthorized');
      return;
    }
  }, [isAuthenticated, user, requiredRole, router]);

  if (!isAuthenticated) {
    return <Loading />;
  }

  if (requiredRole && !requiredRole.includes(user?.role || '')) {
    return <Loading />;
  }

  return <>{children}</>;
}
```

## 📱 Responsive Design

### **Mobile-First Approach**
```typescript
// /hooks/useMediaQuery.ts
import { useState, useEffect } from 'react';

export function useMediaQuery(query: string) {
  const [matches, setMatches] = useState(false);

  useEffect(() => {
    const media = window.matchMedia(query);
    if (media.matches !== matches) {
      setMatches(media.matches);
    }
    
    const listener = () => setMatches(media.matches);
    media.addEventListener('change', listener);
    return () => media.removeEventListener('change', listener);
  }, [matches, query]);

  return matches;
}

// Uso en componentes
export function ResponsiveComponent() {
  const isMobile = useMediaQuery('(max-width: 768px)');
  const isTablet = useMediaQuery('(max-width: 1024px)');

  return (
    <div className={`
      ${isMobile ? 'flex-col space-y-4' : 'flex-row space-x-4'}
      ${isTablet ? 'px-4' : 'px-8'}
    `}>
      {/* Content */}
    </div>
  );
}
```

## 🎯 Best Practices

### **1. Estructura de Archivos**
```
src/
├── app/                    # Next.js App Router pages
├── components/             # Componentes reutilizables
│   ├── ui/                # Componentes base (shadcn/ui)
│   ├── forms/             # Formularios específicos
│   └── layout/            # Componentes de layout
├── hooks/                 # Custom hooks
├── lib/                   # Utilidades y configuraciones
├── store/                 # Zustand stores
├── types/                 # TypeScript types
└── utils/                 # Funciones utilitarias
```

### **2. Naming Conventions**
- **Componentes**: PascalCase (`UserForm.tsx`)
- **Hooks**: camelCase con "use" prefix (`useAuth.ts`)
- **Stores**: camelCase con "use" prefix (`useAuthStore.ts`)
- **Types**: PascalCase (`User.ts`)
- **Utils**: camelCase (`formatDate.ts`)

### **3. Performance Optimizations**
```typescript
// Lazy loading de componentes
const DashboardChart = lazy(() => import('./dashboard-chart'));

// Memoización de componentes costosos
const MemoizedUserList = memo(UserList);

// Optimización de queries
const { data: users } = useQuery({
  queryKey: ['users', filters],
  queryFn: () => fetchUsers(filters),
  staleTime: 5 * 60 * 1000,
  cacheTime: 10 * 60 * 1000,
});
```

### **4. Testing Strategy**
```typescript
// Unit tests con Jest
describe('UserForm', () => {
  test('validates email format', async () => {
    render(<UserForm />);
    // Test implementation
  });
});

// Integration tests con React Testing Library
describe('Login Flow', () => {
  test('redirects to dashboard after successful login', async () => {
    // Test implementation
  });
});
```

## 🚀 Deployment

### **Build Optimization**
```bash
# Desarrollo
pnpm dev

# Producción
pnpm build
pnpm start

# Análisis de bundle
pnpm build --analyze
```

### **Environment Variables**
```env
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_APP_NAME=InmoEcuador
```

---

Esta guía proporciona las bases para trabajar con el stack modernizado del frontend. Para más detalles específicos, consulta la documentación oficial de cada tecnología.