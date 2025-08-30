# 📚 BookMyField API - Postman Testing Guide (Production)

## 🚀 Quick Setup

### 1. Import ke Postman

1. **Import Collection:**

   - Buka Postman
   - Klik **Import**
   - Pilih file `postman_collection_production.json`

2. **Import Environment:**
   - Klik **Import** lagi
   - Pilih file `postman_environment_production.json`
   - Pilih environment **"BookMyField Production Environment"** di dropdown

### 2. Default Login Credentials

**🔑 User Account:**

- Email: `user@user.com`
- Password: `password123`

**👑 Admin Account:**

- Email: `admin@admin.com`
- Password: `password123`

## 🧪 Testing Flow (Recommended Order)

### A. 🏁 Setup Awal

1. **Health Check** - Pastikan API berjalan
2. **Get All Fields** - Lihat field yang tersedia

### B. 🔐 Authentication Testing

1. **Login User** - Login sebagai user biasa
2. **Login Admin** - Login sebagai admin
3. **Register User** - Daftar user baru (optional)
4. **Refresh Token** - Test refresh token functionality
5. **Logout** - Test logout

### C. 🏟️ Field Management

1. **Get All Fields** - Lihat semua lapangan (BUTUH AUTH! ✅)
2. **Get Field by ID** - Lihat detail lapangan (BUTUH AUTH! ✅)
3. **[ADMIN] Create Field** - Buat lapangan baru
4. **[ADMIN] Update Field** - Update lapangan
5. **[ADMIN] Delete Field** - Hapus lapangan

### D. 📅 Booking Management

1. **Create Booking** - Buat booking baru
2. **Get My Bookings** - Lihat booking saya
3. **[ADMIN] Get All Bookings** - Lihat semua booking (admin only)
4. **Cancel Booking** - Batalkan booking

### E. 💳 Payment Testing

1. **Create Checkout Session** - Buat session pembayaran
2. **Stripe Webhook** - Test webhook handler

## 🔧 Environment Variables

Variables yang akan otomatis di-set oleh script:

- `base_url` - URL production server
- `access_token` - JWT token user (auto-set dari login)
- `refresh_token` - Refresh token user (auto-set dari login)
- `admin_access_token` - JWT token admin (auto-set dari admin login)
- `admin_refresh_token` - Refresh token admin (auto-set dari admin login)
- `field_id` - ID field untuk testing (auto-set dari get fields)
- `new_field_id` - ID field baru (auto-set dari create field)
- `booking_id` - ID booking untuk testing (auto-set dari create booking)
- `checkout_url` - URL checkout Stripe (auto-set dari create checkout)

## ⚠️ Important Notes

### 🔐 Authentication Requirements

**SEMUA endpoint fields MEMERLUKAN AUTHENTICATION:**

- ✅ `GET /api/v1/fields/` - **BUTUH Bearer Token**
- ✅ `GET /api/v1/fields/:id` - **BUTUH Bearer Token**

**Endpoint yang TIDAK butuh auth:**

- `GET /api/v1/health` - Health check
- `POST /api/v1/auth/register` - Register user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh token

### 📋 Testing Checklist

1. **Pastikan environment dipilih** - "BookMyField Production Environment"
2. **Login dulu** sebelum test endpoint yang butuh authentication
3. **Token otomatis tersimpan** setelah login berhasil
4. **Admin endpoints** butuh admin token
5. **Field ID** akan otomatis tersimpan dari "Get All Fields"
6. **Date format** untuk booking: `YYYY-MM-DDTHH:MM:SSZ` (ISO 8601)

### 🎯 Quick Test Scenario

**Urutan testing cepat:**

1. **Health Check** ✅
2. **Login User** ✅ (auto-save token)
3. **Get All Fields** ✅ (auto-save field_id)
4. **Get Field by ID** ✅ (gunakan saved field_id)
5. **Create Booking** ✅ (auto-save booking_id)
6. **Get My Bookings** ✅
7. **Create Checkout Session** ✅
8. **Cancel Booking** ✅

### ❌ Common Issues & Solutions

**1. "Authorization header missing"**

- Pastikan sudah login dan token tersimpan
- Check environment variables `access_token`

**2. "Invalid or expired token"**

- Token expired, login ulang
- Atau gunakan refresh token

**3. "Field not found"**

- Pastikan `field_id` valid di environment variables
- Jalankan "Get All Fields" dulu

**4. Format datetime booking salah**

- Gunakan format: `2025-09-15T10:00:00Z`
- Jangan gunakan format: `10:00` (akan error!)

## 🌐 Production URL

Base URL: `https://bookmyfield-production.up.railway.app`

## 📖 API Documentation

Swagger UI: `https://bookmyfield-production.up.railway.app/swagger/index.html`

---

**Happy Testing! 🚀**
