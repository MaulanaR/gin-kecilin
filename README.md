# AssistX Submission – Gin + MongoDB

Ini adalah project submission **AssistX** yang dibangun menggunakan **Gin Framework** (Golang) dan **MongoDB**.

## Catatan
- Jika environment variable **tidak di-set**, maka:
  - **Port** default yang digunakan adalah `8080`
  - **Database** default yang digunakan adalah `cctv_db`

## Modul
Project ini terdiri dari 3 modul utama:
1. **User (Auth)**  
   Modul untuk autentikasi user.
2. **Contacts**  
   Modul untuk mengelola data kontak.
3. **CCTVs**  
   Modul untuk mengelola data kamera CCTV.

## Relasi
- Modul **Contacts** dan **CCTVs** memiliki relasi **one-to-many**.  
- Implementasi relasi dilakukan dengan **MongoDB `$lookup`**:
  - Satu **Contact** dapat memiliki banyak **CCTV**.
  - Data CCTV di-join berdasarkan `contact_id`.

## Teknologi
- Golang + Gin
- MongoDB
- Validator (go-playground/validator)
- Gin middleware & utils custom

---

## 🚀 Dockerized Setup

Project ini sudah di-*containerize* menggunakan **Docker** dan **docker-compose**.

### 1. Prasyarat
- Sudah terinstall **Docker** dan **docker-compose**.
- Pastikan port `8080` dan `27017` tidak digunakan aplikasi lain.

### 2. Jalankan dengan docker-compose
```bash
docker compose up -d
