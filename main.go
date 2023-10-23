package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// convertHexToBinary converts hexadecimal string to binary string
func convertHexToBinary(hexData string) (string, error) {
	bytes, err := hex.DecodeString(hexData)
	if err != nil {
		return "", err
	}

	var binaryData string
	for _, b := range bytes {
		binaryData += fmt.Sprintf("%08b", b)
	}
	return binaryData, nil
}

// extractBits extracts the desired sequence of bits specified by the formula
func extractBits(binaryData string, startBit, length int) (string, error) {
	if startBit+length > len(binaryData) {
		return "", errors.New("out of range")
	}
	return binaryData[startBit : startBit+length], nil
}

// binaryStringToDecimal converts binary string to decimal
func binaryStringToDecimal(binaryString string) (uint64, error) {
	return strconv.ParseUint(binaryString, 2, 64)
}

func ExtractAndDecodeWithFormula(hexData string, formula string) (float64, string, error) {
	// Regular expression to parse the formula
	re := regexp.MustCompile(`(\d+)\|(\d+)@(\d+)\+ \(([^,]+),([^)]+)\) \[([^|]+)\|([^]]+)] "([^"]+)"`)
	matches := re.FindStringSubmatch(formula)

	if len(matches) != 9 {
		return 0, "", fmt.Errorf("invalid formula format: %s", formula)
	}

	startBit, _ := strconv.Atoi(matches[1])
	lengthBits, _ := strconv.Atoi(matches[2])
	scaleFactor, _ := strconv.ParseFloat(matches[4], 64)
	offsetAdjustment, _ := strconv.ParseFloat(matches[5], 64)
	minValue, _ := strconv.ParseFloat(matches[6], 64)
	maxValue, _ := strconv.ParseFloat(matches[7], 64)
	unit := matches[8]

	// Convert hex to binary
	binaryData, err := convertHexToBinary(hexData)
	if err != nil {
		return 0, "", err
	}
	//print(binaryData)
	//confirmed by rapidtables

	// Extract bits
	bits, err := extractBits(binaryData, startBit, lengthBits)
	if err != nil {
		return 0, "", err
	}
	//print(bits)
	//confirmed by rapidtables

	// Convert binary string to decimal
	value, err := binaryStringToDecimal(bits)
	if err != nil {
		return 0, "", err
	}

	// Apply the formula
	decodedValue := float64(value)*scaleFactor + offsetAdjustment

	// Check range
	if decodedValue < minValue || decodedValue > maxValue {
		return 0, "", fmt.Errorf("decoded value out of range: %.2f (expected range %.2f to %.2f)", decodedValue, minValue, maxValue)
	}

	return decodedValue, unit, nil
}

func main() {
	hexData := "7e80641a6000e2a2da"
	formula := "48|20@0+ (0.1,0) [1|429496730] \"km\""
	//formula := "31|32@0+ (0.1,0) [1|429496730] \"km\""

	decodedValue, unit, err := ExtractAndDecodeWithFormula(hexData, formula)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Decoded Value: %.2f %s\n", decodedValue, unit)
	}
}
