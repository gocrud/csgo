# 快速参考卡片

本文档提供常用的速查表和对比表，帮助你快速做出决策。

---

## 📁 项目结构速查

### 目录用途对照表

| 位置 | 用途 | 何时使用 | 示例 |
|------|------|---------|------|
| `apps/*/features/` | 功能模块 | 端特有功能 | users/, products/ |
| `apps/*/internal/` | 应用内部共享 | 应用内多功能共享 | admin/internal/auth/ |
| `features/*/internal/` | 功能内部实现 | 功能内部分层 | reports/internal/business/ |
| `shared/domain/` | 共享模型 | 跨端核心实体 | user.go, product.go |
| `shared/repositories/` | 共享仓储 | 跨端数据访问 | user_repository.go |
| `shared/services/` | 共享服务 | 跨端业务服务 | order_service.go |
| `shared/contracts/dtos/` | 共享 DTO | 跨端 API 一致性 | user_response.go |

---

## 🎨 组织模式速查

### 根据代码量选择模式

| 代码量 | 推荐模式 | 特点 | 目录结构 |
|--------|---------|------|---------|
| < 200 行 | 模式1：单文件 | 简单直接 | `handler.go` + `service_extensions.go` |
| 200-1000 行 | 模式2：按操作拆分 ⭐ | 结构清晰 | `models.go` + `create_xxx.go` + `list_xxx.go` |
| > 1000 行 | 模式3：内部分层 | 职责分明 | `internal/entity/` + `internal/data/` + `internal/business/` |

### 文件命名规范

| 文件类型 | 命名规则 | 示例 |
|---------|---------|------|
| Handler | `{动词}_{名词}.go` | `create_user.go`, `list_products.go` |
| Models | `models.go` | 功能内共享的数据模型 |
| Controller | `controller.go` | 路由映射 |
| DI 注册 | `service_extensions.go` | DI 容器注册 |
| Store | `{名词}_store.go` 或 `store.go` | `user_store.go` |

---

## 💾 Domain 模型速查

### 模型放置决策

| 场景 | 位置 | 理由 |
|------|------|------|
| 跨端共享的核心实体 | `shared/domain/` | User, Product, Order |
| 单端特有实体 | `apps/*/features/*/models/` | AdminLog, UserSession |
| 通用基础字段 | `shared/domain/common/` | BaseEntity, SoftDelete |
| 功能内部实体 | `features/*/internal/entity/` | ReportEntity（内部使用） |

### Domain vs DTO 对比

| 维度 | Domain 模型 | DTO |
|------|------------|-----|
| **位置** | `shared/domain/` 或 `features/*/models/` | `features/*/models.go` |
| **用途** | 内部业务逻辑、数据库映射 | 外部数据传输、API 响应 |
| **字段** | 完整字段（包含敏感） | 只有必要字段 |
| **标签** | `gorm` 标签为主 | `json` + `binding` 标签 |
| **验证** | 业务规则验证 | 输入格式验证 |
| **敏感数据** | 包含 Password, Salt 等 | 不包含敏感字段 |

---

## 📋 DTO 组织速查

### DTO 数量决策表

| DTO 数量 | 推荐方案 | 位置 | 优势 |
|---------|---------|------|------|
| < 3 个 | 与操作一起 | `create_user.go` 中定义 | 功能内聚 |
| 3-10 个 | 功能内 models.go ⭐ | `features/xxx/models.go` | 复用方便 |
| > 10 个 | requests/responses 分类 | `requests/` + `responses/` | 分类清晰 |

### DTO 共享策略

| 共享范围 | 位置 | 适用场景 |
|---------|------|---------|
| 不共享 | `features/xxx/models.go` | 功能独有 DTO |
| 端内共享 | `apps/*/shared/dtos/` | 应用内多功能共享 |
| 跨端共享 | `shared/contracts/dtos/` | 保证 API 一致性 |

### DTO 命名规范

| 类型 | 命名格式 | 示例 |
|------|---------|------|
| Request | `{动词}{实体}Request` | `CreateUserRequest`, `UpdateProductRequest` |
| Response | `{实体}Response` | `UserResponse`, `OrderDetailResponse` |
| List Item | `{实体}ListItem` | `UserListItem`, `ProductListItem` |

---

## 🔒 internal 使用速查

### internal 层次对照

| 层次 | 位置 | 可见性 | 适用场景 | 示例 |
|------|------|--------|---------|------|
| 功能级 | `features/*/internal/` | 仅当前功能 | 功能内部实现 | `reports/internal/business/` |
| 应用级 | `apps/*/internal/` | 仅当前应用 | 应用特定逻辑 | `admin/internal/auth/` |
| 全局 | `shared/` | 所有应用 | 跨端共享 | `shared/domain/` |

### internal 内部结构

| 子目录 | 用途 | 示例 |
|--------|------|------|
| `internal/entity/` | 内部实体/领域对象 | `report_entity.go` |
| `internal/data/` | 数据访问层 | `report_store.go` |
| `internal/business/` | 业务逻辑层 | `report_generator.go` |

---

## 🔄 共享策略速查

### 共享层次决策

```
不共享
    ↓
功能级共享 (features/*/internal/)
    ↓
应用级共享 (apps/*/internal/)
    ↓
全局共享 (shared/)
```

