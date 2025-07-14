# 🧠 Conceptos Fundamentales de Go

## 🎯 Para Desarrolladores Python/Django

Esta guía explica conceptos Go usando comparaciones con Python para facilitar el aprendizaje.

## 1. 📦 Packages vs Modules

### **Python:**
```python
# En Python
from myapp.models import Propiedad
import json
```

### **Go:**
```go
// En Go
package dominio  // Define el package actual

import (
    "encoding/json"                    // Package estándar
    "realty-core/internal/repositorio" // Package local
)
```

**🔑 Diferencias clave:**
- Go: `package` al inicio define el package actual
- Go: Imports van en bloque `import ( ... )`
- Go: Solo puedes usar lo que importas explícitamente

## 2. 🏗️ Structs vs Classes

### **Python Class:**
```python
class Propiedad:
    def __init__(self, titulo, precio):
        self.titulo = titulo
        self.precio = precio
        self.fecha_creacion = datetime.now()
    
    def es_valida(self):
        return self.precio > 0
```

### **Go Struct:**
```go
type Propiedad struct {
    Titulo         string    `json:"titulo"`
    Precio         float64   `json:"precio"`
    FechaCreacion  time.Time `json:"fecha_creacion"`
}

// Constructor (función)
func NuevaPropiedad(titulo string, precio float64) *Propiedad {
    return &Propiedad{
        Titulo:        titulo,
        Precio:        precio,
        FechaCreacion: time.Now(),
    }
}

// Método
func (p *Propiedad) EsValida() bool {
    return p.Precio > 0
}
```

**🔑 Diferencias clave:**
- Go: Struct define campos, métodos se definen por separado
- Go: No hay `__init__`, usas funciones constructoras
- Go: Métodos tienen "receiver" `(p *Propiedad)`

## 3. 🎭 Interfaces vs Abstract Classes

### **Python (ABC):**
```python
from abc import ABC, abstractmethod

class RepositoryInterface(ABC):
    @abstractmethod
    def crear(self, propiedad):
        pass
    
    @abstractmethod
    def obtener_por_id(self, id):
        pass

class PropiedadRepository(RepositoryInterface):
    def crear(self, propiedad):
        # Implementación real
        pass
```

### **Go Interface:**
```go
// Interface define QUÉ métodos debe tener
type PropiedadRepository interface {
    Crear(propiedad *Propiedad) error
    ObtenerPorID(id string) (*Propiedad, error)
}

// Implementación (NO necesita declarar que implementa)
type PropiedadRepositoryPostgres struct {
    db *sql.DB
}

func (r *PropiedadRepositoryPostgres) Crear(propiedad *Propiedad) error {
    // Implementación real
    return nil
}

func (r *PropiedadRepositoryPostgres) ObtenerPorID(id string) (*Propiedad, error) {
    // Implementación real
    return nil, nil
}
```

**🔑 Diferencias clave:**
- Go: Interface automáticamente implementada si tienes los métodos
- Go: "Duck typing" pero verificado en compile time
- Go: No necesitas heredar o declarar explícitamente

## 4. 🔧 Punteros vs Referencias

### **Python (todo son referencias):**
```python
propiedad = Propiedad("Casa", 100000)
otra_referencia = propiedad  # Ambas apuntan al mismo objeto
otra_referencia.precio = 200000
print(propiedad.precio)  # 200000 - cambió!
```

### **Go (valores y punteros explícitos):**
```go
// VALOR - se copia toda la struct
propiedad := Propiedad{Titulo: "Casa", Precio: 100000}
copia := propiedad          // Se copia toda la struct
copia.Precio = 200000
fmt.Println(propiedad.Precio) // 100000 - NO cambió

// PUNTERO - apunta al mismo lugar en memoria
propiedad := &Propiedad{Titulo: "Casa", Precio: 100000}
referencia := propiedad      // Apuntan al mismo lugar
referencia.Precio = 200000
fmt.Println(propiedad.Precio) // 200000 - SÍ cambió
```

**🔑 Cuándo usar cada uno:**
```go
// Usar VALOR cuando:
func CalcularImpuesto(prop Propiedad) float64 {  // Solo lees
    return prop.Precio * 0.1
}

// Usar PUNTERO cuando:
func ActualizarPrecio(prop *Propiedad, nuevoPrecio float64) {  // Modificas
    prop.Precio = nuevoPrecio
}
```

## 5. 🚨 Manejo de Errores

### **Python (excepciones):**
```python
def obtener_propiedad(id):
    try:
        # Buscar en BD
        return propiedad
    except DatabaseError as e:
        raise ValueError(f"Error al buscar: {e}")
```

