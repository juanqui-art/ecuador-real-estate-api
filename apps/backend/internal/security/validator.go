package security

import (
	"net"
	"regexp"
	"strings"
	"unicode"
)

// InputValidator provides security validation for user inputs
type InputValidator struct {
	sqlInjectionPatterns []string
	xssPatterns          []string
	pathTraversalPatterns []string
	maxInputLength       int
}

// NewInputValidator creates a new input validator with security patterns
func NewInputValidator() *InputValidator {
	return &InputValidator{
		sqlInjectionPatterns: []string{
			`(?i)\b(union|select|insert|update|delete|drop|create|alter|exec|execute)\b`,
			`(?i)(\-\-|\#|\/\*|\*\/|;)`,
			`(?i)\b(or|and)\s+\d+\s*=\s*\d+`,
			`(?i)\b(or|and)\s+.*\s*=\s*.*`,
			`(?i)'.*(\s|;|\/\*|\-\-|\#).*'`,
			`(?i)".*(\s|;|\/\*|\-\-|\#).*"`,
		},
		xssPatterns: []string{
			`(?i)<script[^>]*>.*?</script>`,
			`(?i)<iframe[^>]*>.*?</iframe>`,
			`(?i)<object[^>]*>.*?</object>`,
			`(?i)<embed[^>]*>.*?</embed>`,
			`(?i)<link[^>]*>`,
			`(?i)<meta[^>]*>`,
			`(?i)javascript:`,
			`(?i)vbscript:`,
			`(?i)data:text/html`,
			`(?i)on\w+\s*=`,
			`(?i)expression\s*\(`,
			`(?i)@import`,
		},
		pathTraversalPatterns: []string{
			`\.\.\/`,
			`\.\.\\`,
			`\.\.\%2f`,
			`\.\.\%5c`,
			`\/etc\/passwd`,
			`\/proc\/`,
			`\\windows\\`,
			`\\system32\\`,
		},
		maxInputLength: 10000, // 10KB max input length
	}
}

// ValidateInput performs comprehensive input validation
func (v *InputValidator) ValidateInput(input string, fieldName string) *ValidationResult {
	result := &ValidationResult{
		IsValid: true,
		Field:   fieldName,
		Threats: make([]ThreatInfo, 0),
	}

	// Check input length
	if len(input) > v.maxInputLength {
		result.IsValid = false
		result.Threats = append(result.Threats, ThreatInfo{
			Type:        "input_length",
			Severity:    "medium",
			Description: "Input exceeds maximum allowed length",
			Pattern:     "",
		})
	}

	// Check for SQL injection patterns
	if threats := v.checkSQLInjection(input); len(threats) > 0 {
		result.IsValid = false
		result.Threats = append(result.Threats, threats...)
	}

	// Check for XSS patterns
	if threats := v.checkXSS(input); len(threats) > 0 {
		result.IsValid = false
		result.Threats = append(result.Threats, threats...)
	}

	// Check for path traversal patterns
	if threats := v.checkPathTraversal(input); len(threats) > 0 {
		result.IsValid = false
		result.Threats = append(result.Threats, threats...)
	}

	// Check for suspicious characters
	if threats := v.checkSuspiciousCharacters(input); len(threats) > 0 {
		result.IsValid = false
		result.Threats = append(result.Threats, threats...)
	}

	return result
}

// ValidationResult contains the result of input validation
type ValidationResult struct {
	IsValid bool         `json:"is_valid"`
	Field   string       `json:"field"`
	Threats []ThreatInfo `json:"threats"`
}

