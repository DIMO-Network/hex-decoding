package main

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ExtractAndDecodeWithFormula(hexData string, formula string) (float64, string, error) {

	// Define a regex for the formula
	re := regexp.MustCompile(`(\d+)\|(.+)`)

	matches := re.FindStringSubmatch(formula)

	if len(matches) != 3 {
		return 0, "", fmt.Errorf("invalid formula format: %s", formula)
	}

	offset, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, "", err
	}

	formulaParts := strings.Fields(matches[2])

	//formulaParts[0] contains length, formulaParts[1] contains factor parts, formulaParts[2] contains range parts, formulaParts[3] contains the units
	if len(formulaParts) != 4 {
		return 0, "", fmt.Errorf("invalid formula format: %s", formula)
	}

	// Extract length, scale factor, and offset adjustment values
	lengthPart := extractNumericPart(formulaParts[0])
	length, err := strconv.Atoi(lengthPart)
	if err != nil {
		return 0, "", err
	}

	// Parse scaling factor and offset adjustment
	factorPart := strings.Trim(formulaParts[1], " ")

	// Extract scaleFactor and offsetAdjustment from factorPart
	scaleFactorIndex := strings.Index(factorPart, "(")
	offsetAdjustmentIndex := strings.Index(factorPart, ",")
	if scaleFactorIndex == -1 || offsetAdjustmentIndex == -1 {
		return 0, "", fmt.Errorf("invalid factor format: %s", factorPart)
	}

	scaleFactorStr := factorPart[scaleFactorIndex+1 : offsetAdjustmentIndex]
	offsetAdjustmentStr := factorPart[offsetAdjustmentIndex+1 : len(factorPart)-1]

	// Convert scaleFactorStr and offsetAdjustmentStr to float64
	scaleFactor, err := strconv.ParseFloat(strings.Trim(scaleFactorStr, " "), 64)
	if err != nil {
		return 0, "", err
	}

	offsetAdjustment, err := strconv.ParseFloat(strings.Trim(offsetAdjustmentStr, " "), 64)
	if err != nil {
		return 0, "", err
	}

	rangePart := strings.Trim(formulaParts[2], "")
	minValueIndex := strings.Index(rangePart, "[")
	maxValueIndex := strings.Index(rangePart, "|")
	if minValueIndex == -1 || maxValueIndex == -1 {
		return 0, "", fmt.Errorf("invalid factor format: %s", rangePart)
	}
	minValueStr := rangePart[minValueIndex+1 : maxValueIndex]
	maxValueStr := rangePart[maxValueIndex+1 : len(rangePart)-1]

	// Convert minValueStr and maxValueStr to float64
	minValue, err := strconv.ParseFloat(strings.Trim(minValueStr, " "), 64)
	if err != nil {
		return 0, "", err
	}

	maxValue, err := strconv.ParseFloat(strings.Trim(maxValueStr, " "), 64)
	if err != nil {
		return 0, "", err
	}

	unit := formulaParts[3]

	// Decode hex
	hexBytes, err := hex.DecodeString(hexData)
	if err != nil {
		return 0, "", err
	}

	// Check if the byte slice is long enough
	if offset+length > len(hexBytes) {
		return 0, "", fmt.Errorf("formula for offset %d and length %d is out of bounds", offset, length)
	}

	// Extract relevant bytes
	hexValue := hexBytes[offset : offset+length]
	intValue, err := strconv.ParseInt(string(hexValue), 16, 64)
	if err != nil {
		return 0, "", err
	}

	// Apply formula components
	decodedValue := (float64(intValue) + offsetAdjustment) * scaleFactor

	// Check if the decoded value is within specified range
	if decodedValue < minValue || decodedValue > maxValue {
		return 0, "", fmt.Errorf("decoded value is out of range: %.2f", decodedValue)
	}

	return decodedValue, unit, nil
}
func extractNumericPart(s string) string {
	re := regexp.MustCompile(`(\d+)`)
	match := re.FindString(s)
	return match
}

func main() {
	// Example
	hexData := "0123"
	formula := "31|16@0+ (1,0) [0|65535] \"km\" "

	decodedValue, unit, err := ExtractAndDecodeWithFormula(hexData, formula)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Decoded Value: %.2f %s\n", decodedValue, unit)
	}
}
