package config

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"interface":        "eth0",
		"default_lease_time": 600,
		"max_lease_time":     7200,
		"range_start":       "192.168.1.100",
		"range_end":         "192.168.1.200",
		"subnet":           "192.168.1.0",
		"netmask":          "255.255.255.0",
		"dns_servers":      "8.8.8.8, 8.8.4.4",
		"gateway":          "192.168.1.1",
	}
}