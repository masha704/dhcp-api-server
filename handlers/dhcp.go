package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
	
	"dhcp-api-server/logging"
)

type DHCPConfig struct {
	Interface       string `json:"interface"`
	DefaultLeaseTime int    `json:"default_lease_time"`
	MaxLeaseTime    int    `json:"max_lease_time"`
	RangeStart      string `json:"range_start"`
	RangeEnd        string `json:"range_end"`
	Subnet          string `json:"subnet"`
	Netmask         string `json:"netmask"`
	DNSservers      string `json:"dns_servers"`
	Gateway         string `json:"gateway"`
}

func DHCPConfigHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getDHCPConfig(w, r)
	case http.MethodPost:
		updateDHCPConfig(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func DHCPStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем статус службы
	cmd := exec.Command("systemctl", "is-active", "isc-dhcp-server")
	status, _ := cmd.Output()

	response := map[string]string{
		"status": strings.TrimSpace(string(status)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DHCPControlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Action string `json:"action"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var cmd *exec.Cmd
	switch request.Action {
	case "start":
		cmd = exec.Command("systemctl", "start", "isc-dhcp-server")
	case "stop":
		cmd = exec.Command("systemctl", "stop", "isc-dhcp-server")
	case "restart":
		cmd = exec.Command("systemctl", "restart", "isc-dhcp-server")
	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	if err := cmd.Run(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service %sed successfully", request.Action)
}

func getDHCPConfig(w http.ResponseWriter, r *http.Request) {
	// Чтение текущей конфигурации DHCP
	content, err := os.ReadFile("/etc/dhcp/dhcpd.conf")
	if err != nil {
		http.Error(w, "DHCP config file not found", http.StatusInternalServerError)
		return
	}

	// Возвращаем raw конфигурацию
	response := map[string]string{
		"config": string(content),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func updateDHCPConfig(w http.ResponseWriter, r *http.Request) {
	var config DHCPConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Упрощенная валидация
	if config.Interface == "" || config.RangeStart == "" || config.RangeEnd == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Генерация конфигурационного файла
	configContent := generateDHCPConfig(config)
	
	// Запись конфигурации
	if err := os.WriteFile("/etc/dhcp/dhcpd.conf", []byte(configContent), 0644); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Обновление интерфейса в настройках службы
	updateServiceInterface(config.Interface)

	// Перезапуск сервиса
	if err := restartDHCPService(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "DHCP configuration updated successfully")
}

func generateDHCPConfig(config DHCPConfig) string {
	// Используем шаблон для генерации конфигурации
	tmpl := `default-lease-time {{.DefaultLeaseTime}};
max-lease-time {{.MaxLeaseTime}};

subnet {{.Subnet}} netmask {{.Netmask}} {
	range {{.RangeStart}} {{.RangeEnd}};
	option routers {{.Gateway}};
	option domain-name-servers {{.DNSservers}};
}`

	t := template.Must(template.New("dhcpd").Parse(tmpl))
	var result strings.Builder
	t.Execute(&result, config)
	return result.String()
}

func updateServiceInterface(iface string) {
	// Обновляем интерфейс в настройках службы
	serviceConfig := fmt.Sprintf("INTERFACESv4=\"%s\"\nINTERFACESv6=\"\"\n", iface)
	os.WriteFile("/etc/default/isc-dhcp-server", []byte(serviceConfig), 0644)
}

func restartDHCPService() error {
	// Останавливаем службу
	exec.Command("systemctl", "stop", "isc-dhcp-server").Run()
	
	// Проверяем конфигурацию
	cmd := exec.Command("dhcpd", "-t")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("configuration test failed: %v", err)
	}
	
	// Запускаем службу
	return exec.Command("systemctl", "start", "isc-dhcp-server").Run()
}

// Middleware для логирования
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(logger *logging.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		wrapped := &responseWriter{w, http.StatusOK}
		next(wrapped, r)
		
		duration := time.Since(start)
		logger.LogRequest(r.Method, r.URL.Path, r.RemoteAddr, wrapped.status, duration)
	}
}