# 🔄 Go vs Django/Python - Comparación Detallada

## 🎯 Para Desarrolladores Django que Aprenden Go

Esta guía traduce conceptos Django/Python a Go usando nuestro proyecto inmobiliario como ejemplo.

---

## 📋 **Estructura de Proyecto**

### **Django Project:**
```
inmobiliaria/
├── manage.py
├── inmobiliaria/
│   ├── settings.py
│   ├── urls.py
│   └── wsgi.py
├── propiedades/
│   ├── models.py
│   ├── views.py
│   ├── urls.py
│   ├── serializers.py
│   └── admin.py
└── requirements.txt
```

### **Go Project (nuestro):**
```
realty-core/
├── main.go
├── internal/
│   ├── dominio/     # ≈ models.py
│   ├── repositorio/ # ≈ ORM queries
│   ├── servicio/    # ≈ business logic
│   └── web/         # ≈ views.py + urls.py
├── cmd/servidor/    # ≈ manage.py runserver
└── go.mod          # ≈ requirements.txt
```

**🔑 Mapeo conceptual:**
- `models.py` → `internal/dominio/`
- `views.py` → `internal/web/handlers/`
- `urls.py` → `internal/web/routes.go`
- ORM queries → `internal/repositorio/`
- Business logic → `internal/servicio/`

---

## 🏗️ **Modelos vs Structs**

### **Django Model:**
```python
# models.py
from django.db import models

class Propiedad(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    titulo = models.CharField(max_length=255)
    precio = models.DecimalField(max_digits=15, decimal_places=2)
    provincia = models.CharField(max_length=100)
    tipo = models.CharField(max_length=50, choices=TIPO_CHOICES)
    fecha_creacion = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        db_table = 'propiedades'
        ordering = ['-fecha_creacion']
    
    def __str__(self):
        return self.titulo
    
    def es_cara(self):
        return self.precio > 200000
```

### **Go Struct (nuestro dominio):**
```go
// internal/dominio/propiedad.go
type Propiedad struct {
    ID            string    `json:"id" db:"id"`
    Titulo        string    `json:"titulo" db:"titulo"`
    Precio        float64   `json:"precio" db:"precio"`
    Provincia     string    `json:"provincia" db:"provincia"`
    Tipo          string    `json:"tipo" db:"tipo"`
    FechaCreacion time.Time `json:"fecha_creacion" db:"fecha_creacion"`
}

// Constructor (no hay __init__ en Go)
func NuevaPropiedad(titulo string, precio float64, provincia string) *Propiedad {
    return &Propiedad{
        ID:            uuid.New().String(),
        Titulo:        titulo,
        Precio:        precio,
        Provincia:     provincia,
        FechaCreacion: time.Now(),
    }
}

// Método (como método de clase)
func (p *Propiedad) EsCara() bool {
    return p.Precio > 200000
}
```

**🔑 Diferencias principales:**
- Django: Modelo incluye lógica de BD (save, delete, etc.)
- Go: Struct solo datos, BD separada en repositorio
- Django: Meta class para configuración
- Go: Tags en campos para configuración
- Django: __str__ automático
- Go: Métodos explícitos

---

## 🗄️ **ORM vs Repository Pattern**

### **Django ORM:**
```python
# views.py o managers.py
from .models import Propiedad

# Crear
propiedad = Propiedad.objects.create(
    titulo="Casa en Guayaquil",
    precio=150000,
    provincia="Guayas"
)

# Leer
propiedades = Propiedad.objects.all()
propiedad = Propiedad.objects.get(id=1)
caras = Propiedad.objects.filter(precio__gt=200000)

# Actualizar
propiedad.precio = 160000
propiedad.save()

# Eliminar
propiedad.delete()
```

