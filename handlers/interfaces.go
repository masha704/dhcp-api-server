package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type InterfaceInfo struct {
	Name         string   `json:"name"`
	MTU          int      `json:"mtu"`
	HardwareAddr string   `json:"hardware_addr"`
	IPAddresses  []string `json:"ip_addresses"`
	Gateway      string   `json:"gateway,omitempty"`
	RXBytes      uint64   `json:"rx_bytes"`
	TXBytes      uint64   `json:"tx_bytes"`
	RXPackets    uint64   `json:"rx_packets"`
	TXPackets    uint64   `json:"tx_packets"`
	RXErrors     uint64   `json:"rx_errors"`
	TXErrors     uint64   `json:"tx_errors"`
}

type InterfaceStats struct {
	RXBytes   uint64
	TXBytes   uint64
	RXPackets uint64
	TXPackets uint64
	RXErrors  uint64
	TXErrors  uint64
}

func GetInterfaces(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var results []InterfaceInfo

	for _, iface := range interfaces {
		// Пропускаем loopback
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		info := InterfaceInfo{
			Name:         iface.Name,
			MTU:          iface.MTU,
			HardwareAddr: iface.HardwareAddr.String(),
		}

		// Получаем IP адреса
		addrs, err := iface.Addrs()
		if err == nil {
			for _, addr := range addrs {
				info.IPAddresses = append(info.IPAddresses, addr.String())
			}
		}

		// Получаем статистику
		stats := getInterfaceStats(iface.Name)
		info.RXBytes = stats.RXBytes
		info.TXBytes = stats.TXBytes
		info.RXPackets = stats.RXPackets
		info.TXPackets = stats.TXPackets
		info.RXErrors = stats.RXErrors
		info.TXErrors = stats.TXErrors

		// Получаем шлюз
		info.Gateway = getGateway(iface.Name)

		results = append(results, info)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func getInterfaceStats(ifaceName string) InterfaceStats {
	stats := InterfaceStats{}
	
	// Читаем статистику из /proc/net/dev
	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return stats
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, ifaceName+":") {
			fields := strings.Fields(line)
			if len(fields) >= 17 {
				fmt.Sscanf(fields[1], "%d", &stats.RXBytes)
				fmt.Sscanf(fields[2], "%d", &stats.RXPackets)
				fmt.Sscanf(fields[3], "%d", &stats.RXErrors)
				fmt.Sscanf(fields[9], "%d", &stats.TXBytes)
				fmt.Sscanf(fields[10], "%d", &stats.TXPackets)
				fmt.Sscanf(fields[11], "%d", &stats.TXErrors)
			}
			break
		}
	}
	
	return stats
}

func getGateway(ifaceName string) string {
	// Используем ip route для получения шлюза
	cmd := exec.Command("ip", "route", "show", "default")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "dev "+ifaceName) {
			fields := strings.Fields(line)
			for i, field := range fields {
				if field == "via" && i+1 < len(fields) {
					return fields[i+1]
				}
			}
		}
	}
	
	return ""
}