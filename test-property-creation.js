/**
 * Test script for Property Creation
 * Tests the complete flow: Frontend → Server Actions → Backend Go → PostgreSQL
 */

const testPropertyData = {
  // Información básica (requerida)
  title: "Casa moderna en Samborondón con piscina",
  description: "Hermosa casa de 3 pisos con acabados de lujo, piscina, jardín y excelente ubicación en urbanización cerrada",
  price: 285000,
  type: "house",
  status: "available",
  
  // Ubicación (requerida)
  province: "Guayas",
  city: "Samborondón",
  address: "Urbanización La Puntilla, Mz. 15 Villa 7",
  sector: "La Puntilla",
  latitude: -2.1894927,
  longitude: -79.8890853,
  location_precision: "exact",
  
  // Características
  bedrooms: 4,
  bathrooms: 3.5,
  area_m2: 320,
  parking_spaces: 2,
  year_built: 2020,
  floors: 3,
  
  // Precios adicionales
  rent_price: 1800,
  common_expenses: 85,
  price_per_m2: 890.63,
  
  // Multimedia
  main_image: "https://example.com/property1.jpg",
  images: ["https://example.com/property1.jpg", "https://example.com/property2.jpg"],
  video_tour: "https://example.com/video-tour.mp4",
  tour_360: "https://example.com/tour-360",
  
  // Estado y clasificación
  property_status: "new",
  tags: ["piscina", "jardín", "urbanización cerrada", "lujo"],
  featured: true,
  view_count: 0,
  
  // Amenidades
  furnished: false,
  garage: true,
  pool: true,
  garden: true,
  terrace: true,
  balcony: false,
  security: true,
  elevator: false,
  air_conditioning: true,
  
  // Contact info
  contact_phone: "0992345678",
  contact_email: "propietario@ejemplo.com",
  notes: "Propiedad en excelente estado, lista para habitar. Acepta financiamiento bancario."
};

async function testPropertyCreation() {
  try {
    console.log('🏠 Testing Property Creation API...');
    console.log('📊 Property Data:', JSON.stringify(testPropertyData, null, 2));
    
    const response = await fetch('http://localhost:8080/api/properties', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'User-Agent': 'Property-Test-Script/1.0',
      },
      body: JSON.stringify(testPropertyData),
    });

    console.log('📡 Response Status:', response.status);
    console.log('📡 Response Headers:', Object.fromEntries(response.headers));

    if (!response.ok) {
      const errorData = await response.text();
      console.error('❌ API Error Response:', errorData);
      throw new Error(`HTTP ${response.status}: ${errorData}`);
    }

    const createdProperty = await response.json();
    console.log('✅ Property Created Successfully!');
    console.log('🆔 Property ID:', createdProperty.id);
    console.log('📋 Created Property:', JSON.stringify(createdProperty, null, 2));
    
    return createdProperty;
  } catch (error) {
    console.error('❌ Test Failed:', error.message);
    throw error;
  }
}

// Run the test
testPropertyCreation()
  .then(() => {
    console.log('🎉 Test completed successfully!');
    console.log('');
    console.log('Next steps:');
    console.log('1. ✅ Backend API is working');
    console.log('2. 🌐 Test frontend at: http://localhost:3004/create-property');
    console.log('3. 🧪 Fill out the form and verify Server Actions work');
    process.exit(0);
  })
  .catch((error) => {
    console.error('💥 Test failed:', error.message);
    process.exit(1);
  });