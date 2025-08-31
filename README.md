# 🏟️ BookMyField API

**BookMyField** adalah platform backend untuk sistem booking lapangan olahraga yang menyediakan solusi manajemen lapangan yang efisien dan andal untuk pengelola fasilitas olahraga dan pengguna. Dalam proyek ini, kami bertujuan untuk menyediakan API yang kuat dan terukur yang mendukung semua fungsi inti dari operasi booking lapangan modern. API ini mencakup fitur untuk mengelola pengguna, lapangan, booking, dan pembayaran dengan pengalaman yang seamless.

**Fitur:**

- 🔐 Otentikasi dengan JWT (Access & Refresh Token)
- 🔓 Otorisasi dengan RBAC (Role Based Access Control) (ADMIN, USER)
- 🏟️ Manajemen lapangan olahraga (CRUD operations)
- 📅 Sistem booking real-time dengan validasi tanggal/waktu
- 💳 Integrasi pembayaran dengan Stripe Payment Gateway (USD support)
- 🎣 Implementasi webhook Stripe untuk payment status update otomatis
- 💰 Manajemen payment lengkap dengan history dan tracking
- 📊 Payment endpoints untuk admin dan user dashboard
- 📄 Dokumentasi standar OpenAPI dengan Swagger UI
- ✅ Validasi data komprehensif dan error handling
- 🔒 JWT token blacklisting untuk keamanan ekstra
- 🔍 Advanced filtering untuk pencarian lapangan
- 📊 Auto-seeding admin user dan sample data untuk development

**Endpoints:**

- Auth endpoints untuk login, register, refresh token, dan logout
- Manajemen lapangan (Admin only)
- Sistem booking lapangan untuk user
- Track booking history pengguna
- Pembatalan booking dengan refund otomatis
- Payment processing dan webhook handling
- Payment history dan tracking untuk user dan admin
- Success/cancel pages untuk payment redirect
- Health check endpoint

**Default Credentials:**

- Admin: admin@admin.com | password123
- User: user@user.com | password123

**Tech Stack:** Go, Gin Framework, PostgreSQL/SQLite, GORM, Stripe API, JWT, Redis, Swagger

**Repo:** https://github.com/qullDev/BookMyField  
**Swagger:** http://localhost:8080/swagger/index.html  
**Postman Collection:** https://github.com/qullDev/BookMyField/blob/main/postman_collection.json

