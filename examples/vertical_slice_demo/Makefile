.PHONY: build run-admin run-api run-worker clean test

# 编译所有入口
build:
	@echo "编译管理端..."
	@go build -o bin/admin cmd/admin/main.go
	@echo "编译 C 端..."
	@go build -o bin/api cmd/api/main.go
	@echo "编译 Worker..."
	@go build -o bin/worker cmd/worker/main.go
	@echo "✅ 编译完成！"

# 运行管理端
run-admin:
	@go run cmd/admin/main.go

# 运行 C 端
run-api:
	@go run cmd/api/main.go

# 运行 Worker
run-worker:
	@go run cmd/worker/main.go

# 清理编译文件
clean:
	@rm -rf bin/
	@echo "✅ 清理完成！"

# 运行测试
test:
	@go test ./...

# 安装依赖
deps:
	@go mod tidy
	@echo "✅ 依赖安装完成！"

# 查看帮助
help:
	@echo "可用命令:"
	@echo "  make build      - 编译所有入口"
	@echo "  make run-admin  - 运行管理端 (端口 8081)"
	@echo "  make run-api    - 运行 C 端 (端口 8080)"
	@echo "  make run-worker - 运行 Worker"
	@echo "  make clean      - 清理编译文件"
	@echo "  make test       - 运行测试"
	@echo "  make deps       - 安装依赖"