### **Go Repository (nuestro código):**
```go
// internal/repositorio/propiedad.go

// Interface define qué operaciones hay
type PropiedadRepository interface {
    Crear(propiedad *dominio.Propiedad) error
    ObtenerTodas() ([]dominio.Propiedad, error)
    ObtenerPorID(id string) (*dominio.Propiedad, error)
    Actualizar(propiedad *dominio.Propiedad) error
    Eliminar(id string) error
}

// Implementación con SQL directo
func (r *PropiedadRepositoryPostgres) Crear(propiedad *dominio.Propiedad) error {
    query := `INSERT INTO propiedades (id, titulo, precio, provincia) VALUES ($1, $2, $3, $4)`
    _, err := r.db.Exec(query, propiedad.ID, propiedad.Titulo, propiedad.Precio, propiedad.Provincia)
    return err
}

func (r *PropiedadRepositoryPostgres) ObtenerTodas() ([]dominio.Propiedad, error) {
    query := `SELECT id, titulo, precio, provincia FROM propiedades`
    rows, err := r.db.Query(query)
    // ... manejar filas y convertir a struct
}
```

**🔑 Comparación:**
| Aspecto | Django ORM | Go Repository |
|---------|------------|---------------|
| **Consultas** | `Propiedad.objects.filter()` | SQL directo |
| **Migraciones** | `makemigrations` automático | Scripts SQL manuales |
| **Relaciones** | ForeignKey automático | Joins manuales |
| **Validación** | En modelo + forms | En servicio |
| **Caché** | QuerySet lazy | Manual si necesitas |

---

## 🎮 **Views vs Handlers**

### **Django Views (Function-Based):**
```python
# views.py
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.views.decorators.http import require_http_methods
import json

@csrf_exempt
@require_http_methods(["GET", "POST"])
def propiedades_list(request):
    if request.method == 'GET':
        propiedades = Propiedad.objects.all().values()
        return JsonResponse(list(propiedades), safe=False)
    
    elif request.method == 'POST':
        data = json.loads(request.body)
        propiedad = Propiedad.objects.create(**data)
        return JsonResponse({
            'id': propiedad.id,
            'titulo': propiedad.titulo,
            'precio': str(propiedad.precio)
        }, status=201)

@csrf_exempt
def propiedad_detail(request, pk):
    try:
        propiedad = Propiedad.objects.get(pk=pk)
    except Propiedad.DoesNotExist:
        return JsonResponse({'error': 'No encontrada'}, status=404)
    
    if request.method == 'GET':
        return JsonResponse({
            'id': propiedad.id,
            'titulo': propiedad.titulo,
            'precio': str(propiedad.precio)
        })
```

### **Go Handlers (nuestro código):**
```go
// internal/web/handlers/propiedad.go
func (h *PropiedadHandler) ListarPropiedades(w http.ResponseWriter, r *http.Request) {
    // Verificar método HTTP
    if r.Method != http.MethodGet {
        h.responderError(w, http.StatusMethodNotAllowed, "Método no permitido")
        return
    }
    
    // Obtener datos del servicio
    propiedades, err := h.servicio.ListarPropiedades()
    if err != nil {
        h.responderError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    // Responder JSON
    h.responderExito(w, http.StatusOK, propiedades, "Propiedades obtenidas exitosamente")
}

func (h *PropiedadHandler) CrearPropiedad(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        h.responderError(w, http.StatusMethodNotAllowed, "Método no permitido")
        return
    }
    
    // Decodificar JSON
    var req CrearPropiedadRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.responderError(w, http.StatusBadRequest, "JSON inválido")
        return
    }
    
    // Usar servicio
    propiedad, err := h.servicio.CrearPropiedad(req.Titulo, req.Descripcion, req.Provincia, req.Ciudad, req.Tipo, req.Precio)
    if err != nil {
        h.responderError(w, http.StatusBadRequest, err.Error())
        return
    }
    
    h.responderExito(w, http.StatusCreated, propiedad, "Propiedad creada exitosamente")
}
```

**🔑 Diferencias:**
- Django: Decoradores para métodos HTTP
- Go: if statements para verificar métodos
- Django: `request.POST` automático
- Go: `json.NewDecoder` manual
- Django: Retorno directo de JsonResponse
- Go: Helper functions para respuestas

---

## 🛣️ **URLs vs Routes**

