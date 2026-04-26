# Books on Hooks

## Quick Start (Podman)

Run in detached mode:

```bash
podman compose up --build -d
```

Stop everything:

```bash
podman compose down
```

## Scripts

Helper scripts live in `scripts/`:

- `./scripts/admin.sh counts`
- `./scripts/admin.sh reset-database`
- `./scripts/test.sh`
- `./scripts/reload.sh`

## Services and Ports

- API (`api`): `http://localhost:4000`
- API dev/hot-reload (`dev`): `http://localhost:5000`
- Frontend (Vite dev server, `client-vite`): `http://localhost:5173`
- Frontend (built static app behind nginx, `client-nginx`): `http://localhost:8081`
- PostgreSQL main DB (`db`): `localhost:5432`
- PostgreSQL test DB (`test-db`): `localhost:5433`

## Frontend Services

This repo has two frontend runtime modes:

- `client-vite`: development mode (Vite server)
- `client-nginx`: production-style mode (built assets served by nginx)

If you only want one frontend mode, start specific services, for example:

```bash
podman compose up --build api db client-vite
```

or:

```bash
podman compose up --build api db client-nginx
```

## Database Access

Main DB connection settings:

- Host: `localhost`
- Port: `5432`
- Database: `books`
- User: `postgres`
- Password: `secret`

Test DB connection settings:

- Host: `localhost`
- Port: `5433`
- Database: `books_test`
- User: `postgres`
- Password: `secret`

Example with `psql`:

```bash
psql "postgres://postgres:secret@localhost:5432/books?sslmode=disable"
```

Test DB:

```bash
psql "postgres://postgres:secret@localhost:5433/books_test?sslmode=disable"
```

## Image Storage

Uploaded book images are stored in the named volume `books-images`.

That volume is mounted at:

- `/data/images` in the `api` container
- `/data/images` in the `dev` container

This means image files persist across container restarts/rebuilds as long as the `books-images` volume is not removed.

To remove all volumes (including images and DB data):

```bash
podman compose down -v
```

## Packages

Backend package overview (`api/internal`):

- `app`: application wiring, middleware, and route setup.
- `domain`: shared domain models and request/response DTO structs.
- `handlers/book`: HTTP handlers for book endpoints and book metric recording.
- `handlers/machine`: HTTP handlers for machine endpoints and machine metric recording.
- `repository`: database-facing persistence, transactional operations, and file-system interactions for image storage lifecycle.
- `requestx`: shared request parsing and normalization helpers.
- `validator`: reusable field-level validation helpers.
- `httpErrors`: common HTTP error response helpers (including structured validation errors).
- `logger`: logging abstraction used across the API.
- `sqlc`: generated query types and database methods from SQL files.

## Admin Auth

Admin routes use server-side session auth.

Required API env vars:

- `ADMIN_USERNAME`
- `ADMIN_PASSWORD_HASH`
- `COOKIE_SECURE`

Generate a bcrypt password hash with the admin utility:

```bash
cd api
go run ./cmd/admin generate-password-hash "your-password"
```

Use the printed hash as `ADMIN_PASSWORD_HASH`.

Notes:

- `COOKIE_SECURE=false` is appropriate for local HTTP development only.
- Set `COOKIE_SECURE=true` when the app is running behind HTTPS.

## Metrics

### Best-Effort

Machine metrics are recorded as best-effort telemetry and must not break core user flows.

Current behavior:

- If metric query params are invalid (for example, malformed `is_qr` or `source`), the API logs a warning and still returns the machine response.
- If metric insert fails, the API logs an error and still returns the machine response.
- Only core machine lookup errors (missing machine, DB read failures) affect the HTTP response status.

### QR and Session Metrics Flow
Book click metrics are written when the summary endpoint is called with `machine` (and optional metadata):

- `machine`: machine id
- `source`: click source, currently `location-grid`
- `is_qr`: whether the user arrived from a QR URL

Machine page entry supports QR tracking with a URL like:

```text
/location/1?is_qr=true
```

Flow:

1. `metricsSession` helper (`frontend/src/lib/metricsSession.js`) initializes session metrics context:
   - ensures a session-scoped `session_id` exists (UUID in `sessionStorage`)
   - reads optional `is_qr` from the machine page URL and persists it in `sessionStorage`
2. Cover links from the machine grid to summary pages include:
   - `machine=<machine id>`
   - `source=location-grid`
   - `is_qr=<true|false>` (from session state)
   - `session_id=<uuid>` (from session state)
3. `BookLocation` forwards query params to `/api/books/:id/summary`.
4. API inserts the metric row using those values.

Notes:

- The `is_qr` value is session-scoped (tab/session lifetime), so it persists across back/forward and multiple summary clicks during the session.
- This QR behavior is only applied for machine-grid sourced clicks (`source=location-grid`).
