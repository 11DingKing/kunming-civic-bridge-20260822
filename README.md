# 承滇风精神 抒云岭民意

昆明市人民建议征集与办理协同平台后端。系统把线上建议、176 个线下征集点位、专业建议人、部门办理和群众回访连接成可审计的闭环。

## 运行

```bash
cp .env.example .env
go run ./cmd/server
```

服务默认监听 `:8080`，提供 `/healthz`、`/readyz` 和 `/api/v1`。SQLite 数据库会在启动时应用 `migrations/`，可通过 `DB_PATH` 指定路径。

## 角色

`citizen` 可提交和回访建议；`reviewer` 可研判和分派；`operator` 可办理和提交回应；`supervisor` 可审签、督办和归档。

## 验证

```bash
go build ./...
go test ./... -count=1
go test -race ./... -count=1
go vet ./...
```
