package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputValidator_ValidateInput(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		name     string
		input    string
		field    string
		expected bool
	}{
		{
			name:     "Valid input",
			input:    "Hello World",
			field:    "message",
			expected: true,
		},
		{
			name:     "SQL injection attempt",
			input:    "'; DROP TABLE users; --",
			field:    "username",
			expected: false,
		},
		{
			name:     "XSS attempt",
			input:    "<script>alert('xss')</script>",
			field:    "comment",
			expected: false,
		},
		{
			name:     "Path traversal attempt",
			input:    "../../../etc/passwd",
			field:    "filename",
			expected: false,
		},
		{
			name:     "Long input",
			input:    string(make([]byte, 15000)), // Exceeds max length
			field:    "description",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validator.ValidateInput(tc.input, tc.field)
			assert.Equal(t, tc.expected, result.IsValid, "Test case: %s", tc.name)
			assert.Equal(t, tc.field, result.Field)
			
			if !result.IsValid {
				assert.Greater(t, len(result.Threats), 0, "Should have threats when invalid")
			}
		})
	}
}

func TestInputValidator_ValidateEmail(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		email    string
		expected bool
	}{
		{"valid@example.com", true},
		{"user.name@domain.co.uk", true},
		{"user+tag@example.com", true},
		{"invalid-email", false},
		{"@domain.com", false},
		{"user@", false},
		{"user..name@domain.com", false}, // Double dots
		{"", false},
		{string(make([]byte, 300)), false}, // Too long
	}

	for _, tc := range testCases {
		result := validator.ValidateEmail(tc.email)
		assert.Equal(t, tc.expected, result, "Email: %s", tc.email)
	}
}

func TestInputValidator_ValidateURL(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		url      string
		expected bool
	}{
		{"https://example.com", true},
		{"http://subdomain.example.com/path", true},
		{"https://example.com/path?query=value", true},
		{"ftp://example.com", false}, // Not http/https
		{"not-a-url", false},
		{"", false},
		{string(make([]byte, 3000)), false}, // Too long
	}

	for _, tc := range testCases {
		result := validator.ValidateURL(tc.url)
		assert.Equal(t, tc.expected, result, "URL: %s", tc.url)
	}
}

func TestInputValidator_ValidateFilename(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		filename string
		expected bool
	}{
		{"document.pdf", true},
		{"my-file_123.txt", true},
		{"file with spaces.doc", true},
		{"../malicious.exe", false}, // Path traversal
		{"CON", false},              // Reserved name
		{"file<script>.txt", false}, // Invalid characters
		{"", false},
		{string(make([]byte, 300)), false}, // Too long
	}

	for _, tc := range testCases {
		result := validator.ValidateFilename(tc.filename)
		assert.Equal(t, tc.expected, result, "Filename: %s", tc.filename)
	}
}

func TestInputValidator_SanitizeInput(t *testing.T) {
	validator := NewInputValidator()

	testCases := []struct {
		input    string
		expected string
	}{
		{"  Hello World  ", "Hello World"},
		{"Text with\x00null", "Text withnull"},
		{"Normal text\n\r\twith whitespace", "Normal text\n\r\twith whitespace"},
	}

	for _, tc := range testCases {
		result := validator.SanitizeInput(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}

func TestIPValidator_ValidateIP(t *testing.T) {
	validator := NewIPValidator()

	testCases := []struct {
		ip       string
		expected bool
		reason   string
	}{
		{"8.8.8.8", true, ""},
		{"127.0.0.1", false, "IP address in blocked range"},
		{"192.168.1.1", false, "IP address in blocked range"},
		{"invalid-ip", false, "invalid IP format"},
		{"::1", false, "IP address in blocked range"},
	}

	for _, tc := range testCases {
		isValid, reason := validator.ValidateIP(tc.ip)
		assert.Equal(t, tc.expected, isValid, "IP: %s", tc.ip)
		if !isValid {
			assert.NotEmpty(t, reason)
		}
	}
}

func TestIPValidator_IsPrivateIP(t *testing.T) {
	validator := NewIPValidator()

	testCases := []struct {
		ip       string
		expected bool
	}{
		{"192.168.1.1", true},
		{"10.0.0.1", true},
		{"172.16.0.1", true},
		{"8.8.8.8", false},
		{"invalid-ip", false},
	}

	for _, tc := range testCases {
		result := validator.IsPrivateIP(tc.ip)
		assert.Equal(t, tc.expected, result, "IP: %s", tc.ip)
	}
}

func TestPasswordValidator_ValidatePassword(t *testing.T) {
	validator := NewPasswordValidator()

	testCases := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "Strong password",
			password: "MyStr0ng!Pass",
			expected: true,
		},
		{
			name:     "Too short",
			password: "Short1!",
			expected: false,
		},
		{
			name:     "No uppercase",
			password: "lowercase123!",
			expected: false,
		},
		{
			name:     "No lowercase",
			password: "UPPERCASE123!",
			expected: false,
		},
		{
			name:     "No digits",
			password: "NoDigits!",
			expected: false,
		},
		{
			name:     "No special characters",
			password: "NoSpecial123",
			expected: false,
		},
		{
			name:     "Common password",
			password: "password123",
			expected: false,
		},
		{
			name:     "Sequential pattern",
			password: "Abcd1234!",
			expected: false,
		},
		{
			name:     "Too long",
			password: string(make([]byte, 200)),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid, errors := validator.ValidatePassword(tc.password)
			assert.Equal(t, tc.expected, isValid, "Test case: %s", tc.name)
			
			if !isValid {
				assert.Greater(t, len(errors), 0, "Should have error messages when invalid")
			}
		})
	}
}

func TestPasswordValidator_hasSequentialPattern(t *testing.T) {
	validator := NewPasswordValidator()

	testCases := []struct {
		password string
		expected bool
	}{
		{"abc123", true},
		{"xyz789", true},
		{"qwerty", true},
		{"321cba", true}, // Reverse sequence
		{"randompassword", false},
		{"MyStr0ng!Pass", false},
	}

	for _, tc := range testCases {
		result := validator.hasSequentialPattern(tc.password)
		assert.Equal(t, tc.expected, result, "Password: %s", tc.password)
	}
}

// Benchmark tests
func BenchmarkInputValidator_ValidateInput(b *testing.B) {
	validator := NewInputValidator()
	input := "This is a test input string with some content"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateInput(input, "test")
	}
}

func BenchmarkInputValidator_ValidateEmail(b *testing.B) {
	validator := NewInputValidator()
	email := "user@example.com"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateEmail(email)
	}
}

func BenchmarkPasswordValidator_ValidatePassword(b *testing.B) {
	validator := NewPasswordValidator()
	password := "MyStr0ng!Password123"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidatePassword(password)
	}
}