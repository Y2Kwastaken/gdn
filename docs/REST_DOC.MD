# 🐾 Grayva API Documentation

**Base URL:**  
[`https://grayva.miles.sh/`](https://grayva.miles.sh/)  
**API Endpoint Root:**  
[`https://grayva.miles.sh/api/v1`](https://grayva.miles.sh/api/v1)

---

## 📸 GET Endpoints

### `GET /api/v1/photos/<id>`

Retrieve metadata for a specific photo by its ID.

**Example:**  
`GET https://grayva.miles.sh/api/v1/photos/12345`

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

**Headers:**

```
X-API-Key: admin-only-api-key  
Content-Type: multipart/form-data
```

**Forms Content**

```
image=...;type=application/json
metadata=...;type=image/...
```

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

### `DELETE /api/v1/photos?id=<id>`  

Delete a photo entry by ID.

**Example:**  
`DELETE https://grayva.miles.sh/api/v1/photos?id=12345`

**Headers:**

```
X-API-Key: admin-only-api-key
```
