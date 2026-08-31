## API Endpoints

Base URL:

`https://golang-apis-dox2.onrender.com`

### Available APIs

| API            | Method | Endpoint                | Description                                 |
| -------------- | ------ | ----------------------- | ------------------------------------------- |
| Hello API      | `GET`  | `/greet`                | Returns a greeting with the current time.   |
| User API       | `GET`  | `/users`                | Returns mock data of all users.             |
| User API by ID | `GET`  | `/users/:id`            | Returns mock user data using the user's ID. |
| Search Users   | `GET`  | `/users?search=:search` | Searches users by their first name.         |

### 1. Hello API

**Endpoint**

```http
GET /greet
```

**Example**

```http
 https://golang-apis-dox2.onrender.com/greet
```

---

### 2. Get All Users

**Endpoint**

```http
GET /users
```

**Example**

```http
https://golang-apis-dox2.onrender.com/users
```

Returns mock data containing all users.

---

### 3. Get User by ID

**Endpoint**

```http
GET /users/:id
```

**Example**

```http
https://golang-apis-dox2.onrender.com/users/1
```

Replace `:id` with the ID of the user you want to retrieve.

---

### 4. Search Users

**Endpoint**

```http
GET /users?search=:search
```

**Example**

```http
 https://golang-apis-dox2.onrender.com/users?search=sagar
```

The `search` query parameter searches users by their first name.

**Example requests:**

```http
GET /users?search=sagar
GET /users?search=vivaan
GET /users?search=aditya
```

The search is case-insensitive and also supports partial matches.

### Response

The APIs return data in JSON format.

```json
{
  "id": 1,
  "name": {
    "firstname": "Sagar",
    "lastname": "Saini"
  }
}
```

### Testing the APIs

You can test these endpoints using:

* Browser
* Postman
* cURL
* **Go Playground** — the project's frontend API explorer
