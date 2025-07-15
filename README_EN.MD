# 🔗 URL Shortener

A simple and high-performance service for shortening URLs with a REST API.

---

## ✨ Features

- Create short links with custom aliases
- Automatic alias generation (if not specified)
- Redirect to original URLs
- Basic Auth protection
- Operation logging
- SQLite storage
- Flexible configuration via YAML

---

## 📌 Requirements

- **Go** 1.21+
- **SQLite3** (built-in driver used)

---

## 🚀 Quick Start

1. **Clone the repository and install dependencies:**
    ```bash
    git clone https://github.com/SergeyGolang/urlshortener
    cd urlshortener
    go mod download
    ```

2. **Copy the config template:**
    ```bash
    cp config.example.yml config/local.yaml
    ```

3. **Set the environment variable:**
    ```bash
    # Linux/macOS:
    export CONFIG_PATH=config/local.yaml
    ```

4. **Start the server:**
    ```bash
    go run cmd/url-shortener/main.go
    ```

---

## 🔧 Example Config (`config.example.yml`)

```yaml
# Environment: local, dev, prod
env: "local"

# Path to the database file
storage_path: "./storage/storage.db"

http_server:
  # Server address and port
  address: "localhost:8082"
  # Request read and response write timeout
  timeout: 4s
  # Idle connection timeout
  idle_timeout: 60s
  # Username for Basic Auth
  user: "test"
  # Password for Basic Auth
  password: "test"

##📋 Usage

Create a short link:

curl -X POST -u test:test \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com", "alias":"example"}' \
  http://localhost:8082/url

Redirect using a short link:

curl -v http://localhost:8082/example

Or open in your browser:

http://localhost:8082/example

## 🛠 Development & Testing

Run tests:

go test -v ./...

## ⚠️ Important!
Do not commit local.yaml to git!