### **Go (valores de retorno):**
```go
func (r *Repository) ObtenerPropiedad(id string) (*Propiedad, error) {
    propiedad, err := r.db.Query("SELECT * FROM propiedades WHERE id = $1", id)
    if err != nil {
        return nil, fmt.Errorf("error al buscar: %w", err)
    }
    return propiedad, nil
}

// Uso:
propiedad, err := repo.ObtenerPropiedad("123")
if err != nil {
    log.Printf("Error: %v", err)
    return
}
// usar propiedad...
```

**🔑 Principios:**
- Go: Errores son valores, no excepciones
- Go: Siempre manejar el error inmediatamente
- Go: `if err != nil` es el patrón estándar

## 6. 🏷️ Tags y Anotaciones

### **Python (decorators/annotations):**
```python
from dataclasses import dataclass
from typing import Optional

@dataclass
class Propiedad:
    titulo: str
    precio: float
    descripcion: Optional[str] = None
```

### **Go (struct tags):**
```go
type Propiedad struct {
    Titulo      string  `json:"titulo" db:"titulo" validate:"required,min=3"`
    Precio      float64 `json:"precio" db:"precio" validate:"gt=0"`
    Descripcion string  `json:"descripcion,omitempty" db:"descripcion"`
}
```

**🔑 Tags comunes:**
- `json:"titulo"` - Nombre en JSON
- `json:",omitempty"` - Omitir si está vacío
- `db:"titulo"` - Nombre de columna en BD
- `validate:"required"` - Para validación

## 7. 🔤 Tipos de Datos

### **Comparación:**

| Python | Go | Uso |
|--------|-----|-----|
| `str` | `string` | Texto |
| `int` | `int`, `int64` | Números enteros |
| `float` | `float32`, `float64` | Decimales |
| `bool` | `bool` | Verdadero/Falso |
| `list` | `[]string` | Arrays/Slices |
| `dict` | `map[string]int` | Mapas |
| `None` | `nil` | Valor nulo |

### **Ejemplos Go:**
```go
// Básicos
var nombre string = "Juan"
var edad int = 30
var precio float64 = 99.99
var activo bool = true

// Slices (como listas Python)
var nombres []string = []string{"Juan", "María"}
nombres = append(nombres, "Pedro")  // Agregar elemento

// Maps (como diccionarios Python)
var precios map[string]float64 = map[string]float64{
    "casa": 100000,
    "departamento": 50000,
}
```

## 8. 🔄 Control de Flujo

### **Python:**
```python
# If/else
if precio > 100000:
    print("Caro")
elif precio > 50000:
    print("Medio")
else:
    print("Barato")

# For loop
for propiedad in propiedades:
    print(propiedad.titulo)

# List comprehension
precios = [p.precio for p in propiedades if p.tipo == "casa"]
```

### **Go:**
```go
// If/else
if precio > 100000 {
    fmt.Println("Caro")
} else if precio > 50000 {
    fmt.Println("Medio")
} else {
    fmt.Println("Barato")
}

// For loop (única forma de loop en Go)
for _, propiedad := range propiedades {
    fmt.Println(propiedad.Titulo)
}

// Filtrar slice
var precios []float64
for _, p := range propiedades {
    if p.Tipo == "casa" {
        precios = append(precios, p.Precio)
    }
}
```

## 9. 🧪 Testing

### **Python:**
```python
import unittest

class TestPropiedad(unittest.TestCase):
    def test_es_valida(self):
        prop = Propiedad("Casa", 100000)
        self.assertTrue(prop.es_valida())
```

### **Go:**
```go
package dominio

import "testing"

func TestPropiedad_EsValida(t *testing.T) {
    prop := NuevaPropiedad("Casa", "Desc", "Guayas", "Guayaquil", "casa", 100000)
    
    if !prop.EsValida() {
        t.Error("Propiedad debería ser válida")
    }
}
```

## 10. 🚀 Mejores Prácticas

### **Naming Conventions:**
```go
// ✅ Correcto
type PropiedadService struct {}      // PascalCase para públicos
var precioMaximo float64            // camelCase para privados
const TipoCasa = "casa"             // PascalCase para constantes

// ❌ Incorrecto
type propiedad_service struct {}     // No usar snake_case
var PrecioMaximo float64            // No PascalCase para privados
```

### **Error Handling:**
```go
// ✅ Correcto
result, err := DoSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// ❌ Incorrecto  
result, _ := DoSomething()  // Ignorar errores
```

### **Imports:**
```go
// ✅ Correcto - agrupados y ordenados
import (
    // Standard library
    "fmt"
    "time"
    
    // Third party
    "github.com/google/uuid"
    
    // Local
    "realty-core/internal/dominio"
)
```

---

💡 **Recuerda:** Go prioriza simplicidad y claridad sobre "features" avanzados. Es menos "mágico" que Python pero más explícito y predecible.