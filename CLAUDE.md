# CLAUDE.md

本文档记录 LexVeritas 后端项目的开发规范和约定，供 AI 助手和开发者参考。

---

## 项目结构

```
lex-veritas-backend/
├── cmd/server/          # 程序入口
├── internal/            # 内部代码 (不对外暴露)
│   ├── dto/            # 请求/响应结构体
│   ├── service/        # 业务逻辑层
│   ├── repository/     # 数据访问层
│   ├── handler/        # HTTP 处理器
│   ├── middleware/     # 中间件
│   ├── model/          # 数据模型
│   ├── config/         # 配置管理
│   ├── router/         # 路由配置
│   ├── graph/          # Eino Graph 编排
│   ├── client/         # 外部服务客户端
│   └── pkg/            # 通用工具包
│       ├── auth/       # JWT、密码工具
│       ├── cache/      # Redis 缓存
│       ├── database/   # PostgreSQL 数据库
│       ├── errors/     # 错误处理
│       ├── logger/     # 日志
│       └── response/   # HTTP 响应
├── docs/               # 文档
└── config.yaml         # 配置文件
```

## 四层架构

| 层级       | 目录          | 职责                                  |
| ---------- | ------------- | ------------------------------------- |
| DTO        | `dto/`        | 请求/响应结构体定义                   |
| Handler    | `handler/`    | HTTP 参数绑定、调用 service、返回响应 |
| Service    | `service/`    | 业务逻辑、事务协调                    |
| Repository | `repository/` | 数据库 CRUD 操作                      |

---

## 缓存规范 (`internal/pkg/cache`)

### 文件说明

| 文件       | 用途                                      |
| ---------- | ----------------------------------------- |
| `cache.go` | Redis 连接和基础操作 (Get/Set/Delete 等)  |
| `keys.go`  | **所有 Redis 键前缀和配置常量统一定义处** |

### 使用规范

1. **所有 Redis 键前缀必须在 `keys.go` 中定义**，禁止在业务代码中硬编码
2. 使用辅助函数生成完整键名，如 `cache.AuthRefreshTokenKey(tokenHash)`
3. 相关配置常量 (如 TTL、最大尝试次数) 也放在 `keys.go` 中

### 示例

```go
// ✅ 正确用法
key := cache.AuthRefreshTokenKey(tokenHash)
cache.SetObject(ctx, key, data, cache.AuthLockoutDuration)

// ❌ 错误用法 - 禁止硬编码键前缀
key := "auth:refresh:" + tokenHash
```

### 键前缀命名规范

格式：`{模块}:{子类型}:{标识符}`

```
auth:refresh:{tokenHash}     # 刷新令牌
auth:blacklist:{jti}         # Token 黑名单
auth:attempts:{identifier}   # 登录尝试计数
session:user:{userId}        # 用户会话
ratelimit:ip:{ip}            # IP 限流
```

---

## 响应规范 (`internal/pkg/response`)

### 统一响应结构

```go
type Response struct {
    Code    int         `json:"code"`    // 业务错误码
    Message string      `json:"message"` // 错误信息
    Data    interface{} `json:"data"`    // 响应数据
}
```

### 常用响应函数

```go
response.Success(c, data)                    // 成功响应
response.SuccessWithMessage(c, "msg", data)  // 带消息的成功响应
response.BadRequest(c, "参数错误")            // 400 参数错误
response.Unauthorized(c, "未登录")            // 401 未授权
response.Forbidden(c, "权限不足")             // 403 禁止访问
response.NotFound(c)                         // 404 资源不存在
response.InternalError(c, err)               // 500 服务器错误
```

---

## 错误处理规范 (`internal/pkg/errors`)

### 错误码分类

| 范围 | 类型     | 示例                                   |
| ---- | -------- | -------------------------------------- |
| 0    | 成功     | `CodeSuccess`                          |
| 1xxx | 通用错误 | `CodeUnauthorized`, `CodeForbidden`    |
| 2xxx | 用户相关 | `CodeUserNotFound`, `CodeInvalidToken` |
| 3xxx | 业务错误 | `CodeChatSessionNotFound`              |
| 5xxx | 系统错误 | `CodeDatabaseError`, `CodeRedisError`  |

### 使用方式

```go
// 创建错误
errors.New(errors.CodeUserNotFound)
errors.NewWithMessage(errors.CodeInvalidParam, "邮箱格式错误")

// 包装原始错误
errors.Wrap(errors.CodeDatabaseError, err)

// 在 handler 中使用
response.Error(c, errors.New(errors.CodeUserNotFound))
```

---

## 日志规范 (`internal/pkg/logger`)

### 使用方式

```go
import "github.com/lexveritas/lex-veritas-backend/internal/pkg/logger"

// 直接使用
logger.Info("用户登录", zap.String("userId", id))
logger.Error("数据库连接失败", zap.Error(err))

// 带 RequestID
logger.WithRequestID(requestID).Info("处理请求")

// 从 context 获取
logger.WithContext(ctx).Info("处理请求")
```

### 日志级别

