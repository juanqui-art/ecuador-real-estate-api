package main

import (
	"fmt"
	"time"

	"realty-core/internal/cache"
)

func main() {
	fmt.Println("🚀 Probando el sistema de cache LRU...")

	// Crear configuración del cache
	config := cache.ImageCacheConfig{
		Enabled:         true,
		Capacity:        5,           // Solo 5 elementos
		MaxSizeBytes:    1024,        // 1KB máximo
		TTL:             30 * time.Second,
		CleanupInterval: 5 * time.Second,
	}

	// Crear cache de imágenes
	imageCache := cache.NewImageCache(config)

	fmt.Printf("✅ Cache creado: Enabled=%v, Capacity=%d\n", 
		imageCache.IsEnabled(), config.Capacity)

	// Prueba 1: Almacenar y recuperar datos
	fmt.Println("\n📁 Prueba 1: Almacenar y recuperar datos")
	
	testData := []byte("datos de imagen de prueba")
	imageCache.Set("imagen1", testData, "image/jpeg")
	
	data, contentType, found := imageCache.Get("imagen1")
	if found {
		fmt.Printf("✅ Recuperado: %s (%s)\n", string(data), contentType)
	} else {
		fmt.Println("❌ No se encontró la imagen")
	}

	// Prueba 2: Thumbnails
	fmt.Println("\n🖼️  Prueba 2: Cache de thumbnails")
	
	thumbnailData := []byte("thumbnail data")
	imageCache.SetThumbnail("prop1", 150, thumbnailData, "image/jpeg")
	
	thumb, thumbType, found := imageCache.GetThumbnail("prop1", 150)
	if found {
		fmt.Printf("✅ Thumbnail recuperado: %s (%s)\n", string(thumb), thumbType)
	}

	// Prueba 3: Variantes de imagen
	fmt.Println("\n🎨 Prueba 3: Cache de variantes")
	
	variantData := []byte("variant 800x600")
	imageCache.SetVariant("prop1", 800, 600, 85, "jpg", variantData, "image/jpeg")
	
	variant, varType, found := imageCache.GetVariant("prop1", 800, 600, 85, "jpg")
	if found {
		fmt.Printf("✅ Variante recuperada: %s (%s)\n", string(variant), varType)
	}

	// Prueba 4: Estadísticas
	fmt.Println("\n📊 Prueba 4: Estadísticas del cache")
	
	stats := imageCache.Stats()
	fmt.Printf("✅ Estadísticas:\n")
	fmt.Printf("   - Elementos: %d\n", stats.Size)
	fmt.Printf("   - Tamaño: %d bytes\n", stats.CurrentSize)
	fmt.Printf("   - Hits thumbnails: %d\n", stats.ThumbnailHits)
	fmt.Printf("   - Hits variantes: %d\n", stats.VariantHits)
	fmt.Printf("   - Tasa hit thumbnails: %.2f%%\n", stats.ThumbnailRate)
	fmt.Printf("   - Tasa hit variantes: %.2f%%\n", stats.VariantRate)

	// Prueba 5: Llenar cache hasta límite de capacidad
	fmt.Println("\n⚡ Prueba 5: Evicción por capacidad")
	
	for i := 0; i < 7; i++ {
		key := fmt.Sprintf("imagen%d", i+2)
		data := []byte(fmt.Sprintf("datos imagen %d", i+2))
		imageCache.Set(key, data, "image/jpeg")
		fmt.Printf("   Agregado: %s (tamaño cache: %d)\n", key, imageCache.Size())
	}

	// Verificar que los primeros elementos fueron evictados
	_, _, found = imageCache.Get("imagen1")
	fmt.Printf("✅ ¿Imagen1 sigue en cache? %v (debería ser false por evicción)\n", found)

	// Prueba 6: Invalidación de imagen
	fmt.Println("\n🗑️  Prueba 6: Invalidación de imagen")
	
	removedCount := imageCache.InvalidateImage("prop1")
	fmt.Printf("✅ Elementos removidos para prop1: %d\n", removedCount)
	
	// Estadísticas finales
	fmt.Println("\n📈 Estadísticas finales:")
	finalStats := imageCache.Stats()
	fmt.Printf("   - Elementos finales: %d\n", finalStats.Size)
	fmt.Printf("   - Tamaño final: %d bytes\n", finalStats.CurrentSize)

	fmt.Println("\n🎉 ¡Prueba del cache completada exitosamente!")
}