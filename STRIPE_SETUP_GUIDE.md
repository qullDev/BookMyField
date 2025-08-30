# 🚨 STRIPE SETUP GUIDE - Mengatasi Error 401

## ❌ **Error yang Anda alami:**

```
"You did not provide an API key. You need to provide your API key in the Authorization header, using Bearer auth"
```

## ✅ **Solusi:**

### 1. **Dapatkan Stripe Test Keys**

**🌍 PENTING: Indonesia belum supported di Stripe secara langsung!**

**Opsi untuk Developer Indonesia:**

1. **Pilih Singapore** sebagai country (paling dekat dengan Indonesia)
2. **Atau pilih United States** untuk testing/development
3. Untuk production real, pertimbangkan payment gateway lokal seperti:
   - Midtrans
   - Xendit
   - DOKU
   - Dana/OVO/GoPay

**Setup Stripe dengan Singapore/US:**

1. Buka [Stripe Dashboard](https://dashboard.stripe.com/test/apikeys)
2. Saat registrasi, pilih **Singapore** atau **United States**
3. Login ke akun Stripe Anda
4. Pastikan Anda dalam **Test Mode** (toggle di kiri atas)
5. Copy **Secret Key** yang dimulai dengan `sk_test_...`

### 2. **Update Environment Variables**

Edit file `.env` di root project Anda:

```bash
# Stripe Configuration (gunakan test keys dari Stripe Dashboard)
STRIPE_SECRET_KEY=sk_test_YOUR_ACTUAL_SECRET_KEY_HERE
STRIPE_WEBHOOK_SECRET=whsec_YOUR_WEBHOOK_SECRET_HERE
```

**Contoh keys yang valid:**

```bash
STRIPE_SECRET_KEY=sk_test_51ABC123DEF456...
STRIPE_WEBHOOK_SECRET=whsec_1ABC123DEF456...
```

### 3. **Restart Server**

Setelah update `.env`, restart server Anda:

```bash
# Stop server (Ctrl+C)
# Kemudian jalankan lagi
go run cmd/api/main.go
```

### 4. **Verifikasi Setup**

Test dengan endpoint health check dulu:

```bash
curl -X GET http://localhost:8080/api/v1/health
```

### 5. **Test Payment Flow**

Setelah Stripe keys ter-setup, ikuti flow ini:

1. **Login** → Dapatkan access token
2. **Get Fields** → Pilih field untuk booking
3. **Create Booking** → Buat booking baru
4. **Create Payment** → Test payment endpoint

## 🔧 **Production Environment**

Untuk production server (Railway), pastikan environment variables di Railway dashboard sudah di-set dengan keys yang sama.

## 📝 **Notes**

- Test keys aman untuk development
- Jangan pernah commit actual keys ke Git
- Gunakan `.env` untuk local development
- Railway environment variables untuk production

## 🇮🇩 **Untuk Developer Indonesia**

### **Rekomendasi Payment Gateway Lokal:**

1. **[Midtrans](https://midtrans.com/)** - Support IDR, kartu lokal
2. **[Xendit](https://xendit.co/)** - Support e-wallet, virtual account
3. **[DOKU](https://doku.com/)** - Payment gateway lokal
4. **[Moota](https://moota.co/)** - Mutasi bank otomatis

### **Alternative untuk Learning/Testing:**

- Gunakan **Stripe dengan Singapore/US** untuk belajar
- Currency tetap bisa `IDR` di code
- Test cards Stripe tetap bisa digunakan
- Untuk production, ganti ke payment gateway lokal

### **Test Card Numbers (Stripe):**

```
Visa: 4242 4242 4242 4242
Mastercard: 5555 5555 5555 4444
Expiry: Any future date
CVC: Any 3 digits
```
