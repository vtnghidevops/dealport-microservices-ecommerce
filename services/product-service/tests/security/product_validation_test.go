package security

import (
	"strings"
	"testing"
)

// TestInputValidation tests validation of product inputs for security issues
func TestInputValidation(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		isValid     bool
		description string
	}{
		{
			name:        "Valid product name",
			input:       "Genuine Product Name",
			isValid:     true,
			description: "Regular product name with no injection attempts",
		},
		{
			name:        "SQL Injection attempt",
			input:       "Product'; DROP TABLE products; --",
			isValid:     false,
			description: "SQL injection attempt in product name",
		},
		{
			name:        "XSS attempt",
			input:       "Product<script>alert('XSS')</script>",
			isValid:     false,
			description: "XSS attempt in product name",
		},
		{
			name:        "Command injection attempt",
			input:       "Product; rm -rf /;",
			isValid:     false,
			description: "Command injection attempt in product name",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := validateInput(test.input)
			if result != test.isValid {
				t.Errorf("Expected validation result %v for input '%s', got %v",
					test.isValid, test.input, result)
			}
		})
	}
}

// TestPriceValidation tests validation of product prices
func TestPriceValidation(t *testing.T) {
	tests := []struct {
		name        string
		price       float64
		isValid     bool
		description string
	}{
		{
			name:        "Valid price",
			price:       99.99,
			isValid:     true,
			description: "Regular positive price",
		},
		{
			name:        "Zero price",
			price:       0,
			isValid:     true, // Free products can be valid
			description: "Zero price for free product",
		},
		{
			name:        "Negative price",
			price:       -10.50,
			isValid:     false,
			description: "Negative price should be invalid",
		},
		{
			name:        "Extremely high price",
			price:       1000000000,
			isValid:     false,
			description: "Extremely high price that might indicate an attack",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := validatePrice(test.price)
			if result != test.isValid {
				t.Errorf("Expected validation result %v for price %.2f, got %v",
					test.isValid, test.price, result)
			}
		})
	}
}

// TestSanitizeHTML tests sanitization of HTML content in product descriptions
func TestSanitizeHTML(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput string
		description    string
	}{
		{
			name:           "Clean text",
			input:          "This is a clean product description",
			expectedOutput: "This is a clean product description",
			description:    "Clean text should remain unchanged",
		},
		{
			name:           "Text with script tag",
			input:          "Description <script>alert('malicious')</script> here",
			expectedOutput: "Description  here",
			description:    "Script tags should be removed",
		},
		{
			name:           "Text with safe HTML formatting",
			input:          "Description with <b>bold</b> and <i>italic</i> text",
			expectedOutput: "Description with <b>bold</b> and <i>italic</i> text",
			description:    "Safe HTML formatting should be preserved",
		},
		{
			name:           "Text with iframe",
			input:          "Description with <iframe src='malicious.com'></iframe> embedded",
			expectedOutput: "Description with  embedded",
			description:    "iframes should be removed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := sanitizeHTML(test.input)
			if result != test.expectedOutput {
				t.Errorf("Expected '%s' but got '%s'", test.expectedOutput, result)
			}
		})
	}
}

// validateInput is a simple input validator for testing
func validateInput(input string) bool {
	// Check for SQL injection patterns
	sqlPatterns := []string{
		"'", "--", "DROP", "DELETE", "UPDATE", "INSERT", ";",
	}

	for _, pattern := range sqlPatterns {
		if strings.Contains(strings.ToUpper(input), strings.ToUpper(pattern)) {
			return false
		}
	}

	// Check for XSS patterns
	xssPatterns := []string{
		"<script>", "</script>", "javascript:", "onerror=", "onload=",
	}

	for _, pattern := range xssPatterns {
		if strings.Contains(strings.ToLower(input), strings.ToLower(pattern)) {
			return false
		}
	}

	// Check for command injection
	cmdPatterns := []string{
		"rm -rf", "; ", "&&", "||", "`", "$(",
	}

	for _, pattern := range cmdPatterns {
		if strings.Contains(input, pattern) {
			return false
		}
	}

	return true
}

// validatePrice validates a product price for security issues
func validatePrice(price float64) bool {
	// Price should not be negative
	if price < 0 {
		return false
	}

	// Price should not be unreasonably high (example threshold)
	if price > 1000000 {
		return false
	}

	return true
}

// sanitizeHTML performs basic HTML sanitization
func sanitizeHTML(input string) string {
	// This is a simplified version for testing
	// In production, use a proper HTML sanitization library

	// Remove script tags
	input = removeHTMLTag(input, "script")

	// Remove iframe tags
	input = removeHTMLTag(input, "iframe")

	// Remove event handlers
	input = strings.ReplaceAll(input, "javascript:", "")
	input = strings.ReplaceAll(input, "onerror=", "")
	input = strings.ReplaceAll(input, "onload=", "")

	return input
}

// removeHTMLTag removes HTML tags of a specific type
func removeHTMLTag(input, tag string) string {
	openTag := "<" + tag
	closeTag := "</" + tag + ">"

	for {
		startIdx := strings.Index(strings.ToLower(input), strings.ToLower(openTag))
		if startIdx == -1 {
			break
		}

		endIdx := strings.Index(strings.ToLower(input[startIdx:]), strings.ToLower(closeTag))
		if endIdx == -1 {
			break
		}
		endIdx += startIdx + len(closeTag)

		input = input[:startIdx] + input[endIdx:]
	}

	return input
}
