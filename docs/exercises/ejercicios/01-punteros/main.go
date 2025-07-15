package main

import (
	"fmt"
)

// ========== USO 1: MODIFICAR VARIABLES EN FUNCIONES ==========

type Usuario struct {
	nombre string
	edad   int
	activo bool
}

// Sin puntero: NO modifica el original
func activarUsuarioCopia(u Usuario) {
	u.activo = true
	fmt.Println("Dentro función (copia):", u.activo)
}

// Con puntero: SÍ modifica el original
func activarUsuarioPuntero(u *Usuario) {
	u.activo = true // Go automáticamente hace (*u).activo
	fmt.Println("Dentro función (puntero):", u.activo)
}

func ejemplo1() {
	fmt.Println("=== USO 1: MODIFICAR EN FUNCIONES ===")

	usuario := Usuario{nombre: "Juan", edad: 30, activo: false}

	fmt.Println("Original:", usuario.activo)
	activarUsuarioCopia(usuario)
	fmt.Println("Después de copia:", usuario.activo) // Sigue false

	activarUsuarioPuntero(&usuario)
	fmt.Println("Después de puntero:", usuario.activo) // Ahora true

	fmt.Println("---")
}

// ========== USO 2: MÉTODOS QUE MODIFICAN EL RECEIVER ==========

type Contador struct {
	valor int
}

// Método con receiver por valor: NO modifica el original
func (c Contador) IncrementarCopia() {
	c.valor++
}

// Método con receiver por puntero: SÍ modifica el original
func (c *Contador) IncrementarOriginal() {
	c.valor++
}

func ejemplo2() {
	fmt.Println("=== USO 2: MÉTODOS CON RECEIVER ===")

	contador := Contador{valor: 0}

	fmt.Println("Inicial:", contador.valor)

	contador.IncrementarCopia()
	fmt.Println("Después de copia:", contador.valor) // Sigue 0

	contador.IncrementarOriginal()
	fmt.Println("Después de puntero:", contador.valor) // Ahora 1

	fmt.Println("---")
}

// ========== USO 3: VALORES OPCIONALES (COMO NULL) ==========

type Producto struct {
	nombre string
	precio *float64 // Puntero = precio opcional
}

func (p Producto) MostrarPrecio() string {
	if p.precio == nil {
		return fmt.Sprintf("%s: Sin precio", p.nombre)
	}
	return fmt.Sprintf("%s: $%.2f", p.nombre, *p.precio)
}

func ejemplo3() {
	fmt.Println("=== USO 3: VALORES OPCIONALES ===")

	// Producto sin precio
	producto1 := Producto{nombre: "Laptop", precio: nil}
	fmt.Println(producto1.MostrarPrecio())

	// Producto con precio
	precio := 999.99
	producto2 := Producto{nombre: "Mouse", precio: &precio}
	fmt.Println(producto2.MostrarPrecio())

	fmt.Println("---")
}

// ========== USO 4: EVITAR COPIAS COSTOSAS ==========

type DocumentoGrande struct {
	contenido [1000]string // Imagine que esto es muy grande
	titulo    string
}

// Sin puntero: copia todo el documento (lento)
func procesarDocumentoCopia(doc DocumentoGrande) {
	fmt.Printf("Procesando (copia): %s\n", doc.titulo)
}

// Con puntero: solo pasa la dirección (rápido)
func procesarDocumentoPuntero(doc *DocumentoGrande) {
	fmt.Printf("Procesando (puntero): %s\n", doc.titulo)
}

func ejemplo4() {
	fmt.Println("=== USO 4: EVITAR COPIAS COSTOSAS ===")

	doc := DocumentoGrande{titulo: "Manual de 1000 páginas"}

	// Esto copia todo el array de 1000 strings
	procesarDocumentoCopia(doc)

	// Esto solo pasa 8 bytes (la dirección)
	procesarDocumentoPuntero(&doc)

	fmt.Println("---")
}

// ========== USO 5: MÚLTIPLES VALORES DE RETORNO ==========

func dividir(a, b int) (resultado *int, error string) {
	if b == 0 {
		return nil, "División por cero"
	}

	res := a / b
	return &res, ""
}

func ejemplo5() {
	fmt.Println("=== USO 5: MÚLTIPLES VALORES DE RETORNO ===")

	// División válida
	if resultado, err := dividir(10, 2); err == "" {
		fmt.Printf("10 / 2 = %d\n", *resultado)
	} else {
		fmt.Println("Error:", err)
	}

	// División por cero
	if resultado, err := dividir(10, 0); err == "" {
		fmt.Printf("10 / 0 = %d\n", *resultado)
	} else {
		fmt.Println("Error:", err)
	}

	fmt.Println("---")
}

