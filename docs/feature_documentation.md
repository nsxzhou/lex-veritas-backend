# LexVeritas 功能文档

> **项目定位**：基于区块链可信存证与检索增强生成（vRAG）的法律智能问答系统

---

## 1. 项目核心理念

**"无依据，不回答"（No Citation, No Answer）**

LexVeritas 追求的不是 AI 的自由发挥，而是严谨的法律回答。系统确保：

- 每一条法律建议都源自经过确权的权威数据
- 每一个引用都可追溯到区块链存证
- 彻底消除数据篡改和 AI 幻觉的风险

---

## 2. 功能模块概览

| 模块           | 描述           | 核心能力                     |
| :------------- | :------------- | :--------------------------- |
| **用户管理**   | 多角色用户系统 | 注册/登录、OAuth、额度管理   |
| **智能问答**   | 法律 AI 对话   | 会话管理、流式回复、法条引用 |
| **知识库管理** | 法律文档处理   | 上传、分块、向量化、存证     |
| **区块链存证** | 可信验证层     | Merkle Tree、链上验证        |
| **审计日志**   | 系统监控       | 操作追踪、安全告警           |
| **Token 计量** | 资源管理       | 用量统计、配额控制           |

---

## 3. 用户管理模块

### 3.1 数据模型

```
User（用户）
├── 基础信息：Email、Phone、Name、Avatar
├── 认证：PasswordHash
├── 角色：user / admin / super_admin
├── 状态：active / inactive / banned
├── 额度：TokenQuota（配额）、TokenUsed（已用）
└── 关联：OAuthAccounts（第三方登录）
```

### 3.2 功能清单

| 功能              | 描述                                |
| :---------------- | :---------------------------------- |
| **本地注册/登录** | 邮箱 + 密码认证                     |
| **OAuth 登录**    | 支持第三方平台（Google, GitHub 等） |
| **角色权限**      | 普通用户、管理员、超级管理员        |
| **Token 额度**    | 默认 10 万 Tokens，可调整           |
| **账户状态管理**  | 激活、禁用、封禁                    |

---

## 4. 智能问答模块

### 4.1 数据模型

```
ChatSession（会话）
├── 用户归属：UserID / GuestSessionID（支持匿名）
├── 会话信息：Title、Summary
├── 统计：MessageCount、TokensUsed
└── 时间戳：LastMessageAt

ChatMessage（消息）
├── 角色：user / assistant / system
├── 内容：Content
├── Token 统计：TokensIn、TokensOut
├── 引用：Citations（法条引用列表）
└── 元数据：Metadata（JSONB）

MessageCitation（法条引用）
├── 引用内容：Text、Source、ArticleNumber
├── 法律层级：LawHierarchy
└── 区块链验证：ChunkHash、VerificationID、BlockNumber、Verified
```

### 4.2 功能清单

| 功能           | 描述                       |
| :------------- | :------------------------- |
| **多轮对话**   | 支持上下文连续对话         |
| **匿名访问**   | 游客无需登录即可体验       |
| **法条引用**   | 回答自动标注法律来源       |
| **链上验证**   | 引用内容可追溯至区块链存证 |
| **Token 计量** | 精确统计输入/输出 Token    |
| **会话管理**   | 创建、查看、删除会话       |

### 4.3 问答流程

```
用户提问
    ↓
语义检索（Milvus 向量搜索）
    ↓
引用扩展（加载关联法条）
    ↓
区块链验证（校验数据完整性）
    ↓
LLM 生成回答（严格 Prompt 约束）
    ↓
引用一致性校验
    ↓
返回带引用的回答
```

---

## 5. 知识库管理模块

### 5.1 数据模型

```
Document（文档）
├── 基础信息：Name、Type、Size、FilePath
├── 法律元数据：LawName、LawType、EffectiveDate、PublishOrg
├── 处理状态：pending / processing / indexed / minted / error
├── 分块统计：ChunkCount
└── 存证信息：IsMinted、MintedAt、VersionID

DocumentChunk（文档分块）
├── 内容：Content、ContentHash
├── 法律结构：ChunkOrder、LawHierarchy、ArticleNumber
├── 引用关系：References（JSONB）
├── Merkle 验证：MerkleIndex、MerkleProof
├── 版本：VersionID
└── 向量化状态：IsEmbedded、EmbeddedAt
```

