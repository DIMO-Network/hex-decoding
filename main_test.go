package main

import (
	"math"
	"testing"
)

func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestExtractAndDecodeWithFormula(t *testing.T) {
	const tolerance = 1e-5

	var tests = []struct {
		hexData  string
		formula  string
		expected float64
		unit     string
		err      string
	}{
		{"7e80641a6000e2a2da", "48|20@0+ (0.1,0) [1|429496730] \"km\"", 92830.1, "km", ""},
		{"invalidhex", "31|32@0+ (0.1,0) [1|429496730] \"km\"", 0, "", "encoding/hex: invalid byte: U+0069 'i'"},
		{"7e80641a6000e2a2da", "31|32@0+ (0.1,0) [1|2] \"km\"", 0, "", "decoded value out of range: 80533537.70 (expected range 1.00 to 2.00)"},
	}

	for _, test := range tests {
		decoded, unit, err := ExtractAndDecodeWithFormula(test.hexData, test.formula)

		if err != nil {
			if err.Error() != test.err {
				t.Errorf("Expected error \"%v\" but got \"%v\"", test.err, err)
			}
		} else if test.err != "" {
			t.Errorf("Expected error \"%v\" but got nil", test.err)
		} else if !almostEqual(decoded, test.expected, tolerance) || unit != test.unit {
			t.Errorf("ExtractAndDecodeWithFormula(%q, %q): expected %v %v, actual %v %v",
				test.hexData, test.formula, test.expected, test.unit, decoded, unit)
		}
	}
}
