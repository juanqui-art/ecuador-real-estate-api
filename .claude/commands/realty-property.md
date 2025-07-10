# Realty Property Management

Manage real estate properties for: $ARGUMENTS

## Context - Property Domain
Our Property struct includes:
- Basic info: ID, Title, Description, Price, Province, City, Type, Status
- Location: Latitude, Longitude, Sector, Address, LocationPrecision  
- Details: Bedrooms, Bathrooms, AreaM2, YearBuilt, Floors
- Features: Furnished, Garage, Pool, Garden, Terrace, Balcony, Security, Elevator
- Media: MainImage, Images[], VideoTour, Tour360
- SEO: Slug, Tags[], Featured, ViewCount
- Timestamps: CreatedAt, UpdatedAt

## Property Operations:
1. **Struct modifications:**
   - Add new fields with proper Go types and tags
   - Include validation tags (validate, db, json)
   - Generate SQL migrations for new fields
   - Update repository methods

2. **Business logic:**
   - Property validation (Ecuador provinces, cities)
   - Price calculations (per m2, rent vs sale)
   - Status transitions (available → sold → rented)
   - Feature combinations validation

3. **Data operations:**
   - CRUD operations with PostgreSQL
   - Full-text search integration
   - Filtering and pagination
   - Related data (images, companies)

## Ecuador-specific validations:
- Provinces: Azuay, Bolívar, Cañar, Carchi, Chimborazo, Cotopaxi, El Oro, Esmeraldas, Galápagos, Guayas, Imbabura, Loja, Los Ríos, Manabí, Morona Santiago, Napo, Orellana, Pastaza, Pichincha, Santa Elena, Santo Domingo, Sucumbíos, Tungurahua, Zamora Chinchipe
- Property types: casa, apartamento, terreno, comercial
- Price ranges by province/city
- Address format validation

## Code patterns to follow:
- Use pointer receivers for methods that modify state
- Implement builder pattern for complex property creation
- Add proper error handling with domain-specific errors
- Use context.Context for all operations
- Include proper logging and metrics

## Database patterns:
- Use transactions for related operations
- Implement soft deletes with deleted_at
- Add proper indexes for search fields
- Use JSON fields for flexible data (images, tags)

## Testing requirements:
- Table-driven tests for validation
- Mock repository for service tests
- Integration tests for complex operations
- Benchmark tests for search/filter operations

## Common use cases:
- "add elevator field to Property struct"
- "create property search by price range and city"
- "validate property price against market average"
- "implement property status workflow"
- "add property comparison functionality"
- "optimize property listing query performance"

## Integration points:
- Image system: property.MainImage, property.Images
- Cache system: property search results, property details
- FTS system: property.Title, property.Description search
- User system: property ownership, favorites (future)

## Performance considerations:
- Eager load related data when needed
- Use pagination for large result sets
- Cache frequent searches
- Optimize database queries with proper indexes
- Use connection pooling for database access

## Output format:
- Complete Go code with proper error handling
- SQL migrations if schema changes needed
- Unit tests for new functionality
- Documentation for new features
- Performance benchmarks if applicable