### 5.2 支持的文档类型

| 类型     | 格式     | 说明               |
| :------- | :------- | :----------------- |
| PDF      | `.pdf`   | 主流法律文档格式   |
| DOCX     | `.docx`  | Word 文档          |
| TXT      | `.txt`   | 纯文本             |
| Markdown | `.md`    | 结构化文本（推荐） |
| URL      | 网页链接 | 在线法规抓取       |

### 5.3 文档处理流程

```
上传文档
    ↓
格式转换（→ Markdown）
    ↓
结构化分块（以法条为原子单位）
    ↓
元数据提取（法律层级、引用关系）
    ↓
内容哈希计算
    ↓
向量化嵌入（OpenAI Embeddings）
    ↓
存入向量数据库（Milvus）
    ↓
构建 Merkle Tree
    ↓
上链存证（可选）
```

### 5.4 结构化分块策略

采用**法律结构感知分块**，而非固定字符分块：

| 层级   | 正则模式          | 说明             |
| :----- | :---------------- | :--------------- |
| 编     | `第[一二三...]编` | 最高层级         |
| 章     | `第[一二三...]章` | 章节划分         |
| 节     | `第[一二三...]节` | 小节划分         |
| **条** | `第[一二三...]条` | **最小原子单位** |

**优势**：

- 100% 保持法条完整性
- 保留完整的法律层级上下文
- 精准的引用关系追踪

---

## 6. 区块链存证模块

### 6.1 数据模型

```
KnowledgeVersion（知识库版本）
├── Merkle Tree：MerkleRoot、ChunkCount
├── 版本描述：Description
├── 链上信息：TxHash、BlockNumber
├── 状态：pending / confirmed
└── 关联文档：Documents（多对多）

ProofRecord（验证记录）
├── 验证主体：ChunkID、DocumentID
├── 验证数据：LeafHash、ComputedRoot、OnChainRoot
├── 验证结果：Verified、VersionID
├── 链上信息：BlockNumber、TxHash
└── 触发来源：TriggerType、TriggerBy
```

### 6.2 核心机制

#### Merkle Tree 结构

```
                    Merkle Root (上链)
                   /              \
            Hash(AB)              Hash(CD)
           /       \            /        \
      Hash(A)     Hash(B)   Hash(C)     Hash(D)
         |           |         |           |
      Chunk 1    Chunk 2    Chunk 3    Chunk 4
```

#### 验证流程

1. **存证**：计算所有 Chunk 的 Merkle Root → 上链存储
2. **验证**：
   - 计算目标 Chunk 的哈希
   - 使用 Merkle Proof 计算根哈希
   - 与链上 Merkle Root 比对
   - 一致则验证通过 ✅

### 6.3 功能清单

| 功能                 | 描述                         |
| :------------------- | :--------------------------- |
| **版本管理**         | 知识库多版本追踪             |
| **Merkle Tree 构建** | 自动计算文档集的 Merkle Root |
| **链上存证**         | 将 Merkle Root 写入智能合约  |
| **实时验证**         | 问答时自动验证引用数据完整性 |
| **验证记录**         | 保存所有验证历史             |

---

## 7. 审计日志模块

### 7.1 数据模型

```
AuditLog（审计日志）
├── 事件类型：tamper / access / verify / upload / mint / system / auth
├── 严重程度：low / medium / high
├── 事件描述：Message、Details（JSONB）
├── 来源信息：Source、SourceType
├── 用户信息：UserID、UserIP、UserAgent
└── 处理状态：unresolved / investigating / resolved
```

### 7.2 审计事件类型

| 类型     | 描述         | 严重程度  |
| :------- | :----------- | :-------- |
| `tamper` | 数据篡改检测 | 🔴 High   |
| `auth`   | 认证相关事件 | 🟡 Medium |
| `verify` | 验证操作记录 | 🟢 Low    |
| `upload` | 文档上传记录 | 🟢 Low    |
| `mint`   | 存证操作记录 | 🟡 Medium |
| `access` | 访问记录     | 🟢 Low    |
| `system` | 系统事件     | 🟡 Medium |

