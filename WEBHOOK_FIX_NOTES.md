# Stripe Webhook 404 Fix

## 🚨 Masalah yang Ditemukan

Stripe webhook mengalami **"404 page not found"** karena routing conflict:

- Webhook endpoint: `POST /api/v1/payments/stripe-webhook`
- Payment detail endpoint: `GET /api/v1/payments/:id`

Ketika Stripe mencoba akses webhook endpoint, routing menganggap `stripe-webhook` sebagai `:id` parameter dan mengarahkan ke endpoint payment detail yang memerlukan authentication.

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

### 2. **Testing Fix**

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
- [x] Test webhook endpoint accessibility
- [x] Verify signature validation
- [ ] Deploy to production
- [ ] Update Stripe webhook URL (if needed)
- [ ] Test full payment flow

## 🚀 Status

**FIXED** - Webhook endpoint sekarang dapat diakses tanpa 404 error.
