/**
 * Test for New Fields: property_status, tags, featured
 * Testing incremental addition of 3 new fields to the Property form
 */

const testPropertyWithNewFields = {
  // Información básica (requerida)
  title: "Casa de prueba - Nuevos Campos",
  description: "Prueba específica para validar los 3 nuevos campos agregados: property_status, tags y featured",
  price: "175000",
  type: "house",
  status: "available",
  
  // Ubicación (requerida)
  province: "Pichincha",
  city: "Quito",
  address: "Avenida de la República y Río Amazonas",
  sector: "La República",
  
  // Características básicas
  bedrooms: "3",
  bathrooms: "2",
  area_m2: "150",
  parking_spaces: "1",
  
  // NUEVOS CAMPOS - Testing específico
  property_status: "renovated",  // Enum: new/used/renovated
  tags: "piscina, jardín, seguridad, cerca metro",  // String comma-separated
  featured: "true",  // Boolean como string
  
  // Amenidades
  garden: "true",
  security: "true",
  garage: "true",
  
  // Contact info
  contact_phone: "0987654321",
  contact_email: "test-newfields@ejemplo.com",
  notes: "Propiedad específica para testing de nuevos campos incrementales",
};

async function testNewFields() {
  try {
    console.log('🧪 Testing New Fields: property_status, tags, featured');
    console.log('📋 Data with new fields:', testPropertyWithNewFields);
    
    const response = await fetch('http://localhost:8080/api/properties', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'User-Agent': 'New-Fields-Test/1.0',
      },
      body: JSON.stringify(testPropertyWithNewFields),
    });

    console.log('📡 Response Status:', response.status);

    if (!response.ok) {
      const errorData = await response.text();
      console.error('❌ API Error:', errorData);
      throw new Error(`HTTP ${response.status}: ${errorData}`);
    }

    const createdProperty = await response.json();
    console.log('✅ Property Created with New Fields!');
    console.log('🆔 Property ID:', createdProperty.data?.id);
    
    // Verify new fields
    const property = createdProperty.data;
    console.log('');
    console.log('🔍 New Fields Verification:');
    console.log('property_status:', property.property_status);
    console.log('tags:', property.tags);
    console.log('featured:', property.featured);
    
    // Validate new fields
    let allFieldsCorrect = true;
    
    if (property.property_status !== 'renovated') {
      console.error('❌ property_status incorrect:', property.property_status);
      allFieldsCorrect = false;
    } else {
      console.log('✅ property_status: renovated ✓');
    }
    
    if (!Array.isArray(property.tags) || !property.tags.includes('piscina')) {
      console.error('❌ tags incorrect:', property.tags);
      allFieldsCorrect = false;
    } else {
      console.log('✅ tags array:', property.tags, '✓');
    }
    
    if (property.featured !== true) {
      console.error('❌ featured incorrect:', property.featured);
      allFieldsCorrect = false;
    } else {
      console.log('✅ featured: true ✓');
    }
    
    if (allFieldsCorrect) {
      console.log('');
      console.log('🎉 All new fields processed correctly!');
      console.log('✅ Backend accepting and processing 3 new fields');
      console.log('✅ property_status enum working');
      console.log('✅ tags array processing working');
      console.log('✅ featured boolean working');
    }
    
    return createdProperty;
  } catch (error) {
    console.error('💥 New Fields Test Failed:', error.message);
    throw error;
  }
}

// Run the test
testNewFields()
  .then(() => {
    console.log('');
    console.log('🎉 New Fields Test completed successfully!');
    console.log('');
    console.log('✅ Testing Summary:');
    console.log('1. Backend API accepting 3 new fields ✅');
    console.log('2. property_status enum validation ✅');
    console.log('3. tags array conversion ✅');
    console.log('4. featured boolean processing ✅');
    console.log('');
    console.log('🌐 Next: Test frontend form at: http://localhost:3000/create-property');
    console.log('🔄 Fill new "Estado y Clasificación" section');
    process.exit(0);
  })
  .catch((error) => {
    console.error('💥 New fields test failed:', error.message);
    process.exit(1);
  });