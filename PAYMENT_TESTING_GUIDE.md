# 💳 Payment Testing Guide - BookMyField API

## 🔄 **Complete Payment Flow**

### 1. **Prerequisites**

- ✅ User sudah login (punya access token)
- ✅ Ada booking yang sudah dibuat
- ✅ Booking dalam status "pending"

### 2. **Step-by-step Payment Testing**

#### Step 1: Login User

```json
POST {{base_url}}/api/v1/auth/login
{
    "email": "user@user.com",
    "password": "password123"
}
```

**Response:**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "refresh_token": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### Step 2: Get All Fields (untuk dapat field_id)

```
GET {{base_url}}/api/v1/fields/
Authorization: Bearer {{access_token}}
```

#### Step 3: Create Booking

```json
POST {{base_url}}/api/v1/bookings/
Authorization: Bearer {{access_token}}
{
    "field_id": "{{field_id}}",
    "start_time": "2025-09-20T10:00:00Z",
    "end_time": "2025-09-20T12:00:00Z",
    "notes": "Booking untuk latihan tim"
}
```

**Response:**

```json
{
  "id": "booking-uuid-here",
  "user_id": "user-uuid",
  "field_id": "field-uuid",
  "start_time": "2025-09-20T10:00:00Z",
  "end_time": "2025-09-20T12:00:00Z",
  "status": "pending",
  "notes": "Booking untuk latihan tim",
  "created_at": "2025-08-30T15:30:00Z"
}
```

#### Step 4: Create Checkout Session

```json
POST {{base_url}}/api/v1/checkout
Authorization: Bearer {{access_token}}
{
    "booking_id": "{{booking_id}}"
}
```

**Expected Response:**

```json
{
  "checkout_url": "https://checkout.stripe.com/pay/cs_test_...",
  "session_id": "cs_test_..."
}
```

#### Step 5: Simulate Payment Success (Webhook)

```json
POST {{base_url}}/api/v1/webhook
Content-Type: application/json
Stripe-Signature: test_signature
{
    "type": "checkout.session.completed",
    "data": {
        "object": {
            "id": "cs_test_123",
            "metadata": {
                "booking_id": "{{booking_id}}"
            },
            "payment_status": "paid"
        }
    }
}
```

## 🔧 **Postman Testing Setup**

### Environment Variables:

```
base_url = https://bookmyfield-production.up.railway.app
access_token = (auto-set dari login)
field_id = (auto-set dari get fields)
booking_id = (auto-set dari create booking)
checkout_url = (auto-set dari create checkout)
```

### Test Scripts untuk Auto-save:

#### Login User - Test Script:

```javascript
pm.test("Login successful", function () {
  var jsonData = pm.response.json();
  pm.environment.set("access_token", jsonData.access_token);
  pm.environment.set("refresh_token", jsonData.refresh_token);
});
```

#### Get All Fields - Test Script:

```javascript
pm.test("Get fields successful", function () {
  var jsonData = pm.response.json();
  if (jsonData.length > 0) {
    pm.environment.set("field_id", jsonData[0].id);
  }
});
```

#### Create Booking - Test Script:

```javascript
pm.test("Booking created", function () {
  var jsonData = pm.response.json();
  pm.environment.set("booking_id", jsonData.id);
});
```

#### Create Checkout - Test Script:

```javascript
pm.test("Checkout session created", function () {
  var jsonData = pm.response.json();
  pm.environment.set("checkout_url", jsonData.checkout_url);
});
```

## ⚠️ **Common Issues & Solutions**

### 1. **401 Unauthorized pada /checkout**

**Penyebab:** Token expired atau tidak ada
**Solusi:**

- Login ulang untuk mendapatkan token baru
- Pastikan Authorization header format: `Bearer {{access_token}}`

### 2. **"Booking not found"**

**Penyebab:** Booking ID salah atau tidak dimiliki user
**Solusi:**

- Pastikan booking_id benar (dari response create booking)
- Pastikan booking dibuat oleh user yang sama

### 3. **"Stripe not configured"**

**Penyebab:** Environment variable Stripe tidak ada
**Solusi:**

- Check STRIPE_SECRET_KEY di environment production

### 4. **Webhook testing**

**Catatan:**

- Webhook endpoint tidak perlu authentication
- Biasanya dipanggil oleh Stripe server
- Untuk testing manual, gunakan sample payload

## 🎯 **Quick Testing Scenario**

1. **Health Check** ✅
2. **Login User** ✅ (save token)
3. **Get All Fields** ✅ (save field_id)
4. **Create Booking** ✅ (save booking_id)
5. **Create Checkout Session** ✅ (save checkout_url)
6. **[Optional] Test Webhook** ✅

## 📝 **Expected Responses**

### Create Checkout Session Success:

```json
{
  "checkout_url": "https://checkout.stripe.com/pay/cs_test_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0",
  "session_id": "cs_test_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0"
}
```

### After Webhook (Check booking status):

```json
{
  "id": "booking-uuid",
  "status": "confirmed", // changed from "pending"
  "payments": [
    {
      "id": "payment-uuid",
      "stripe_session_id": "cs_test_...",
      "amount": 10000,
      "status": "completed"
    }
  ]
}
```

## 🌐 **Production URLs**

- **Base URL:** `https://bookmyfield-production.up.railway.app`
- **Checkout:** `POST /api/v1/checkout`
- **Webhook:** `POST /api/v1/webhook`

---

**Happy Payment Testing! 💳✨**
