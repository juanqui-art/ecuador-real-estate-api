# 🧪 Ejercicios Prácticos - Realty Core

## 🎯 Objetivo
Ejercicios graduales para dominar cada parte del código. Haz UN ejercicio a la vez.

---

## 📁 **Ejercicios: internal/dominio/propiedad.go**

### **🟢 Ejercicio 1: Agregar Campo Básico**
**Nivel:** Principiante  
**Tiempo:** 15 minutos

**Tarea:** Agregar campo `Garage bool` a la struct Propiedad.

**Pasos:**
1. Abrir `internal/dominio/propiedad.go`
2. Agregar `Garage bool` después del campo `AreaM2`
3. Actualizar función `NuevaPropiedad` para incluir garage
4. Compilar: `go build ./...`

**Código esperado:**
```go
type Propiedad struct {
    // ... campos existentes
    AreaM2  float64 `json:"area_m2" db:"area_m2"`
    Garage  bool    `json:"garage" db:"garage"`     // ← Nuevo campo
    // ... resto de campos
}
```

**✅ Verificación:** El código compila sin errores.

### **🟡 Ejercicio 2: Método de Validación**
**Nivel:** Intermedio  
**Tiempo:** 20 minutos

**Tarea:** Crear método `TieneComodidades()` que retorne true si tiene más de 2 baños Y garage.

**Código esperado:**
```go
func (p *Propiedad) TieneComodidades() bool {
    return p.Banos > 2 && p.Garage
}
```

**Prueba:**
```go
// En main.go temporal, agregar:
prop := dominio.NuevaPropiedad("Casa", "Desc", "Guayas", "Guayaquil", "casa", 100000)
prop.Banos = 3
prop.Garage = true
fmt.Println(prop.TieneComodidades()) // Debe imprimir: true
```

### **🟠 Ejercicio 3: Constante Nueva**
**Nivel:** Principiante  
**Tiempo:** 10 minutos

**Tarea:** Agregar nueva constante para estado "mantenimiento".

**Código esperado:**
```go
const (
    EstadoDisponible = "disponible"
    EstadoVendida    = "vendida"
    EstadoAlquilada  = "alquilada"
    EstadoReservada  = "reservada"
    EstadoMantenimiento = "mantenimiento" // ← Nueva constante
)
```

---

## 📁 **Ejercicios: internal/repositorio/propiedad.go**

### **🟢 Ejercicio 4: Entender Query SQL**
**Nivel:** Principiante  
**Tiempo:** 30 minutos

**Tarea:** Leer y explicar qué hace cada query en el repositorio.

**Preguntas para responder:**
1. ¿Qué hace `$1, $2, $3` en las queries?
2. ¿Por qué usamos `QueryRow` vs `Query`?
3. ¿Qué hace `rows.Close()`?

**Actividad:**
Agregar comentarios explicativos en el código:
```go
// Esta query inserta una nueva propiedad en la tabla
query := `INSERT INTO propiedades (...) VALUES ($1, $2, $3...)`

// $1, $2, $3 son placeholders para prevenir SQL injection
_, err := r.db.Exec(query, propiedad.ID, propiedad.Titulo, ...)
```

### **🟡 Ejercicio 5: Nueva Query Simple**
**Nivel:** Intermedio  
**Tiempo:** 45 minutos

**Tarea:** Agregar método `ContarPropiedades()` que retorne el número total.

**Pasos:**
1. Agregar método a la interface:
```go
type PropiedadRepository interface {
    // ... métodos existentes
    ContarPropiedades() (int, error)
}
```

2. Implementar en struct:
```go
func (r *PropiedadRepositoryPostgres) ContarPropiedades() (int, error) {
    query := `SELECT COUNT(*) FROM propiedades`
    
    var count int
    err := r.db.QueryRow(query).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("error al contar propiedades: %w", err)
    }
    
    return count, nil
}
```

**✅ Verificación:** Compilar y probar que retorna número correcto.

### **🔴 Ejercicio 6: Query con Filtro**
**Nivel:** Avanzado  
**Tiempo:** 60 minutos

**Tarea:** Crear método `BuscarPorPrecio(min, max float64)` que filtre por rango.

**Pista:**
```sql
SELECT * FROM propiedades WHERE precio >= $1 AND precio <= $2
```

---

## 📁 **Ejercicios: internal/servicio/propiedad.go**

### **🟢 Ejercicio 7: Nueva Validación**
**Nivel:** Principiante  
**Tiempo:** 20 minutos

**Tarea:** Agregar validación que el área debe ser mayor a 10 m².

**En función `validarDatosCreacion`:**
```go
if areaM2 <= 10 {
    return fmt.Errorf("área debe ser mayor a 10 m²")
}
```

### **🟡 Ejercicio 8: Método de Servicio**
**Nivel:** Intermedio  
**Tiempo:** 30 minutos

