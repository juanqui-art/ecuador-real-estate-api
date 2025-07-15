# 📚 Guía de Estudio: Dominio (internal/dominio/propiedad.go)

## 🎯 Objetivo de Esta Sesión
Entender completamente el archivo `propiedad.go` - el corazón de nuestro modelo de datos.

## 📖 Conceptos Clave a Aprender

### 1. **¿Qué es el Dominio?**
- Representa las **reglas de negocio** de una inmobiliaria
- Define **QUÉ es** una propiedad y **cómo se comporta**
- NO sabe nada de base de datos ni HTTP

### 2. **Comparación: Django vs Go**

#### En Django (Python):
```python
# models.py
class Propiedad(models.Model):
    titulo = models.CharField(max_length=255)
    precio = models.DecimalField(max_digits=15, decimal_places=2)
    provincia = models.CharField(max_length=100)
    
    def es_valida(self):
        return self.precio > 0 and self.titulo != ""
```

#### En Go (nuestro código):
```go
// propiedad.go
type Propiedad struct {
    Titulo    string  `json:"titulo" db:"titulo"`
    Precio    float64 `json:"precio" db:"precio"`
    Provincia string  `json:"provincia" db:"provincia"`
}

func (p *Propiedad) EsValida() bool {
    return p.Precio > 0 && p.Titulo != ""
}
```

## 🔍 Análisis Línea por Línea

### **Líneas 1-9: Package e Imports**
```go
package dominio

import (
    "regexp"
    "strings" 
    "time"
    "github.com/google/uuid"
)
```

**🤔 Preguntas para reflexionar:**
- ¿Por qué se llama `package dominio`?
- ¿Qué hace cada import?
- ¿Por qué usamos `github.com/google/uuid`?

**📝 Respuestas:**
- `package dominio` = Agrupa código relacionado con reglas de negocio
- `regexp` = Para limpiar texto en slugs
- `strings` = Para manipular strings
- `time` = Para fechas
- `uuid` = Para generar IDs únicos

### **Líneas 13-43: Struct Propiedad**
```go
type Propiedad struct {
    ID string `json:"id" db:"id"`
    Slug string `json:"slug" db:"slug"`
    Titulo string `json:"titulo" db:"titulo"`
    // ... más campos
}
```

**🤔 Preguntas importantes:**
1. ¿Qué son los `tags` como `json:"id"`?
2. ¿Por qué usamos `string` para ID y no `int`?
3. ¿Qué diferencia hay entre `float64` y `float32`?

**📝 Respuestas detalladas:**

#### **Tags en Go:**
```go
Titulo string `json:"titulo" db:"titulo"`
//             ↑               ↑
//         Para JSON       Para Base de Datos
```
- `json:"titulo"` = Cuando convertimos a JSON, usar "titulo"
- `db:"titulo"` = En base de datos, columna se llama "titulo"

#### **Tipos de datos:**
```go
ID          string    // Texto (UUID: "abc-123-def")
Precio      float64   // Número decimal (285000.50)
Dormitorios int       // Número entero (4)
Banos       float32   // Decimal más pequeño (3.5)
```

### **Líneas 46-69: Constructor NuevaPropiedad**
```go
func NuevaPropiedad(titulo, descripcion, provincia, ciudad, tipo string, precio float64) *Propiedad {
    id := uuid.New().String()
    slug := GenerarSlug(titulo, id)
    
    return &Propiedad{
        ID:     id,
        Slug:   slug,
        Titulo: titulo,
        // ...
    }
}
```

**🤔 Preguntas clave:**
1. ¿Por qué retorna `*Propiedad` y no `Propiedad`?
2. ¿Qué hace `uuid.New().String()`?
3. ¿Por qué usamos `&Propiedad{}` con &?

**📝 Respuestas:**

#### **Punteros vs Valores:**
```go
// Retornar VALOR (copia toda la struct)
func CrearPropiedad() Propiedad { ... }

// Retornar PUNTERO (solo dirección de memoria)
func NuevaPropiedad() *Propiedad { ... }  // ← Más eficiente
```

#### **El símbolo &:**
```go
propiedad := Propiedad{Titulo: "Casa"}  // Valor
puntero := &Propiedad{Titulo: "Casa"}   // Puntero al valor
```

### **Líneas 72-75: Método ActualizarFecha**
```go
func (p *Propiedad) ActualizarFecha() {
    p.FechaActualizacion = time.Now()
}
```

**🤔 Preguntas:**
1. ¿Qué significa `(p *Propiedad)` antes del nombre?
2. ¿Por qué `*Propiedad` y no `Propiedad`?

**📝 Respuestas:**

#### **Métodos en Go:**
```go
// Método = función que "pertenece" a una struct
func (p *Propiedad) ActualizarFecha() { ... }
//   ↑ 
//   "receiver" - como "self" en Python
```

#### **Receiver por Puntero:**
```go
func (p *Propiedad) ActualizarFecha() { ... }  // ← Modifica el original
func (p Propiedad) ActualizarFecha() { ... }   // ← Modifica una copia
```

### **Líneas 78-83: Método EsValida**
```go
func (p *Propiedad) EsValida() bool {
    return p.Titulo != "" &&
           p.Precio > 0 &&
           p.Provincia != "" &&
           p.Ciudad != "" &&
           p.Tipo != ""
}
```

**🤔 Preguntas:**
1. ¿Por qué las validaciones van aquí y no en el servicio?
2. ¿Qué hace el operador `&&`?

### **Líneas 86-101: Constantes**
```go
const (
    TipoCasa         = "casa"
    TipoDepartamento = "departamento"
    TipoTerreno      = "terreno"
    TipoComercial    = "comercial"
)
```

**🤔 Preguntas:**
1. ¿Por qué usar constantes en lugar de strings directos?
2. ¿Cuál es la ventaja?

**📝 Ventajas de constantes:**
```go
// ❌ Malo - fácil de escribir mal
if propiedad.Tipo == "casa" { ... }

// ✅ Bueno - el IDE autocompleta y detecta errores
if propiedad.Tipo == TipoCasa { ... }
```

## 🧪 Ejercicios Prácticos

### **Ejercicio 1: Agregar Campo**
Agregar un campo `Garage bool` a la struct:

1. Añadir el campo a la struct
2. Actualizar el constructor `NuevaPropiedad`
3. Probar que compila

### **Ejercicio 2: Nueva Validación**
Crear método `TieneGarage()` que retorne true si tiene garage.

### **Ejercicio 3: Nueva Constante**
Agregar constante para un nuevo estado: `EstadoMantenimiento = "mantenimiento"`

## ❓ Preguntas para Auto-evaluación

1. **¿Qué es una struct en Go?**
2. **¿Para qué sirven los tags como `json:"titulo"`?**
3. **¿Cuál es la diferencia entre un método y una función?**
4. **¿Por qué usamos punteros (`*Propiedad`)?**
5. **¿Qué ventaja tienen las constantes?**

## 🎯 Conceptos Dominados

Al terminar esta sesión debes entender:
- ✅ Structs y campos
- ✅ Tags para JSON y DB
- ✅ Constructores de structs
- ✅ Métodos vs funciones
- ✅ Punteros vs valores
- ✅ Constantes
- ✅ Validaciones básicas

## 📚 Próxima Sesión
**Archivo:** `internal/repositorio/propiedad.go`
**Conceptos:** Interfaces, SQL, database/sql

---

💡 **Tip:** Tomate el tiempo que necesites. Es mejor entender profundamente cada concepto que avanzar rápido sin comprenderlo.