[![Go Version](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![API Documentation](https://img.shields.io/badge/API-Swagger-orange.svg)](http://localhost:8080/swagger/index.html)

## 📚 API Documentation

Dokumentasi API lengkap tersedia melalui Swagger UI:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### 📸 Screenshots Swagger UI

![Swagger Documentation](./images/swagger%201.png)
![Swagger API Endpoints](./images/swagger%202.png)

## ✨ Features

- **🔐 User Authentication**: Secure registration and login with JWT tokens
- **🏟️ Field Management**: Admins can create, update, and delete field information
- **📅 Booking System**: Users can book fields, view their booking history, and cancel bookings
- **💳 Payment Integration**: Seamless payment processing via Stripe with USD support
- **💰 Payment Management**: Complete payment tracking, history, and admin oversight
- **🎯 Payment Success/Cancel Pages**: User-friendly redirect pages after payment
- **👥 Role-Based Access Control**: Differentiated access for regular users and administrators
- **🔒 Security**: JWT blacklisting, input validation, and proper error handling
- **📊 Filtering**: Advanced field filtering by location and price range

## 🛠️ Tech Stack

- **Backend**: Go (Golang) v1.24+ with Gin framework
- **Database**: PostgreSQL with GORM (SQLite fallback for development)
- **Cache**: Redis for JWT token blacklisting and session management
- **Payments**: Stripe API with webhook support (USD currency)
- **API Documentation**: Swagger/OpenAPI 3.0 with interactive UI
- **Authentication**: JWT with access & refresh tokens
- **Validation**: Gin validator with comprehensive input validation
- **Development Tools**: Air for live reloading

---

## 🚀 Getting Started

Follow these steps to get the backend server running on your local machine.

### 📋 Prerequisites

- Go (version 1.18 or higher)
- PostgreSQL (or SQLite will be used automatically for development)
- Redis
- [Air](https://github.com/cosmtrek/air) for live reloading (optional, but recommended)

### 🔧 Installation

1.  **Clone the repository:**

    ```sh
    git clone https://github.com/qullDev/BookMyField.git
    cd BookMyField
    ```

2.  **Install Go dependencies:**

    ```sh
    go mod tidy
    ```

3.  **Set up Environment Variables:**
    Create a `.env` file in the root of the project by copying the example file:
    ```sh
    cp .env.example .env
    ```
    Now, open the `.env` file and fill in the required values for your local setup (Database, Redis, JWT Secret, Stripe keys).

### ▶️ Running the Application

- The application will automatically create the necessary database tables on startup (`AutoMigrate`).
- It will also seed the database with an admin user, a regular user, and some initial field data.
- If no `DATABASE_URL` is provided, the app will automatically use SQLite for development.

**Sample seeded data includes:**

- Admin user and regular user accounts
- 3 sample fields: Futsal (Jakarta), Basket (Bandung), Badminton (Surabaya)

To run the server:

```sh
go run ./cmd/api/main.go
```

For development with live-reloading (requires `air`):

```sh
air
```

The API server will start on `http://localhost:8080`.

### 👤 Default Accounts

After running the application, you can use these default accounts for testing:

**Admin Account:**

- Email: `admin@admin.com`
- Password: `password123`

**Regular User Account:**

- Email: `user@user.com`
- Password: `password123`

---

## 📖 API Reference

**Base URL**: `/api/v1`

### 🔐 Authentication

Endpoints for user registration and login.

#### 1. User Registration

- **Endpoint**: `POST /api/v1/auth/register`
- **Description**: Registers a new user.
- **Request Body**:

  ```json
  {
    "name": "John Doe",
    "email": "john.doe@example.com",
    "password": "password123"
  }
  ```

  **Validation:**

  - `name`: Required, minimum 2 characters
  - `email`: Required, valid email format
  - `password`: Required, minimum 6 characters

- **Success Response** (`201 Created`):
  ```json
  {
    "message": "User registered successfully"
  }
  ```
- **Error Response** (`400 Bad Request`):
  ```json
  {
    "error": "Email already registered"
  }
  ```
  **Other error responses:**
  - `400`: Validation errors (e.g., "Key: 'name' Error:Field validation for 'name' failed on the 'min' tag")
  - `500`: "Failed to hash password" or "Failed to create user"

#### 2. User Login

- **Endpoint**: `POST /api/v1/auth/login`
- **Description**: Authenticates a user and returns access and refresh tokens.
- **Request Body**:

  ```json
  {
    "email": "admin@admin.com",
    "password": "password123"
  }
  ```

  **Validation:**

  - `email`: Required, valid email format
  - `password`: Required

- **Success Response** (`200 OK`):
  ```json
  {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600,
    "refresh_token": "550e8400-e29b-41d4-a716-446655440000"
  }
  ```
- **Error Response** (`401 Unauthorized`):
  ```json
  {
    "error": "Invalid email or password"
  }
  ```
  **Other error responses:**
  - `400`: Validation errors (invalid email format, missing fields)
  - `500`: "Failed to generate access token" or "Failed to store refresh token"

#### 3. User Logout

- **Endpoint**: `POST /api/v1/auth/logout`
- **Authorization**: `Bearer <access_token>`
- **Description**: Logs the user out by blacklisting the access token and deleting the refresh token.
- **Request Body**:
  ```json
  {
    "refresh_token": "550e8400-e29b-41d4-a716-446655440000"
  }
  ```
- **Success Response** (`200 OK`):
  ```json
  {
    "message": "Logged out"
  }
  ```
- **Error Response** (`400 Bad Request`):
  ```json
  {
    "error": "Invalid body"
  }
  ```
  **Other error responses:**
  - `401`: Unauthorized (invalid or missing access token)

#### 4. Refresh Access Token

- **Endpoint**: `POST /api/v1/auth/refresh`
- **Description**: Issues a new access token using a valid refresh token and rotates the refresh token.
- **Request Body**:

  ```json
  {
    "refresh_token": "550e8400-e29b-41d4-a716-446655440000"
  }
  ```

  **Validation:**

  - `refresh_token`: Required

- **Success Response** (`200 OK`):
  ```json
  {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600,
    "refresh_token": "new-refresh-token-here"
  }
  ```
- **Error Response** (`401 Unauthorized`):
  ```json
  {
    "error": "Refresh token revoked or expired"
  }
  ```
  **Other error responses:**
  - `400`: "Invalid body" (missing refresh_token)
  - `401`: "User not found" (user associated with token doesn't exist)
  - `500`: "Failed to generate access token" or "Failed to store new refresh token"

### 🏟️ Fields

Endpoints for retrieving and managing field information.

#### 1. Get All Fields

- **Endpoint**: `GET /api/v1/fields`
- **Description**: Retrieves a list of all available fields, with optional filters.
- **Query Parameters**:

  - `location` (string, optional): Filter fields by location (case-insensitive search using ILIKE)
  - `min_price` (number, optional): Filter for fields with a price greater than or equal to this value
  - `max_price` (number, optional): Filter for fields with a price less than or equal to this value

- **Example URL**: `GET /api/v1/fields?location=jakarta&min_price=100000&max_price=300000`

- **Success Response** (`200 OK`):
  ```json
  [
    {
      "id": "c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d",
      "name": "Lapangan Futsal A",
      "location": "Jakarta",
      "price": 200000,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
  ```
- **Error Response** (`500 Internal Server Error`):
  ```json
  {
    "error": "Failed to retrieve fields"
  }
  ```

#### 2. Get Field by ID

- **Endpoint**: `GET /api/v1/fields/:id`
- **Description**: Retrieves details for a specific field.
- **Path Parameters**:

  - `id` (string, required): UUID of the field

- **Example URL**: `GET /api/v1/fields/c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d`

- **Success Response** (`200 OK`):
  ```json
  {
    "id": "c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d",
    "name": "Lapangan Futsal A",
    "location": "Jakarta",
    "price": 200000,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
  ```
- **Error Response** (`404 Not Found`):
  ```json
  {
    "error": "Field not found"
  }
  ```

#### 3. Create Field (Admin Only)

- **Endpoint**: `POST /api/v1/fields/admin`
- **Authorization**: `Bearer <admin_access_token>`
- **Description**: Creates a new field. Requires admin privileges.
- **Request Body**:

  ```json
  {
    "name": "Lapangan Tennis Baru",
    "location": "Bandung",
    "price": 150000
  }
  ```

  **Validation:**

  - `name`: Required
  - `location`: Required
  - `price`: Required, must be a number

- **Success Response** (`201 Created`):
  ```json
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Lapangan Tennis Baru",
    "location": "Bandung",
    "price": 150000,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
  ```
- **Error Response** (`403 Forbidden`):
  ```json
  {
    "error": "Forbidden"
  }
  ```
  **Other error responses:**
  - `400`: Validation errors (e.g., missing required fields)
  - `401`: Unauthorized (invalid or missing access token)
  - `500`: Database error

#### 4. Update Field (Admin Only)

- **Endpoint**: `PUT /api/v1/fields/admin/:id`
- **Authorization**: `Bearer <admin_access_token>`
- **Description**: Updates an existing field. Requires admin privileges.
- **Path Parameters**:

  - `id` (string, required): UUID of the field to update

- **Request Body**:
  ```json
  {
    "name": "Lapangan Tennis Updated",
    "location": "Bandung Barat",
    "price": 175000
  }
  ```
- **Success Response** (`200 OK`):
  ```json
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Lapangan Tennis Updated",
    "location": "Bandung Barat",
    "price": 175000,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z"
  }
  ```
- **Error Response** (`404 Not Found`):
  ```json
  {
    "error": "Field not found"
  }
  ```
  **Other error responses:**
  - `400`: Validation errors
  - `401`: Unauthorized (invalid or missing access token)
  - `403`: Forbidden (not admin)
  - `500`: Database error

#### 5. Delete Field (Admin Only)

- **Endpoint**: `DELETE /api/v1/fields/admin/:id`
- **Authorization**: `Bearer <admin_access_token>`
- **Description**: Deletes a field. Requires admin privileges.
- **Path Parameters**:

  - `id` (string, required): UUID of the field to delete

- **Success Response** (`200 OK`):
  ```json
  {
    "message": "Field deleted successfully"
  }
  ```
- **Error Response** (`404 Not Found`):
  ```json
  {
    "error": "Field not found"
  }
  ```
  **Other error responses:**
  - `401`: Unauthorized (invalid or missing access token)
  - `403`: Forbidden (not admin)
  - `500`: Database error

### 📅 Bookings

Endpoints for creating and managing user bookings.

#### 1. Get All Bookings (Admin Only)

- **Endpoint**: `GET /api/v1/bookings`
- **Authorization**: `Bearer <admin_access_token>`
- **Description**: Retrieves a list of all bookings from all users. Requires admin privileges.
- **Success Response** (`200 OK`):
  ```json
  [
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "user_id": "550e8400-e29b-41d4-a716-446655440003",
      "field_id": "c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d",
      "start_time": "2024-09-15T10:00:00Z",
      "end_time": "2024-09-15T12:00:00Z",
      "status": "confirmed",
      "payments": [],
      "created_at": "2024-09-15T09:00:00Z",
      "updated_at": "2024-09-15T09:30:00Z"
    }
  ]
  ```
- **Error Response** (`403 Forbidden`):
  ```json
  {
    "error": "Forbidden"
  }
  ```
  **Other error responses:**
  - `401`: Unauthorized (invalid or missing access token)
  - `500`: Database error

#### 2. Create a Booking

- **Endpoint**: `POST /api/v1/bookings`
- **Authorization**: `Bearer <user_access_token>`
- **Description**: Creates a new booking for a field.
- **Request Body**:

  ```json
  {
    "field_id": "c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d",
    "start_time": "2024-09-15T10:00:00Z",
    "end_time": "2024-09-15T12:00:00Z"
  }
  ```

  **Validation:**

  - `field_id`: Required, must be valid UUID
  - `start_time`: Required, must be valid ISO 8601 datetime, cannot be in the past
  - `end_time`: Required, must be valid ISO 8601 datetime, must be after start_time

- **Success Response** (`201 Created`):
  ```json
  {
    "id": "550e8400-e29b-41d4-a716-446655440004",
    "user_id": "550e8400-e29b-41d4-a716-446655440003",
    "field_id": "c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d",
    "start_time": "2024-09-15T10:00:00Z",
    "end_time": "2024-09-15T12:00:00Z",
    "status": "pending",
    "payments": [],
    "created_at": "2024-09-15T09:00:00Z",
    "updated_at": "2024-09-15T09:00:00Z"
  }
  ```
- **Error Response** (`409 Conflict`):
  ```json
  {
    "error": "Field is already booked for this time slot"
  }
  ```
  **Other error responses:**
  - `400`: "Invalid field_id format", "Invalid booking time"
  - `401`: Unauthorized
  - `404`: "Field not found"
  - `500`: "Failed to create booking"

#### 3. Get My Bookings

- **Endpoint**: `GET /api/v1/bookings/me`
- **Authorization**: `Bearer <user_access_token>`
- **Description**: Retrieves all bookings for the currently authenticated user with related field and payment data.
- **Success Response** (`200 OK`):
  ```json
  [
    {
      "id": "550e8400-e29b-41d4-a716-446655440004",
      "user_id": "550e8400-e29b-41d4-a716-446655440003",
      "field_id": "c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d",
      "start_time": "2024-09-15T10:00:00Z",
      "end_time": "2024-09-15T12:00:00Z",
      "status": "confirmed",
      "created_at": "2024-09-15T09:00:00Z",
      "updated_at": "2024-09-15T09:30:00Z",
      "field": {
        "id": "c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d",
        "name": "Lapangan Futsal A",
        "location": "Jakarta",
        "price": 200000,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      },
      "payments": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440005",
          "booking_id": "550e8400-e29b-41d4-a716-446655440004",
          "amount": 200000,
          "currency": "idr",
          "status": "succeeded",
          "stripe_ref_id": "cs_test_...",
          "created_at": "2024-09-15T09:15:00Z",
          "updated_at": "2024-09-15T09:20:00Z"
        }
      ]
    }
  ]
  ```
- **Error Response** (`401 Unauthorized`):
  ```json
  {
    "error": "Unauthorized"
  }
  ```
  **Other error responses:**
  - `500`: "Failed to fetch bookings"

#### 4. Cancel a Booking

- **Endpoint**: `DELETE /api/v1/bookings/:id` atau `DELETE /api/v1/bookings/:id/cancel`
- **Authorization**: `Bearer <user_access_token>`
- **Description**: Cancels a user's booking. If already paid, it will attempt to refund via Stripe.
- **Success Response** (`200 OK`):
  ```json
  {
    "message": "Booking cancelled and payment refunded successfully",
    "refund_id": "re_...",
    "refund_status": "succeeded"
  }
  ```
- **Error Response** (`404 Not Found`):
  ```json
  {
    "error": "Booking not found"
  }
  ```

### 💳 Payments

Endpoints for handling payments with Stripe integration.

#### 1. Create Stripe Checkout Session

- **Endpoint**: `POST /api/v1/payments/create-checkout-session`
- **Authorization**: `Bearer <user_access_token>`
- **Description**: Creates a Stripe checkout session for a pending booking.
- **Request Body**:

  ```json
  {
    "booking_id": "c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d"
  }
  ```

  **Validation:**

  - `booking_id`: Required, must be valid UUID

- **Success Response** (`200 OK`):
  ```json
  {
    "session_id": "cs_test_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
    "session_url": "https://checkout.stripe.com/pay/cs_test_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
  }
  ```
- **Error Response** (`400 Bad Request`):
  ```json
  {
    "error": "Payment already exists for this booking"
  }
  ```
  **Other error responses:**
  - `400`: "Booking is not in pending status"
  - `401`: "Unauthorized"
  - `404`: "Booking not found or not authorized"
  - `500`: Stripe API errors or "Failed to create payment record"

#### 2. Stripe Webhook

- **Endpoint**: `POST /api/v1/payments/stripe-webhook`
- **Description**: Handles Stripe webhook events to update payment and booking statuses. This endpoint is called by Stripe, not the frontend.
- **Headers**:
  - `Stripe-Signature`: Required for webhook verification
- **Webhook Events Handled**:

  - `checkout.session.completed`: Updates payment status to "succeeded" and booking status to "confirmed"
  - `checkout.session.expired`: Updates payment status to "failed"
  - `checkout.session.async_payment_failed`: Updates payment status to "failed"

- **Success Response** (`200 OK`):
  ```json
  {
    "status": "received"
  }
  ```
- **Error Response** (`400 Bad Request`):

  ```json
  {
    "error": "Invalid webhook signature"
  }
  ```

  **Other error responses:**

  - `400`: "Failed to read payload"
  - `500`: "Failed to update payment status"

- **Note**: This endpoint requires proper Stripe webhook signature verification and should be configured in your Stripe dashboard.

#### 2b. Stripe Webhook Test (Development)

- **Endpoint**: `POST /api/v1/payments/stripe-webhook-test`
- **Description**: Test webhook endpoint without signature validation for development/testing purposes.
- **Authorization**: None required
- **Request Body**: Same format as production webhook but without signature validation
- **Example**:
  ```json
  {
    "type": "checkout.session.completed",
    "data": {
      "object": {
        "id": "cs_test_session_123",
        "metadata": {
          "booking_id": "booking-uuid-here"
        }
      }
    }
  }
  ```
- **Response**:
  - `200`: `{"status": "test webhook processed", "event_type": "checkout.session.completed"}`
  - `400`: "Invalid JSON payload" or "Missing event type"

#### 3. Get All Payments (Admin Only)

- **Endpoint**: `GET /api/v1/payments/`
- **Authorization**: `Bearer <admin_access_token>`
- **Description**: Retrieves all payments in the system with complete booking, user, and field details. Admin access only.
- **Success Response** (`200 OK`):
  ```json
  [
    {
      "id": "550e8400-e29b-41d4-a716-446655440005",
      "booking_id": "550e8400-e29b-41d4-a716-446655440004",
      "amount": 10000,
      "currency": "idr",
      "status": "succeeded",
      "stripe_ref_id": "cs_test_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
      "created_at": "2024-09-15T09:15:00Z",
      "updated_at": "2024-09-15T09:20:00Z",
      "Booking": {
        "id": "550e8400-e29b-41d4-a716-446655440004",
        "user": {
          "id": "550e8400-e29b-41d4-a716-446655440003",
          "name": "John Doe",
          "email": "john@example.com",
          "role": "user"
        },
        "field": {
          "id": "c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d",
          "name": "Lapangan Futsal A",
          "location": "Jakarta",
          "price": 10000
        },
        "start_time": "2024-09-15T10:00:00Z",
        "end_time": "2024-09-15T12:00:00Z",
        "status": "confirmed",
        "notes": "Booking untuk latihan tim"
      }
    }
  ]
  ```
- **Error Response** (`403 Forbidden`):
  ```json
  {
    "error": "Forbidden"
  }
  ```

#### 4. Get My Payments

- **Endpoint**: `GET /api/v1/payments/me`
- **Authorization**: `Bearer <user_access_token>`
- **Description**: Retrieves all payments for the authenticated user with booking and field details.
- **Success Response** (`200 OK`):
  ```json
  [
    {
      "id": "550e8400-e29b-41d4-a716-446655440005",
      "booking_id": "550e8400-e29b-41d4-a716-446655440004",
      "amount": 10000,
      "currency": "idr",
      "status": "succeeded",
      "stripe_ref_id": "cs_test_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
      "created_at": "2024-09-15T09:15:00Z",
      "updated_at": "2024-09-15T09:20:00Z",
      "Booking": {
        "user": {
          "id": "550e8400-e29b-41d4-a716-446655440003",
          "name": "John Doe",
          "email": "john@example.com",
          "role": "user"
        },
        "field": {
          "id": "c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d",
          "name": "Lapangan Futsal A",
          "location": "Jakarta",
          "price": 10000
        },
        "start_time": "2024-09-15T10:00:00Z",
        "end_time": "2024-09-15T12:00:00Z",
        "status": "confirmed"
      }
    }
  ]
  ```
- **Error Response** (`401 Unauthorized`):
  ```json
  {
    "error": "Unauthorized"
  }
  ```

#### 5. Get Payment by ID

- **Endpoint**: `GET /api/v1/payments/:id`
- **Authorization**: `Bearer <access_token>`
- **Description**: Retrieves specific payment details. Users can only access their own payments, admins can access all payments.
- **Path Parameters**:
  - `id` (string, required): UUID of the payment
- **Success Response** (`200 OK`):
  ```json
  {
    "id": "550e8400-e29b-41d4-a716-446655440005",
    "booking_id": "550e8400-e29b-41d4-a716-446655440004",
    "amount": 10000,
    "currency": "idr",
    "status": "succeeded",
    "stripe_ref_id": "cs_test_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
    "created_at": "2024-09-15T09:15:00Z",
    "updated_at": "2024-09-15T09:20:00Z",
    "Booking": {
      "user": {
        "id": "550e8400-e29b-41d4-a716-446655440003",
        "name": "John Doe",
        "email": "john@example.com",
        "role": "user"
      },
      "field": {
        "id": "c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d",
        "name": "Lapangan Futsal A",
        "location": "Jakarta",
        "price": 10000
      },
      "start_time": "2024-09-15T10:00:00Z",
      "end_time": "2024-09-15T12:00:00Z",
      "status": "confirmed"
    }
  }
  ```
- **Error Response** (`404 Not Found`):
  ```json
  {
    "error": "Payment not found"
  }
  ```

#### 6. Payment Success Page

- **Endpoint**: `GET /success?session_id={session_id}`
- **Description**: Redirect page after successful Stripe payment. Returns success message with session details.
- **Query Parameters**:
  - `session_id` (string, optional): Stripe checkout session ID
- **Success Response** (`200 OK`):
  ```json
  {
    "status": "success",
    "message": "🎉 Payment completed successfully!",
    "session_id": "cs_test_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
    "next_steps": "Your booking has been confirmed. Check your bookings in the app."
  }
  ```

#### 7. Payment Cancel Page

- **Endpoint**: `GET /cancel`
- **Description**: Redirect page when payment is cancelled by user.
- **Success Response** (`200 OK`):
  ```json
  {
    "status": "cancelled",
    "message": "❌ Payment was cancelled. You can try again anytime."
  }
  ```

---

## 🏗️ Project Structure

```
BookMyField/
├── cmd/
│   └── api/
│       ├── main.go                 # Application entry point
│       ├── docs/                   # Swagger documentation
│       └── tmp/                    # Air build artifacts
├── internal/
│   ├── config/
│   │   ├── db.go                   # Database configuration
│   │   ├── jwt.go                  # JWT utilities
│   │   ├── redis.go                # Redis configuration
│   │   └── stripe.go               # Stripe configuration
│   ├── controllers/
│   │   ├── auth_controller.go      # Authentication endpoints
│   │   ├── booking_controller.go   # Booking management
│   │   ├── field_controller.go     # Field management
│   │   └── payment_controller.go   # Payment processing
│   ├── middlewares/
│   │   ├── jwt.go                  # JWT authentication middleware
│   │   └── role.go                 # Role-based access control
│   ├── models/
│   │   ├── booking.go              # Booking model
│   │   ├── field.go                # Field model
│   │   ├── payment.go              # Payment model
│   │   └── user.go                 # User model
│   ├── routes/
│   │   ├── auth.go                 # Authentication routes
│   │   ├── booking.go              # Booking routes
│   │   ├── field.go                # Field routes
│   │   └── payment.go              # Payment routes
│   └── seed/
│       ├── admin.go                # Admin user seeding
│       ├── field.go                # Field data seeding
│       └── user.go                 # Regular user seeding
├── images/                         # Documentation images
├── .env.example                    # Environment variables template
├── go.mod                          # Go module file
├── go.sum                          # Go dependencies checksum
└── README.md                       # This file
```

## 🔧 Configuration

### Environment Variables

Copy `.env.example` to `.env` and configure the following variables:

```bash
# Database Configuration (PostgreSQL)
# If not provided, SQLite will be used automatically for development
DATABASE_URL=postgres://username:password@host:port/dbname

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key_here_make_it_long_and_random

# Redis Configuration
REDIS_URL=redis://default:password@host:port

# Stripe Configuration
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key_here
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret_here

# Server Configuration
PORT=8080
```

## 🔐 Authentication

All protected endpoints require a JWT token in the Authorization header:

```bash
Authorization: Bearer <your_jwt_token>
```

### Get Access Token

```bash
curl -X POST "https://bookmyfield-production.up.railway.app/api/auth/login"
  -H "Content-Type: application/json"
  -d '{
    "email": "admin@test.com",
    "password": "admin123"
  }'
```

## 🚨 Error Handling

````

## ��� Error Handling

The API uses conventional HTTP response codes to indicate the success or failure of requests:

- **200**: Success
- **201**: Created
- **400**: Bad Request - Invalid request format or missing required fields
- **401**: Unauthorized - Invalid or missing authentication token
- **403**: Forbidden - Insufficient permissions
- **404**: Not Found - Resource not found
- **409**: Conflict - Resource conflict (e.g., duplicate booking)
- **500**: Internal Server Error - Server-side error

### Common Error Response Format

```json
{
  "error": "Descriptive error message"
}
````

## ��� Security Features

- **JWT Authentication**: Secure token-based authentication
- **Token Blacklisting**: Logout functionality with token invalidation
- **Refresh Token Rotation**: Enhanced security with token rotation
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: GORM ORM with prepared statements
- **CORS Support**: Cross-origin resource sharing configuration
- **Role-Based Access**: Admin and user role differentiation

## ��� Testing

### Manual Testing with cURL

**Register a new user:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Login with admin:**

```bash
curl -X POST "https://bookmyfield-production.up.railway.app/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "admin123"
  }'
```

**Create a booking:**

```bash
curl -X POST "https://bookmyfield-production.up.railway.app/api/bookings" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "field_id": "field-uuid-here",
    "date": "2024-01-15",
    "start_time": "10:00",
    "end_time": "12:00",
    "notes": "Weekly practice session"
  }'
```

**Create payment checkout session:**

```bash
curl -X POST "https://bookmyfield-production.up.railway.app/api/payments/create-checkout-session" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "booking_id": "booking-uuid-here"
  }'
```

**Get user payment history:**

```bash
curl -X GET "https://bookmyfield-production.up.railway.app/api/payments/my-payments" \
  -H "Authorization: Bearer <your_jwt_token>"
```

**Admin - Get all payments:**

```bash
curl -X GET "https://bookmyfield-production.up.railway.app/api/admin/payments" \
  -H "Authorization: Bearer <admin_jwt_token>"
```

curl -X POST http://localhost:8080/api/v1/auth/login \
 -H "Content-Type: application/json" \
 -d '{
"email": "admin@admin.com",
"password": "password123"
}'

````

**Login with regular user:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@user.com",
    "password": "password123"
  }'
````

**Get all fields:**

```bash
curl -X GET http://localhost:8080/api/v1/fields
```

**Get fields with filters:**

```bash
curl -X GET "http://localhost:8080/api/v1/fields?location=jakarta&min_price=100000&max_price=300000"
```

**Create a booking (requires user authentication):**

### Using Postman

1. Import the Postman collection from `postman_collection_production.json`
2. Import the environment from `postman_environment_production.json`
3. Run the "Login" request to authenticate
4. The access token will be automatically saved for subsequent requests

For detailed testing guides, see:

- [Payment Testing Guide](PAYMENT_TESTING_GUIDE.md)
- [Postman Production Guide](POSTMAN_PRODUCTION_GUIDE.md)
- [Payment Endpoints Guide](PAYMENT_ENDPOINTS_GUIDE.md)

## 📚 Additional Documentation

- [Payment Testing Guide](PAYMENT_TESTING_GUIDE.md) - Complete payment flow testing
- [Stripe Setup Guide](STRIPE_SETUP_GUIDE.md) - Stripe configuration
- [Postman Guide](POSTMAN_GUIDE.md) - API testing with Postman
- [Payment Endpoints Guide](PAYMENT_ENDPOINTS_GUIDE.md) - Payment API reference

## 💳 Payment Flow

1. **Create Booking** - User creates a booking for a field
2. **Initiate Payment** - User creates a Stripe checkout session
3. **Complete Payment** - User completes payment on Stripe
4. **Webhook Processing** - Stripe webhook updates payment status
5. **View History** - User can view payment history and details

## 🚀 Production Deployment

The API is deployed on Railway at:
**https://bookmyfield-production.up.railway.app**

### Database Schema

The application uses PostgreSQL with the following main tables:

- `users` - User accounts with role-based access
- `fields` - Sports field listings with pricing
- `bookings` - Field reservations with time slots
- `payments` - Payment records with Stripe integration

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new features
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📊 Database Schema

### Users Table

- `id` (UUID, Primary Key)
- `name` (VARCHAR, Not Null)
- `email` (VARCHAR, Not Null, Unique)
- `password` (VARCHAR, Not Null, Hashed with bcrypt)
- `role` (VARCHAR, Not Null, Default: 'user') - Values: 'user' | 'admin'
- `creadted_at` (TIMESTAMP) - Note: typo in model, should be created_at

### Fields Table

- `id` (UUID, Primary Key)
- `name` (VARCHAR(100), Not Null)
- `location` (VARCHAR(255), Not Null)
- `price` (DECIMAL/FLOAT, Not Null) - Price in IDR
- `created_at`, `updated_at` (TIMESTAMP)

### Bookings Table

- `id` (UUID, Primary Key)
- `user_id` (UUID, Foreign Key → users.id)
- `field_id` (UUID, Foreign Key → fields.id)
- `start_time` (TIMESTAMP)
- `end_time` (TIMESTAMP)
- `status` (VARCHAR) - Values: 'pending' | 'confirmed' | 'cancelled'
- `created_at`, `updated_at` (TIMESTAMP)

### Payments Table

- `id` (UUID, Primary Key)
- `booking_id` (UUID, Foreign Key → bookings.id)
- `amount` (DECIMAL/FLOAT) - Amount in IDR
- `currency` (VARCHAR) - Default: 'idr'
- `status` (VARCHAR) - Values: 'pending' | 'succeeded' | 'failed' | 'refunded'
- `stripe_ref_id` (VARCHAR) - Stripe session ID or payment intent ID
- `created_at`, `updated_at` (TIMESTAMP)

### Relationships

- **User → Bookings**: One-to-Many (a user can have multiple bookings)
- **Field → Bookings**: One-to-Many (a field can have multiple bookings)
- **Booking → Payments**: One-to-Many (a booking can have multiple payment attempts)
- **Booking → User**: Many-to-One (with Preload support)
- **Booking → Field**: Many-to-One (with Preload support)

## ��� Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature')
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## ��� License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ���‍��� Author

**qullDev**

- GitHub: [@qullDev](https://github.com/qullDev)

---

**Happy Coding! ���**
