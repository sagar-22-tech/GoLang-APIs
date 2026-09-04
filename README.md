# API Endpoints

Base URL:

**`https://golang-apis-dox2.onrender.com`**

## Available APIs

| API                | Method    | Endpoint                    | Description                                 |
| ------------------ | --------- | --------------------------- | ------------------------------------------- |
| **Hello API**      | **`GET`** | **`/greet`**                | Returns a greeting with the current time.   |
| **User API**       | **`GET`** | **`/users`**                | Returns mock data of all users.             |
| **User API by ID** | **`GET`** | **`/users/:id`**            | Returns mock user data using the user's ID. |
| **Search Users**   | **`GET`** | **`/users?search=:search`** | Searches users by their first name.         |

> **Redis Caching:** Redis has been integrated into selected User API endpoints to cache frequently requested data and improve API performance. During testing, Redis caching reduced API response time by approximately **50%**.

---

## 1. Hello API

### Endpoint

```http
GET /greet
```

### Example

```text
https://golang-apis-dox2.onrender.com/greet
```

Returns a greeting along with the current time and HTTP method.

---

## 2. Get All Users

### Endpoint

```http
GET /users
```

### Example

```text
https://golang-apis-dox2.onrender.com/users
```

Returns mock data containing all users.

Redis caching is used to cache frequently requested user data and reduce repeated data lookups.

---

## 3. Get User by ID

### Endpoint

```http
GET /users/:id
```

### Example

```text
https://golang-apis-dox2.onrender.com/users/1
```

Replace **`:id`** with the ID of the user you want to retrieve.

Redis caching is used to cache user data by ID, allowing repeated requests for the same user to be served directly from Redis.

---

## 4. Search Users

### Endpoint

```http
GET /users?search=:search
```

### Example

```text
https://golang-apis-dox2.onrender.com/users?search=sagar
```

The **`search`** query parameter searches users by their first name.

### Example Requests

```http
GET /users?search=sagar

GET /users?search=vivaan

GET /users?search=aditya
```

The search is:

* Case-insensitive
* Supports partial matches
* Cached using Redis for repeated searches

---

# Redis Caching

Redis is integrated as a caching layer for selected User API endpoints.

The purpose of Redis is to avoid repeatedly reading and searching the JSON data source when the same data is requested multiple times.

## How It Works

```text
                    Client Request
                          │
                          ▼
                    Check Redis
                          │
                 ┌────────┴────────┐
                 │                 │
             CACHE HIT         CACHE MISS
                 │                 │
                 ▼                 ▼
          Return cached        Read db.json
              data                 │
                                   ▼
                              Store in Redis
                                   │
                                   ▼
                              Return data
```

## Cache Hit

When the requested data already exists in Redis:

```text
Client
  ↓
Go API
  ↓
Redis
  ↓
Cached Data
  ↓
Client
```

The API can return the cached data without repeatedly reading and processing the JSON file.

## Cache Miss

When the requested data does not exist in Redis:

```text
Client
  ↓
Go API
  ↓
Redis → CACHE MISS
  ↓
Read db.json
  ↓
Search / Retrieve User
  ↓
Store Result in Redis
  ↓
Return Response
```

The next request for the same data can then be served directly from Redis.

---

## Performance Improvement

Redis caching resulted in an **approximately 50% reduction in API response time during testing**.

### Improvements

* Implemented Redis caching for frequently accessed API endpoints.
* Reduced repeated reads and searches of the JSON data source.
* Cached user data by ID.
* Cached repeated user search results.
* Used Redis Cloud for the deployed API.
* Observed approximately **50% faster response times** during testing.

> **Performance Result:** API response time decreased by approximately **50%** after implementing Redis caching, based on development testing.

---

# Response Format

The APIs return data in JSON format.

### Example User Response

```json
{
  "id": 1,
  "name": {
    "firstname": "Sagar",
    "lastname": "Saini"
  },
  "username": "sagar_saini",
  "email": "sagar.saini@example.com",
  "address": {
    "street": "MG Road",
    "city": "Delhi",
    "zipcode": "110001",
    "geo": {
      "lat": "28.6139",
      "lng": "77.2090"
    }
  },
  "phone": "+91-9876543201",
  "website": "aaravsharma.dev",
  "company": {
    "name": "TechNova Solutions",
    "post": "Software Engineer",
    "salary": "75000"
  }
}
```

---

# Testing the APIs

You can test these endpoints using:

* **Browser**
* **Postman**
* **cURL**
* **Go Playground** — an interactive API explorer for testing the project's endpoints.

### Go Playground

```text
https://go-playground-weld.vercel.app/
```

---

# Technologies Used

* **Go**
* **net/http**
* **Redis**
* **Redis Cloud**
* **Docker**
* **Render**
* **JSON**
* **CORS**

---

# Architecture

```text
                    ┌──────────────────┐
                    │      Client      │
                    └────────┬─────────┘
                             │
                             ▼
                    ┌──────────────────┐
                    │    Go REST API   │
                    │     (Render)     │
                    └────────┬─────────┘
                             │
                    ┌────────┴─────────┐
                    │                  │
                    ▼                  ▼
             ┌─────────────┐    ┌─────────────┐
             │ Redis Cloud │    │   db.json   │
             │    Cache    │    │ Data Source │
             └─────────────┘    └─────────────┘
```

Redis acts as a caching layer between the Go API and the underlying JSON data source, allowing frequently requested data to be served faster.
