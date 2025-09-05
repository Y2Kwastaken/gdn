# 🐾 Grayva API Documentation

**Base URL:**  
[`https://grayva.miles.sh/`](https://grayva.miles.sh/)  
**API Endpoint Root:**  
[`https://grayva.miles.sh/api/v1`](https://grayva.miles.sh/api/v1)

---

## 📸 GET Endpoints

### `GET /api/v1/photo/<id>`
Retrieve metadata for a specific photo by its ID.

**Example:**  
`GET https://grayva.miles.sh/api/v1/photo/12345`

---

### `GET /api/v1/id?amount=<amount>&random=<boolean>`
Fetch a list of photo IDs.

- `amount`: Number of IDs to return.
- `random`: If `true`, returns random IDs.

**Example:**  
`GET https://grayva.miles.sh/api/v1/id?amount=10&random=true`

---

### `GET /api/v1/id/count`
Returns the total count of photo entries in the system.

**Example:**  
`GET https://grayva.miles.sh/api/v1/id/count`

---

## ✏️ PUT Endpoint

### `PUT /api/v1/photo`  
Create or update a photo entry.

**Headers:**
```
X-API-Key: admin-only-api-key  
Content-Type: application/json
```

**Body:**
```json
{
  "title": "Cat!",
  "description": "Cutie Pie",
  "tags": ["belly", "gray"]
}
```

---

## 🗑️ DELETE Endpoint

### `DELETE /api/v1/photo?id=<id>`  
Delete a photo entry by ID.

**Example:**  
`DELETE https://grayva.miles.sh/api/v1/photo?id=12345`

**Headers:**
```
X-API-Key: admin-only-api-key
```