### 共享策略对比表

| 层次 | 位置 | 可见性 | 适用场景 | 示例 |
|------|------|--------|---------|------|
| 功能私有 | `features/xxx/` | 仅当前功能 | 功能特有代码 | `create_order.go` |
| 功能级共享 | `features/xxx/internal/` | 仅当前功能 | 功能内部复用 | `report_generator.go` |
| 应用级共享 | `apps/*/internal/` | 仅当前应用 | 应用特定逻辑 | `admin_auth.go` |
| 全局共享 | `shared/` | 所有应用 | 核心业务代码 | `user.go` |

---

## 🎯 常见场景速查

### 场景1：新建简单 CRUD 功能

```
选择：模式1（单文件实现）
目录：features/categories/
文件：handler.go + service_extensions.go
```

### 场景2：新建中等复杂功能（多个操作）

```
选择：模式2（按操作拆分）⭐
目录：features/orders/
文件：models.go + create_order.go + list_orders.go + controller.go
```

### 场景3：新建复杂功能（需要内部分层）

```
选择：模式3（内部分层）
目录：features/reports/
结构：internal/entity/ + internal/data/ + internal/business/ + handlers
```

### 场景4：Domain 模型放置

| 问题 | 答案 |
|------|------|
| User 实体多端使用？ | → `shared/domain/user.go` |
| AdminLog 只有管理端用？ | → `apps/admin/features/logs/models/admin_log.go` |
| BaseEntity 所有表都用？ | → `shared/domain/common/base_entity.go` |

### 场景5：DTO 放置

| 问题 | 答案 |
|------|------|
| CreateUserRequest 只在创建用户用？ | → `features/users/models.go` |
| UserResponse 多端返回相同格式？ | → `shared/contracts/dtos/user_response.go` |
| Pagination 端内多功能用？ | → `apps/admin/shared/dtos/pagination.go` |

### 场景6：何时使用 internal

| 问题 | 答案 |
|------|------|
| 报表功能需要内部分层？ | → `features/reports/internal/` |
| Admin 端多功能需要权限检查？ | → `apps/admin/internal/auth/` |
| 核心用户实体需要共享？ | → `shared/domain/user.go`（不用 internal） |

---

## 📊 决策流程图（简化版）

### 代码放置决策

```
1. 需要跨端共享？
   是 → shared/
   否 → 继续

2. 需要应用内共享？
   是 → apps/*/internal/
   否 → 继续

3. 需要功能内复用？
   是 → features/*/internal/
   否 → features/*/（直接放）
```

### DTO 组织决策

```
1. DTO 数量？
   < 3 个 → 与操作一起
   3-10 个 → models.go ⭐
   > 10 个 → requests/responses/

2. 需要共享？
   跨端 → shared/contracts/dtos/
   端内 → apps/*/shared/dtos/
   不需要 → 功能内
```

---

## ⚡ 快速命令

### 创建新功能（模式2）

```bash
mkdir -p apps/admin/features/tags
cd apps/admin/features/tags

# 创建文件
touch models.go
touch create_tag.go
touch list_tags.go
touch update_tag.go
touch delete_tag.go
touch controller.go
touch service_extensions.go
```

### 创建内部分层功能（模式3）

```bash
mkdir -p apps/admin/features/reports/internal/{entity,data,business}
cd apps/admin/features/reports

# 创建目录结构
touch models.go
touch internal/entity/report_entity.go
touch internal/data/report_store.go
touch internal/business/report_generator.go
touch generate_report.go
touch controller.go
touch service_extensions.go
```

---

## 🚨 常见错误速查

| 错误 | 正确做法 |
|------|---------|
| ❌ 一开始就放 shared | ✅ 先放 features，确认需要时再提取 |
| ❌ 所有东西都放 internal | ✅ 只有内部实现放 internal，DTO/Handler 在外面 |
| ❌ Domain 模型暴露敏感字段 | ✅ 通过 DTO 转换，隐藏敏感信息 |
| ❌ DTO 命名不清晰 | ✅ 使用 CreateXxxRequest, XxxResponse 格式 |
| ❌ 简单功能过度分层 | ✅ 根据复杂度选择合适的模式 |

---

## 🔍 术语对照表

| 术语 | 英文 | 解释 |
|------|------|------|
| 垂直切片 | Vertical Slice | 按功能垂直组织代码 |
| 领域模型 | Domain Model | 业务实体的抽象，对应数据库表 |
| 数据传输对象 | DTO (Data Transfer Object) | API 层传输数据的对象 |
| 仓储 | Repository | 数据访问层的抽象 |
| 聚合根 | Aggregate Root | 一组关联实体的根实体 |
| 依赖注入 | DI (Dependency Injection) | 控制反转的一种实现 |

---

## 📚 相关文档链接

- [项目结构组织](../01-project-structure.md)
- [Domain 模型组织](../02-domain-models.md)
- [DTO 模型组织](../03-dto-models.md)
- [internal 目录使用](../04-internal-directory.md)
- [共享策略指南](../05-sharing-strategies.md)
- [完整项目示例](./complete-example.md)
- [检查清单](./checklist.md)
- [决策树汇总](./decision-trees.md)

---

**提示：** 将本页面加入书签，以便快速查阅！
