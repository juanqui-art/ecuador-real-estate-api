# Troubleshooting Guide

## 🔧 Problemas Comunes y Soluciones

### **Autenticación**

#### **❌ Error: "missing or invalid authorization header" en logout**
**Síntomas:**
- Logout devuelve 401
- Headers de autorización no se envían
- Problema específico con Server Actions

**Solución:**
```typescript
// Problema: ApiClient no captura token antes de logout
// Solución implementada en /lib/api-client.ts

private async request<T>(endpoint: string, config: RequestConfig = {}): Promise<ApiResponse<T>> {
  // Manejo especial para logout - capturar token antes de que se limpie
  if (endpoint === '/auth/logout') {
    const accessToken = this.getAccessToken();
    if (accessToken) {
      config.headers = {
        ...config.headers,
        'Authorization': `Bearer ${accessToken}`
      };
    }
  }
  
  const headers = await this.prepareHeaders(config.headers);
  // ... resto del código
}
```

#### **❌ Error: "Token already expired" en auto-refresh**
**Síntomas:**
- Refresh token expirado
- Usuario debe hacer login manual
- Aplicación no maneja gracefully

**Solución:**
```typescript
// /lib/api-client.ts
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
    this.setTokens(data.access_token, data.refresh_token);
    return data.access_token;
  } catch (error) {
    // Limpiar tokens y redirigir
    this.clearTokens();
    if (typeof window !== 'undefined') {
      window.location.href = '/login';
    }
    return null;
  }
}
```

#### **❌ Error: "Usuario no autenticado" después de refresh**
**Síntomas:**
- Store state no sync con localStorage
- Pérdida de auth state en refresh

**Solución:**
```typescript
// /store/auth.ts - Agregar onRehydrateStorage
export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      // ... state
    }),
    {
      name: 'auth-storage',
      onRehydrateStorage: () => (state) => {
        if (state) {
          // Sync localStorage con store state
          if (state.access_token) {
            localStorage.setItem('access_token', state.access_token);
          }
          if (state.refresh_token) {
            localStorage.setItem('refresh_token', state.refresh_token);
          }
        }
      },
    }
  )
);
```

### **Formularios**

#### **❌ Error: "Validation error" en TanStack Form**
**Síntomas:**
- Validación no funciona correctamente
- Errores no se muestran

**Solución:**
```typescript
// Verificar que zodValidator esté configurado
import { zodValidator } from '@tanstack/zod-form-adapter';

const form = useForm({
  defaultValues: {
    email: '',
    password: '',
  },
  validatorAdapter: zodValidator, // ¡Importante!
  onSubmit: async ({ value }) => {
    // Handle submit
  },
});

// Verificar validación en campos
<form.Field
  name="email"
  validators={{
    onChange: z.string().email(), // Validación correcta
  }}
  children={(field) => (
    <div>
      <Input
        value={field.state.value}
        onChange={(e) => field.handleChange(e.target.value)}
      />
      {field.state.meta.errors.map((error) => (
        <p key={error} className="text-red-500">{error}</p>
      ))}
    </div>
  )}
/>
```

#### **❌ Error: "Form not submitting" en TanStack Form**
**Síntomas:**
- Button click no triggers submit
- onSubmit no se ejecuta

**Solución:**
```typescript
// Verificar que el form handler esté correcto
<form onSubmit={(e) => {
  e.preventDefault();
  e.stopPropagation();
  form.handleSubmit();
}}>
  {/* Form fields */}
  <Button type="submit">Submit</Button>
</form>

// Verificar que el button sea type="submit"
<Button type="submit" disabled={form.state.isSubmitting}>
  {form.state.isSubmitting ? 'Enviando...' : 'Enviar'}
</Button>
```

### **API y Networking**

#### **❌ Error: "Network request failed" en fetch**
**Síntomas:**
- Requests fallan sin razón aparente
- CORS errors

**Solución:**
```typescript
// Verificar headers correctos
const defaultHeaders = {
  'Content-Type': 'application/json',
  'User-Agent': 'RealtyCore-Dashboard/1.0 (Next.js)',
};

// Verificar URL base
const baseURL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Verificar CORS en backend
// Go backend debe incluir:
// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
```

#### **❌ Error: "Request timeout" en API calls**
**Síntomas:**
- Requests muy lentos
- Timeout errors

