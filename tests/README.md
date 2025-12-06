# 测试目录

本目录包含项目所有模块的单元测试，按模块分类组织。

## 目录结构

```
test/
├── configuration/      # Configuration 模块测试
│   └── configuration_test.go
├── errors/            # Errors 模块测试
│   └── builder_test.go
└── validation/        # Validation 模块测试
    ├── collection_rules_test.go  # 集合验证规则测试
    ├── number_rules_test.go      # 数字验证规则测试
    ├── result_test.go            # 验证结果测试
    ├── string_builder_test.go    # 字符串构建器测试
    └── validator_test.go         # 验证器核心测试
```

## 包命名约定

所有测试文件使用独立的测试包名（`*_test` 后缀），遵循 Go 语言最佳实践：

- `configuration_test` - Configuration 模块测试
- `errors_test` - Errors 模块测试
- `validation_test` - Validation 模块测试

这种方式的优点：
- 测试代码与源代码完全分离
- 只能访问公开的 API，确保黑盒测试
- 避免测试依赖内部实现细节

## 运行测试

### 运行所有测试
```bash
go test ./test/...
```

### 运行特定模块测试
```bash
go test ./test/validation/...
go test ./test/errors/...
go test ./test/configuration/...
```

### 运行特定测试
```bash
go test ./test/validation/... -run TestNotEmptySlice
```

### 查看详细输出
```bash
go test ./test/... -v
```

## 测试覆盖率

查看测试覆盖率：
```bash
go test ./test/... -cover
```

生成覆盖率报告：
```bash
go test ./test/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```
