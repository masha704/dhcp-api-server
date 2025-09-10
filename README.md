# 🌐 DHCP REST API Server

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![API Documentation](https://img.shields.io/badge/API-Documentation-ff69b4)](docs/API.md)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

Мощный REST API сервер на Go для управления сетевыми интерфейсами Linux и ISC DHCP сервером с комплексным логированием и валидацией.

## ✨ Возможности

- 📊 **Мониторинг сетевых интерфейсов** - Полная информация о системных интерфейсах
- ⚙️ **Управление DHCP сервером** - Полный контроль над конфигурацией ISC DHCP
- 🔄 **Управление службой** - Запуск, остановка и перезагрузка DHCP через API
- 📝 **Логирование запросов** - Детальное логирование всех API запросов
- ✅ **Валидация данных** - Надежная проверка параметров конфигурации
- 🐳 **Docker поддержка** - Готовность к контейнеризации
- 🔐 **Безопасность** - Валидация и обработка ошибок

## 🚀 Быстрый старт

### Предварительные требования
- Go 1.21+
- ISC DHCP Server
- Linux окружение

### Установка

```bash
# Клонирование репозитория
git clone https://github.com/masha704/dhcp-api-server.git
cd dhcp-api-server

# Установка зависимостей
go mod download

# Сборка приложения
go build -o dhcp-api-server

# Запуск сервера (требует sudo для управления DHCP)
sudo ./dhcp-api-server
```
📡 API Эндпоинты

GET	/api/interfaces	Информация о сетевых интерфейсах
GET	/api/dhcp/config	Текущая конфигурация DHCP
POST	/api/dhcp/config	Обновление конфигурации DHCP
GET	/api/dhcp/status	Статус службы DHCP
POST	/api/dhcp/control	Управление службой DHCP

💡 Примеры использования
Получение информации о интерфейсах
```bash
curl http://localhost:8080/api/interfaces | jq
```
Получение статуса DHCP
```bash
curl http://localhost:8080/api/dhcp/status
```
Конфигурация DHCP сервера
```bash
curl -X POST http://localhost:8080/api/dhcp/config \
  -H "Content-Type: application/json" \
  -d '{
    "interface": "eth0",
    "default_lease_time": 600,
    "max_lease_time": 7200,
    "range_start": "192.168.1.100",
    "range_end": "192.168.1.200",
    "subnet": "192.168.1.0",
    "netmask": "255.255.255.0",
    "dns_servers": "8.8.8.8, 8.8.4.4",
    "gateway": "192.168.1.1"
  }'
```
Управление службой
```bash
# Перезагрузка службы
curl -X POST http://localhost:8080/api/dhcp/control \
  -H "Content-Type: application/json" \
  -d '{"action": "restart"}'
```
🏗️ Структура проекта
```text
dhcp-api-server/
├── src/                    # Исходный код
│   ├── main.go            # Основной файл
│   ├── handlers/          # Обработчики API
│   └── utils/             # Вспомогательные утилиты
├── docs/                  # Документация
│   ├── API.md            # API документация
│   └── SETUP.md          # Инструкции по установке
├── scripts/               # Скрипты развертывания
├── tests/                 # Тесты
├── go.mod                # Go модули
├── go.sum                # Зависимости
├── Dockerfile            # Конфигурация Docker
├── docker-compose.yml    # Docker Compose
└── README.md             # Этот файл

```
🛠️ Технологии
Go - Основной язык программирования

ISC DHCP Server - DHCP сервер

REST API - Архитектурный стиль

Linux Networking - Работа с сетевыми интерфейсами

📊 Ответ API
Успешный ответ
```json
{
  "status": "success",
  "data": {...},
  "message": "Operation completed successfully"
}
```
Ошибка
```json
{
  "status": "error",
  "error": "Invalid configuration parameters",
  "details": "Range start must be valid IP address"
}
```

📄 Лицензия
Этот проект лицензирован под MIT License - смотрите файл LICENSE для деталей.

