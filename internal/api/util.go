package api

import "strings"

func validateAndFormatEthereumAddress(address string) (string, bool) {
	address = strings.ToLower(address)
	// Check if the address starts with "0x"
	if address[:2] != "0x" {
		address = "0x" + address
	}

	// Check if the address has the correct length
	if len(address) != 42 {
		return "", false
	}

	// Check if the remaining characters are valid hexadecimal characters
	for _, char := range address[2:] {
		if (char < '0' || char > '9') && (char < 'a' || char > 'f') {
			return "", false
		}
	}

	return address, true
}
