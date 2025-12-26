# numx

数字类型扩展包，提供 JSON 安全的大数字类型，解决 JavaScript 精度丢失问题。

## 问题背景

JavaScript 的 `Number.MAX_SAFE_INTEGER` 为 `9007199254740991`（2^53 - 1），超过这个范围的整数会丢失精度。后端常用的 `int64` 类型范围远大于此，导致前后端数据交互时出现精度问题。

## 解决方案

`numx` 包提供了几种类型，在 JSON 序列化时自动转换为字符串，避免精度丢失：

- **ID**: int64 类型别名，用于数据库主键
- **BigInt**: 通用的 int64 包装
- **BigUint**: uint64 包装
- **Timestamp**: int64 时间戳（毫秒）

## 安装

```bash
go get github.com/whl/gocrud/csgo/numx
```

## 使用示例

### 基本使用

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/whl/gocrud/csgo/numx"
)

type User struct {
    ID        numx.ID        `json:"id"`
    Score     numx.BigInt    `json:"score"`
    Views     numx.BigUint   `json:"views"`
    CreatedAt numx.Timestamp `json:"created_at"`
}

func main() {
    user := User{
        ID:        9007199254740992, // 超过 JS 安全整数范围
        Score:     -123456789012345,
        Views:     18446744073709551615, // uint64 最大值
        CreatedAt: numx.Now(),
    }

    // JSON 序列化
    data, _ := json.Marshal(user)
    fmt.Println(string(data))
    // 输出: {"id":"9007199254740992","score":"-123456789012345","views":"18446744073709551615","created_at":"1703234567890"}

    // JSON 反序列化（支持 string 和 number 两种格式）
    jsonStr := `{"id":"123","score":-456,"views":"789","created_at":"1703234567890"}`
    var u User
    json.Unmarshal([]byte(jsonStr), &u)
    fmt.Printf("%+v\n", u)
}
```

### 数据库操作

所有类型都实现了 `sql.Scanner` 和 `driver.Valuer` 接口，可直接用于数据库操作：

```go
type Article struct {
    ID        numx.ID        `db:"id" json:"id"`
    AuthorID  numx.ID        `db:"author_id" json:"author_id"`
    Views     numx.BigUint   `db:"views" json:"views"`
    CreatedAt numx.Timestamp `db:"created_at" json:"created_at"`
}

// 数据库查询
var article Article
db.QueryRow("SELECT id, author_id, views, created_at FROM articles WHERE id = ?", 123).
    Scan(&article.ID, &article.AuthorID, &article.Views, &article.CreatedAt)

// 数据库插入
db.Exec("INSERT INTO articles (id, author_id, views, created_at) VALUES (?, ?, ?, ?)",
    article.ID, article.AuthorID, article.Views, article.CreatedAt)
```

### 类型转换

```go
id := numx.ID(123456)

// 转换为基础类型
i := id.Int64()        // int64
s := id.String()       // "123456"
isZero := id.IsZero()  // false

// Timestamp 特有方法
ts := numx.Now()
t := ts.Time()         // time.Time
ts2 := numx.FromTime(time.Now())
```

### Web API 示例

```go
// 请求
type CreateUserRequest struct {
    Name      string     `json:"name"`
    ParentID  numx.ID    `json:"parent_id"`  // 前端传 "123" 或 123
}

// 响应
type UserResponse struct {
    ID        numx.ID        `json:"id"`         // 输出 "9007199254740992"
    Name      string         `json:"name"`
    CreatedAt numx.Timestamp `json:"created_at"` // 输出 "1703234567890"
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // req.ParentID 自动解析为 int64
    user := User{
        ID:        numx.ID(generateID()),
        Name:      req.Name,
        ParentID:  req.ParentID,
        CreatedAt: numx.Now(),
    }
    
    // 自动序列化为字符串
    json.NewEncoder(w).Encode(UserResponse{
        ID:        user.ID,
        Name:      user.Name,
        CreatedAt: user.CreatedAt,
    })
}
```

## API 文档

### ID

```go
type ID int64

// JSON 序列化方法
func (id ID) MarshalJSON() ([]byte, error)
func (id *ID) UnmarshalJSON(data []byte) error

// 数据库方法
func (id ID) Value() (driver.Value, error)
func (id *ID) Scan(value interface{}) error

// 辅助方法
func (id ID) Int64() int64
func (id ID) String() string
func (id ID) IsZero() bool
```

### BigInt

```go
type BigInt int64

// 方法同 ID，不含 IsZero
```

### BigUint

```go
type BigUint uint64

func (b BigUint) Uint64() uint64
func (b BigUint) String() string
// 其他方法同上
```

### Timestamp

```go
type Timestamp int64

// 额外方法
func (t Timestamp) Time() time.Time
func Now() Timestamp
func FromTime(t time.Time) Timestamp
```

## JSON 序列化示例

| Go 类型 | Go 值 | JSON 输出 | JSON 输入（支持） |
|---------|-------|-----------|------------------|
| `numx.ID` | `123` | `"123"` | `"123"` 或 `123` |
| `numx.BigInt` | `-456` | `"-456"` | `"-456"` 或 `-456` |
| `numx.BigUint` | `789` | `"789"` | `"789"` 或 `789` |
| `numx.Timestamp` | `1703234567890` | `"1703234567890"` | `"1703234567890"` 或 `1703234567890` |

## 注意事项

1. **反序列化兼容性**: 所有类型都支持从 JSON string 和 number 两种格式反序列化，确保前后端兼容
2. **数据库兼容性**: 在数据库中存储为原生的 int64/uint64，无额外开销
3. **类型安全**: 保留所有原生类型的运算和比较能力
4. **零值判断**: ID 类型提供 `IsZero()` 方法，方便判断是否为默认值

## 测试

```bash
cd numx
go test -v
```

## License

MIT