### 7.3 功能清单

| 功能         | 描述                       |
| :----------- | :------------------------- |
| **事件记录** | 自动记录关键操作           |
| **篡改告警** | 检测到数据不一致时立即告警 |
| **日志查询** | 按时间、用户、类型筛选     |
| **问题追踪** | 标记、调查、解决问题       |

---

## 8. Token 计量模块

### 8.1 数据模型

```
TokenUsage（使用记录）
├── 归属：UserID、SessionID
├── API 细分：APIEndpoint、APIMethod
├── Token 统计：TokensIn、TokensOut、TotalTokens
├── 模型信息：Model、ModelType
├── 费用估算：EstimatedCost
└── 请求信息：RequestID、Latency、Success

TokenUsageDaily（日统计）
├── 维度：UserID、APIEndpoint、Model、Date
└── 聚合：RequestCount、TokensIn、TokensOut、TotalCost
```

### 8.2 功能清单

| 功能          | 描述                           |
| :------------ | :----------------------------- |
| **实时计量**  | 记录每次 API 调用的 Token 消耗 |
| **用户配额**  | 限制单用户 Token 使用量        |
| **费用估算**  | 根据模型定价估算成本           |
| **日报表**    | 按日汇总统计数据               |
| **Dashboard** | 可视化用量展示                 |

---

## 9. 系统配置

### 9.1 数据模型

```
SystemConfig（系统配置）
├── 配置项：Key（唯一）、Value
├── 类型：Type（string/int/bool/json）
├── 分类：Category
├── 描述：Description
└── 更新信息：UpdatedBy、UpdatedAt
```

### 9.2 配置分类

| 分类         | 示例配置项                 |
| :----------- | :------------------------- |
| `llm`        | 模型名称、温度、最大 Token |
| `retrieval`  | Top-K 数量、相似度阈值     |
| `blockchain` | 合约地址、网络配置         |
| `quota`      | 默认用户配额、VIP 配额     |
| `security`   | JWT 过期时间、速率限制     |

---

## 10. 技术架构

### 10.1 技术栈

| 层级       | 技术                  |
| :--------- | :-------------------- |
| **后端**   | Go + Gin + Eino       |
| **数据库** | PostgreSQL (业务数据) |
| **向量库** | Milvus (语义检索)     |
| **缓存**   | Redis (会话、验证码)  |
| **区块链** | Polygon (智能合约)    |
| **LLM**    | OpenAI / 其他兼容 API |

### 10.2 目录结构

```
lex-veritas-backend/
├── cmd/server/          # 入口
│   └── main.go
├── internal/
│   ├── config/          # 配置管理
│   ├── handler/         # HTTP 处理器
│   ├── service/         # 业务逻辑
│   ├── repository/      # 数据访问
│   ├── dto/             # 数据传输对象
│   ├── model/           # 数据模型
│   ├── graph/           # Eino RAG 编排
│   └── client/          # 外部服务客户端
├── docs/                # 文档
└── config.yaml          # 配置文件
```

---

## 11. 核心价值

| 特性           | 传统法律 AI  | LexVeritas            |
| :------------- | :----------- | :-------------------- |
| **数据来源**   | 可能被篡改   | 区块链存证，不可篡改  |
| **AI 幻觉**    | 可能捏造法条 | 严格引用验证，零幻觉  |
| **可追溯性**   | 来源不明     | 每条引用可追溯至链上  |
| **法条完整性** | 可能截断     | 结构化分块，100% 完整 |
| **上下文理解** | 可能丢失     | 保留完整法律层级      |

---

## 12. 未来规划

- [ ] 图数据库集成（Neo4j）- 复杂引用关系
- [ ] 法条版本溯源 - 时间点查询
- [ ] 法律冲突检测 - 自动识别矛盾法条
- [ ] 多语言支持 - 中英双语法律咨询
- [ ] 移动端适配 - iOS/Android 应用
