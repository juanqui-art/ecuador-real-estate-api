package repositorio

import (
	"database/sql"
	"testing"

	"realty-core/internal/dominio"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MockRepository implementa PropiedadRepository para testing
type MockRepository struct {
	propiedades map[string]*dominio.Propiedad
	shouldError bool
}

// NewMockRepository crea un nuevo mock repository
func NewMockRepository() *MockRepository {
	return &MockRepository{
		propiedades: make(map[string]*dominio.Propiedad),
		shouldError: false,
	}
}

// Implementar la interfaz PropiedadRepository
func (m *MockRepository) Crear(propiedad *dominio.Propiedad) error {
	if m.shouldError {
		return assert.AnError
	}
	m.propiedades[propiedad.ID] = propiedad
	return nil
}

func (m *MockRepository) ObtenerPorID(id string) (*dominio.Propiedad, error) {
	if m.shouldError {
		return nil, assert.AnError
	}
	propiedad, exists := m.propiedades[id]
	if !exists {
		return nil, sql.ErrNoRows
	}
	return propiedad, nil
}

func (m *MockRepository) ObtenerPorSlug(slug string) (*dominio.Propiedad, error) {
	if m.shouldError {
		return nil, assert.AnError
	}
	for _, propiedad := range m.propiedades {
		if propiedad.Slug == slug {
			return propiedad, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *MockRepository) ObtenerTodas() ([]dominio.Propiedad, error) {
	if m.shouldError {
		return nil, assert.AnError
	}

	var propiedades []dominio.Propiedad
	for _, propiedad := range m.propiedades {
		propiedades = append(propiedades, *propiedad)
	}
	return propiedades, nil
}

func (m *MockRepository) Actualizar(propiedad *dominio.Propiedad) error {
	if m.shouldError {
		return assert.AnError
	}
	if _, exists := m.propiedades[propiedad.ID]; !exists {
		return sql.ErrNoRows
	}
	m.propiedades[propiedad.ID] = propiedad
	return nil
}

func (m *MockRepository) Eliminar(id string) error {
	if m.shouldError {
		return assert.AnError
	}
	if _, exists := m.propiedades[id]; !exists {
		return sql.ErrNoRows
	}
	delete(m.propiedades, id)
	return nil
}

// Métodos helper para testing
func (m *MockRepository) SetShouldError(shouldError bool) {
	m.shouldError = shouldError
}

func (m *MockRepository) AddPropiedad(propiedad *dominio.Propiedad) {
	m.propiedades[propiedad.ID] = propiedad
}

func (m *MockRepository) GetPropiedadCount() int {
	return len(m.propiedades)
}

// TestSuite usando testify suite para agrupar tests relacionados
type PropiedadRepositoryTestSuite struct {
	suite.Suite
	mockRepo *MockRepository
}

// SetupTest se ejecuta antes de cada test
func (suite *PropiedadRepositoryTestSuite) SetupTest() {
	suite.mockRepo = NewMockRepository()
}

// Test de creación de propiedad
func (suite *PropiedadRepositoryTestSuite) TestCrear() {
	propiedad := dominio.NuevaPropiedad(
		"Casa en Quito",
		"Hermosa casa",
		"Pichincha",
		"Quito",
		"casa",
		150000,
	)

	// Test exitoso
	err := suite.mockRepo.Crear(propiedad)
	suite.NoError(err)
	suite.Equal(1, suite.mockRepo.GetPropiedadCount())

	// Test con error
	suite.mockRepo.SetShouldError(true)
	err = suite.mockRepo.Crear(propiedad)
	suite.Error(err)
}

// Test de obtener por ID
func (suite *PropiedadRepositoryTestSuite) TestObtenerPorID() {
	propiedad := dominio.NuevaPropiedad(
		"Casa en Quito",
		"Hermosa casa",
		"Pichincha",
		"Quito",
		"casa",
		150000,
	)

	// Agregar propiedad al mock
	suite.mockRepo.AddPropiedad(propiedad)

	// Test exitoso
	resultado, err := suite.mockRepo.ObtenerPorID(propiedad.ID)
	suite.NoError(err)
	suite.NotNil(resultado)
	suite.Equal(propiedad.ID, resultado.ID)
	suite.Equal(propiedad.Titulo, resultado.Titulo)

	// Test propiedad no encontrada
	_, err = suite.mockRepo.ObtenerPorID("id-inexistente")
	suite.Error(err)
	suite.Equal(sql.ErrNoRows, err)

	// Test con error
	suite.mockRepo.SetShouldError(true)
	_, err = suite.mockRepo.ObtenerPorID(propiedad.ID)
	suite.Error(err)
}

// Test de obtener por slug
func (suite *PropiedadRepositoryTestSuite) TestObtenerPorSlug() {
	propiedad := dominio.NuevaPropiedad(
		"Casa en Quito",
		"Hermosa casa",
		"Pichincha",
		"Quito",
		"casa",
		150000,
	)

	suite.mockRepo.AddPropiedad(propiedad)

	// Test exitoso
	resultado, err := suite.mockRepo.ObtenerPorSlug(propiedad.Slug)
	suite.NoError(err)
	suite.NotNil(resultado)
	suite.Equal(propiedad.Slug, resultado.Slug)

	// Test slug no encontrado
	_, err = suite.mockRepo.ObtenerPorSlug("slug-inexistente")
	suite.Error(err)
	suite.Equal(sql.ErrNoRows, err)
}

// Test de obtener todas las propiedades
func (suite *PropiedadRepositoryTestSuite) TestObtenerTodas() {
	// Agregar múltiples propiedades
	propiedades := []*dominio.Propiedad{
		dominio.NuevaPropiedad("Casa 1", "Desc 1", "Pichincha", "Quito", "casa", 100000),
		dominio.NuevaPropiedad("Casa 2", "Desc 2", "Guayas", "Guayaquil", "casa", 200000),
		dominio.NuevaPropiedad("Depto 1", "Desc 3", "Azuay", "Cuenca", "departamento", 80000),
	}

	for _, propiedad := range propiedades {
		suite.mockRepo.AddPropiedad(propiedad)
	}

	// Test exitoso
	resultado, err := suite.mockRepo.ObtenerTodas()
	suite.NoError(err)
	suite.Len(resultado, 3)

	// Verificar que todas las propiedades están presentes
	ids := make(map[string]bool)
	for _, propiedad := range resultado {
		ids[propiedad.ID] = true
	}

	for _, propiedadOriginal := range propiedades {
		suite.True(ids[propiedadOriginal.ID], "Propiedad %s debería estar en el resultado", propiedadOriginal.ID)
	}

	// Test con error
	suite.mockRepo.SetShouldError(true)
	_, err = suite.mockRepo.ObtenerTodas()
	suite.Error(err)
}

// Test de actualizar propiedad
func (suite *PropiedadRepositoryTestSuite) TestActualizar() {
	propiedad := dominio.NuevaPropiedad(
		"Casa Original",
		"Descripción original",
		"Pichincha",
		"Quito",
		"casa",
		150000,
	)

	suite.mockRepo.AddPropiedad(propiedad)

	// Modificar propiedad
	propiedad.Titulo = "Casa Actualizada"
	propiedad.Precio = 180000

	// Test exitoso
	err := suite.mockRepo.Actualizar(propiedad)
	suite.NoError(err)

	// Verificar que se actualizó
	resultado, err := suite.mockRepo.ObtenerPorID(propiedad.ID)
	suite.NoError(err)
	suite.Equal("Casa Actualizada", resultado.Titulo)
	suite.Equal(180000.0, resultado.Precio)

	// Test propiedad no encontrada
	propiedadInexistente := dominio.NuevaPropiedad("Test", "Test", "Pichincha", "Quito", "casa", 100000)
	err = suite.mockRepo.Actualizar(propiedadInexistente)
	suite.Error(err)
	suite.Equal(sql.ErrNoRows, err)

	// Test con error
	suite.mockRepo.SetShouldError(true)
	err = suite.mockRepo.Actualizar(propiedad)
	suite.Error(err)
}

// Test de eliminar propiedad
func (suite *PropiedadRepositoryTestSuite) TestEliminar() {
	propiedad := dominio.NuevaPropiedad(
		"Casa a Eliminar",
		"Esta casa será eliminada",
		"Pichincha",
		"Quito",
		"casa",
		150000,
	)

	suite.mockRepo.AddPropiedad(propiedad)
	suite.Equal(1, suite.mockRepo.GetPropiedadCount())

	// Test exitoso
	err := suite.mockRepo.Eliminar(propiedad.ID)
	suite.NoError(err)
	suite.Equal(0, suite.mockRepo.GetPropiedadCount())

	// Verificar que ya no existe
	_, err = suite.mockRepo.ObtenerPorID(propiedad.ID)
	suite.Error(err)
	suite.Equal(sql.ErrNoRows, err)

	// Test eliminar propiedad inexistente
	err = suite.mockRepo.Eliminar("id-inexistente")
	suite.Error(err)
	suite.Equal(sql.ErrNoRows, err)

	// Test con error
	propiedad2 := dominio.NuevaPropiedad("Casa 2", "Desc", "Pichincha", "Quito", "casa", 100000)
	suite.mockRepo.AddPropiedad(propiedad2)
	suite.mockRepo.SetShouldError(true)
	err = suite.mockRepo.Eliminar(propiedad2.ID)
	suite.Error(err)
}

// Ejecutar la suite de tests
func TestPropiedadRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PropiedadRepositoryTestSuite))
}

