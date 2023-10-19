package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
)

func ExtractAndDecodeWithFormula(hexData string, formula string) (float64, string, error) {
	// Define a regex for the formula
	re := regexp.MustCompile(`(\d+)\|(\d+)@(\d+)\+ \(([^,]+),([^)]+)\) \[([^|]+)\|([^]]+)] "([^"]+)"`)
	matches := re.FindStringSubmatch(formula)

	if len(matches) != 9 {
		return 0, "", fmt.Errorf("invalid formula format: %s", formula)
	}

	startBit, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, "", err
	}
	lengthBits, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, "", err
	}
	endianness, err := strconv.Atoi(matches[3])
	if err != nil {
		return 0, "", err
	}
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

	// Decode hex
	hexBytes, err := hex.DecodeString(hexData)
	if err != nil {
		return 0, "", err
	}

	// Calculate start byte, end byte
	startByte := startBit / 8
	lengthBytes := (lengthBits + 7) / 8 //
	endByte := startByte + lengthBytes
	if endByte > len(hexBytes) {
		return 0, "", fmt.Errorf("data out of range: data length is %d, but endByte is %d", len(hexBytes), endByte)
	}

	// Extract data and convert to uint32 with endianness
	var dataUint32 uint32
	dataBytes := hexBytes[startByte:endByte]
	if endianness == 1 {
		dataUint32 = binary.BigEndian.Uint32(append(make([]byte, 4-len(dataBytes)), dataBytes...)) // big-endian
	} else {
		dataUint32 = binary.LittleEndian.Uint32(append(dataBytes, make([]byte, 4-len(dataBytes))...)) // little-endian
	}

	// Create mask and shift bits
	mask := (1<<lengthBits - 1) << (startBit % 8)
	value := (dataUint32 & uint32(mask)) >> (startBit % 8)

	// Apply formula components
	decodedValue := (float64(value) * scaleFactor) + offsetAdjustment

	// Check if the decoded value is within specified range
	if decodedValue < minValue || decodedValue > maxValue {
		return 0, "", fmt.Errorf("decoded value out of range: %.2f (expected range %.2f to %.2f)", decodedValue, minValue, maxValue)
	}

	return decodedValue, unit, nil
}

func main() {
	// Example odometer
	hexData := "7e80641a6000e2a2"
	formula := "31|32@0+ (0.1,0) [1|429496730] \"km\""

	decodedValue, unit, err := ExtractAndDecodeWithFormula(hexData, formula)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Decoded Value: %.2f %s\n", decodedValue, unit)
	}
}
