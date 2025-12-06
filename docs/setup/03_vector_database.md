# 向量数据库配置指南

> Milvus 部署与 Schema 设计

---

## 🐳 Milvus 部署

### Docker Compose (推荐)

`docker-compose.yml`:

```yaml
version: "3.5"

services:
  etcd:
    image: quay.io/coreos/etcd:v3.5.5
    environment:
      - ETCD_AUTO_COMPACTION_MODE=revision
      - ETCD_AUTO_COMPACTION_RETENTION=1000

  minio:
    image: minio/minio:RELEASE.2023-03-20T20-16-18Z
    command: minio server /minio_data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin

  milvus:
    image: milvusdb/milvus:v2.3.0
    command: ["milvus", "run", "standalone"]
    environment:
      ETCD_ENDPOINTS: etcd:23 79
      MINIO_ADDRESS: minio:9000
    ports:
      - "19530:19530"
      - "9091:9091"
    depends_on:
      - etcd
      - minio
```

启动:

```bash
docker-compose up -d
```

---

## 📊 Collection Schema

### Python 完整定义

```python
from pymilvus import CollectionSchema, FieldSchema, DataType, Collection, connections

# 连接 Milvus
connections.connect("default", host="localhost", port="19530")

# 定义字段
fields = [
    # 基础字段
    FieldSchema(name="chunk_id", dtype=DataType.VARCHAR, max_length=64, is_primary=True),
    FieldSchema(name="chunk_text", dtype=DataType.VARCHAR, max_length=8192),
    FieldSchema(name="chunk_hash", dtype=DataType.VARCHAR, max_length=64),
    FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=1536),

    # 区块链验证字段
    FieldSchema(name="merkle_proof", dtype=DataType.VARCHAR, max_length=4096),
    FieldSchema(name="merkle_index", dtype=DataType.INT64),
    FieldSchema(name="version_id", dtype=DataType.INT32),

    # 来源字段
    FieldSchema(name="source_file", dtype=DataType.VARCHAR, max_length=256),
    FieldSchema(name="chunk_order", dtype=DataType.INT32),

    # 法律专用字段
    FieldSchema(name="law_hierarchy", dtype=DataType.VARCHAR, max_length=512),
    FieldSchema(name="references", dtype=DataType.VARCHAR, max_length=2048),
    FieldSchema(name="article_number", dtype=DataType.VARCHAR, max_length=32),
    FieldSchema(name="law_name", dtype=DataType.VARCHAR, max_length=128),
    FieldSchema(name="law_type", dtype=DataType.VARCHAR, max_length=64),
    FieldSchema(name="effective_date", dtype=DataType.VARCHAR, max_length=32),
    FieldSchema(name="priority_score", dtype=DataType.INT32),
]

# 创建 Schema
schema = CollectionSchema(
    fields=fields,
    description="LexVeritas Legal Knowledge Base"
)

# 创建 Collection
collection = Collection(name="legal_knowledge_base", schema=schema)
```

---

## 🔍 索引策略

### HNSW (生产环境推荐)

```python
index_params = {
    "metric_type": "L2",
    "index_type": "HNSW",
    "params": {
        "M": 16,              # 邻居数
        "efConstruction": 256 # 构建深度
    }
}

collection.create_index(
    field_name="embedding",
    index_params=index_params
)

# 加载到内存
collection.load()
```

---

## 📝 数据插入

### 批量插入示例

```python
import json

def ingest_to_milvus(chunks, version_id):
    """批量插入 Chunks 到 Milvus"""

    data = [
        [chunk['chunk_id'] for chunk in chunks],
        [chunk['chunk_text'] for chunk in chunks],
        [chunk['chunk_hash'] for chunk in chunks],
        [chunk['embedding'] for chunk in chunks],
        [json.dumps(chunk['merkle_proof']) for chunk in chunks],
        [chunk['merkle_index'] for chunk in chunks],
        [version_id for _ in chunks],
        [chunk['source_file'] for chunk in chunks],
        [chunk['chunk_order'] for chunk in chunks],
        [chunk['law_hierarchy'] for chunk in chunks],
        [json.dumps(chunk['references']) for chunk in chunks],
        [chunk['article_number'] for chunk in chunks],
        [chunk['law_name'] for chunk in chunks],
        [chunk['law_type'] for chunk in chunks],
        [chunk['effective_date'] for chunk in chunks],
        [chunk['priority_score'] for chunk in chunks],
    ]

    collection.insert(data)
    collection.flush()  # 确保持久化

    print(f"✅ Inserted {len(chunks)} chunks")
```

---

## 🔎 检索示例

### 向量检索 + 版本过滤

```python
def search_chunks(query_vector, top_k=10, version_id=None):
    """检索最相关的 Chunks"""

    search_params = {
        "metric_type": "L2",
        "params": {"ef": 64}  # HNSW 搜索参数
    }

    # 版本过滤表达式
    expr = f"version_id == {version_id}" if version_id else None

    results = collection.search(
        data=[query_vector],
        anns_field="embedding",
        param=search_params,
        limit=top_k,
        expr=expr,
        output_fields=["chunk_text", "chunk_hash", "merkle_proof",
                      "law_hierarchy", "references", "article_number"]
    )

    return results[0]
```

---

## ⚡ 性能优化

### 1. 批量插入

```python
# 每次插入 1000 条
batch_size = 1000
for i in range(0, len(all_chunks), batch_size):
    batch = all_chunks[i:i+batch_size]
    ingest_to_milvus(batch, version_id)
```

### 2. 索引参数调优

| 参数           | 开发环境 | 生产环境 |
| :------------- | :------- | :------- |
| M              | 8        | 16       |
| efConstruction | 128      | 256      |
| ef (search)    | 32       | 64       |

### 3. 内存管理

```python
# 定期 compact
collection.compact()

# 查看统计信息
collection.get_stats()
```

---

## ❓ 常见问题

**Q: 连接 Milvus 失败?**

```bash
docker-compose logs milvus
# 检查: curl localhost:9091/healthz
```

**Q: 插入速度慢?**

- 使用批量插入 (batch_size=1000)
- 增加 Milvus 内存配置

**Q: 检索结果不准确?**

- 检查 `metric_type` 是否与训练时一致
- 尝试增大 `ef` 参数

---

## 🔗 相关文档

- [← 区块链层实现](./02_blockchain_layer.md)
- [数据处理管道 →](./04_data_pipeline.md)
