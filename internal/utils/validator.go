package utils

import (
	"regexp"
	"unicode"
)

func ValidatePhoneNumber(phone string) bool {
	// Indonesian phone number validation
	re := regexp.MustCompile(`^(\+62|62|0)8[1-9][0-9]{6,9}$`)
	return re.MatchString(phone)
}

func ValidatePIN(pin string) bool {
	if len(pin) != 6 {
		return false
	}
	
	for _, char := range pin {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}

func ValidateName(name string) bool {
	return len(name) > 0 && len(name) <= 50
}

func ValidateAmount(amount int64) bool {
	return amount > 0
}
