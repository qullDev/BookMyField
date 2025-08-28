# BookMyField API Documentation

Welcome to the BookMyField API! This document provides a comprehensive guide for frontend developers to interact with the backend services for a field booking application.

## Features

-   **User Authentication**: Secure registration and login for users.
-   **Field Management**: Admins can create, update, and delete field information.
-   **Booking System**: Users can book fields, view their booking history, and cancel bookings.
-   **Payment Integration**: Seamless payment processing via Stripe.
-   **Role-Based Access Control**: Differentiated access for regular users and administrators.

## Tech Stack

-   **Backend**: Go (Golang) with Gin framework
-   **Database**: PostgreSQL
-   **Cache**: Redis for token management
-   **Payments**: Stripe

---

## Getting Started

Follow these steps to get the backend server running on your local machine.

### Prerequisites

-   Go (version 1.18 or higher)
-   PostgreSQL
-   Redis
-   [Air](https://github.com/cosmtrek/air) for live reloading (optional, but recommended)

### Installation

1.  **Clone the repository:**
    ```sh
    git clone <your-repo-url>
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

### Running the Application

-   The application will automatically create the necessary database tables on startup (`AutoMigrate`).
-   It will also seed the database with an admin user, a regular user, and some initial field data.

To run the server:

```sh
go run ./cmd/api/main.go
```

For development with live-reloading (requires `air`):

```sh
air
```

The API server will start on `http://localhost:8080`.

---

## API Reference

**Base URL**: `/api/v1`

### Authentication

Endpoints for user registration and login.

#### 1. User Registration

-   **Endpoint**: `POST /auth/register`
-   **Description**: Registers a new user.
-   **Request Body**:
    ```json
    {
      "name": "John Doe",
      "email": "john.doe@example.com",
      "password": "password123"
    }
    ```
-   **Success Response** (`201 Created`):
    ```json
    {
      "message": "User registered successfully"
    }
    ```
-   **Error Response** (`400 Bad Request`):
    ```json
    {
      "error": "Email already registered"
    }
    ```

#### 2. User Login

-   **Endpoint**: `POST /auth/login`
-   **Description**: Authenticates a user and returns access and refresh tokens.
-   **Request Body**:
    ```json
    {
      "email": "john.doe@example.com",
      "password": "password123"
    }
    ```
-   **Success Response** (`200 OK`):
    ```json
    {
      "access_token": "ey...",
      "expires_in": 3600,
      "refresh_token": "..."
    }
    ```
-   **Error Response** (`401 Unauthorized`):
    ```json
    {
      "error": "Invalid email or password"
    }
    ```

### Fields

Endpoints for retrieving and managing field information.

#### 1. Get All Fields

-   **Endpoint**: `GET /fields`
-   **Description**: Retrieves a list of all available fields, with optional filters.
-   **Query Parameters**:
    -   `location` (string, optional): Filter fields by location (case-insensitive search).
    -   `min_price` (number, optional): Filter for fields with a price greater than or equal to this value.
    -   `max_price` (number, optional): Filter for fields with a price less than or equal to this value.
-   **Success Response** (`200 OK`):
    ```json
    [
      {
        "ID": "c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d",
        "Name": "Main Soccer Field",
        "Location": "Central Park",
        "Price": 50.00
      }
    ]
    ```

#### 2. Get Field by ID

-   **Endpoint**: `GET /fields/:id`
-   **Description**: Retrieves details for a specific field.
-   **Success Response** (`200 OK`):
    ```json
    {
      "ID": "c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d",
      "Name": "Main Soccer Field",
      "Location": "Central Park",
      "Price": 50.00
    }
    ```
-   **Error Response** (`404 Not Found`):
    ```json
    {
      "error": "Field not found"
    }
    ```

#### 3. Create Field (Admin Only)

-   **Endpoint**: `POST /fields/admin`
-   **Authorization**: `Bearer <admin_access_token>`
-   **Request Body**:
    ```json
    {
      "name": "New Tennis Court",
      "location": "Westside Club",
      "price": 75.50
    }
    ```
-   **Success Response** (`201 Created`):
    ```json
    {
      "ID": "...",
      "Name": "New Tennis Court",
      "Location": "Westside Club",
      "Price": 75.50
    }
    ```

### Bookings

Endpoints for creating and managing user bookings.

#### 1. Create a Booking

-   **Endpoint**: `POST /bookings`
-   **Authorization**: `Bearer <user_access_token>`
-   **Request Body**:
    ```json
    {
      "field_id": "c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d",
      "start_time": "2024-09-15T10:00:00Z",
      "end_time": "2024-09-15T12:00:00Z"
    }
    ```
-   **Success Response** (`201 Created`):
    ```json
    {
      "ID": "...",
      "UserID": "...",
      "FieldID": "c1f8e4d9-8a2b-4b6e-9c1d-5a8f8c7b6a5d",
      "StartTime": "2024-09-15T10:00:00Z",
      "EndTime": "2024-09-15T12:00:00Z",
      "Status": "pending",
      "Payments": null
    }
    ```

#### 2. Get My Bookings

-   **Endpoint**: `GET /bookings/me`
-   **Authorization**: `Bearer <user_access_token>`
-   **Description**: Retrieves all bookings for the currently authenticated user.
-   **Success Response** (`200 OK`):
    ```json
    [
      {
        "ID": "...",
        "UserID": "...",
        "FieldID": "...",
        "StartTime": "...",
        "EndTime": "...",
        "Status": "confirmed",
        "Field": {
          "ID": "...",
          "Name": "Main Soccer Field",
          "Location": "Central Park",
          "Price": 50.00
        },
        "Payments": [
          {
            "ID": "...",
            "BookingID": "...",
            "Amount": 100.00,
            "Status": "succeeded",
            "PaymentIntentID": "pi_..."
          }
        ]
      }
    ]
    ```

#### 3. Cancel a Booking

-   **Endpoint**: `DELETE /bookings/:id/cancel`
-   **Authorization**: `Bearer <user_access_token>`
-   **Description**: Cancels a user's booking. If already paid, it will attempt to refund via Stripe.
-   **Success Response** (`200 OK`):
    ```json
    {
      "message": "Booking cancelled and payment refunded"
    }
    ```
-   **Error Response** (`404 Not Found`):
    ```json
    {
      "error": "Booking not found"
    }
    ```

### Payments

Endpoints for handling payments.

#### 1. Create Stripe Checkout Session

-   **Endpoint**: `POST /checkout`
-   **Authorization**: `Bearer <user_access_token>`
-   **Description**: Creates a Stripe checkout session for a pending booking.
-   **Request Body**:
    ```json
    {
      "booking_id": "..."
    }
    ```
-   **Success Response** (`200 OK`):
    ```json
    {
      "session_id": "cs_test_...",
      "url": "https://checkout.stripe.com/pay/cs_test_..."
    }
    ```

#### 2. Stripe Webhook

-   **Endpoint**: `POST /webhook`
-   **Description**: Listens for events from Stripe to update payment and booking statuses. This is intended for Stripe to call, not the frontend.