### **Django URLs:**
```python
# propiedades/urls.py
from django.urls import path
from . import views

urlpatterns = [
    path('api/propiedades/', views.propiedades_list, name='propiedades-list'),
    path('api/propiedades/<uuid:pk>/', views.propiedad_detail, name='propiedad-detail'),
    path('api/propiedades/filtrar/', views.filtrar_propiedades, name='filtrar'),
]

# proyecto/urls.py
from django.contrib import admin
from django.urls import path, include

urlpatterns = [
    path('admin/', admin.site.urls),
    path('', include('propiedades.urls')),
]
```

### **Go Routes (nuestro código):**
```go
// internal/web/routes.go
func ConfigurarRutas(propiedadHandler *handlers.PropiedadHandler) *http.ServeMux {
    mux := http.NewServeMux()
    
    // Ruta para listar/crear propiedades
    mux.HandleFunc("/api/propiedades", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            propiedadHandler.ListarPropiedades(w, r)
        case http.MethodPost:
            propiedadHandler.CrearPropiedad(w, r)
        default:
            http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        }
    })
    
    // Ruta para operaciones específicas
    mux.HandleFunc("/api/propiedades/", propiedadHandler.ObtenerPropiedad)
    mux.HandleFunc("/api/propiedades/filtrar", propiedadHandler.FiltrarPropiedades)
    
    return mux
}
```

**🔑 Diferencias:**
- Django: Un handler por URL con decoradores
- Go: Switch en método HTTP en mismo handler
- Django: Parámetros automáticos en URL `<uuid:pk>`
- Go: Parsing manual de URL path
- Django: `reverse()` para URLs nombradas
- Go: URLs hardcoded (puedes hacer helpers)

---

## ⚙️ **Settings vs Configuration**

### **Django Settings:**
```python
# settings.py
import os
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent.parent

DEBUG = True
ALLOWED_HOSTS = []

DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.postgresql',
        'NAME': 'inmobiliaria_db',
        'USER': 'postgres',
        'PASSWORD': 'password',
        'HOST': 'localhost',
        'PORT': '5432',
    }
}

REST_FRAMEWORK = {
    'DEFAULT_PAGINATION_CLASS': 'rest_framework.pagination.PageNumberPagination',
    'PAGE_SIZE': 20
}
```

### **Go Configuration (nuestro código):**
```go
// cmd/servidor/main.go
func main() {
    // Cargar variables de entorno
    if err := godotenv.Load(); err != nil {
        log.Println("Archivo .env no encontrado")
    }
    
    // Leer configuración
    databaseURL := obtenerVariable("DATABASE_URL", "postgresql://localhost/inmobiliaria_db")
    puerto := obtenerVariable("PORT", "8080")
    logLevel := obtenerVariable("LOG_LEVEL", "info")
    
    // Configurar base de datos
    db, err := repositorio.ConectarBaseDatos(databaseURL)
    if err != nil {
        log.Fatalf("Error conectando BD: %v", err)
    }
    defer db.Close()
}

func obtenerVariable(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

**🔑 Diferencias:**
- Django: `settings.py` centralizado y automático
- Go: Variables de entorno manuales
- Django: Configuración por proyecto/entorno automática
- Go: Logic manual para cada variable
- Django: Muchas configuraciones built-in
- Go: Todo explícito

---

## 🧪 **Testing**

### **Django Tests:**
```python
# tests.py
from django.test import TestCase, Client
from django.urls import reverse
from .models import Propiedad
import json

class PropiedadTestCase(TestCase):
    def setUp(self):
        self.client = Client()
        self.propiedad = Propiedad.objects.create(
            titulo="Casa Test",
            precio=100000,
            provincia="Guayas"
        )
    
    def test_listar_propiedades(self):
        response = self.client.get(reverse('propiedades-list'))
        self.assertEqual(response.status_code, 200)
        data = response.json()
        self.assertEqual(len(data), 1)
    
    def test_crear_propiedad(self):
        data = {
            'titulo': 'Casa Nueva',
            'precio': 150000,
            'provincia': 'Pichincha'
        }
        response = self.client.post(
            reverse('propiedades-list'),
            data=json.dumps(data),
            content_type='application/json'
        )
        self.assertEqual(response.status_code, 201)
