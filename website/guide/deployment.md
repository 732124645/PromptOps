# Deployment

## Docker Compose

The repository root ships a `Dockerfile` (multi-stage: build frontend → build
backend → runtime image) and a `docker-compose.yml`. One command builds and
starts everything:

```bash
docker compose up --build
```

Then open `http://localhost:8080` — the Go server also serves the built Web UI.

Default `docker-compose.yml`:

```yaml
services:
  promptops:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PROMPTOPS_TOKEN=promptops-dev-token
      - PROMPTOPS_DB=/app/data/promptops.db
    volumes:
      - ./data:/app/data
```

Data is persisted through the mounted `./data` volume.

## Single-binary deployment

The backend uses a pure-Go SQLite driver, so it can be compiled to a single
file with CGO disabled:

```bash
cd web && npm install && npm run build      # produces web/dist
cd ../server && CGO_ENABLED=0 go build -o promptops .
```

Deploy the `promptops` binary together with `web/dist`. If
`web/dist/index.html` exists in the working directory, the server serves the
frontend automatically.

## Production notes

- Set a strong random admin token via `PROMPTOPS_TOKEN`.
- A default `admin` / `admin` account is created on first run — change it in
  "User Management" as soon as possible.
- Put the SQLite file on a persistent volume and back it up regularly.
