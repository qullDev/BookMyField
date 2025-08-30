# 📊 Payment Endpoints Documentation

## 🔗 **Available Payment Endpoints**

### 1. **Create Checkout Session**

```
POST /api/v1/payments/create-checkout-session
```

- **Description:** Create Stripe checkout session for payment
- **Auth:** Required (User/Admin)
- **Body:** `{"booking_id": "uuid"}`
- **Response:** Stripe checkout URL

### 2. **Get All Payments (Admin Only)**

```
GET /api/v1/payments/
```

- **Description:** Get all payments in system
- **Auth:** Required (Admin only)
- **Response:** Array of payments with booking/user/field details

### 3. **Get My Payments**

```
GET /api/v1/payments/me
```

- **Description:** Get user's own payments
- **Auth:** Required (User/Admin)
- **Response:** Array of user's payments with details

### 4. **Get Payment by ID**

```
GET /api/v1/payments/{payment_id}
```

- **Description:** Get specific payment details
- **Auth:** Required (User can only access own payments, Admin can access all)
- **Response:** Payment details with booking/user/field data

### 5. **Stripe Webhook**

```
POST /api/v1/payments/stripe-webhook
```

- **Description:** Handle Stripe payment events
- **Auth:** None (Called by Stripe)
- **Body:** Stripe event payload

### 6. **Payment Success Page**

```
GET /success?session_id={session_id}
```

- **Description:** Redirect page after successful payment
- **Auth:** None
- **Response:** Success message with session details

### 7. **Payment Cancel Page**

```
GET /cancel
```

- **Description:** Redirect page when payment is cancelled
- **Auth:** None
- **Response:** Cancel message

## 📝 **Payment Data Structure**

```json
{
  "id": "uuid",
  "booking_id": "uuid",
  "amount": 10000,
  "currency": "idr",
  "status": "pending|succeeded|failed",
  "stripe_ref_id": "cs_test_...",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "Booking": {
    "id": "uuid",
    "user": {
      "id": "uuid",
      "name": "User Name",
      "email": "user@email.com",
      "role": "user"
    },
    "field": {
      "id": "uuid",
      "name": "Field Name",
      "location": "Location",
      "price": 10000
    },
    "start_time": "timestamp",
    "end_time": "timestamp",
    "status": "pending|confirmed|cancelled",
    "notes": "booking notes"
  }
}
```

## 🎯 **Complete Payment Flow**

1. **User Login** → Get access token
2. **Get Fields** → Choose field for booking
3. **Create Booking** → Get booking ID
4. **Create Checkout Session** → Get Stripe payment URL
5. **User Pays** → Stripe processes payment
6. **Payment Success** → User redirected to success page
7. **Get My Payments** → View payment history

## 🚀 **Testing Credentials**

### User Account:

- Email: `user@user.com`
- Password: `password123`

### Admin Account:

- Email: `admin@admin.com`
- Password: `password123`

## 📱 **Postman Collection Updated**

The production Postman collection now includes:

- ✅ Get All Payments (Admin)
- ✅ Get My Payments
- ✅ Get Payment by ID
- ✅ Auto-save payment IDs for testing
- ✅ Proper authorization headers

**API is now fully complete with comprehensive payment management! 🎉**
