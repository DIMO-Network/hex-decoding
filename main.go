package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ExtractAndDecodeWithFormula extracts the data following the PID and applies the formula for decoding
func ExtractAndDecodeWithFormula(hexData, pid, formula string) (float64, string, error) {
	// Parse formula
	re := regexp.MustCompile(`(\d+)\|(\d+)@(\d+)\+ \(([^,]+),([^)]+)\) \[([^|]+)\|([^]]+)] "([^"]+)"`)
	matches := re.FindStringSubmatch(formula)

	if len(matches) != 9 {
		return 0, "", fmt.Errorf("invalid formula format: %s", formula)
	}

	lengthBits, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, "", err
	}
	numBytes := lengthBits / 8

	// Find the index of PID in the hex string, strings.Index returns first occurrence
	pidIndex := strings.Index(hexData, pid)
	if pidIndex == -1 {
		return 0, "", errors.New("PID not found")
	}

	// Calculate the start index of the data after the PID
	dataStartIndex := pidIndex + len(pid)

	if dataStartIndex+numBytes*2 > len(hexData) {
		return 0, "", errors.New("not enough data")
	}

	// Extract the relevant portion of the hex data
	valueHex := hexData[dataStartIndex : dataStartIndex+numBytes*2]

	value, err := strconv.ParseUint(valueHex, 16, 64)
	if err != nil {
		return 0, "", err
	}

	// Parse the formula parameters
	scaleFactor, err := strconv.ParseFloat(matches[4], 64)
	if err != nil {
		return 0, "", err
	}
	offsetAdjustment, err := strconv.ParseFloat(matches[5], 64)
	if err != nil {
		return 0, "", err
	}
	minValue, err := strconv.ParseFloat(matches[6], 64)
	if err != nil {
		return 0, "", err
	}
	maxValue, err := strconv.ParseFloat(matches[7], 64)
	if err != nil {
		return 0, "", err
	}
	unit := matches[8]

	// Apply the formula
	decodedValue := float64(value)*scaleFactor + offsetAdjustment

	// Validate the range
	if decodedValue < minValue || decodedValue > maxValue {
		return 0, "", fmt.Errorf("decoded value out of range: %.2f (expected range %.2f to %.2f)", decodedValue, minValue, maxValue)
	}

	return decodedValue, unit, nil
}

func main() {
	hexData := "7e80641a60008b24200"
	pid := "a6"
	formula := "31|32@0+ (0.1,0) [1|4294967295] \"km\""

	decodedValue, unit, err := ExtractAndDecodeWithFormula(hexData, pid, formula)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Decoded Value: %.2f %s\n", decodedValue, unit)
}
