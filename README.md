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

### ✅ Halaman yang Sudah Selesai

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

## Instalasi dan Menjalankan Project

### Prerequisites

- Docker dan Docker Compose
- Go 1.21+ (untuk development)
- Air (untuk live reload development)

### 1. Clone Repository

```bash
git clone <repository-url>
cd e-wallet-api
```

### 2. Setup Environment
```bash
# Copy environment file
cp .env.example .env

# Edit konfigurasi jika diperlukan
nano .env
```

### 3. Jalankan dengan Docker Compose
```bash
# Build dan jalankan semua services
docker-compose up --build

# Atau jalankan di background
docker-compose up --build -d
```

### 4. Verifikasi Installation
```bash
# Cek status container
docker-compose ps

# Test API health check
curl http://localhost:8080/health
```
### 5. Development Mode(Optional)
```bash
# Install Air untuk live reload
go install github.com/cosmtrek/air@latest

# Jalankan database services saja
docker-compose up postgres rabbitmq -d

# Jalankan aplikasi dengan live reload
air
```
### 6. Stop Services
```bash
# Stop semua services
docker-compose down

# Stop dan hapus data (HATI-HATI)
docker-compose down -v
```
