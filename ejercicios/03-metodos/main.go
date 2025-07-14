package main

import (
	"fmt"
	"time"
)

// 🏗️ Struct de ejemplo
type Propiedad struct {
	ID          string
	Titulo      string
	Precio      float64
	Provincia   string
	FechaCreada time.Time
}

func main() {
	// 🎯 EJEMPLO 1: Métodos básicos
	fmt.Println("=== EJEMPLO 1: Métodos Básicos ===")
	
	prop := &Propiedad{
		ID:          "prop-001",
		Titulo:      "Casa en Guayaquil",
		Precio:      150000.0,
		Provincia:   "Guayas",
		FechaCreada: time.Now(),
	}
	
	fmt.Println("Título:", prop.Titulo)
	fmt.Println("Es cara:", prop.EsCara())
	fmt.Println("Descripción:", prop.ObtenerDescripcion())
	fmt.Println()

	// 🎯 EJEMPLO 2: Métodos que modifican (receiver por puntero)
	fmt.Println("=== EJEMPLO 2: Métodos que Modifican ===")
	
	fmt.Println("Precio antes:", prop.Precio)
	prop.AplicarDescuento(10.0) // 10% descuento
	fmt.Println("Precio después del descuento:", prop.Precio)
	
	fmt.Println("Título antes:", prop.Titulo)
	prop.CambiarTitulo("Casa Renovada en Guayaquil")
	fmt.Println("Título después:", prop.Titulo)
	fmt.Println()

	// 🎯 EJEMPLO 3: Diferencia entre receiver por valor vs puntero
	fmt.Println("=== EJEMPLO 3: Receiver Valor vs Puntero ===")
	
	// Crear dos propiedades iguales
	prop1 := Propiedad{Titulo: "Casa", Precio: 100000.0}
	prop2 := Propiedad{Titulo: "Casa", Precio: 100000.0}
	
	fmt.Println("Precios iniciales:")
	fmt.Println("prop1:", prop1.Precio)
	fmt.Println("prop2:", prop2.Precio)
	
	// Método con receiver por VALOR (no modifica original)
	prop1.IntentarCambiarPrecio(200000.0)
	fmt.Println("prop1 después de IntentarCambiarPrecio:", prop1.Precio) // No cambia
	
	// Método con receiver por PUNTERO (modifica original)
	prop2.CambiarPrecio(200000.0)
	fmt.Println("prop2 después de CambiarPrecio:", prop2.Precio) // Sí cambia
	fmt.Println()

	// 🎯 EJEMPLO 4: Métodos encadenados
	fmt.Println("=== EJEMPLO 4: Métodos Encadenados ===")
	
	propNueva := &Propiedad{
		ID:          "prop-002",
		Titulo:      "Departamento",
		Precio:      80000.0,
		Provincia:   "Pichincha",
		FechaCreada: time.Now(),
	}
	
	// Encadenar métodos que retornan *Propiedad
	propNueva.CambiarTitulo("Departamento Moderno").
		AplicarDescuento(5.0).
		CambiarTitulo("Departamento Moderno en Oferta")
	
	fmt.Println("Resultado final:", propNueva.Titulo)
	fmt.Println("Precio final:", propNueva.Precio)
	fmt.Println()

	// 🎯 EJEMPLO 5: Métodos vs Funciones
	fmt.Println("=== EJEMPLO 5: Métodos vs Funciones ===")
	
	// Método - se llama en la instancia
	esCara := prop.EsCara()
	fmt.Println("Método - Es cara:", esCara)
	
	// Función - se pasa la instancia como parámetro
	esCaraFuncion := EsPropiedadCara(prop)
	fmt.Println("Función - Es cara:", esCaraFuncion)
	
	// Función - calcula impuesto sin modificar la propiedad
	impuesto := CalcularImpuesto(prop.Precio)
	fmt.Println("Impuesto (función):", impuesto)
}

// ========== MÉTODOS (con receiver) ==========

// Método con receiver por PUNTERO - puede modificar
func (p *Propiedad) AplicarDescuento(porcentaje float64) *Propiedad {
	p.Precio = p.Precio * (1 - porcentaje/100)
	return p // Retorna el puntero para encadenar métodos
}

// Método con receiver por PUNTERO - puede modificar
func (p *Propiedad) CambiarTitulo(nuevoTitulo string) *Propiedad {
	p.Titulo = nuevoTitulo
	return p // Retorna el puntero para encadenar métodos
}

// Método con receiver por PUNTERO - puede modificar
func (p *Propiedad) CambiarPrecio(nuevoPrecio float64) {
	p.Precio = nuevoPrecio
}

// Método con receiver por VALOR - NO puede modificar
func (p Propiedad) IntentarCambiarPrecio(nuevoPrecio float64) {
	p.Precio = nuevoPrecio // Esto solo cambia la copia, no el original
}

// Método con receiver por VALOR - solo lee, no modifica
func (p Propiedad) EsCara() bool {
	return p.Precio > 200000.0
}

// Método con receiver por VALOR - solo lee, no modifica
func (p Propiedad) ObtenerDescripcion() string {
	return fmt.Sprintf("%s en %s por $%.2f", p.Titulo, p.Provincia, p.Precio)
}

// ========== FUNCIONES (sin receiver) ==========

// Función que recibe propiedad como parámetro
func EsPropiedadCara(p *Propiedad) bool {
	return p.Precio > 200000.0
}

// Función que solo necesita el precio
func CalcularImpuesto(precio float64) float64 {
	return precio * 0.12 // 12% IVA
}