package main

import (
	"fmt"
	"strings"
	"testing"
)

// =============================================================================
// TABLE-DRIVEN TESTS TUTORIAL
// =============================================================================

// Table-driven tests are a Go testing pattern where you define test cases
// in a data structure (table) and iterate through them in a single test function.
// This makes tests more readable, maintainable, and comprehensive.

// =============================================================================
// 1. BASIC TABLE-DRIVEN TEST
// =============================================================================

// Function to test: Simple string reversal
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Basic table-driven test
func TestReverseString(t *testing.T) {
	// Define test cases in a table
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single character",
			input:    "a",
			expected: "a",
		},
		{
			name:     "simple word",
			input:    "hello",
			expected: "olleh",
		},
		{
			name:     "palindrome",
			input:    "racecar",
			expected: "racecar",
		},
		{
			name:     "with spaces",
			input:    "hello world",
			expected: "dlrow olleh",
		},
		{
			name:     "unicode characters",
			input:    "café",
			expected: "éfac",
		},
	}

	// Iterate through test cases
	for _, tt := range tests {
		// Use t.Run to create sub-tests
		t.Run(tt.name, func(t *testing.T) {
			result := reverseString(tt.input)
			if result != tt.expected {
				t.Errorf("reverseString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// 2. ADVANCED TABLE-DRIVEN TEST WITH ERROR CASES
// =============================================================================

// Function to test: Division with error handling
func divideNumbers(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// Advanced table-driven test with error cases
func TestDivide(t *testing.T) {
	tests := []struct {
		name        string
		a, b        int
		expected    int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "normal division",
			a:           10,
			b:           2,
			expected:    5,
			expectError: false,
		},
		{
			name:        "division by zero",
			a:           10,
			b:           0,
			expected:    0,
			expectError: true,
			errorMsg:    "division by zero",
		},
		{
			name:        "negative numbers",
			a:           -10,
			b:           2,
			expected:    -5,
			expectError: false,
		},
		{
			name:        "zero dividend",
			a:           0,
			b:           5,
			expected:    0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := divideNumbers(tt.a, tt.b)

			// Check error cases
			if tt.expectError {
				if err == nil {
					t.Errorf("divide(%d, %d) expected error but got none", tt.a, tt.b)
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("divide(%d, %d) error = %v, want error containing %q", tt.a, tt.b, err, tt.errorMsg)
				}
				return
			}

			// Check success cases
			if err != nil {
				t.Errorf("divide(%d, %d) unexpected error: %v", tt.a, tt.b, err)
				return
			}

			if result != tt.expected {
				t.Errorf("divide(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// 3. TABLE-DRIVEN TEST WITH COMPLEX STRUCTURES
// =============================================================================

// Complex data structure to test
type User struct {
	Name  string
	Age   int
	Email string
}

// Function to test: Validate user data
func validateUser(user User) []string {
	var errors []string

	if user.Name == "" {
		errors = append(errors, "name is required")
	}
	if user.Age < 0 || user.Age > 150 {
		errors = append(errors, "age must be between 0 and 150")
	}
	if user.Email == "" {
		errors = append(errors, "email is required")
	} else if !strings.Contains(user.Email, "@") {
		errors = append(errors, "email must contain @")
	}

	return errors
}

// Table-driven test with complex structures
func TestValidateUser(t *testing.T) {
	tests := []struct {
		name           string
		user           User
		expectedErrors []string
	}{
		{
			name: "valid user",
			user: User{
				Name:  "John Doe",
				Age:   30,
				Email: "john@example.com",
			},
			expectedErrors: []string{},
		},
		{
			name: "missing name",
			user: User{
				Name:  "",
				Age:   30,
				Email: "john@example.com",
			},
			expectedErrors: []string{"name is required"},
		},
		{
			name: "invalid age",
			user: User{
				Name:  "John Doe",
				Age:   -5,
				Email: "john@example.com",
			},
			expectedErrors: []string{"age must be between 0 and 150"},
		},
		{
			name: "invalid email",
			user: User{
				Name:  "John Doe",
				Age:   30,
				Email: "invalid-email",
			},
			expectedErrors: []string{"email must contain @"},
		},
		{
			name: "multiple errors",
			user: User{
				Name:  "",
				Age:   200,
				Email: "",
			},
			expectedErrors: []string{
				"name is required",
				"age must be between 0 and 150",
				"email is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validateUser(tt.user)

			// Check if we have the expected number of errors
			if len(errors) != len(tt.expectedErrors) {
				t.Errorf("validateUser() returned %d errors, want %d", len(errors), len(tt.expectedErrors))
				t.Errorf("got errors: %v", errors)
				t.Errorf("want errors: %v", tt.expectedErrors)
				return
			}

			// Check if all expected errors are present
			for _, expectedError := range tt.expectedErrors {
				found := false
				for _, actualError := range errors {
					if actualError == expectedError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("validateUser() missing expected error: %q", expectedError)
				}
			}
		})
	}
}

// =============================================================================
// 4. BENCHMARK TABLE-DRIVEN TESTS
// =============================================================================

// Function to benchmark: String concatenation methods
func concatenateWithPlus(strs []string) string {
	result := ""
	for _, s := range strs {
		result += s
	}
	return result
}

func concatenateWithBuilder(strs []string) string {
	var builder strings.Builder
	for _, s := range strs {
		builder.WriteString(s)
	}
	return builder.String()
}

// Table-driven benchmark
func BenchmarkConcatenation(b *testing.B) {
	benchmarks := []struct {
		name string
		fn   func([]string) string
	}{
		{
			name: "string_plus",
			fn:   concatenateWithPlus,
		},
		{
			name: "strings_builder",
			fn:   concatenateWithBuilder,
		},
	}

	testData := []string{"hello", "world", "golang", "testing", "benchmark"}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bm.fn(testData)
			}
		})
	}
}

// =============================================================================
// 5. HELPER FUNCTIONS FOR TABLE-DRIVEN TESTS
// =============================================================================

// Helper function to compare slices
func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Function to test: Split string by delimiter
func splitString(s, delimiter string) []string {
	if delimiter == "" {
		return []string{s}
	}
	return strings.Split(s, delimiter)
}

// Table-driven test using helper functions
func TestSplitString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		delimiter string
		expected  []string
	}{
		{
			name:      "simple split",
			input:     "a,b,c",
			delimiter: ",",
			expected:  []string{"a", "b", "c"},
		},
		{
			name:      "empty delimiter",
			input:     "hello",
			delimiter: "",
			expected:  []string{"hello"},
		},
		{
			name:      "no delimiter found",
			input:     "hello world",
			delimiter: ",",
			expected:  []string{"hello world"},
		},
		{
			name:      "multiple delimiters",
			input:     "a,,b,,c",
			delimiter: ",",
			expected:  []string{"a", "", "b", "", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitString(tt.input, tt.delimiter)
			if !compareStringSlices(result, tt.expected) {
				t.Errorf("splitString(%q, %q) = %v, want %v", tt.input, tt.delimiter, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// 6. DEMO FUNCTION TO SHOW TABLE-DRIVEN TESTS
// =============================================================================

// Demo function to show how table-driven tests work
func tableDrivenTestsDemo() {
	fmt.Println("🧪 TABLE-DRIVEN TESTS TUTORIAL")
	fmt.Println("===============================")

	fmt.Println("\n1. Basic Table-Driven Test Structure:")
	fmt.Println("   - Define test cases in a slice of structs")
	fmt.Println("   - Each struct contains input, expected output, and test name")
	fmt.Println("   - Use t.Run() to create sub-tests")
	fmt.Println("   - Iterate through test cases with a for loop")

	fmt.Println("\n2. Benefits of Table-Driven Tests:")
	fmt.Println("   ✅ Easy to add new test cases")
	fmt.Println("   ✅ Clear and readable test structure")
	fmt.Println("   ✅ Comprehensive coverage with minimal code")
	fmt.Println("   ✅ Easy to maintain and modify")
	fmt.Println("   ✅ Better error messages with test case names")

	fmt.Println("\n3. Common Patterns:")
	fmt.Println("   - Basic: input → expected output")
	fmt.Println("   - Error cases: input → expected error")
	fmt.Println("   - Complex validation: input → expected validation results")
	fmt.Println("   - Benchmarks: different implementations")

	fmt.Println("\n4. Best Practices:")
	fmt.Println("   - Use descriptive test case names")
	fmt.Println("   - Include edge cases and error conditions")
	fmt.Println("   - Use helper functions for complex comparisons")
	fmt.Println("   - Keep test data close to the test function")
	fmt.Println("   - Use t.Run() for better test organization")

	fmt.Println("\n✅ Table-driven tests make your tests more maintainable and comprehensive!")
}
