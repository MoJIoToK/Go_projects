# Limits Normalization Service

Service for normalizing client limits data and storing the result in a database.

## Description

The service accepts raw limits data either as plain text or as a file,
applies normalization rules, stores normalized data in a SQLite database
and provides an API to retrieve the normalized result.

## Business Rules

- Limits are grouped by `CLIENT_CODE` and `SECCODE`
- Each position must contain limit kinds `{0, 1, 2, 365}`
- Positions with `OPEN_BALANCE = 0` and `OPEN_LIMIT = 0` for all limit kinds
  are removed
- If all positions of a client are removed, a fallback position is created
- Validation errors do not stop processing and are returned as warnings

## API

### POST /limits/normalize

Accepts raw limits data in one of the following formats:

#### Plain text

- `Content-Type: text/plain`
- Request body contains raw limits data

#### File upload

- `Content-Type: multipart/form-data`
- File is read from the request body

The endpoint:

- parses input data
- applies normalization rules
- saves normalized data to the database

Response:

- `200 OK` on success
- validation warnings are logged but do not stop processing

---

### GET /limits

Returns normalized limits data in JSON format.

Response:

- `200 OK`
- JSON with clients, positions and limits

## Storage

- SQLite database
- File: `./db/limits.db`
- Tables: `clients`, `positions`, `limits`

## Run

```bash
go run main.go
```



