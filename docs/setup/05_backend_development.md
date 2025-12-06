# 后端开发指南

> Go + Eino RAG 编排

---

## 📁 项目结构

```
backend/
├── cmd/server/main.go
├── internal/
│   ├── config/config.go
│   ├── handler/chat_handler.go
│   ├── graph/
│   │   ├── builder.go
│   │   └── nodes/
│   │       ├── retriever.go
│   │       ├── verifier.go
│   │       └── llm.go
│   └── client/
│       ├── milvus_client.go
│       ├── blockchain_client.go
│       └── llm_client.go
├── go.mod
└── config.yaml
```

---

## 🚀 三层检索 Chain

```go
package graph

import (
    "context"
    "github.com/cloudwego/eino/compose"
    "github.com/cloudwego/eino/schema"
)

func BuildHierarchicalRAGChain(
    ctx context.Context,
    vectorStore compose.Retriever,
    chatModel compose.ChatModel,
) (*compose.Chain, error) {

    // 第二层: 引用扩展
    expandReferences := compose.InvokableLambda(
        func(ctx context.Context, docs []*schema.Document) ([]*schema.Document, error) {
            expanded := make([]*schema.Document, 0)
            visited := make(map[string]bool)

            for _, doc := range docs {
                expanded = append(expanded, doc)
                visited[doc.ID] = true

                // 加载引用的法条
                refs, ok := doc.MetaData["references"].([]string)
                if !ok {
                    continue
                }

                for _, refID := range refs {
                    if !visited[refID] {
                        refDoc := loadDocumentByID(ctx, vectorStore, refID)
                        if refDoc != nil {
                            expanded = append(expanded, refDoc)
                            visited[refID] = true
                        }
                    }
                }
            }

            return expanded, nil
        },
    )

    // 第三层: 上下文补全
    enrichContext := compose.InvokableLambda(
        func(ctx context.Context, docs []*schema.Document) (map[string]any, error) {
            var contextBuilder strings.Builder
            for i, doc := range docs {
                hierarchy := doc.MetaData["law_hierarchy"]
                contextBuilder.WriteString(fmt.Sprintf(
                    "[法条%d] %s\n%s\n\n",
                    i+1, hierarchy, doc.Content,
                ))
            }

            return map[string]any{
                "legal_context": contextBuilder.String(),
            }, nil
        },
    )

    // 组装 Chain
    return compose.NewChain[string, *schema.Message]().
        AppendRetriever(vectorStore).      // (一层: 语义检索
        AppendLambda(expandReferences).    // 第二层: 引用扩展
        AppendLambda(enrichContext).       // 第三层: 上下文补全
        AppendChatModel(chatModel).        // LLM 生成
        Compile(ctx)
}
```

---

## 🔐 Merkle 验证节点

```go
package nodes

import (
    "crypto/sha256"
    "encoding/hex"
)

type MerkleVerifier struct {
    blockchainClient *BlockchainClient
}

func (mv *MerkleVerifier) Verify(chunk *schema.Document) (bool, error) {
    // 1. 计算 Chunk 哈希
    chunkHash := sha256.Sum256([]byte(chunk.Content))

    // 2. 获取 Merkle Proof
    proofData := chunk.MetaData["merkle_proof"].(string)
    var proof MerkleProof
    json.Unmarshal([]byte(proofData), &proof)

    // 3. 本地验证
    computedRoot := mv.computeRoot(chunkHash[:], proof.Proof, proof.Index)

    // 4. 与链上 Root 比对
    onChainRoot, err := mv.blockchainClient.GetCurrentRoot(context.Background())
    if err != nil {
        return false, err
    }

    return bytes.Equal(computedRoot, onChainRoot), nil
}

func (mv *MerkleVerifier) computeRoot(leaf []byte, proof []string, index int) []byte {
    computedHash := leaf

    for _, proofElement := range proof {
        proofBytes, _ := hex.DecodeString(proofElement)

        if index%2 == 0 {
            computedHash = sha256Hash(append(computedHash, proofBytes...))
        } else {
            computedHash = sha256Hash(append(proofBytes, computedHash...))
        }

        index = index / 2
    }

    return computedHash
}
```

---

## 🌐 API 接口

```go
package handler

import (
    "github.com/gin-gonic/gin"
)

type ChatHandler struct {
    ragChain *compose.Chain
}

func (h *ChatHandler) HandleChat(c *gin.Context) {
    var req ChatRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 执行 RAG Chain
    result, err := h.ragChain.Invoke(c.Request.Context(), req.Question)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, ChatResponse{
        Answer:     result.Content,
        References: extractReferences(result),
    })
}
```

---

## ⚙️ 配置文件

`config.yaml`:

```yaml
milvus:
  host: localhost
  port: 19530
  collection: legal_knowledge_base

blockchain:
  rpc_url: https://rpc-amoy.polygon.technology
  contract_address: "0x..."

llm:
  provider: openai
  model: gpt-4
  api_key: ${OPENAI_API_KEY}

server:
  port: 8080
```

---

## 🔗 相关文档

- [← 数据处理管道](./04_data_pipeline.md)
- [部署与运维 →](./06_deployment.md)
