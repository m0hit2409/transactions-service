# Transactions Service

A small Go service that manages cardholder **accounts** and the **transactions**
posted against them. Purchases and withdrawals are stored as negative amounts and
credit vouchers as positive — the caller always sends a positive amount and the
**server applies the sign** based on the operation type.

The service uses an embedded **SQLite** database (via
[go-sqlite3](https://github.com/mattn/go-sqlite3)) so there are no external
dependencies to spin up. Schema migrations run automatically on startup.

---

## Quick start

### Option A — Docker (recommended)

**Prerequisites:** a Docker-compatible runtime (Docker Desktop, or see
[Colima](#docker-runtime-on-macos-colima) below).

1. From the `transactions-service` directory, build and start everything with
   one command:

   ```bash
   ./run          # build + start in the foreground (Ctrl-C to stop)
   ./run -d       # or: start in the background
   ```

   This builds the image and starts the API container. SQLite migrations run
   automatically on startup — there is nothing else to set up.

2. Confirm it's up:

   ```bash
   curl -s localhost:8080/health
   # {"status":"ok"}
   ```

The API is now available at <http://localhost:8080> and the interactive docs
at <http://localhost:8080/docs/index.html>.

To stop it: `Ctrl-C` (foreground) or `docker compose down` (if started with
`-d`).

#### Docker runtime on macOS (Colima)

If you don't have Docker Desktop, [Colima](https://github.com/abiosoft/colima)
provides a Docker-compatible runtime:

```bash
brew install colima docker docker-compose
colima start
./run
```

### Option B — Without Docker

**Prerequisites:**
- Go 1.25+ (`go version`)
- A C compiler (`gcc`/`clang`) on your `PATH` — `go-sqlite3` uses CGO, so this
  is required even though there's no separate database service to install.
  - macOS: install the Xcode Command Line Tools — `xcode-select --install`
  - Debian/Ubuntu: `sudo apt-get install build-essential`

1. From the `transactions-service` directory, fetch and verify dependencies:

   ```bash
   go mod download
   go mod tidy      # optional sanity check that go.mod/go.sum are in sync
   ```

2. (Optional) Copy the example env file and adjust if you want non-default
   settings:

   ```bash
   cp .env.example .env.local
   ```

   By default the app listens on port `8080` and stores its SQLite file at
   `data/pismo.db` (created automatically) — no env file is required to run
   it.

3. Run the service:

   ```bash
   go run ./cmd/api
   ```

   Or build a binary first and run that:

   ```bash
   go build -o bin/api ./cmd/api
   ./bin/api
   ```

4. In another terminal, confirm it's up:

   ```bash
   curl -s localhost:8080/health
   # {"status":"ok"}
   ```

The API is now available at <http://localhost:8080> and the interactive docs
at <http://localhost:8080/docs/index.html>.

To stop it: `Ctrl-C` in the terminal running the process.

---

## Endpoints

| Method | Path                    | Description           |
| ------ | ----------------------- | --------------------- |
| POST   | `/accounts`             | Create an account     |
| GET    | `/accounts/{accountId}` | Retrieve an account   |
| POST   | `/transactions`         | Create a transaction  |
| GET    | `/health`               | Liveness check        |
| GET    | `/docs/index.html`      | Swagger UI            |

### Examples

```bash
# Create an account
curl -s -X POST localhost:8080/accounts \
  -H 'Content-Type: application/json' \
  -d '{"document_number":"12345678900"}'
# 201 -> {"account_id":1,"document_number":"12345678900", ...}

# Fetch it
curl -s localhost:8080/accounts/1

# A purchase (operation_type_id 1) — amount is stored negative
curl -s -X POST localhost:8080/transactions \
  -H 'Content-Type: application/json' \
  -d '{"account_id":1,"operation_type_id":1,"amount":123.45}'
# 201 -> {"transaction_id":1,"amount":"-123.45", ...}

# A credit voucher (operation_type_id 4) — amount is stored positive
curl -s -X POST localhost:8080/transactions \
  -H 'Content-Type: application/json' \
  -d '{"account_id":1,"operation_type_id":4,"amount":60}'
# 201 -> {"transaction_id":2,"amount":"60", ...}
```

### Operation types

| ID | Description                | Sign |
| -- | -------------------------- | ---- |
| 1  | Normal Purchase            | −    |
| 2  | Purchase with installments | −    |
| 3  | Withdrawal                 | −    |
| 4  | Credit Voucher             | +    |

### Status codes

| Situation                              | Status |
| -------------------------------------- | ------ |
| Created                                | 201    |
| Malformed body / non-positive amount   | 400    |
| Account not found                      | 404    |
| Duplicate document number              | 409    |
| Unknown operation type                 | 422    |

Errors share one envelope: `{"error":"message"}`.

---

## Configuration

| Variable       | Default          | Description                   |
| -------------- | ---------------- | ----------------------------- |
| `PORT`         | `8080`           | HTTP listen port              |
| `DATABASE_URL` | `data/pismo.db`  | SQLite file path              |

Copy `.env.example` to `.env.local` and override as needed.
