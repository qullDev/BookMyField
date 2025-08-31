# Webhook Testing & Troubleshooting Guide

## 🚨 **Masalah "Bad Request" di Webhook**

### **Penyebab Umum:**

1. **Invalid Stripe Signature** - Webhook asli memerlukan signature valid dari Stripe
2. **Format payload salah** - Structure JSON tidak sesuai
3. **Headers tidak lengkap** - Missing required headers

## ✅ **Solusi Testing**

### **1. Endpoint Testing (Tanpa Signature Validation)**

**URL untuk testing:** `POST /api/v1/payments/stripe-webhook-test`

```bash
curl -X POST "https://bookmyfield-production.up.railway.app/api/v1/payments/stripe-webhook-test" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "checkout.session.completed",
    "data": {
        "object": {
            "id": "cs_test_session_123",
            "payment_status": "paid",
            "metadata": {
                "booking_id": "your-booking-id"
            }
        }
    }
}'
```

**Response berhasil:**

```json
{
  "status": "test webhook processed",
  "event_type": "checkout.session.completed",
  "session_id": "cs_test_session_123"
}
```

### **2. Event Types yang Didukung**

#### ✅ **checkout.session.completed**

- **Fungsi:** Update payment status ke "succeeded"
- **Fungsi:** Update booking status ke "confirmed"
- **Payload:**

```json
{
  "type": "checkout.session.completed",
  "data": {
    "object": {
      "id": "cs_session_id",
      "metadata": {
        "booking_id": "booking-uuid"
      }
    }
  }
}
```

#### ✅ **checkout.session.expired**

- **Fungsi:** Update payment status ke "failed"
- **Payload:**

```json
{
  "type": "checkout.session.expired",
  "data": {
    "object": {
      "id": "cs_session_id"
    }
  }
}
```

#### ✅ **checkout.session.async_payment_failed**

- **Fungsi:** Update payment status ke "failed"
- **Payload:** Same as expired

## 🔧 **Testing Flow**

### **Step 1: Buat Payment Session**

```bash
POST /api/v1/payments/create-checkout-session
{
  "booking_id": "your-booking-id"
}
```

**Response:** Dapatkan `stripe_session_id`

### **Step 2: Test Webhook**

```bash
POST /api/v1/payments/stripe-webhook-test
{
  "type": "checkout.session.completed",
  "data": {
    "object": {
      "id": "stripe_session_id_dari_step1",
      "metadata": {
        "booking_id": "your-booking-id"
      }
    }
  }
}
```

### **Step 3: Verify Payment Status**

```bash
GET /api/v1/payments/{{payment_id}}
```

**Expected result:** Payment status = "succeeded"

## 📋 **Postman Testing**

### **Collection Updates:**

- ✅ Added "Stripe Webhook Test (Development)" request
- ✅ Uses `{{stripe_session_id}}` and `{{booking_id}}` variables
- ✅ No signature validation required

### **Testing Steps:**

1. Import updated collection
2. Run "Create Checkout Session" to get session ID
3. Run "Stripe Webhook Test" to simulate payment completion
4. Run "Get Payment by ID" to verify status

## 🚀 **Production Webhook**

Untuk production, gunakan endpoint asli dengan signature validation:

- **URL:** `POST /api/v1/payments/stripe-webhook`
- **Requires:** Valid Stripe-Signature header
- **Setup:** Configure di Stripe Dashboard

### **Stripe Dashboard Configuration:**

```
Endpoint URL: https://bookmyfield-production.up.railway.app/api/v1/payments/stripe-webhook
Events to send:
- checkout.session.completed
- checkout.session.expired
- checkout.session.async_payment_failed
```

## 🐛 **Debugging Tips**

1. **Check logs** untuk error messages
2. **Verify payload structure** sesuai format yang diharapkan
3. **Test dengan endpoint -test** terlebih dahulu
4. **Pastikan environment variables** (STRIPE_WEBHOOK_SECRET) tersedia

## ⚠️ **Important Notes**

- **Test endpoint** hanya untuk development
- **Production webhook** memerlukan signature validation
- **Payment status** akan otomatis terupdate
- **Booking status** akan berubah ke "confirmed" saat payment berhasil
