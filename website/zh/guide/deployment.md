# 部署

## Docker Compose

仓库根目录提供了 `Dockerfile`(多阶段:构建前端 → 构建后端 → 运行镜像)与
`docker-compose.yml`。一条命令即可构建并启动:

```bash
docker compose up --build
```

启动后访问 `http://localhost:8080` —— Go 服务端会同时托管已构建的 Web UI。

`docker-compose.yml` 默认配置:

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

数据通过挂载的 `./data` 卷持久化。

## 单二进制部署

后端使用纯 Go 的 SQLite 驱动,可在禁用 CGO 的情况下编译为单文件:

```bash
cd web && npm install && npm run build      # 产出 web/dist
cd ../server && CGO_ENABLED=0 go build -o promptops .
```

将 `promptops` 二进制与 `web/dist` 一起部署即可。若运行目录下存在
`web/dist/index.html`,服务端会自动托管前端。

## 生产建议

- 通过 `PROMPTOPS_TOKEN` 设置一个强随机的管理员 Token。
- 首次启动会自动创建 `admin` / `admin` 账号,请尽快在「用户管理」中修改。
- 将 SQLite 文件放在持久化卷上并定期备份。