```

### **Go Tests:**
```go
// internal/dominio/propiedad_test.go
package dominio

import "testing"

func TestPropiedad_EsValida(t *testing.T) {
    // Crear propiedad válida
    propiedad := NuevaPropiedad("Casa Test", "Descripción", "Guayas", "Guayaquil", "casa", 100000)
    
    if !propiedad.EsValida() {
        t.Error("Propiedad debería ser válida")
    }
    
    // Probar propiedad inválida
    propiedad.Precio = 0
    if propiedad.EsValida() {
        t.Error("Propiedad con precio 0 no debería ser válida")
    }
}

// internal/web/handlers/propiedad_test.go
func TestPropiedadHandler_ListarPropiedades(t *testing.T) {
    // Setup
    mockRepo := &MockRepository{}
    servicio := servicio.NuevoPropiedadService(mockRepo)
    handler := handlers.NuevoPropiedadHandler(servicio)
    
    // Test
    req := httptest.NewRequest("GET", "/api/propiedades", nil)
    w := httptest.NewRecorder()
    
    handler.ListarPropiedades(w, req)
    
    // Verificar
    if w.Code != http.StatusOK {
        t.Errorf("Esperaba 200, obtuvo %d", w.Code)
    }
}
```

**🔑 Diferencias:**
- Django: Base de datos en memoria automática
- Go: Mocks manuales
- Django: TestCase con setUp/tearDown
- Go: Funciones individuales
- Django: Cliente HTTP integrado
- Go: httptest package
- Django: Fixtures y factorías
- Go: Datos de prueba manuales

---

## 🚀 **Deployment**

### **Django Deployment:**
```python
# gunicorn_config.py
bind = "0.0.0.0:8000"
workers = 4
worker_class = "gevent"

# requirements.txt
Django==4.2.0
psycopg2-binary==2.9.5
djangorestframework==3.14.0
gunicorn==20.1.0

# Dockerfile
FROM python:3.11
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY . .
CMD ["gunicorn", "inmobiliaria.wsgi:application"]
```

### **Go Deployment:**
```go
// main.go ya es ejecutable
// go build -o servidor cmd/servidor/main.go

# Dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o servidor cmd/servidor/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/servidor .
CMD ["./servidor"]
```

**🔑 Diferencias:**
- Django: Necesitas WSGI server (gunicorn)
- Go: Binario standalone
- Django: Interprete Python en producción
- Go: No dependencies en runtime
- Django: requirements.txt + virtual env
- Go: go.mod + compilación estática

---

## 📊 **Cuándo Usar Cada Uno**

### **Usar Django cuando:**
- ✅ Desarrollo rápido (RAD)
- ✅ Admin interface necesaria
- ✅ ORM complex queries
- ✅ Muchas librerías disponibles
- ✅ Prototipado rápido
- ✅ Equipo Python existente

### **Usar Go cuando:**
- ✅ Performance crítico
- ✅ Microservicios
- ✅ APIs simples y rápidas
- ✅ Deployment simple
- ✅ Concurrencia pesada
- ✅ Equipos nuevos (menos curva aprendizaje)

---

## 🎯 **Resumen para Desarrollador Django**

### **Lo que extrañarás de Django:**
- ORM potente con QuerySets
- Admin interface automática
- Migraciones automáticas
- Django Rest Framework
- Ecosystem gigante

### **Lo que amarás de Go:**
- Velocidad de ejecución
- Deployment super simple
- Menos "magia", más control
- Concurrencia nativa
- Binarios standalone
- Menos memory footprint

### **Curva de aprendizaje:**
1. **Semana 1-2:** Sintaxis Go básica
2. **Semana 3-4:** Structs, interfaces, punteros
3. **Semana 5-6:** HTTP handlers, JSON
4. **Semana 7-8:** Patterns (repository, service)
5. **Mes 2+:** Go idioms y optimizaciones

---

💡 **Conclusión:** Go y Django resuelven problemas similares de formas muy diferentes. Django es más "mágico" y rápido para prototipar, Go es más explícito y eficiente en runtime. ¡Ambos tienen su lugar!