**Solución:**
```typescript
// Implementar timeout en requests
const controller = new AbortController();
const timeoutId = setTimeout(() => controller.abort(), 10000); // 10s

try {
  const response = await fetch(url, {
    ...options,
    signal: controller.signal,
  });
  clearTimeout(timeoutId);
  return response;
} catch (error) {
  if (error.name === 'AbortError') {
    throw new Error('Request timeout');
  }
  throw error;
}
```

### **State Management**

#### **❌ Error: "Hydration mismatch" en Zustand**
**Síntomas:**
- SSR/Client state mismatch
- Console warnings sobre hydration

**Solución:**
```typescript
// /hooks/useHydration.ts
import { useState, useEffect } from 'react';

export function useHydrated() {
  const [hydrated, setHydrated] = useState(false);

  useEffect(() => {
    setHydrated(true);
  }, []);

  return hydrated;
}

// Uso en componentes
export function AuthenticatedComponent() {
  const hydrated = useHydrated();
  const { isAuthenticated } = useAuthStore();

  if (!hydrated) {
    return <Loading />;
  }

  return isAuthenticated ? <Dashboard /> : <Login />;
}
```

#### **❌ Error: "Store not updating" en Zustand**
**Síntomas:**
- State changes no reflejan en UI
- Components no re-render

**Solución:**
```typescript
// Verificar que set() se llame correctamente
const useStore = create((set) => ({
  value: 0,
  increment: () => set((state) => ({ value: state.value + 1 })),
  // ❌ Incorrecto: set({ value: state.value + 1 })
  // ✅ Correcto: set((state) => ({ value: state.value + 1 }))
}));

// Verificar subscriptions
const { value, increment } = useStore(); // ✅ Correcto
// const value = useStore(state => state.value); // También correcto
```

### **Routing y Navigation**

#### **❌ Error: "Page not found" en Next.js App Router**
**Síntomas:**
- 404 errors en rutas válidas
- Routing no funciona

**Solución:**
```typescript
// Verificar estructura de archivos
app/
├── page.tsx          // /
├── login/
│   └── page.tsx      // /login
├── dashboard/
│   └── page.tsx      // /dashboard
└── layout.tsx        // Root layout

// Verificar que componentes sean default exports
// ❌ Incorrecto
export function LoginPage() { ... }

// ✅ Correcto
export default function LoginPage() { ... }
```

#### **❌ Error: "useRouter not working" en navigation**
**Síntomas:**
- Router push no funciona
- Navigation errors

**Solución:**
```typescript
// Verificar import correcto
import { useRouter } from 'next/navigation'; // ✅ App Router
// import { useRouter } from 'next/router'; // ❌ Pages Router

// Verificar uso correcto
const router = useRouter();

// Para navigation
router.push('/dashboard');
router.replace('/login');

// Para refresh
router.refresh();
```

### **Build y Deployment**

#### **❌ Error: "Module not found" en build**
**Síntomas:**
- Build fails con import errors
- Módulos no encontrados

**Solución:**
```bash
# Limpiar cache
rm -rf .next
rm -rf node_modules
pnpm install

# Verificar imports
# ❌ Incorrecto
import { Button } from '../../../components/ui/button';

# ✅ Correcto
import { Button } from '@/components/ui/button';
```

#### **❌ Error: "Environment variables not loaded"**
**Síntomas:**
- Variables de entorno undefined
- API URLs no funcionan

**Solución:**
```typescript
// Verificar que variables tengan prefix NEXT_PUBLIC_
// .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_APP_NAME=InmoEcuador

// Verificar uso correcto
const apiUrl = process.env.NEXT_PUBLIC_API_URL; // ✅ Correcto
const apiUrl = process.env.API_URL; // ❌ No funciona en client

// Para server-side
const secret = process.env.JWT_SECRET; // ✅ Correcto en server
```

### **Performance Issues**

#### **❌ Error: "Slow page loads" en desarrollo**
**Síntomas:**
- Páginas cargan lentamente
- Bundle size grande

**Solución:**
```typescript
// Implementar lazy loading
const DashboardChart = lazy(() => import('./dashboard-chart'));

// Usar Suspense
<Suspense fallback={<Loading />}>
  <DashboardChart />
</Suspense>

// Optimizar imports
// ❌ Incorrecto
import * as Icons from 'lucide-react';

// ✅ Correcto
import { Home, User, Settings } from 'lucide-react';
```

