package utils

import (
	"net"
)

// IsValidIP проверяет валидность IP адреса
func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// ValidateIPRange проверяет валидность диапазона IP адресов
func ValidateIPRange(start, end string) bool {
	startIP := net.ParseIP(start)
	endIP := net.ParseIP(end)
	
	if startIP == nil || endIP == nil {
		return false
	}
	
	// Проверяем, что start <= end
	for i := 0; i < len(startIP); i++ {
		if startIP[i] > endIP[i] {
			return false
		}
		if startIP[i] < endIP[i] {
			break
		}
	}
	
	return true
}