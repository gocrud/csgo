# 最佳实践

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

构建高质量API的最佳实践。

## API设计原则

### 1. RESTful 设计

```
GET    /api/users      - 获取列表
GET    /api/users/:id  - 获取单个
POST   /api/users      - 创建
PUT    /api/users/:id  - 更新
DELETE /api/users/:id  - 删除
```

### 2. 统一响应格式

```json
{
  "success": true,
  "data": {}
}
```

### 3. 错误处理

```json
{
  "success": false,
  "error": {
    "code": "USER.NOT_FOUND",
    "message": "用户不存在"
  }
}
```

### 4. 版本控制

```
/api/v1/users
/api/v2/users
```

### 5. 分页和过滤

```
GET /api/users?page=1&size=20&sort=created_at&order=desc
```

## 代码组织

```
project/
├── controllers/  # 控制器
├── services/     # 业务逻辑
├── models/       # 数据模型
├── validators/   # 验证器
└── main.go
```

## 下一步

继续学习：[实践项目：CRUD API](project-crud-api.md) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

