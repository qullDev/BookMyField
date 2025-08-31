# Stripe Webhook 404 Fix

## 🚨 Masalah yang Ditemukan

Stripe webhook mengalami **"404 page not found"** karena dua masalah:

### 1. **Routing Conflict (FIXED)**

- Webhook endpoint: `POST /api/v1/payments/stripe-webhook`
- Payment detail endpoint: `GET /api/v1/payments/:id`

Gin router menganggap `stripe-webhook` sebagai `:id` parameter dan mengarahkan ke endpoint payment detail yang memerlukan authentication.

### 2. **URL Salah di Postman Collection (FIXED)**

- URL yang salah: `{{base_url}}/api/v1/webhook` ❌
- URL yang benar: `{{base_url}}/api/v1/payments/stripe-webhook` ✅

## ✅ Solusi yang Diterapkan

### 1. **Perubahan Routing**

**Sebelum (Bermasalah):**

```go
func PaymentRoutes(api *gin.RouterGroup) {
    payment := api.Group("/payments")
    {
        payment.GET("/:id", middlewares.AuthMiddleware(), controllers.GetPaymentByID)
        payment.POST("/stripe-webhook", controllers.StripeWebhook) // ❌ Conflict dengan /:id
    }
}
```

**Sesudah (Fixed):**

```go
func PaymentRoutes(api *gin.RouterGroup) {
    // Webhook endpoint di luar group untuk menghindari conflict
    api.POST("/payments/stripe-webhook", controllers.StripeWebhook) // ✅ No conflict

    payment := api.Group("/payments")
    {
        payment.GET("/:id", middlewares.AuthMiddleware(), controllers.GetPaymentByID)
    }
}
```

### 2. **Fix URL di Postman Collections**

**Postman Collection (Salah):**

```json
"url": {
  "raw": "{{base_url}}/api/v1/webhook", // ❌ URL salah
  "path": ["api", "v1", "webhook"]
}
```

**Postman Collection (Fixed):**

```json
"url": {
  "raw": "{{base_url}}/api/v1/payments/stripe-webhook", // ✅ URL benar
  "path": ["api", "v1", "payments", "stripe-webhook"]
}
```

### 3. **Testing Fix**

```bash
# Test webhook endpoint
curl -X POST "https://bookmyfield-production.up.railway.app/api/v1/payments/stripe-webhook" \
  -H "Content-Type: application/json" \
  -H "Stripe-Signature: test" \
  -d '{"test": "webhook"}'

# Response: {"error":"Invalid webhook signature"} ✅ (Expected - webhook works!)
```

## 🔧 Konfigurasi Stripe Dashboard

Pastikan webhook URL di Stripe Dashboard adalah:

```
https://bookmyfield-production.up.railway.app/api/v1/payments/stripe-webhook
```

Events yang perlu di-subscribe:

- `checkout.session.completed`
- `checkout.session.expired`
- `checkout.session.async_payment_failed`

## 📋 Checklist Deployment

- [x] Fix routing conflict
- [x] Fix URL di Postman collections
- [x] Test webhook endpoint accessibility
- [x] Verify signature validation
- [x] Update documentation
- [ ] Deploy to production (if needed)
- [ ] Update Stripe webhook URL (if needed)
- [ ] Test full payment flow

## 🚀 Status

**FIXED** - Webhook endpoint sekarang dapat diakses tanpa 404 error.

### ✅ Verification Results:

```bash
# Test URL yang salah
curl -X POST "https://bookmyfield-production.up.railway.app/api/v1/webhook"
# Response: 404 page not found ❌

# Test URL yang benar
curl -X POST "https://bookmyfield-production.up.railway.app/api/v1/payments/stripe-webhook"
# Response: {"error":"Invalid webhook signature"} ✅ (Expected - webhook works!)
```

**Postman Collection URLs sudah diperbaiki:**

- ✅ `postman_collection.json`
- ✅ `postman_collection_production.json`
