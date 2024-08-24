# Project Documentation

## Project Overview

This project is a web application that allows users to manage services, add items to their cart, make payments, and track their orders. The application uses a combination of PostgreSQL for relational data and MongoDB for non-relational data, integrating Razorpay for payment processing.

### Features

- **User Management**: Create and manage users.
- **Service Management**: Create and manage services.
- **Cart Management**: Add, view, and remove items from the cart.
- **Order Management**: Process payments, move items to orders, and track delivery status.

## Endpoints

### 1. Cart Management

#### Get User Cart
- **Endpoint**: `GET /cart`
- **Description**: Fetch all items in the user's cart.
- **Function**: `GetUserCart`

#### Add Item to Cart
- **Endpoint**: `POST /cart/:id`
- **Description**: Add a new item to the user's cart. The `id` parameter specifies the service ID.
- **Function**: `AddItemToCart`

#### Delete Item from Cart
- **Endpoint**: `DELETE /cart/:id`
- **Description**: Remove an item from the user's cart by its ID.
- **Function**: `DeleteItemFromCart`

### 2. Service Management

#### Create Service
- **Endpoint**: `POST /services`
- **Description**: Create a new service.
- **Function**: `CreateService`

#### Get Services
- **Endpoint**: `GET /services`
- **Description**: Fetch all services.
- **Function**: `GetService`

### 3. User Management

#### Create User
- **Endpoint**: `POST /users`
- **Description**: Create a new user.
- **Function**: `CreateUser`

#### Get User
- **Endpoint**: `GET /users/:id`
- **Description**: Fetch a user by ID.
- **Function**: `GetUser`

### 4. Order Management

#### Payment
- **Endpoint**: `POST /payment/:user_id`
- **Description**: Handle payment processing and move items from the cart to orders. The `user_id` parameter specifies the user.
- **Function**: `Payment`

#### Update Order Status
- **Endpoint**: `PUT /orders/:order_id/status`
- **Description**: Update the delivery status of an order. The `order_id` parameter specifies the order.
- **Function**: `UpdateOrderStatus`

#### Get Order Status
- **Endpoint**: `GET /orders/:order_id`
- **Description**: Fetch the status and details of an order. The `order_id` parameter specifies the order.
- **Function**: `GetOrderStatus`

## Models

### User
```go
type User struct {
    ID    uint    `json:"id"`
    Name  string  `json:"name"`
    Email *string `json:"email"`
}