| 级别  | 用途                        |
| ----- | --------------------------- |
| Debug | 调试信息，生产环境关闭      |
| Info  | 关键业务事件 (登录、注册等) |
| Warn  | 异常但可恢复的情况          |
| Error | 错误，需要关注              |
| Fatal | 致命错误，程序将退出        |

---

## 数据模型规范 (`internal/model`)

### 文件组织

| 文件          | 内容                                   |
| ------------- | -------------------------------------- |
| `types.go`    | 枚举类型定义 (UserRole, UserStatus 等) |
| `user.go`     | 用户和 OAuth 账户模型                  |
| `chat.go`     | 聊天会话和消息模型                     |
| `document.go` | 文档模型                               |
| `system.go`   | 审计日志、系统配置                     |

### GORM 标签规范

```go
type User struct {
    ID        string         `json:"id" gorm:"type:varchar(36);primaryKey"`
    Email     string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
    Status    UserStatus     `json:"status" gorm:"type:varchar(20);default:'active'"`
    CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`  // 软删除
}
```

---

## Handler 规范 (`internal/handler`)

### 结构

```go
type AuthHandler struct {
    authSvc auth.AuthService  // 依赖注入
}

func NewAuthHandler(authSvc auth.AuthService) *AuthHandler {
    return &AuthHandler{authSvc: authSvc}
}
```

### 方法命名

- `Login` - 登录
- `Register` - 注册
- `GetByID` - 按 ID 获取
- `List` - 列表查询
- `Create` - 创建
- `Update` - 更新
- `Delete` - 删除

---

## 中间件规范 (`internal/middleware`)

### 现有中间件

| 文件           | 功能                         |
| -------------- | ---------------------------- |
| `auth.go`      | JWT 认证、角色验证           |
| `guest.go`     | 匿名用户限制 (1 会话/5 对话) |
| `quota.go`     | 额度检查、RequireAdmin       |
| `cors.go`      | CORS 跨域                    |
| `logger.go`    | 请求日志                     |
| `ratelimit.go` | 限流                         |
| `recovery.go`  | Panic 恢复                   |

### 从上下文获取用户信息

```go
userID := middleware.GetUserID(c)
role := middleware.GetRole(c)
claims := middleware.GetClaims(c)
guestID := middleware.GetGuestID(c)  // 匿名用户
```

---

## 用户分级系统

### 三级用户模型

| 类型         | 存储         | 会话 | 对话      | 查看统计 | 上传文档 |
| ------------ | ------------ | ---- | --------- | -------- | -------- |
| **匿名用户** | Cookie+Redis | 1 个 | 5 次/会话 | ❌       | ❌       |
| **登录用户** | 数据库       | 不限 | 按额度    | 自己     | ❌       |
| **管理员**   | 数据库       | 不限 | 不限      | 全部     | ✅       |

### User 模型额度字段

```go
type User struct {
    // ...
    TokenQuota int64 `gorm:"default:100000"` // 默认额度
    TokenUsed  int64 `gorm:"default:0"`      // 已使用
}
```

### 中间件使用

```go
// 匿名用户限制
router.Use(middleware.GuestLimit())

// 额度检查 (登录用户)
router.Use(middleware.QuotaCheck())

// 管理员权限
router.Use(middleware.RequireAdmin())
```

### 额度服务 (`internal/user/quota.go`)

```go
quotaSvc := user.NewQuotaService()

// 检查/消费额度
quotaSvc.CheckQuota(ctx, userID, 1000)
quotaSvc.ConsumeTokens(ctx, userID, 500)

// 调整额度 (admin)
quotaSvc.AdjustQuota(ctx, userID, 200000)

// 获取统计
quotaSvc.GetUsage(ctx, userID)
quotaSvc.GetAllUsage(ctx, page, pageSize)  // admin
```

---

## 认证模块规范 (`internal/auth`)

### 双令牌策略

| Token         | 存储   | 有效期  | 用途              |
| ------------- | ------ | ------- | ----------------- |
| Access Token  | 客户端 | 15 分钟 | 无状态请求认证    |
| Refresh Token | Redis  | 7 天    | 刷新 Access Token |

### 文件结构

| 文件          | 功能                   |
| ------------- | ---------------------- |
| `jwt.go`      | JWT 生成/验证          |
| `password.go` | 密码哈希 (bcrypt)      |
| `service.go`  | AuthService 接口和实现 |

---

## 配置规范 (`config.yaml`)

### 环境变量覆盖

格式：`LEXVERITAS_{SECTION}_{KEY}`

```bash
LEXVERITAS_DATABASE_PASSWORD=xxx
LEXVERITAS_JWT_SECRET=xxx
LEXVERITAS_REDIS_PASSWORD=xxx
```

---

## 代码风格

1. **YAGNI** - 不写当前不需要的代码
2. **KISS** - 选择最简单的方案
3. **DRY** - 不重复自己，提取公共逻辑
4. **中文注释** - 关键业务逻辑使用中文注释
5. **错误处理** - 始终检查并处理错误，不忽略
