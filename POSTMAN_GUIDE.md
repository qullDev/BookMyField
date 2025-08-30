# BookMyField API - Postman Collection Guide

## 📁 File yang Sudah Dibuat:

- `postman_collection.json` - Collection lengkap semua API endpoints
- `postman_environment.json` - Environment variables untuk testing

## 🚀 Cara Import ke Postman:

### 1. Import Collection:

1. Buka Postman
2. Click **Import** button
3. Pilih file `postman_collection.json`
4. Collection "BookMyField API Collection" akan muncul

### 2. Import Environment:

1. Click **Environment** tab di Postman
2. Click **Import**
3. Pilih file `postman_environment.json`
4. Pilih environment "BookMyField Environment"

## 🧪 Testing Flow Recommendations:

### A. Setup Awal:

1. **Health Check** - Pastikan API berjalan
2. **Get All Fields** - Lihat field yang tersedia, copy ID field untuk testing

### B. Authentication Testing:

1. **Login Admin** - Login sebagai admin (admin@admin.com / password123)
2. **Login User** - Login sebagai user (user@user.com / password123)
3. **Register User** - Daftar user baru
4. **Refresh Token** - Test refresh token functionality
5. **Logout** - Test logout

### C. Field Management (Admin Only):

1. **[ADMIN] Create Field** - Buat field baru
2. **[ADMIN] Update Field** - Update field yang sudah ada
3. **Get Field by ID** - Lihat detail field
4. **[ADMIN] Delete Field** - Hapus field (hati-hati!)

### D. Booking Flow (User):

1. **Create Booking** - Buat booking baru
2. **Get My Bookings** - Lihat booking saya
3. **[ADMIN] Get All Bookings** - Admin lihat semua booking
4. **Cancel Booking** - Cancel booking yang sudah dibuat

### E. Payment Flow:

1. **Create Checkout Session** - Buat sesi pembayaran
2. **Stripe Webhook** - Simulasi webhook dari Stripe

## 🔧 Environment Variables:

Variables ini akan otomatis di-set oleh script:

- `access_token` - JWT token user
- `admin_access_token` - JWT token admin
- `refresh_token` - Refresh token user
- `admin_refresh_token` - Refresh token admin
- `field_id` - ID field untuk testing
- `booking_id` - ID booking untuk testing
- `payment_id` - ID payment untuk testing
- `checkout_url` - URL checkout Stripe

## 📝 Default Login Credentials:

**Admin Account:**

- Email: admin@admin.com
- Password: password123

**User Account:**

- Email: user@user.com
- Password: password123

## ⚠️ Important Notes:

1. **Pastikan server berjalan** di `http://localhost:8080`
2. **Login dulu** sebelum test endpoint yang butuh authentication
3. **Token otomatis tersimpan** setelah login berhasil
4. **Admin endpoints** butuh admin token
5. **Field ID** bisa didapat dari "Get All Fields"
6. **Booking Date** format: YYYY-MM-DD
7. **Time** format: HH:MM (24 hour)

## 🎯 Quick Test Scenario:

1. Health Check ✅
2. Login Admin ✅
3. Get All Fields ✅ (copy field_id)
4. Login User ✅
5. Create Booking ✅
6. Create Checkout Session ✅
7. Get My Bookings ✅
8. Cancel Booking ✅

## 🔍 Response Codes:

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict (booking collision)
- `500` - Internal Server Error

Happy Testing! 🚀
