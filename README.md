# Excel Import CRUD (Go + Gin + MySQL + Redis)

A small, self-contained Gin API that:

1. Accepts an `.xlsx` upload, validates it, and imports it **asynchronously**.
2. Stores the rows in MySQL.
3. Caches reads in Redis (5-minute TTL).
4. Lets you list, view, edit, and delete the imported records (CRUD), keeping
   MySQL and Redis in sync.

The whole app is deliberately kept in a single `main` package with one file
per concern (`db.go`, `cache.go`, `excel.go`, `handlers.go`, ...) instead of
a layered/DI architecture, so the flow of a request is easy to follow
top-to-bottom.

## Project layout

```
main.go          entry point: config, connections, routes
config.go        env-var based configuration
models.go        Customer / UpdateCustomerInput / JobStatus structs
db.go            MySQL connection + table creation
repository.go    SQL queries (insert/list/get/update/delete)
cache.go         Redis read/write helpers (5 min TTL)
excel.go         Excel parsing + header/row validation
jobs.go          in-memory tracker for async import jobs
handlers.go      HTTP handlers (upload, jobs, CRUD)
response.go      consistent JSON response helpers
postman_collection.json  ready-to-import Postman collection
```

## Requirements

- Go 1.22+
- MySQL running locally (or update `MYSQL_DSN`)
- Redis running locally (or update `REDIS_ADDR`)

## Setup

```bash
cp .env.example .env      # edit if your MySQL/Redis differ from the defaults
export $(cat .env | xargs)

go mod tidy                # downloads gin, mysql driver, redis client, excelize
go run .
```

The app creates the `customers` table automatically on startup (`CREATE
TABLE IF NOT EXISTS`), so no manual migration is needed — just make sure the
database named in `MYSQL_DSN` (default `assignment_db`) already exists:

```sql
CREATE DATABASE IF NOT EXISTS assignment_db;
```

> **Note:** this environment could not reach the public Go module proxy, so
> `go.sum` is not included. Running `go mod tidy` once with normal internet
> access will fetch the dependencies and generate it.

## Excel file format

The first row must be a header row containing at least: `first_name`,
`last_name`, `email` (case-insensitive, spaces are treated as underscores).
These optional columns are also read if present: `company_name`, `address`,
`city`, `county`, `postal`, `phone`, `web`. Rows missing a required field are
skipped and counted, not failed — the import still completes for the valid
rows.

## API

| Method | Path                  | Description                                   |
|--------|-----------------------|------------------------------------------------|
| POST   | `/api/upload`         | Upload an `.xlsx` file (form field: `file`). Returns a `job_id` immediately; import runs in the background. |
| GET    | `/api/jobs/:id`       | Check the status of an import (`processing` / `completed` / `failed`). |
| GET    | `/api/customers`      | List all customers. Served from Redis if cached, else MySQL (and re-caches). |
| GET    | `/api/customers/:id`  | Get one customer, same cache-aside behavior.  |
| PUT    | `/api/customers/:id`  | Update any subset of fields (JSON body). Updates MySQL then refreshes the Redis cache. |
| DELETE | `/api/customers/:id`  | Delete a customer from MySQL and Redis.       |

### Example: upload

```bash
curl -F "file=@customers.xlsx" http://localhost:8080/api/upload
```

### Example: update

```bash
curl -X PUT http://localhost:8080/api/customers/1 \
  -H "Content-Type: application/json" \
  -d '{"city": "New City", "phone": "1234567890"}'
```

A ready-to-import `postman_collection.json` is included with all of the
above requests set up.

## Design notes

- **Async import**: the upload handler validates the file synchronously
  (format, headers) but hands the row-by-row insert off to a goroutine, so
  the client isn't blocked while thousands of rows are written. Progress can
  be polled via `/api/jobs/:id`.
- **Batch inserts**: rows are inserted in chunks of 200 in a single `INSERT
  ... VALUES (...), (...), ...` statement instead of one query per row.
- **Cache-aside reads**: `GET` endpoints check Redis first and only hit
  MySQL on a miss, then repopulate the cache (5 minute TTL).
- **Write-through updates**: `PUT` updates MySQL, refreshes the single-record
  cache entry, and drops the now-stale "all customers" list cache so the
  next list read rebuilds it.
- **Horizontal scaling**: the app is stateless except for the in-memory job
  tracker (which is only used for progress reporting, not correctness) — you
  can run multiple instances behind a load balancer pointed at the same
  MySQL/Redis, and connection pool limits are set on the MySQL client to
  avoid one instance exhausting the database.
