# 🧪 Ejercicios de Aprendizaje Go

Esta carpeta contiene ejercicios prácticos para aprender Go paso a paso, usando conceptos del proyecto inmobiliario.

## 📚 Orden de Estudio

### 1. **01-punteros/** - Conceptos de Punteros
- **Conceptos:** `&`, `*`, diferencia valor vs puntero
- **Tiempo:** 30 minutos
- **Ejecutar:** `go run 01-punteros/main.go`

### 2. **02-structs/** - Estructuras de Datos
- **Conceptos:** `struct`, constructores, campos, comparaciones
- **Tiempo:** 45 minutos
- **Ejecutar:** `go run 02-structs/main.go`

### 3. **03-metodos/** - Métodos y Receivers
- **Conceptos:** `func (receiver)`, métodos vs funciones, encadenamiento
- **Tiempo:** 60 minutos
- **Ejecutar:** `go run 03-metodos/main.go`

## 🚀 Cómo Ejecutar

### **Desde la terminal:**
```bash
cd ejercicios
go run 01-punteros/main.go
go run 02-structs/main.go
go run 03-metodos/main.go
```

### **Desde GoLand:**
1. Abrir cualquier archivo `main.go`
2. Click derecho → "Run 'go build main.go'"
3. O usar el botón ▶️ verde

## 📋 Checklist de Aprendizaje

### 01-punteros
- [ ] Entiendo qué es un puntero (`&variable`)
- [ ] Entiendo cómo desreferenciar (`*puntero`)
- [ ] Entiendo la diferencia entre pasar por valor vs puntero
- [ ] Puedo explicar cuándo usar cada uno

### 02-structs
- [ ] Puedo crear structs básicas
- [ ] Entiendo los tags (`json:"campo"`)
- [ ] Puedo crear constructores
- [ ] Entiendo la diferencia entre `Struct{}` y `&Struct{}`

### 03-metodos
- [ ] Entiendo qué es un receiver
- [ ] Entiendo la diferencia entre `(s Struct)` y `(s *Struct)`
- [ ] Puedo crear métodos que modifican vs que solo leen
- [ ] Entiendo la diferencia entre métodos y funciones

## 🎯 Ejercicios Adicionales

### **Ejercicio A: Crear tu propia struct**
1. Crear struct `Persona` con campos: Nombre, Edad, Email
2. Crear constructor `NuevaPersona()`
3. Crear métodos `EsMayorDeEdad()` y `CambiarEmail()`

### **Ejercicio B: Modificar Propiedad**
1. Agregar campo `Garage bool` a la struct Propiedad
2. Crear método `TieneGarage()` 
3. Crear método `AgregarGarage()` que cambie el valor

### **Ejercicio C: Lista de Propiedades**
1. Crear slice `[]Propiedad`
2. Crear función que filtre propiedades por precio
3. Crear función que calcule precio promedio

## 💡 Consejos

1. **Ejecuta cada ejemplo** - No solo leas, ejecuta y ve el resultado
2. **Modifica el código** - Cambia valores y ve qué pasa
3. **Pregúntate "¿por qué?"** - Entiende el propósito de cada línea
4. **Compara con Python** - Piensa cómo harías lo mismo en Python/Django
5. **Usa fmt.Println** - Agrega más prints para debug

## 🆘 Si te Atascas

1. Lee los comentarios en el código
2. Ejecuta línea por línea mentalmente
3. Usa `fmt.Printf("%T %v", variable, variable)` para ver tipo y valor
4. Revisa la documentación: `CONCEPTOS_GO_FUNDAMENTALES.md`

---

💪 **¡Tómate tu tiempo! Es mejor entender profundamente que avanzar rápido.**