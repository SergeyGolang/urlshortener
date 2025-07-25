# 🔗 Сокращатель URL

Простой и производительный сервис для сокращения ссылок с REST API.

---

## ✨ Возможности

- Создание коротких ссылок с кастомными алиасами
- Автоматическая генерация алиасов (если не указан)
- Редирект на оригинальные URL
- Защита через Basic Auth
- Логирование операций
- Хранение в SQLite
- Гибкая конфигурация через YAML

---

## 📌 Требования

- **Go** 1.21+
- **SQLite3** (используется встроенный драйвер)

---

## 🚀 Быстрый старт

1. Клонируйте репозиторий и установите зависимости:
    ```bash
    git clone https://github.com/SergeyGolang/urlshortener
    cd urlshortener
    go mod download
    ```

2. Создайте папку для конфигов и скопируйте шаблон:
    ```bash
    cp config.example.yml config/local.yaml
    ```

3. Установите переменную окружения:
    ```bash
    # Linux/macOS:
    export CONFIG_PATH=config/local.yaml
    ```

4. Запустите сервер:
    ```bash
    go run cmd/url-shortener/main.go
    ```

---

## 🔧 Пример конфига (`config.example.yml`)

```yaml
# Режим работы: local, dev, prod
env: "local"

# Путь к файлу базы данных
storage_path: "./storage/storage.db"

http_server:
  # Адрес и порт сервера
  address: "localhost:8082"
  # Время на чтение запроса и отправку ответа
  timeout: 4s
  # Время открытого соединения с клиентом
  idle_timeout: 60s
  # Имя пользователя для Basic Auth
  user: "test"
  # Пароль для Basic Auth
  password: "test"
```

## 📋 Использование 

Cоздание короткой ссылки:

curl -X POST -u test:test \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com", "alias":"example"}' \
  http://localhost:8082/url


Переход по короткой ссылке:

curl -v http://localhost:8082/example

Или откройте в браузере: 

http://localhost:8082/example


## 🛠 Разработка и тестирование

Запуск тестов: 

go test -v ./...


## ⚠️ Важно!

Не коммитьте local.yaml в git!

## Разработчик

https://github.com/SergeyGolang

