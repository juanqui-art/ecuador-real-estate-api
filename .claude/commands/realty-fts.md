# Realty Full-Text Search

Optimize search functionality for: $ARGUMENTS

## Context - Current FTS Implementation
PostgreSQL Full-Text Search with:
- **Spanish language support** with custom dictionaries
- **Weighted search vectors** for title, description, location
- **Ranking system** using ts_rank_cd
- **Suggestions** with autocomplete
- **Advanced filters** combined with FTS

## FTS Configuration:
```sql
-- Search vector with weights
UPDATE properties SET search_vector = 
    setweight(to_tsvector('spanish', title), 'A') ||
    setweight(to_tsvector('spanish', description), 'B') ||
    setweight(to_tsvector('spanish', city || ' ' || province), 'C');

-- GIN index for performance
CREATE INDEX idx_properties_search_vector ON properties USING GIN(search_vector);
```

## Search Operations:
1. **Basic search:**
   - Full-text query parsing
   - Stemming and normalization
   - Relevance ranking
   - Phrase search support

2. **Advanced search:**
   - Combine FTS with filters
   - Geographic search
   - Price range filtering
   - Property type filtering

3. **Autocomplete:**
   - Suggest property titles
   - Suggest locations
   - Suggest property types
   - Real-time suggestions

## Search patterns:
```go
// Basic FTS query
func (r *Repository) SearchProperties(query string, limit int) ([]PropertySearchResult, error) {
    sql := `
        SELECT id, title, price, province, city, ts_rank_cd(search_vector, plainto_tsquery('spanish', $1)) as rank
        FROM properties 
        WHERE search_vector @@ plainto_tsquery('spanish', $1)
        ORDER BY rank DESC
        LIMIT $2
    `
    // Execute query
}

// Advanced search with filters
func (r *Repository) AdvancedSearch(params AdvancedSearchParams) ([]PropertySearchResult, error) {
    // Build dynamic query with FTS and filters
}
```

## Common use cases:
- "improve search ranking for Ecuador locations"
- "add autocomplete for property types"
- "optimize search query performance"
- "implement search result highlighting"
- "add search filters for amenities"
- "create search analytics"

## Performance optimizations:
- Proper indexing strategies
- Query optimization
- Result caching
- Search result pagination
- Precomputed popular searches

## Output format:
- FTS query optimization
- Index recommendations
- Search result ranking
- Autocomplete implementation
- Performance improvements