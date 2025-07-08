# Ejemplos de Uso - Geolocalización y Relaciones

## Geolocalización de Propiedades

### Coordenadas de Ciudades Principales de Ecuador

```go
// Quito
latitudQuito := -0.1807
longitudQuito := -78.4678

// Guayaquil  
latitudGuayaquil := -2.1469
longitudGuayaquil := -79.6485

// Cuenca
latitudCuenca := -2.9001
longitudCuenca := -79.0059
```

### Crear Propiedad con Ubicación

```go
package main

import (
    "fmt"
    "realty-core/internal/dominio"
)

func main() {
    // 1. Crear propiedad básica
    propiedad := dominio.NuevaPropiedad(
        "Casa moderna en Samborondón",
        "Hermosa casa con piscina y jardín",
        "Guayas",
        "Samborondón", 
        "casa",
        285000,
    )
    
    // 2. Configurar ubicación GPS (cerca de Guayaquil)
    err := propiedad.ConfigurarUbicacion(
        -2.1469,  // Latitud Guayaquil
        -79.6485, // Longitud Guayaquil
        dominio.PrecisionExacta,
    )
    
    if err != nil {
        fmt.Printf("Error configurando ubicación: %v\n", err)
        return
    }
    
    // 3. Verificar ubicación
    if propiedad.TieneUbicacionConfigurada() {
        fmt.Printf("Propiedad ubicada en: %.6f, %.6f\n", 
            propiedad.Latitud, propiedad.Longitud)
    }
    
    fmt.Printf("Propiedad creada: %s\n", propiedad.Titulo)
}
```

## Relaciones con Inmobiliarias

### Crear Inmobiliaria

```go
// 1. Crear inmobiliaria
inmobiliaria := dominio.NuevaInmobiliaria(
    "InmoEcuador S.A.",
    "+593987654321",
    "contacto@inmoecuador.com",
    "Pichincha",
    "Quito",
)

// 2. Configurar información adicional
err := inmobiliaria.ConfigurarLicencia("1234567890001") // RUC
if err != nil {
    fmt.Printf("Error: %v\n", err)
}

err = inmobiliaria.ConfigurarSitioWeb("https://www.inmoecuador.com")
if err != nil {
    fmt.Printf("Error: %v\n", err)
}

fmt.Printf("Inmobiliaria creada: %s\n", inmobiliaria.Nombre)
```

### Asociar Propiedad con Inmobiliaria

```go
// Método 1: Crear propiedad ya asociada
propiedadConInmo := dominio.NuevaPropiedadConInmobiliaria(
    "Departamento en La Carolina",
    "Moderno departamento de 2 dormitorios",
    "Pichincha",
    "Quito",
    "departamento", 
    95000,
    inmobiliaria.ID,
)

// Método 2: Asociar después de crear
propiedad := dominio.NuevaPropiedad(...)
propiedad.AsociarInmobiliaria(inmobiliaria.ID)

// Verificar asociación
if propiedad.TieneInmobiliaria() {
    fmt.Printf("Propiedad manejada por inmobiliaria ID: %s\n", 
        propiedad.ObtenerInmobiliariaID())
}
```

## Ejemplo Completo para Frontend

### JSON que recibe el Frontend

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "slug": "casa-moderna-samborondon-123e4567",
  "titulo": "Casa moderna en Samborondón", 
  "descripcion": "Hermosa casa con piscina y jardín",
  "precio": 285000,
  "provincia": "Guayas",
  "ciudad": "Samborondón",
  "sector": "Vía a la Costa Km 15",
  "direccion": "Urbanización Los Ceibos, Mz 5, Villa 12",
  "latitud": -2.1469,
  "longitud": -79.6485, 
  "precision_ubicacion": "exacta",
  "tipo": "casa",
  "estado": "disponible",
  "dormitorios": 4,
  "banos": 3.5,
  "area_m2": 320.5,
  "inmobiliaria_id": "456e7890-e89b-12d3-a456-426614174001",
  "fecha_creacion": "2024-01-15T10:30:00Z",
  "fecha_actualizacion": "2024-01-15T10:30:00Z"
}
```

### Código Frontend (JavaScript) para Mapa

```javascript
// Mostrar propiedad en mapa (usando Google Maps)
function mostrarPropiedadEnMapa(propiedad) {
    if (propiedad.latitud === 0 && propiedad.longitud === 0) {
        console.log("Propiedad sin ubicación GPS");
        return;
    }
    
    const mapa = new google.maps.Map(document.getElementById('mapa'), {
        center: { lat: propiedad.latitud, lng: propiedad.longitud },
        zoom: 15
    });
    
    const marcador = new google.maps.Marker({
        position: { lat: propiedad.latitud, lng: propiedad.longitud },
        map: mapa,
        title: propiedad.titulo
    });
    
    // Ventana de información
    const infoWindow = new google.maps.InfoWindow({
        content: `
            <h3>${propiedad.titulo}</h3>
            <p>Precio: $${propiedad.precio.toLocaleString()}</p>
            <p>Ubicación: ${propiedad.precision_ubicacion}</p>
        `
    });
    
    marcador.addListener('click', () => {
        infoWindow.open(mapa, marcador);
    });
}
```

## Validaciones Específicas para Ecuador

### Coordenadas Válidas

```go
// Estas coordenadas están dentro de Ecuador
coordenadasValidas := []struct{
    ciudad string
    lat, lng float64
}{
    {"Quito", -0.1807, -78.4678},
    {"Guayaquil", -2.1469, -79.6485}, 
    {"Cuenca", -2.9001, -79.0059},
    {"Manta", -0.9677, -80.7089},
    {"Loja", -3.9969, -79.2067},
}

for _, coord := range coordenadasValidas {
    valida := dominio.EsCoordenadasValidasEcuador(coord.lat, coord.lng)
    fmt.Printf("%s: %v\n", coord.ciudad, valida) // Todas serán true
}

// Coordenadas fuera de Ecuador (serán false)
invalidasBogota := dominio.EsCoordenadasValidasEcuador(4.7110, -74.0721) // false
invalidasLima := dominio.EsCoordenadasValidasEcuador(-12.0464, -77.0428)  // false
```

### Teléfonos Ecuatorianos

```go
telefonos := []string{
    "0987654321",      // Móvil válido
    "022345678",       // Fijo Quito válido  
    "+593987654321",   // Con código país
    "0423456789",      // Fijo Guayaquil válido
}

for _, tel := range telefonos {
    valido := dominio.EsTelefonoValidoEcuador(tel)
    formateado := dominio.FormatearTelefono(tel)
    fmt.Printf("Tel: %s -> Válido: %v, Formateado: %s\n", tel, valido, formateado)
}
```