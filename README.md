# Books on Hooks

## Quick Start (Podman)

From the repo root:

```bash
podman compose up --build
```

Run in detached mode:

```bash
podman compose up --build -d
```

Stop everything:

```bash
podman compose down
```

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

- `/usr/local/bin/images` in the `api` container
- `/usr/src/app/api/images` in the `dev` container

This means image files persist across container restarts/rebuilds as long as the `books-images` volume is not removed.

To remove all volumes (including images and DB data):

```bash
podman compose down -v
```
