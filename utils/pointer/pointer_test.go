package pointer

import (
	"testing"
)

func TestPointerAndValue(t *testing.T) {
	// Test Pointer function
	t.Run("Test Pointer", func(t *testing.T) {
		value := 42
		p := Pointer(value)

		if p == nil {
			t.Errorf("Pointer: expected a non-nil pointer, got nil")
		}

		if *p != value {
			t.Errorf("Pointer: expected %d, got %d", value, *p)
		}
	})

	// Test Value function with non-nil pointer
	t.Run("Test Value with non-nil pointer", func(t *testing.T) {
		value := 42
		p := &value
		result := Value(p)

		if result != value {
			t.Errorf("Value: expected %d, got %d", value, result)
		}
	})

	// Test Value function with nil pointer
	t.Run("Test Value with nil pointer", func(t *testing.T) {
		var p *int
		result := Value(p)

		// Since p is nil, the result should be the zero value of int
		if result != 0 {
			t.Errorf("Value: expected 0, got %d", result)
		}
	})

	// Test Value function with nil pointer of a custom struct
	t.Run("Test Value with nil pointer of a custom struct", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		var p *Person
		result := Value(p)

		// Since p is nil, the result should be the zero value of Person
		zeroPerson := Person{}
		if result != zeroPerson {
			t.Errorf("Value: expected %v, got %v", zeroPerson, result)
		}
	})

	// Test Pointer function with a string
	t.Run("Test Pointer with string", func(t *testing.T) {
		str := "Hello, World!"
		p := Pointer(str)

		if p == nil {
			t.Errorf("Pointer: expected a non-nil pointer, got nil")
		}

		if *p != str {
			t.Errorf("Pointer: expected %s, got %s", str, *p)
		}
	})

	// Test Value function with nil pointer of a string
	t.Run("Test Value with nil pointer of a string", func(t *testing.T) {
		var p *string
		result := Value(p)

		// Since p is nil, the result should be an empty string
		if result != "" {
			t.Errorf("Value: expected '', got %s", result)
		}
	})
}
