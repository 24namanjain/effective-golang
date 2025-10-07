## Nakama step-by-step (manual)

Goal: understand Nakama by building one tiny Go runtime with two RPCs:
- open RPC issues a session token (JWT) using internal auth
- secured RPC requires a valid session and returns caller info

### 0) Prerequisites
- Docker and Docker Compose
- curl
- Go (optional locally; we’ll build with Docker anyway)

### 1) Create project structure
Run these commands manually (or create folders in your editor):
```bash
mkdir -p nakama-step-by-step/runtime/modules
```

### 2) Create minimal Go runtime module
Create a file at `nakama-step-by-step/runtime/modules/main.go` with a minimal skeleton. It must export `InitModule` and register two RPCs. You can start empty and fill later; the important part is the function signature:
```go
package main

import (
    "context"
    "database/sql"
    "github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
    // later: initializer.RegisterRpc("jwt_open", JwtOpenRPC)
    // later: initializer.RegisterRpc("jwt_secured", JwtSecuredRPC)
    logger.Info("Module loaded")
    return nil
}
```

Initialize a Go module inside `runtime/modules`:
```bash
cd nakama-step-by-step/runtime/modules
go mod init example.com/nakama-step-by-step
go get github.com/heroiclabs/nakama-common@latest
```

```result
```

Add RPCs (later):
- `jwt_open`: use `nk.AuthenticateCustom(ctx, customID, username, true)` to mint a session token (JWT)
- `jwt_secured`: check `ctx.Value(runtime.RUNTIME_CTX_USER_ID)`; if empty => 401

### 3) Build the plugin for Linux (for Docker image)
From repo root, build with a Linux toolchain using Docker:
```bash
docker run --rm -v "$PWD/nakama-step-by-step":/workspace -w /workspace/runtime/modules golang:1.21 \
  bash -lc 'GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ../nakama_go_jwt.so'
```

This should create `nakama-step-by-step/runtime/nakama_go_jwt.so`.

### 4) Configure Nakama to load the plugin
Create `nakama-step-by-step/nakama.yml`:
```yaml
runtime:
  path: "./data/modules"
  go:
    - nakama_go_jwt.so
```

### 5) Run Nakama via Docker Compose
Create `docker-compose.yml` next to this README:
```yaml
services:
  nakama:
    image: heroiclabs/nakama:latest
    container_name: nakama
    restart: unless-stopped
    ports:
      - "7350:7350"
      - "7351:7351"
    volumes:
      - ./nakama-step-by-step/nakama.yml:/nakama/data/nakama.yml:ro
      - ./nakama-step-by-step/runtime:/nakama/data:ro
    command: ["--config", "/nakama/data/nakama.yml", "--logger.level", "INFO"]
```

Start it:
```bash
docker compose up -d nakama
```

### 6) Implement and test RPCs
After you implement the two RPCs in `main.go` and rebuild the `.so` (repeat step 3), restart the container.

- Open RPC (get a session token):
```bash
curl -s -X POST "http://127.0.0.1:7350/v2/rpc/jwt_open" \
  -H "Content-Type: application/json" \
  -d '{"id":"demo-user","username":"demo"}'
```
Copy the value of `token` from the JSON response.

- Secured RPC (requires token):
```bash
curl -s -X POST "http://127.0.0.1:7350/v2/rpc/jwt_secured" \
  -H "Authorization: Bearer <paste-token-here>" \
  -H "Content-Type: application/json" \
  -d '{}'
```

You should see JSON like:
```json
{"status":"ok","user_id":"...","username":"demo"}
```

### Mental model: what’s happening
- Nakama loads your Go plugin and calls `InitModule`.
- `RegisterRpc` makes endpoints available at `/v2/rpc/<id>`.
- `AuthenticateCustom` creates/returns a session and JWT using server keys.
- For secured RPCs, Nakama validates the `Authorization: Bearer <token>` header and places `user_id` in the context.

### Reference
- Nakama Console/API docs: [Nakama Console API](https://heroiclabs.com/docs/nakama/api/console/)
