# 部署与运维指南

> Docker 容器化部署

---

## 🐳 Docker 部署

### Dockerfile

```dockerfile
# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server cmd/server/main.go

# 运行阶段
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/server .
COPY config.yaml .

EXPOSE 8080
CMD ["./server"]
```

### docker-compose.yml

```yaml
version: "3.8"

services:
  # 向量数据库栈
  etcd:
    image: quay.io/coreos/etcd:v3.5.5
    environment:
      - ETCD_AUTO_COMPACTION_MODE=revision
      - ETCD_AUTO_COMPACTION_RETENTION=1000
    volumes:
      - etcd_data:/etcd

  minio:
    image: minio/minio:RELEASE.2023-03-20T20-16-18Z
    command: minio server /minio_data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - minio_data:/minio_data

  milvus:
    image: milvusdb/milvus:v2.3.0
    command: ["milvus", "run", "standalone"]
    environment:
      ETCD_ENDPOINTS: etcd:2379
      MINIO_ADDRESS: minio:9000
    ports:
      - "19530:19530"
    depends_on:
      - etcd
      - minio
    volumes:
      - milvus_data:/var/lib/milvus

  # LexVeritas 后端
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - MILVUS_HOST=milvus
      - MILVUS_PORT=19530
      - BLOCKCHAIN_RPC=${AMOY_RPC_URL}
      - CONTRACT_ADDRESS=${CONTRACT_ADDRESS}
      - OPENAI_API_KEY=${OPENAI_API_KEY}
    depends_on:
      - milvus

  # 前端
  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    depends_on:
      - backend

volumes:
  etcd_data:
  minio_data:
  milvus_data:
```

### 启动

```bash
# 设置环境变量
cp .env.example .env
# 编辑 .env

# 启动全部服务
docker-compose up -d

# 查看日志
docker-compose logs -f backend
```

---

## 📊 监控

### 健康检查

```go
// backend/internal/handler/health.go
func (h *HealthHandler) Check(c *gin.Context) {
    status := gin.H{
        "status": "ok",
        "milvus": h.checkMilvus(),
        "blockchain": h.checkBlockchain(),
    }

    c.JSON(200, status)
}
```

### Prometheus Metrics

```yaml
# docker-compose.yml 添加
prometheus:
  image: prom/prometheus
  ports:
    - "9090:9090"
  volumes:
    - ./prometheus.yml:/etc/prometheus/prometheus.yml
```

---

## 🔧 故障排查

### Milvus 连接失败

```bash
# 检查容器状态
docker-compose ps

# 查看 Milvus 日志
docker-compose logs milvus

# 重启服务
docker-compose restart milvus
```

### 内存溢出

```yaml
# docker-compose.yml 限制内存
services:
  milvus:
    deploy:
      resources:
        limits:
          memory: 4G
```

---

## 🔗 相关文档

- [← 后端开发指南](./05_backend_development.md)
- [快速开始](./01_quick_start.md)