**Tarea:** Crear método `ObtenerPropiedadesCaras()` que retorne propiedades > $200,000.

**Estructura:**
```go
func (s *PropiedadService) ObtenerPropiedadesCaras() ([]dominio.Propiedad, error) {
    // 1. Obtener todas las propiedades del repositorio
    // 2. Filtrar las que cuestan > 200000
    // 3. Retornar la lista filtrada
}
```

---

## 📁 **Ejercicios: internal/web/handlers/propiedad.go**

### **🟢 Ejercicio 9: Entender HTTP**
**Nivel:** Principiante  
**Tiempo:** 30 minutos

**Tarea:** Agregar logs para ver qué requests llegan.

**En cada handler, agregar:**
```go
log.Printf("Recibida petición %s %s", r.Method, r.URL.Path)
```

### **🟡 Ejercicio 10: Nuevo Endpoint Simple**
**Nivel:** Intermedio  
**Tiempo:** 45 minutos

**Tarea:** Crear endpoint `GET /api/propiedades/count` que retorne el total.

**Pasos:**
1. Agregar método al handler:
```go
func (h *PropiedadHandler) ContarPropiedades(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        h.responderError(w, http.StatusMethodNotAllowed, "Método no permitido")
        return
    }
    
    count, err := h.servicio.ContarPropiedades() // Necesitas crear este método
    if err != nil {
        h.responderError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    resultado := map[string]int{"total": count}
    h.responderExito(w, http.StatusOK, resultado, "Total de propiedades")
}
```

2. Agregar ruta en `internal/web/routes.go`

---

## 📁 **Ejercicios: cmd/servidor/main.go**

### **🟢 Ejercicio 11: Configuración**
**Nivel:** Principiante  
**Tiempo:** 15 minutos

**Tarea:** Agregar nueva variable de entorno `API_VERSION`.

**Pasos:**
1. En `.env` agregar: `API_VERSION=1.0.0`
2. En `main.go` leer la variable
3. Usar en algún endpoint

### **🟡 Ejercicio 12: Logging Mejorado**
**Nivel:** Intermedio  
**Tiempo:** 30 minutos

**Tarea:** Agregar timestamp a todos los logs.

**Antes:**
```
2025/07/04 19:12:22 Servidor iniciado...
```

**Después:**
```
[2025-07-04 19:12:22] INFO: Servidor iniciado en puerto 8080
```

---

## 🎓 **Proyecto Final: Nueva Funcionalidad Completa**

### **🔴 Ejercicio 13: Favoritos (End-to-End)**
**Nivel:** Avanzado  
**Tiempo:** 2-3 horas

**Tarea:** Implementar sistema de propiedades favoritas.

**Requiere cambios en TODOS los archivos:**

1. **Dominio:** Agregar campo `Favorita bool`
2. **Repositorio:** Métodos para marcar/desmarcar favorita
3. **Servicio:** Lógica para gestionar favoritos
4. **Handlers:** Endpoints PUT para marcar favorita
5. **Main:** Configurar nuevas rutas

**Endpoints objetivo:**
- `PUT /api/propiedades/{id}/favorita` - Marcar como favorita
- `DELETE /api/propiedades/{id}/favorita` - Quitar de favoritas
- `GET /api/propiedades/favoritas` - Listar solo favoritas

---

## 📊 **Plan de Estudio Recomendado**

### **Semana 1:**
- Ejercicios 1-3 (Dominio)
- Leer CONCEPTOS_GO_FUNDAMENTALES.md

### **Semana 2:**
- Ejercicios 4-6 (Repositorio)
- Entender interfaces y SQL

### **Semana 3:**
- Ejercicios 7-8 (Servicio)
- Lógica de negocio

### **Semana 4:**
- Ejercicios 9-10 (Handlers)
- HTTP y JSON

### **Semana 5:**
- Ejercicios 11-12 (Main)
- Ejercicio 13 (Proyecto final)

---

## 🎯 **Consejos para Hacer Ejercicios**

1. **Un ejercicio a la vez** - No saltes al siguiente hasta dominar el actual
2. **Compila frecuentemente** - `go build ./...` después de cada cambio
3. **Usa fmt.Println** para debug - Ve qué valores tienen las variables
4. **Lee los errores** - Go tiene mensajes de error muy descriptivos
5. **Pregunta cuando te atasques** - Es normal tener dudas

## ✅ **Cómo Saber si Dominaste un Concepto**

- Puedes explicarlo con tus propias palabras
- Puedes modificar el código sin romperlo
- Entiendes por qué se hizo de esa manera
- Puedes encontrar y arreglar errores simples

---

💡 **Recuerda:** El objetivo es ENTENDER, no memorizar. Tomate el tiempo que necesites.