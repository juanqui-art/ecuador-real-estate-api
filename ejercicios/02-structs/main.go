package main

import (
	"fmt"
	"time"
)

// 🏗️ EJEMPLO 1: Struct básica
type Persona struct {
	Nombre string
	Edad   int
	Email  string
}

// 🏗️ EJEMPLO 2: Struct con tags (como en el proyecto)
type Propiedad struct {
	ID          string    `json:"id"`
	Titulo      string    `json:"titulo"`
	Precio      float64   `json:"precio"`
	Provincia   string    `json:"provincia"`
	FechaCreada time.Time `json:"fecha_creada"`
}

func main() {
	// 🎯 EJEMPLO 1: Crear struct básica
	fmt.Println("=== EJEMPLO 1: Struct Básica ===")
	
	// Forma 1: Crear struct directamente
	persona1 := Persona{
		Nombre: "Juan",
		Edad:   30,
		Email:  "juan@email.com",
	}
	
	fmt.Println("Persona 1:", persona1)
	fmt.Println("Nombre:", persona1.Nombre)
	fmt.Println("Edad:", persona1.Edad)
	fmt.Println()

	// 🎯 EJEMPLO 2: Crear struct con constructor
	fmt.Println("=== EJEMPLO 2: Struct con Constructor ===")
	
	propiedad := NuevaPropiedad("Casa en Guayaquil", 150000.0, "Guayas")
	fmt.Println("Propiedad creada:", propiedad)
	fmt.Println("Título:", propiedad.Titulo)
	fmt.Println("Precio:", propiedad.Precio)
	fmt.Println("ID:", propiedad.ID)
	fmt.Println()

	// 🎯 EJEMPLO 3: Modificar struct
	fmt.Println("=== EJEMPLO 3: Modificar Struct ===")
	
	fmt.Println("Precio original:", propiedad.Precio)
	propiedad.Precio = 160000.0
	fmt.Println("Precio modificado:", propiedad.Precio)
	fmt.Println()

	// 🎯 EJEMPLO 4: Struct con punteros
	fmt.Println("=== EJEMPLO 4: Struct con Punteros ===")
	
	// Crear puntero a struct
	ptrPropiedad := &Propiedad{
		ID:          "prop-001",
		Titulo:      "Departamento en Quito",
		Precio:      80000.0,
		Provincia:   "Pichincha",
		FechaCreada: time.Now(),
	}
	
	fmt.Println("Título desde puntero:", ptrPropiedad.Titulo)
	
	// Modificar a través del puntero
	ptrPropiedad.Precio = 85000.0
	fmt.Println("Precio modificado desde puntero:", ptrPropiedad.Precio)
	fmt.Println()

	// 🎯 EJEMPLO 5: Diferencia entre valor y puntero
	fmt.Println("=== EJEMPLO 5: Valor vs Puntero ===")
	
	// Crear struct por valor
	casa := Propiedad{Titulo: "Casa", Precio: 100000.0}
	
	// Pasar por valor (se copia toda la struct)
	casaCopia := casa
	casaCopia.Precio = 120000.0
	
	fmt.Println("Casa original:", casa.Precio)    // 100000
	fmt.Println("Casa copia:", casaCopia.Precio)  // 120000
	
	// Pasar por puntero (se comparte)
	casaPuntero := &casa
	casaPuntero.Precio = 110000.0
	
	fmt.Println("Casa original después del puntero:", casa.Precio) // 110000
	fmt.Println()

	// 🎯 EJEMPLO 6: Comparar structs
	fmt.Println("=== EJEMPLO 6: Comparar Structs ===")
	
	prop1 := Propiedad{ID: "1", Titulo: "Casa", Precio: 100000.0}
	prop2 := Propiedad{ID: "1", Titulo: "Casa", Precio: 100000.0}
	prop3 := Propiedad{ID: "2", Titulo: "Casa", Precio: 100000.0}
	
	fmt.Println("prop1 == prop2:", prop1 == prop2) // true - mismos valores
	fmt.Println("prop1 == prop3:", prop1 == prop3) // false - ID diferente
}

// Constructor para Propiedad (como en el proyecto real)
func NuevaPropiedad(titulo string, precio float64, provincia string) *Propiedad {
	return &Propiedad{
		ID:          generarID(),
		Titulo:      titulo,
		Precio:      precio,
		Provincia:   provincia,
		FechaCreada: time.Now(),
	}
}

// Función auxiliar para generar ID simple
func generarID() string {
	return fmt.Sprintf("prop-%d", time.Now().UnixNano())
}