package main

import (
	"fmt"
	"log"
	"net/http"
)

// Импортируем обработчики из handlers package
import (
	"dhcp-api-server/handlers"
)

func main() {
	// Настройка маршрутов
	http.HandleFunc("/api/interfaces", handlers.GetInterfaces)
	http.HandleFunc("/api/dhcp/config", handlers.DHCPConfigHandler)
	http.HandleFunc("/api/dhcp/status", handlers.DHCPStatusHandler)
	http.HandleFunc("/api/dhcp/control", handlers.DHCPControlHandler)

	// Запуск сервера
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}