// Tests unitarios adicionales para funciones específicas

// TestNuevoPropiedadRepositoryPostgres verifica la creación del repositorio
func TestNuevoPropiedadRepositoryPostgres(t *testing.T) {
	// Crear una conexión mock (nil en este caso)
	var db *sql.DB

	repo := NuevoPropiedadRepositoryPostgres(db)

	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

// TestRepositoryInterface verifica que MockRepository implementa la interfaz
func TestRepositoryInterface(t *testing.T) {
	var repo PropiedadRepository
	repo = NewMockRepository()

	assert.NotNil(t, repo)

	// Verificar que podemos usar todos los métodos de la interfaz
	propiedad := dominio.NuevaPropiedad("Test", "Test", "Pichincha", "Quito", "casa", 100000)

	err := repo.Crear(propiedad)
	assert.NoError(t, err)

	resultado, err := repo.ObtenerPorID(propiedad.ID)
	assert.NoError(t, err)
	assert.Equal(t, propiedad.ID, resultado.ID)

	propiedades, err := repo.ObtenerTodas()
	assert.NoError(t, err)
	assert.Len(t, propiedades, 1)

	err = repo.Actualizar(propiedad)
	assert.NoError(t, err)

	err = repo.Eliminar(propiedad.ID)
	assert.NoError(t, err)
}

// TestConcurrentAccess prueba el acceso concurrente al mock (ejemplo básico)
func TestConcurrentAccess(t *testing.T) {
	repo := NewMockRepository()

	// Crear varias propiedades concurrentemente
	propiedades := make([]*dominio.Propiedad, 10)
	for i := 0; i < 10; i++ {
		propiedades[i] = dominio.NuevaPropiedad(
			"Casa "+string(rune(i+'0')),
			"Descripción",
			"Pichincha",
			"Quito",
			"casa",
			100000+float64(i*10000),
		)
	}

	// Agregar todas las propiedades
	for _, propiedad := range propiedades {
		err := repo.Crear(propiedad)
		assert.NoError(t, err)
	}

	// Verificar que todas se crearon
	todasLasPropiedades, err := repo.ObtenerTodas()
	assert.NoError(t, err)
	assert.Len(t, todasLasPropiedades, 10)
}

// Ejemplo de test con table-driven para diferentes tipos de errores
func TestMockRepositoryErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		shouldError bool
		expectError bool
	}{
		{
			name:        "sin error",
			shouldError: false,
			expectError: false,
		},
		{
			name:        "con error",
			shouldError: true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			repo.SetShouldError(tt.shouldError)

			propiedad := dominio.NuevaPropiedad("Test", "Test", "Pichincha", "Quito", "casa", 100000)

			err := repo.Crear(propiedad)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// BenchmarkMockRepository - ejemplo de benchmark para el mock
func BenchmarkMockRepository_Crear(b *testing.B) {
	repo := NewMockRepository()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		propiedad := dominio.NuevaPropiedad(
			"Casa Benchmark",
			"Descripción benchmark",
			"Pichincha",
			"Quito",
			"casa",
			100000,
		)
		_ = repo.Crear(propiedad)
	}
}

// BenchmarkMockRepository_ObtenerTodas - benchmark para obtener todas las propiedades
func BenchmarkMockRepository_ObtenerTodas(b *testing.B) {
	repo := NewMockRepository()

	// Preparar datos
	for i := 0; i < 1000; i++ {
		propiedad := dominio.NuevaPropiedad(
			"Casa "+string(rune(i)),
			"Descripción",
			"Pichincha",
			"Quito",
			"casa",
			100000,
		)
		_ = repo.Crear(propiedad)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.ObtenerTodas()
	}
}
