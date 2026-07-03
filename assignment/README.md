# Assignment Import API

A simple Go + Gin REST API for importing assignment user data from CSV into MySQL and caching the results in Redis.

## What this project does

- Uploads a CSV file and imports valid rows into MySQL
- Supports CRUD operations for assignment users
- Uses Redis to cache list and detail reads
- Refreshes the Redis list cache after import/create/update/delete operations

## Tech stack

- Go
- Gin web framework
- sqlx + MySQL
- Redis via go-redis
- Zerolog for logging

## Prerequisites

- Go 1.22+
- MySQL running locally or in Docker
- Redis running locally or in Docker

## Configuration

Update the database and Redis settings in [services/core/app/config.yaml](services/core/app/config.yaml).

Example defaults:
- Server port: 9025
- MySQL: localhost:3306, database assignmentdb
- Redis: localhost:6379

## Database setup

### 1) Start MySQL
If MySQL is not already running, start it locally or using Docker.

### 2) Create the database
Run the following in MySQL:

```sql
CREATE DATABASE assignmentdb;
```

### 3) Verify the connection
The app will automatically create the required table named `assignment_users` when it starts.

## Redis setup

### 1) Start Redis
If Redis is not already running, start it locally or using Docker.

### 2) Verify Redis is available
You can test it with:

```bash
redis-cli ping
```

If Redis is running correctly, it should return:

```text
PONG
```

## Run the application

```bash
cd /home/scalent/Downloads/assignment
go run ./services/core/app
```

The server will start at:

```text
http://localhost:9025
```

## API endpoints

Base URL:

```text
http://localhost:9025/tms-core/assignment-user
```

| Method | Endpoint | Description |
| --- | --- | --- |
| POST | /upload-csv | Upload and import a CSV file |
| POST | / | Create a new assignment user |
| PATCH | /:id | Partial update an assignment user |
| DELETE | /:id | Soft delete an assignment user |
| GET | /:id | Get one assignment user |
| GET | / | List all assignment users |

## CSV format

The CSV file must include these headers:

```text
first_name,last_name,company_name,address,city,county,postal,phone,email,web
```

## Example curl commands

### Upload CSV

```bash
curl -X POST "http://localhost:9025/tms-core/assignment-user/upload-csv" \
  -F "csvFile=@/path/to/users.csv"
```

### Create a user

```bash
curl -X POST "http://localhost:9025/tms-core/assignment-user/" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "John",
    "lastName": "Doe",
    "companyName": "Acme Inc",
    "email": "john@example.com"
  }'
```

### Get all users

```bash
curl -X GET "http://localhost:9025/tms-core/assignment-user/"
```

### Get one user

```bash
curl -X GET "http://localhost:9025/tms-core/assignment-user/1"
```

### Update a user

```bash
curl -X PATCH "http://localhost:9025/tms-core/assignment-user/1" \
  -H "Content-Type: application/json" \
  -d '{
    "companyName": "Updated Company",
    "phone": "555-9876"
  }'
```

### Delete a user

```bash
curl -X DELETE "http://localhost:9025/tms-core/assignment-user/1"
```

## Redis behavior

- Reads are served from Redis when available
- The list cache key is: assignment_users:list
- Detail cache keys are stored as: assignment_users:detail:<id>
- Cache entries expire after 5 minutes
- Import/create/update/delete operations refresh the list cache from MySQL