#### **❌ Error: "Memory leaks" en React Query**
**Síntomas:**
- Memoria aumenta con el tiempo
- App se vuelve lenta

**Solución:**
```typescript
// Configurar cache times correctamente
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutos
      cacheTime: 10 * 60 * 1000, // 10 minutos
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

// Limpiar queries cuando necesario
const queryClient = useQueryClient();
queryClient.removeQueries({ queryKey: ['old-data'] });
```

### **Database y Backend**

#### **❌ Error: "Connection refused" a PostgreSQL**
**Síntomas:**
- Backend no conecta a DB
- Connection errors

**Solución:**
```bash
# Verificar que PostgreSQL esté corriendo
brew services list | grep postgresql

# Verificar conexión
psql -h localhost -p 5433 -U juanquizhpi -d inmobiliaria_db

# Verificar variables de entorno
# .env (backend)
DATABASE_URL=postgresql://juanquizhpi@localhost:5433/inmobiliaria_db
```

#### **❌ Error: "SQL syntax error" en queries**
**Síntomas:**
- Queries fallan
- Database errors

**Solución:**
```go
// Verificar sintaxis PostgreSQL
// ❌ Incorrecto
query := "SELECT * FROM users WHERE role = ?"

// ✅ Correcto (PostgreSQL usa $1, $2, etc.)
query := "SELECT * FROM users WHERE role = $1"

// Verificar escaping
query := "SELECT * FROM users WHERE name = $1"
db.Query(query, userName) // ✅ Correcto
```

## 🔍 Debugging Tools

### **Browser DevTools**
```javascript
// Verificar estado de auth
console.log(useAuthStore.getState());

// Verificar tokens
console.log(localStorage.getItem('access_token'));

// Verificar network requests
// Network tab -> Filter by "XHR/Fetch"
```

### **React DevTools**
```typescript
// Instalar React DevTools extension
// Components tab -> Buscar componente
// Profiler tab -> Performance analysis
```

### **Next.js Debug Mode**
```bash
# Desarrollo con debug
DEBUG=* pnpm dev

# Específico para Next.js
DEBUG=next:* pnpm dev
```

### **Database Debugging**
```bash
# Logs de PostgreSQL
tail -f /usr/local/var/log/postgresql.log

# Query debugging en Go
log.Printf("Executing query: %s with params: %v", query, params)
```

## 📝 Logging Strategy

### **Frontend Logging**
```typescript
// /lib/logger.ts
interface LogEvent {
  level: 'info' | 'warn' | 'error';
  message: string;
  context?: Record<string, any>;
  timestamp: string;
}

export function log(level: LogEvent['level'], message: string, context?: Record<string, any>) {
  const event: LogEvent = {
    level,
    message,
    context,
    timestamp: new Date().toISOString(),
  };

  // Console para desarrollo
  if (process.env.NODE_ENV === 'development') {
    console[level](message, context);
  }

  // Enviar a servicio de logging en producción
  if (process.env.NODE_ENV === 'production') {
    // sendToLoggingService(event);
  }
}

// Uso
log('info', 'User logged in', { userId: user.id });
log('error', 'API request failed', { endpoint: '/users', error: error.message });
```

### **Error Boundary Logging**
```typescript
// /components/ui/error-boundary.tsx
public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
  log('error', 'Error boundary caught error', {
    error: error.message,
    stack: error.stack,
    errorInfo,
  });
}
```

## 🚨 Common Pitfalls

### **1. State Management**
- No mutar state directamente
- Usar set() function correctamente
- Verificar dependencies en useEffect

### **2. API Calls**
- Manejar loading states
- Implementar error handling
- Verificar response status

### **3. Forms**
- Validar en cliente y servidor
- Manejar loading states
- Limpiar form después de submit

### **4. Authentication**
- Verificar tokens antes de requests
- Manejar token expiration
- Limpiar state en logout

### **5. Performance**
- Implementar lazy loading
- Optimizar re-renders
- Usar React.memo cuando apropiado

---

Este troubleshooting guide cubre los problemas más comunes encontrados durante el desarrollo. Mantén este documento actualizado conforme encuentres nuevos issues y sus soluciones.