// ========== USO 6: MODIFICAR SLICES Y MAPS ==========

func modificarSlice(slice []int) {
	// Los slices YA son referencias, pero veamos el concepto
	slice[0] = 999
}

func modificarSlicePuntero(slice *[]int) {
	// Para modificar el slice mismo (no solo su contenido)
	*slice = append(*slice, 100, 200)
}

func ejemplo6() {
	fmt.Println("=== USO 6: MODIFICAR SLICES ===")

	numeros := []int{1, 2, 3}
	fmt.Println("Original:", numeros)

	modificarSlice(numeros)
	fmt.Println("Después de modificar contenido:", numeros)

	modificarSlicePuntero(&numeros)
	fmt.Println("Después de modificar slice:", numeros)

	fmt.Println("---")
}

// ========== USO 7: CONSTRUCTORES DE OBJETOS ==========

func NuevoUsuario(nombre string) *Usuario {
	return &Usuario{
		nombre: nombre,
		edad:   0,
		activo: true,
	}
}

func ejemplo7() {
	fmt.Println("=== USO 7: CONSTRUCTORES ===")

	// Constructor devuelve puntero para eficiencia
	usuario := NuevoUsuario("Ana")
	fmt.Printf("Usuario creado: %+v\n", *usuario)

	// Podemos modificar directamente
	usuario.edad = 25
	fmt.Printf("Usuario modificado: %+v\n", *usuario)

	fmt.Println("---")
}

// ========== USO 8: LINKED LISTS Y ESTRUCTURAS RECURSIVAS ==========

type Nodo struct {
	valor     int
	siguiente *Nodo // Puntero al siguiente nodo
}

func crearLista() *Nodo {
	// Crear lista: 1 -> 2 -> 3 -> nil
	nodo3 := &Nodo{valor: 3, siguiente: nil}
	nodo2 := &Nodo{valor: 2, siguiente: nodo3}
	nodo1 := &Nodo{valor: 1, siguiente: nodo2}

	return nodo1
}

func imprimirLista(nodo *Nodo) {
	current := nodo
	for current != nil {
		fmt.Printf("%d -> ", current.valor)
		current = current.siguiente
	}
	fmt.Println("nil")
}

func ejemplo8() {
	fmt.Println("=== USO 8: LINKED LISTS ===")

	lista := crearLista()
	imprimirLista(lista)

	fmt.Println("---")
}

// ========== USO 9: COMPARACIÓN DE RENDIMIENTO ==========

func ejemploRendimiento() {
	fmt.Println("=== USO 9: RENDIMIENTO ===")

	// Imagina un struct muy grande
	type DatosGigantes struct {
		datos [10000]int
		info  string
	}

	datos := DatosGigantes{info: "Datos muy grandes"}

	// Método 1: Pasar por valor (copia todo)
	func(d DatosGigantes) {
		fmt.Printf("Por valor: %s (costoso)\n", d.info)
	}(datos)

	// Método 2: Pasar por puntero (solo dirección)
	func(d *DatosGigantes) {
		fmt.Printf("Por puntero: %s (eficiente)\n", d.info)
	}(&datos)

	fmt.Println("---")
}

// ========== RESUMEN DE CUÁNDO USAR ==========

func resumenUsos() {
	fmt.Println("=== RESUMEN: CUÁNDO USAR PUNTEROS ===")

	fmt.Println("✅ USA PUNTEROS CUANDO:")
	fmt.Println("  1. Necesites modificar el valor original")
	fmt.Println("  2. Quieras evitar copias costosas")
	fmt.Println("  3. Necesites valores opcionales (nil)")
	fmt.Println("  4. Construyas estructuras recursivas")
	fmt.Println("  5. Implementes métodos que modifican el receiver")

	fmt.Println("\n❌ NO USES PUNTEROS CUANDO:")
	fmt.Println("  1. Solo necesites leer el valor")
	fmt.Println("  2. Trabajes con tipos pequeños (int, bool)")
	fmt.Println("  3. No necesites modificar el original")
	fmt.Println("  4. Quieras simplicidad sobre eficiencia")
}

func main() {
	ejemplo1()
	ejemplo2()
	ejemplo3()
	ejemplo4()
	ejemplo5()
	ejemplo6()
	ejemplo7()
	ejemplo8()
	ejemploRendimiento()
	resumenUsos()
}
