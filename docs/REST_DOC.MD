# 🐾 Grayva API Documentation

**Base URL:**  
[`https://domain.com/`](https://domain.com/)  
**API Endpoint Root:**  
[`https://domain.com/api/v1`](https://domain.com/api/v1)

---

## 📸 GET Endpoints

### `GET /api/v1/photos/<id>`

Retrieve metadata for a specific photo by its ID.

**Example:**  
`GET https://domain.com/api/v1/photos/12345`

### `GET /api/v1/photos`

parameters

| Feature | Example    | Max | Min | Note                             |
| ------- | ---------- | --- | --- | -------------------------------- |
| limit   | ?limit=n   | 20  | 1   | Number of results to return      |
| offset  | ?offset=n  | n   | 1   | How much to offset the search by |
| entries | ?entries=1 | 1   | 1   | Return total entry number        |

---

## ✏️ PUT Endpoint

### `PUT /api/v1/photos`  

Create or update a photo entry.

#### Headers

```example
X-API-Key: admin-only-api-key  
Content-Type: multipart/form-data
```

#### Forms Content

```example
metadata=...;type=application/json
image=...;type=image/...
```

Note the json metadata must be sent first in the series for validation, prior to the image, otherwise your response will be rejected

**Json Body:**

```json
{
  "title": "Cat!",
  "description": "Cutie Pie",
  "tags": ["belly", "gray"]
}
```

---

## 🗑️ DELETE Endpoint

### `DELETE /api/v1/photos/<id>`  

Delete a photo entry by ID.

**Example:**  
`DELETE https://domain.com/api/v1/photos/12345`

**Headers:**

```curl
X-API-Key: admin-only-api-key
```

## Authentication Endpoint

### `GET /api/v1/auth`

Verify your authentication status

**Example**
`GET https://domain.com/api/v1/auth`

**Headers:**

```curl
X-API-Key: admin-only-api key
```

**Response:**

```json
{
  "valid": "true"
}
```
