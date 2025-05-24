# E-Wallet API

Sistem backend RESTful API untuk simulasi e-wallet sederhana dengan fitur register, login, top-up, pembayaran, transfer, dan laporan transaksi.

## Tech Stack

- **Golang** dengan Framework **Gin**
- **PostgreSQL** sebagai database
- **GORM** sebagai ORM
- **JWT** untuk autentikasi
- **RabbitMQ** sebagai message broker
- **Docker Compose** untuk deployment
- **Air** untuk live reload development

## Fitur

### ✅ API yang Sudah Selesai

1. **API Register** - `POST /register`
   - Mendaftarkan pengguna baru dengan UUID
   - Validasi nomor telepon unik
   - Enkripsi PIN menggunakan bcrypt

2. **API Login** - `POST /login`
   - Login dengan nomor telepon dan PIN
   - Generate JWT access token dan refresh token

3. **API Top Up** - `POST /topup`
   - Menambahkan saldo ke akun user
   - Memerlukan autentikasi JWT
   - Mencatat transaksi

4. **API Payment** - `POST /payments`
   - Melakukan pembayaran dengan saldo
   - Validasi saldo mencukupi
   - Mencatat transaksi

5. **API Transfer** - `POST /transfers`
   - Transfer saldo antar user secara asynchronous
   - Menggunakan RabbitMQ queue system
   - Background processing

6. **API Report Transactions** - `GET /transactions`
   - Menampilkan histori transaksi user
   - Terurut berdasarkan tanggal terbaru

7. **API Update Profile** - `PUT /update-profile`
   - Mengubah profil user (kecuali nomor telepon dan PIN)
   - Validasi input

### 🔧 Teknologi Tambahan

- **Air** - Live reload untuk development
- **bcrypt** - Enkripsi PIN
- **UUID** - Generate unique identifier
- **CORS** - Cross-Origin Resource Sharing
- **Environment Variables** - Konfigurasi aplikasi
- **Database Transactions** - Konsistensi data
- **Background Workers** - Asynchronous processing


## 📦 Installation & Dependencies

## Prerequisites

- **Docker** dan **Docker Compose**
- **Go 1.21+** (untuk development)
- **Air** (untuk live reload development)
- **PostgreSQL** (untuk local development)
- **RabbitMQ** (untuk local development)

## Quick Start

### 1. Clone Repository

```bash
git clone <repository-url>
cd e-wallet-backend
```

### 2. Setup Environment

```bash
# Copy environment file
cp .env.example .env

# Edit konfigurasi jika diperlukan
nano .env
```

### 3. Start Application
```bash
# Quick start dengan Docker (Recommended)
make docker-up

# Atau manual
docker-compose up -d --build
Ï
```
### 4. Verify Installation
```bash
# Cek status services
make status

# Test API health check
curl http://localhost:8080/health
Ï
```

## Cara Menjalankan Aplikasi

### Full Docker Environment (Recommended for Production)
```bash
# Start semua services
make docker-up
# atau
docker-compose up -d --build

# Lihat logs real-time
make docker-logs
# atau
docker-compose logs -f

# Lihat status services
make status
# atau
docker-compose ps

# Stop services
make docker-down
# atau
docker-compose down
Ï
```

### Access Point
```bash
API: http://localhost:8080
RabbitMQ Management: http://localhost:15672 (guest/guest)
PostgreSQL: localhost:5433
```

### Full Local Development
```bash
# Prerequisites: Install PostgreSQL dan RabbitMQ di local

# MacOS dengan Homebrew
brew install postgresql rabbitmq
brew services start postgresql
brew services start rabbitmq

# Ubuntu/Debian
sudo apt-get install postgresql rabbitmq-server
sudo systemctl start postgresql
sudo systemctl start rabbitmq-server

# Windows (dengan Chocolatey)
choco install postgresql rabbitmq

# Jalankan aplikasi
make dev-local

# atau
air
```
