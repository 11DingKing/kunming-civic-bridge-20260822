# BENZHI_README

## 项目说明

- 项目：11DingKing/kunming-civic-bridge-20260822
- 项目用途：昆明市人民建议征集与办理协同平台后端。系统把线上建议、176 个线下征集点位、专业建议人、部门办理和群众回访连接成可审计的闭环。
- Go 工具链：`golang:1.23.0`
- 前端工具链：无

## 标准构建、运行和测试命令

进入容器后执行：

```bash
# 编译
cd '/app' && GOTOOLCHAIN=local go build ./...

# 启动
cd '/app' && GOTOOLCHAIN=local go run ./cmd/server

# 测试
cd '/app' && GOTOOLCHAIN=local go test ./...
```

## Docker 构建和进入容器

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh benzhi-task-65-amd64 linux/amd64
./build_benzhi_docker.sh benzhi-task-65-arm64 linux/arm64
docker run -it benzhi-task-65-amd64:latest
docker run -it --platform linux/arm64 benzhi-task-65-arm64:latest
```

## 题目验证命令

1. 预期退出码 0：`go test ./internal/repository -run '^TestRunningJobCannotBeClaimedAgain$' -count=1`
2. 预期退出码 0：`go test ./...`
3. 预期退出码 0：`GOTOOLCHAIN=local go build -buildvcs=false ./... && GOTOOLCHAIN=local go vet ./...`
