# Realty Ecuador Localization

Implement Ecuador-specific validations for: $ARGUMENTS

## Context - Ecuador Market
Real estate system optimized for Ecuador with:
- **24 provinces** with specific cities
- **Property types:** casa, apartamento, terreno, comercial
- **Price ranges** by region
- **Address formats** for Ecuador
- **Legal requirements** for property transactions

## Ecuador Data:
**Provinces:** Azuay, Bolívar, Cañar, Carchi, Chimborazo, Cotopaxi, El Oro, Esmeraldas, Galápagos, Guayas, Imbabura, Loja, Los Ríos, Manabí, Morona Santiago, Napo, Orellana, Pastaza, Pichincha, Santa Elena, Santo Domingo, Sucumbíos, Tungurahua, Zamora Chinchipe

**Major Cities:**
- Pichincha: Quito, Cayambe, Mejía
- Guayas: Guayaquil, Durán, Samborondón
- Azuay: Cuenca, Gualaceo, Paute

## Validation Patterns:
```go
// Province validation
func ValidateProvince(province string) error {
    validProvinces := []string{
        "Azuay", "Bolívar", "Cañar", "Carchi", "Chimborazo", "Cotopaxi",
        "El Oro", "Esmeraldas", "Galápagos", "Guayas", "Imbabura", "Loja",
        "Los Ríos", "Manabí", "Morona Santiago", "Napo", "Orellana", "Pastaza",
        "Pichincha", "Santa Elena", "Santo Domingo", "Sucumbíos", "Tungurahua", "Zamora Chinchipe",
    }
    
    for _, valid := range validProvinces {
        if province == valid {
            return nil
        }
    }
    return fmt.Errorf("invalid province: %s", province)
}

// Price validation by region
func ValidatePriceRange(price float64, province string, propertyType string) error {
    // Region-specific price validation
}
```

## Common use cases:
- "validate Quito postal codes"
- "add province-city relationship validation"
- "implement currency formatting for Ecuador"
- "create address validation for Ecuador"
- "add property tax calculations"

## Market-specific features:
- Property price trends by region
- Legal document requirements
- Tax calculations
- Market analysis tools
- Local regulations compliance

## Output format:
- Ecuador-specific validation functions
- Province and city catalogs
- Price range validations
- Address formatting
- Legal compliance checks