// ThreatInfo contains information about detected threats
type ThreatInfo struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`    // low, medium, high, critical
	Description string `json:"description"`
	Pattern     string `json:"pattern"`
}

// checkSQLInjection checks for SQL injection patterns
func (v *InputValidator) checkSQLInjection(input string) []ThreatInfo {
	threats := make([]ThreatInfo, 0)
	
	for _, pattern := range v.sqlInjectionPatterns {
		if matched, _ := regexp.MatchString(pattern, input); matched {
			threats = append(threats, ThreatInfo{
				Type:        "sql_injection",
				Severity:    "high",
				Description: "Potential SQL injection attempt detected",
				Pattern:     pattern,
			})
		}
	}
	
	return threats
}

// checkXSS checks for Cross-Site Scripting patterns
func (v *InputValidator) checkXSS(input string) []ThreatInfo {
	threats := make([]ThreatInfo, 0)
	
	for _, pattern := range v.xssPatterns {
		if matched, _ := regexp.MatchString(pattern, input); matched {
			threats = append(threats, ThreatInfo{
				Type:        "xss",
				Severity:    "high",
				Description: "Potential XSS attempt detected",
				Pattern:     pattern,
			})
		}
	}
	
	return threats
}

// checkPathTraversal checks for path traversal patterns
func (v *InputValidator) checkPathTraversal(input string) []ThreatInfo {
	threats := make([]ThreatInfo, 0)
	
	for _, pattern := range v.pathTraversalPatterns {
		if matched, _ := regexp.MatchString(pattern, input); matched {
			threats = append(threats, ThreatInfo{
				Type:        "path_traversal",
				Severity:    "medium",
				Description: "Potential path traversal attempt detected",
				Pattern:     pattern,
			})
		}
	}
	
	return threats
}

// checkSuspiciousCharacters checks for suspicious character patterns
func (v *InputValidator) checkSuspiciousCharacters(input string) []ThreatInfo {
	threats := make([]ThreatInfo, 0)
	
	// Check for excessive special characters (possible obfuscation)
	specialCharCount := 0
	for _, r := range input {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
			specialCharCount++
		}
	}
	
	if len(input) > 0 && float64(specialCharCount)/float64(len(input)) > 0.5 {
		threats = append(threats, ThreatInfo{
			Type:        "suspicious_characters",
			Severity:    "low",
			Description: "High ratio of special characters detected",
			Pattern:     "special_char_ratio",
		})
	}
	
	// Check for control characters
	for _, r := range input {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			threats = append(threats, ThreatInfo{
				Type:        "control_characters",
				Severity:    "medium",
				Description: "Control characters detected in input",
				Pattern:     "control_chars",
			})
			break
		}
	}
	
	return threats
}

// SanitizeInput removes potentially dangerous characters from input
func (v *InputValidator) SanitizeInput(input string) string {
	// Remove control characters except basic whitespace
	result := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return -1
		}
		return r
	}, input)
	
	// Trim whitespace
	result = strings.TrimSpace(result)
	
	return result
}

// ValidateEmail validates email format with additional security checks
func (v *InputValidator) ValidateEmail(email string) bool {
	if len(email) == 0 || len(email) > 254 {
		return false
	}
	
	// Basic RFC 5322 regex (simplified)
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	if !matched {
		return false
	}
	
	// Additional security checks
	if strings.Contains(email, "..") {
		return false
	}
	
	// Check for suspicious patterns in email
	result := v.ValidateInput(email, "email")
	return result.IsValid
}

// ValidateURL validates URL format and checks for suspicious patterns
func (v *InputValidator) ValidateURL(url string) bool {
	if len(url) == 0 || len(url) > 2048 {
		return false
	}
	
	// Check for basic URL format
	urlRegex := `^https?:\/\/([\w\-]+(\.[\w\-]+)+)([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?$`
	matched, _ := regexp.MatchString(urlRegex, url)
	if !matched {
		return false
	}
	
	// Additional security checks
	result := v.ValidateInput(url, "url")
	return result.IsValid
}

// ValidateFilename validates filename for safe file operations
func (v *InputValidator) ValidateFilename(filename string) bool {
	if len(filename) == 0 || len(filename) > 255 {
		return false
	}
	
	// Check for valid filename characters
	filenameRegex := `^[a-zA-Z0-9._\-\s]+$`
	matched, _ := regexp.MatchString(filenameRegex, filename)
	if !matched {
		return false
	}
	
	// Check for reserved names (Windows)
	reservedNames := []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}
	
	upperFilename := strings.ToUpper(filename)
	for _, reserved := range reservedNames {
		if upperFilename == reserved {
			return false
		}
	}
	
	// Check for path traversal
	result := v.ValidateInput(filename, "filename")
	return result.IsValid
}

// IPValidator provides IP address validation and security checks
type IPValidator struct {
	blockedRanges []net.IPNet
}

// NewIPValidator creates a new IP validator with blocked ranges
func NewIPValidator() *IPValidator {
	validator := &IPValidator{
		blockedRanges: make([]net.IPNet, 0),
	}
	
	// Add common blocked ranges (private networks, localhost, etc.)
	blockedCIDRs := []string{
		"127.0.0.0/8",   // Localhost
		"10.0.0.0/8",    // Private
		"172.16.0.0/12", // Private
		"192.168.0.0/16", // Private
		"169.254.0.0/16", // Link-local
		"224.0.0.0/4",   // Multicast
		"::1/128",       // IPv6 localhost
		"fc00::/7",      // IPv6 private
		"fe80::/10",     // IPv6 link-local
	}
	
	for _, cidr := range blockedCIDRs {
		_, network, err := net.ParseCIDR(cidr)
		if err == nil {
			validator.blockedRanges = append(validator.blockedRanges, *network)
		}
	}
	
	return validator
}

// ValidateIP validates IP address and checks against blocked ranges
func (v *IPValidator) ValidateIP(ipStr string) (bool, string) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, "invalid IP format"
	}
	
	// Check against blocked ranges
	for _, blocked := range v.blockedRanges {
		if blocked.Contains(ip) {
			return false, "IP address in blocked range"
		}
	}
	
	return true, ""
}

// IsPrivateIP checks if an IP address is in a private range
func (v *IPValidator) IsPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}
	
	for _, cidr := range privateRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err == nil && network.Contains(ip) {
			return true
		}
	}
	
	return false
}

// PasswordValidator provides password strength validation
type PasswordValidator struct {
	minLength       int
	requireUpper    bool
	requireLower    bool
	requireDigits   bool
	requireSpecial  bool
	bannedPasswords []string
}

// NewPasswordValidator creates a new password validator
func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		minLength:      8,
		requireUpper:   true,
		requireLower:   true,
		requireDigits:  true,
		requireSpecial: true,
		bannedPasswords: []string{
			"password", "123456", "123456789", "qwerty", "abc123",
			"password123", "admin", "letmein", "welcome", "monkey",
			"1234567890", "password1", "qwerty123", "123qwe",
		},
	}
}

// ValidatePassword validates password strength
func (v *PasswordValidator) ValidatePassword(password string) (bool, []string) {
	var errors []string
	
	// Check minimum length
	if len(password) < v.minLength {
		errors = append(errors, "Password must be at least 8 characters long")
	}
	
	// Check maximum length (prevent DoS)
	if len(password) > 128 {
		errors = append(errors, "Password must be less than 128 characters")
	}
	
	// Check character requirements
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}
	
	if v.requireUpper && !hasUpper {
		errors = append(errors, "Password must contain at least one uppercase letter")
	}
	
	if v.requireLower && !hasLower {
		errors = append(errors, "Password must contain at least one lowercase letter")
	}
	
	if v.requireDigits && !hasDigit {
		errors = append(errors, "Password must contain at least one digit")
	}
	
	if v.requireSpecial && !hasSpecial {
		errors = append(errors, "Password must contain at least one special character")
	}
	
	// Check against banned passwords
	lowerPassword := strings.ToLower(password)
	for _, banned := range v.bannedPasswords {
		if lowerPassword == banned {
			errors = append(errors, "Password is too common")
			break
		}
	}
	
	// Check for sequential patterns
	if v.hasSequentialPattern(password) {
		errors = append(errors, "Password contains sequential patterns")
	}
	
	return len(errors) == 0, errors
}

// hasSequentialPattern checks for sequential character patterns
func (v *PasswordValidator) hasSequentialPattern(password string) bool {
	sequences := []string{
		"abcdefghijklmnopqrstuvwxyz",
		"0123456789",
		"qwertyuiop",
		"asdfghjkl",
		"zxcvbnm",
	}
	
	lowerPassword := strings.ToLower(password)
	
	for _, seq := range sequences {
		for i := 0; i <= len(seq)-3; i++ {
			if strings.Contains(lowerPassword, seq[i:i+3]) {
				return true
			}
			// Check reverse sequence
			reverse := reverseString(seq[i:i+3])
			if strings.Contains(lowerPassword, reverse) {
				return true
			}
		}
	}
	
	return false
}

// reverseString reverses a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}