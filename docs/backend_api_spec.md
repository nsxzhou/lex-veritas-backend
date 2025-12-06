# LexVeritas 后端开发 API 规范

> 本文档定义了前后端对接的 API 接口规范，按开发优先级排列

---

## 技术栈建议

| 组件       | 技术方案                      | 说明                |
| :--------- | :---------------------------- | :------------------ |
| Web 框架   | Go + Gin                      | 高性能 HTTP 服务    |
| RAG 编排   | Eino (CloudWeGo)              | Graph 式 Chain 编排 |
| 向量数据库 | Milvus 2.3+                   | 法条语义检索        |
| 区块链     | Polygon Amoy                  | Merkle Root 存证    |
| Embedding  | OpenAI text-embedding-3-small | 1536 维向量         |
| LLM        | OpenAI GPT-4                  | 答案生成            |

---

## P0 - MVP 核心 API

### 1. 用户认证 `/api/auth`

#### POST `/api/auth/login`

邮箱密码登录

**请求体：**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**响应：**

```json
{
  "token": "jwt_token_here",
  "user": {
    "id": "user_123",
    "name": "张三",
    "email": "user@example.com",
    "role": "user",
    "status": "active"
  }
}
```

#### POST `/api/auth/login/phone`

手机验证码登录

**请求体：**

```json
{
  "phone": "13800138000",
  "code": "123456"
}
```

#### POST `/api/auth/send-code`

发送验证码

**请求体：**

```json
{
  "phone": "13800138000"
}
```

#### POST `/api/auth/register`

用户注册

**请求体：**

```json
{
  "name": "张三",
  "email": "user@example.com",
  "password": "password123"
}
```

#### GET `/api/auth/me`

获取当前用户信息（需认证）

---

### 2. 智能问答 `/api/chat`

#### POST `/api/chat` (SSE 流式)

发送问题，返回流式响应

**请求体：**

```json
{
  "question": "民法典关于物权的规定是什么？",
  "sessionId": "session_123"
}
```

**SSE 响应事件：**

```
data: {"type": "text", "content": "根据《"}
data: {"type": "text", "content": "中华人民共和国民法典》"}
data: {"type": "citation", "citation": {...}}
data: {"type": "done"}
```

**Citation 结构：**

```json
{
  "id": "cit_1",
  "text": "第一百二十三条 民事主体依法享有物权...",
  "source": "中华人民共和国民法典",
  "verificationId": "0x7f...3a2b",
  "blockNumber": 18234567,
  "timestamp": "2023-10-27T10:00:00Z",
  "metadata": {
    "law_name": "中华人民共和国民法典",
    "part": "第一编 总则",
    "chapter": "第五章 民事权利",
    "article_number": "第123条"
  }
}
```

#### GET `/api/chat/sessions`

获取会话列表

#### GET `/api/chat/sessions/:id`

获取会话详情（含历史消息）

#### POST `/api/chat/sessions`

创建新会话

#### DELETE `/api/chat/sessions/:id`

删除会话

---

### 3. 知识库管理 `/api/knowledge`

#### GET `/api/knowledge/documents`

获取文档列表

**查询参数：**

- `type`: pdf | docx | txt | url
- `status`: indexed | processing | error | minted
- `search`: 搜索关键词

**响应：**

```json
[
  {
    "id": "doc_1",
    "name": "劳动合同法_2024修订版.pdf",
    "type": "pdf",
    "size": "2.4 MB",
    "uploadDate": "2024-03-15",
    "status": "indexed",
    "uploadedBy": "Admin"
  }
]
```

#### POST `/api/knowledge/documents`

上传文档（multipart/form-data）

**表单字段：**

- `file`: 文件
- `type`: pdf | docx | txt | url
- `url`: (可选) URL 类型时使用

#### DELETE `/api/knowledge/documents/:id`

删除文档

#### POST `/api/knowledge/documents/mint`

批量上链

**请求体：**

```json
{
  "documentIds": ["doc_1", "doc_2"]
}
```

**响应：**

```json
{
  "merkleRoot": "0x8f3...2a9",
  "txHash": "0x7e2...1b8",
  "blockNumber": 18234567,
  "versionId": 1
}
```

---

## P1 - 管理功能 API

### 4. Dashboard 统计 `/api/stats`

#### GET `/api/stats/overview`

获取总览数据

**响应：**

```json
{
  "totalQueries": 128930,
  "activeUsers": 8200,
  "knowledgeItems": 1200000,
  "systemHealth": 99.9
}
```

#### GET `/api/stats/query-volume`

获取提问量趋势

**查询参数：**

- `period`: 24h | 7d | 30d

#### GET `/api/stats/recent-activity`

获取最近活动

---

### 5. 区块链存证 `/api/proof`

#### GET `/api/proof/merkle-tree`

获取 Merkle Tree 结构

**响应：**

```json
{
  "root": "0x8f3...2a9",
  "levels": [
    [{"hash": "0x8f3...2a9", "type": "root", "status": "verified"}],
    [{"hash": "0x4a1...b2c", "type": "node", "status": "verified"}, {...}],
    [{"hash": "0x1b2...3c4", "type": "leaf", "status": "verified", "data": "Doc A"}, {...}]
  ]
}
```

#### GET `/api/proof/recent`

获取最近存证记录

#### POST `/api/proof/verify/:chunkId`

验证 Chunk 链上完整性

**响应：**

```json
{
  "verified": true,
  "chunkHash": "0x1b2...3c4",
  "onChainRoot": "0x8f3...2a9",
  "computedRoot": "0x8f3...2a9",
  "blockNumber": 18234567,
  "txHash": "0x7f...3a2b"
}
```

#### GET `/api/proof/stats`

获取存证统计

---

### 6. 审计日志 `/api/audit`

#### GET `/api/audit/logs`

获取审计日志

**查询参数：**

- `page`: 页码
- `limit`: 每页数量
- `type`: tamper | access | verify | system
- `severity`: high | medium | low
- `status`: unresolved | investigating | resolved

**响应：**

```json
{
  "logs": [
    {
      "id": 1,
      "type": "tamper",
      "severity": "high",
      "message": "检测到数据篡改尝试",
      "source": "Doc D (0x3c4...5d6)",
      "timestamp": "2024-03-20 15:45:22",
      "status": "unresolved"
    }
  ],
  "total": 100,
  "page": 1,
  "limit": 20
}
```

#### POST `/api/audit/logs/:id/resolve`

标记日志已处理

---

## P2 - 扩展功能 API

### 7. 用户管理 `/api/users`

#### GET `/api/users`

获取用户列表

**查询参数：**

- `page`, `limit`, `role`, `status`, `search`

#### GET `/api/users/:id`

获取用户详情

#### PUT `/api/users/:id`

更新用户信息

#### DELETE `/api/users/:id`

删除用户

#### GET `/api/users/:id/history`

获取用户咨询历史

#### GET `/api/users/:id/token-usage`

获取 Token 使用趋势

---

### 8. 系统设置 `/api/settings`

#### GET `/api/settings`

获取系统配置

#### PUT `/api/settings`

更新系统配置

---

## 通用规范

### 认证

- 所有需认证接口在 Header 中携带 `Authorization: Bearer <token>`

### 错误响应

```json
{
  "error": "错误描述",
  "code": "ERROR_CODE"
}
```

### 状态码

- `200` 成功
- `201` 创建成功
- `400` 请求参数错误
- `401` 未认证
- `403` 权限不足
- `404` 资源不存在
- `500` 服务器错误
