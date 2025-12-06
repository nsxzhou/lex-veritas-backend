# 法律文档 RAG 优化方案：结构化分块与引用增强

> **摘要**:本文档详细阐述了针对法律领域 RAG (Retrieval-Augmented Generation) 系统的优化策略。针对法律文档层级严谨、引用复杂的特点,我们摒弃了传统的固定字符分块模式,采用**结构化感知分块 (Structure-Aware Chunking)** 和 **引用增强 (Citation Enrichment)** 技术,旨在实现 100% 的检索准确性和上下文完整性。本文档基于 **eino 框架** 提供具体实施方案。

---

## 0. 文档格式选择策略

### 0.1 三格式并存策略

针对"应该使用 PDF、DOCX 还是 Markdown 格式"的问题,我们推荐采用**三格式并存**策略:

| 格式         | 主要用途       | 优势                                                                        | 在 RAG 系统中的角色 |
| :----------- | :------------- | :-------------------------------------------------------------------------- | :------------------ |
| **Markdown** | RAG 处理主格式 | • 结构清晰,层级分明<br>• 机器可读性极强<br>• 易于正则解析<br>• LLM 处理友好 | **主要处理格式**    |
| **PDF**      | 存证与归档     | • 格式固定,不可篡改<br>• 法律效力认可<br>• 跨平台一致性                     | **存证备份**        |
| **DOCX**     | 编辑与协作     | • 易于编辑修改<br>• 协作功能丰富<br>• XML 结构化存储                        | **编辑暂存**        |

### 0.2 推荐工作流程

```
法律文档源文件 (PDF/DOCX)
    ↓ 转换
Markdown 格式 (结构化)
    ↓ RAG处理
分块 → 向量化 → 检索 → 生成
    ↓
同时保留 PDF 原文作为区块链存证
```

**关键理由**:

- **Markdown** 的简单语法(#、##、###)天然对应法律层级(编、章、节),极易提取结构
- PDF/DOCX 可通过工具(如 `pandoc`)无损转换为 Markdown
- 保留 PDF 原文可满足法律存证和区块链哈希验证需求

### 0.3 转换工具建议

```bash
# PDF → Markdown (保留结构)
pandoc input.pdf -o output.md --extract-media=./media

# DOCX → Markdown
pandoc input.docx -o output.md -t markdown --wrap=none
```

---

## 1. 背景与挑战

在构建法律领域的 RAG 系统时，直接使用通用的文本分块器（如 LangChain 的 `RecursiveCharacterTextSplitter`）会面临以下严重问题：

### 1.1 语义断裂 (Semantic Fragmentation)

法律条款（Article）是法律逻辑的最小完整单元。固定字符分块（如每 500 字符切一刀）极易将一条完整的法律规定从中间切断，导致 LLM 无法获取完整的规则描述，从而产生误读。

### 1.2 上下文丢失 (Context Loss)

法律文档具有极强的层级结构（编 -> 章 -> 节 -> 条）。

- **例子**：单独的一条"第十条 这里的规定..."，如果失去了"第一编 总则 - 第二章 自然人"这个父级上下文，其适用范围和法律效力就变得模糊不清。
- **后果**：检索时可能匹配到错误的"第十条"（因为不同法律或不同章节都可能有第十条）。

### 1.3 引用失效 (Broken Citations)

法律条文之间存在复杂的交叉引用（如"依照本法第四十五条规定"）。

- **问题**：如果"第四十五条"不在当前的 Chunk 中，LLM 就无法验证引用的具体内容，导致"幻觉"或无法回答。

---

## 2. 核心解决方案：结构化感知分块 (Structure-Aware Chunking)

我们提出一种基于规则的、感知文档结构的 ETL (Extract, Transform, Load) 流程。

### 2.1 分块策略：以"法条"为原子单位

我们不再关注字符数量，而是关注**文档结构**。

- **分割逻辑**：利用正则表达式精准识别法条边界。
- **正则示例**：`^第[零一二三四五六七八九十百]+条\s+`
- **处理流程**：
  1. 扫描文档，识别所有"编"、"章"、"节"的标题，构建文档树（Document Tree）。
  2. 识别所有"条"（Article），将其作为叶子节点。
  3. 每个 Chunk **严格包含且仅包含**一条完整的法条（对于超长法条可进行子句分割，但必须保持元数据关联）。

### 2.2 数据结构对比

| 特性         | 传统分块 (Fixed-size) | 结构化分块 (Structure-Aware) |
| :----------- | :-------------------- | :--------------------------- |
| **边界**     | 随机（字符数限制）    | 精准（法条结尾）             |
| **完整性**   | 低（常被截断）        | 高（100% 完整）              |
| **上下文**   | 仅包含前后文          | 包含完整的层级路径           |
| **检索粒度** | 文本片段              | 法律条款                     |

---

## 3. 上下文增强 (Context Enrichment)

为了解决上下文丢失问题，我们在向量化（Embedding）之前，对文本进行**显式增强**。

### 3.1 元数据注入 (Metadata Injection)

在存入向量数据库（Milvus）时，每个 Chunk 都携带丰富的结构化元数据：

```json
{
  "chunk_id": "civil_code_123",
  "content": "第123条 民事主体...",
  "metadata": {
    "law_name": "中华人民共和国民法典",
    "part": "第一编 总则",
    "chapter": "第五章 民事权利",
    "section": "第一节 物权",
    "article_number": "123"
  }
}
```

### 3.2 语义嵌入增强 (Embedding Enrichment)

在计算向量时，我们不仅仅计算 `content` 的向量，而是计算 **"上下文 + 内容"** 的组合向量。

- **原始文本**：
  > "第一百二十三条 民事主体依法享有物权。"
- **增强后的 Embedding 输入**：
  > "法律：中华人民共和国民法典；层级：第一编 总则 > 第五章 民事权利 > 第一节 物权；条文：第一百二十三条 民事主体依法享有物权。"

**优势**：这样即使用户查询"民法典关于物权的规定"，即使条文中没有出现"民法典"字样，也能通过增强的上下文被精准召回。

---

## 4. 引用处理机制 (Citation Handling)

针对法条中的引用关系，我们引入**图谱思维**。

### 4.1 引用提取 (Reference Extraction)

在 Ingestion 阶段，使用 NLP 或正则提取条文中的引用指针。

- **输入**："依照本法第四十五条规定..."
- **提取**：`ref: ["Article 45"]`

### 4.2 链接存储

将提取到的引用存入 Metadata：

```json
"references": ["article_45", "article_90"]
```

### 4.3 检索策略 (Retrieval Strategy)

在 RAG 检索阶段，采用**两步检索**：

1.  **Step 1**：根据用户 Query 检索出最相关的 Top-K 法条（如 Article A）。
2.  **Step 2**：读取 Article A 的 `references` 字段，自动加载其引用的 Article 45 和 Article 90。
3.  **Synthesis**：将 Article A + Article 45 + Article 90 一并送入 LLM 上下文。

---

## 5. 引用关系图谱化存储 (Graph-Enhanced Citation Storage)

为了应对复杂的法律引用网络（包括直接引用、间接引用、司法解释等），我们在元数据中引入**图结构**。

### 5.1 增强的元数据结构

```json
{
  "chunk_id": "civil_code_article_123",
  "content": "第123条 民事主体依法享有物权...",
  "metadata": {
    "law_name": "中华人民共和国民法典",
    "part": "第一编 总则",
    "chapter": "第五章 民事权利",
    "section": "第一节 物权",
    "article_number": "123",
    "references": {
      "direct_refs": ["article_45", "article_90"],
      "referenced_by": ["article_200", "article_305"],
      "related": ["article_121", "article_122", "article_124"],
      "interpretations": ["supreme_court_2024_05"],
      "amendments": ["2020_amendment", "2024_amendment"]
    },
    "law_hierarchy": {
      "type": "法律",
      "jurisdiction": "全国",
      "priority_score": 10
    }
  }
}
```

### 5.2 优势

- **完整的引用链**：避免遗漏 A→B→C 的间接引用
- **反向索引**：快速找到"哪些法条引用了当前法条"
- **司法解释关联**：自动关联最高法、最高检的权威解释

---

## 6. 三层检索策略 (Hierarchical Retrieval)

为确保 100% 的引用完整性，我们将两步检索扩展为**三层检索**。

### 6.1 检索架构

```
第一层：语义检索 (Semantic Search)
  ↓ 召回 Top-K 最相关法条

第二层：引用扩展 (Citation Expansion)
  ↓ 自动加载直接引用的法条（避免引用断链）

第三层：上下文补全 (Context Completion)
  ↓ 补充同章节的上下文法条（提供完整的法律逻辑）
```

### 6.2 实现逻辑

```python
def enhanced_retrieval(query, top_k=5):
    """三层检索实现"""

    # 第一层：语义检索
    primary_results = vector_search(query, top_k)

    # 第二层：引用扩展
    expanded_results = []
    for chunk in primary_results:
        expanded_results.append(chunk)
        # 加载所有直接引用的法条
        for ref_id in chunk.metadata.references.direct_refs:
            expanded_results.append(load_chunk(ref_id))

    # 第三层：上下文补全
    final_results = []
    for chunk in primary_results:
        # 加载该章节的标题和前置条款
        parent_context = load_parent_context(
            law_name=chunk.metadata.law_name,
            chapter=chunk.metadata.chapter
        )
        final_results.append(parent_context)
        final_results.extend(expanded_results)

    return deduplicate(final_results)
```

### 6.3 效果对比

| 场景                        | 两步检索        | 三层检索        |
| :-------------------------- | :-------------- | :-------------- |
| 简单查询（单一法条）        | ✅ 足够         | ✅ 足够         |
| 复杂引用（A 引用 B 引用 C） | ⚠️ 可能遗漏 C   | ✅ 完整加载     |
| 需要上下文理解              | ⚠️ 缺少章节背景 | ✅ 自动补充章节 |

---

## 7. 引用一致性校验 (Citation Consistency Check)

在 LLM 生成答案后，增加**后处理验证层**，确保不会引用未经检索验证的法条。

### 7.1 校验流程

```python
def verify_citation_consistency(llm_answer, retrieved_chunks):
    """
    检查LLM生成的答案中是否包含未经验证的法条引用
    """
    # 1. 提取LLM答案中的所有法条引用
    mentioned_articles = extract_article_references(llm_answer)
    # 例如：["第45条", "第123条", "第200条"]

    # 2. 检查这些引用是否都在retrieved_chunks中
    verified_articles = {
        chunk.metadata.article_number
        for chunk in retrieved_chunks
    }

    # 3. 如果LLM引用了未检索到的法条，标记为"需人工复核"
    unverified = mentioned_articles - verified_articles

    if unverified:
        return {
            "status": "NEEDS_REVIEW",
            "warning": f"答案引用了未验证的法条：{unverified}",
            "action": "自动补充检索或拒绝回答"
        }

    return {"status": "VERIFIED"}
```

### 7.2 防御性生成 (Defensive Generation)

更新 Prompt，强制 LLM 只引用已提供的法条：

```python
STRICT_PROMPT = """
你是 LexVeritas 法律 AI。严格遵守以下规则：

1. 你只能引用下方【可用法条】中列出的条款。
2. 如果需要引用【可用法条】之外的法律依据，必须明确说明"需要进一步检索第X条"，而不是直接引用。
3. 在回答中，每一句话后必须标注 [Ref: Article X]。

【可用法条】：
{verified_chunks_with_ids}

【用户问题】：
{question}
"""
```

---

## 8. 法条版本溯源 (Temporal Versioning)

法律文档会修订，需要支持**时间点查询**和**版本对比**。

### 8.1 版本化数据结构

```json
{
  "chunk_id": "civil_code_article_123",
  "versions": [
    {
      "version_id": "v1.0",
      "effective_date": "2021-01-01",
      "content": "第123条 民事主体依法享有物权。",
      "merkle_root": "0xabc123...",
      "source": "2020-05-28 全国人大通过"
    },
    {
      "version_id": "v2.0",
      "effective_date": "2024-07-01",
      "content": "第123条 民事主体依法享有物权。物权包括数字资产。",
      "merkle_root": "0xdef456...",
      "source": "2024-03-15 修订",
      "amendment_note": "增加数字资产保护条款"
    }
  ],
  "current_version": "v2.0"
}
```

### 8.2 时间点查询逻辑

```python
def retrieve_by_date(article_id, query_date):
    """根据时间点查询生效的法条版本"""
    versions = load_all_versions(article_id)

    # 找到在 query_date 时生效的版本
    for version in reversed(versions):
        if version.effective_date <= query_date:
            return version

    return None  # 该时间点法条尚未生效
```

### 8.3 应用场景

- **用户查询**："2022 年签订的合同是否符合民法典规定？" → 使用 2022-01-01 生效的版本
- **历史案例分析**：自动匹配案件发生时的法律文本

---

## 9. 法律冲突检测与优先级排序 (Conflict Detection)

当检索到多个相关法条时，系统需要自动识别并处理法律冲突。

### 9.1 法律层级体系

```python
LAW_HIERARCHY = {
    "宪法": 100,
    "法律": 90,          # 如民法典、刑法
    "行政法规": 80,      # 国务院制定
    "地方性法规": 70,    # 省级人大制定
    "部门规章": 60,      # 部委制定
    "地方政府规章": 50
}

CONFLICT_RESOLUTION_RULES = [
    "特别法优于一般法",
    "新法优于旧法",
    "上位法优于下位法"
]
```

### 9.2 冲突检测算法

```python
def detect_conflicts(retrieved_chunks):
    """检测检索结果中是否存在法律冲突"""
    conflicts = []

    for i, chunk_a in enumerate(retrieved_chunks):
        for chunk_b in retrieved_chunks[i+1:]:
            # 检查是否涉及同一法律问题但规定不同
            if semantic_similarity(chunk_a, chunk_b) > 0.9:
                if chunk_a.content != chunk_b.content:
                    conflicts.append({
                        "articles": [chunk_a.id, chunk_b.id],
                        "resolution": resolve_conflict(chunk_a, chunk_b)
                    })

    return conflicts

def resolve_conflict(chunk_a, chunk_b):
    """根据法律原则解决冲突"""
    # 规则1：特别法优于一般法
    if chunk_a.metadata.law_type == "特别法":
        return chunk_a

    # 规则2：新法优于旧法
    if chunk_a.version.effective_date > chunk_b.version.effective_date:
        return chunk_a

    # 规则3：上位法优于下位法
    if LAW_HIERARCHY[chunk_a.law_hierarchy.type] > LAW_HIERARCHY[chunk_b.law_hierarchy.type]:
        return chunk_a

    return chunk_b
```

---

## 10. 基于 eino 框架的实施路径

### 10.1 eino 框架核心能力

[eino](https://github.com/cloudwego/eino) 是字节跳动开源的 Golang AI 应用开发框架,专为构建 LLM 应用设计:

- ✅ **RecursiveSplitter**: 递归分块器,支持自定义分隔符和层级
- ✅ **多格式解析器**: 原生支持 PDF、DOCX、XLSX 解析
- ✅ **向量化组件**: 集成多种 Embedding 模型
- ✅ **Chain 编排**: 可构建复杂的 RAG 检索链
- ✅ **Stream 支持**: 流式处理大规模文档

### 10.2 法律文档专用分块器实现

**文件**: `lex-veritas-ingestion/chunkers/legal_chunker.go`

```go
package chunkers

import (
    "context"
    "regexp"
    "github.com/cloudwego/eino/schema"
    "github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
)

// 法律文档分块器配置
type LegalChunkerConfig struct {
    // 层级分隔符(按优先级:
    // 1. 编 → 2. 章 → 3. 节 → 4. 条 → 5. 段落)
    Separators []string

    // 是否注入层级上下文
    EnrichContext bool

    // 引用提取正则
    ReferencePattern *regexp.Regexp
}

// 创建法律文档分块器
func NewLegalChunker(ctx context.Context, cfg *LegalChunkerConfig) (*LegalChunker, error) {
    // 默认配置
    if cfg == nil {
        cfg = &LegalChunkerConfig{
            Separators: []string{
                "\\n第[零一二三四五六七八九十百千]+编",  // 编
                "\\n第[零一二三四五六七八九十百千]+章",  // 章
                "\\n第[零一二三四五六七八九十百千]+节",  // 节
                "\\n第[零一二三四五六七八九十百千]+条",  // 条(最小单位)
                "\\n\\n",                            // 段落
            },
            EnrichContext: true,
            ReferencePattern: regexp.MustCompile("第[零一二三四五六七八九十百千]+条"),
        }
    }

    // 使用 eino 的 RecursiveSplitter 作为底层
    splitter, err := recursive.NewSplitter(ctx, &recursive.Config{
        ChunkSize:   2000,  // 法条通常较长
        OverlapSize: 200,   // 保留上下文重叠
        Separators:  cfg.Separators,
    })
    if err != nil {
        return nil, err
    }

    return &LegalChunker{
        splitter: splitter,
        config:   cfg,
    }, nil
}

// 分块并增强元数据
func (lc *LegalChunker) Transform(ctx context.Context, docs []*schema.Document) ([]*schema.Document, error) {
    // 1. 基础分块
    chunks, err := lc.splitter.Transform(ctx, docs)
    if err != nil {
        return nil, err
    }

    // 2. 元数据增强
    for _, chunk := range chunks {
        // 提取层级路径
        hierarchy := lc.extractHierarchy(chunk.Content)
        chunk.MetaData["law_hierarchy"] = hierarchy

        // 提取引用关系
        refs := lc.config.ReferencePattern.FindAllString(chunk.Content, -1)
        chunk.MetaData["references"] = refs

        // 上下文注入(可选)
        if lc.config.EnrichContext {
            chunk.Content = lc.enrichWithContext(chunk.Content, hierarchy)
        }
    }

    return chunks, nil
}

// 提取层级路径: "中华人民共和国民法典 > 第一编 总则 > 第五章 民事权利 > 第123条"
func (lc *LegalChunker) extractHierarchy(content string) string {
    // 使用正则依次匹配编、章、节、条
    // ...(实现略)
    return "第一编 总则 > 第五章 民事权利 > 第123条"
}

// 上下文增强: 在 Embedding 前注入层级信息
func (lc *LegalChunker) enrichWithContext(content, hierarchy string) string {
    return fmt.Sprintf("【层级】%s\n【内容】%s", hierarchy, content)
}
```

### 10.3 三层检索链实现

**文件**: `lex-veritas-backend/retrieval/hierarchical_retriever.go`

```go
package retrieval

import (
    "context"
    "github.com/cloudwego/eino/compose"
    "github.com/cloudwego/eino/schema"
)

// 构建三层检索链
func BuildHierarchicalRAGChain(
    ctx context.Context,
    vectorStore compose.Retriever,
    chatModel compose.ChatModel,
) (*compose.Chain, error) {

    // 第二层: 引用扩展 Lambda
    expandReferences := compose.InvokableLambda(
        func(ctx context.Context, docs []*schema.Document) ([]*schema.Document, error) {
            expanded := make([]*schema.Document, 0, len(docs)*2)
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

    // 第三层: 上下文补全 Lambda
    enrichContext := compose.InvokableLambda(
        func(ctx context.Context, docs []*schema.Document) (map[string]any, error) {
            // 构建包含完整上下文的 Prompt
            contextBuilder := strings.Builder{}
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

    // 严格 Prompt 模板
    strictPrompt, _ := prompt.FromMessages(
        schema.FString,
        schema.SystemMessage(
            "你是专业法律AI助手。严格规则:\n"+
            "1. 仅引用下方【可用法条】中的内容\n"+
            "2. 每句话后标注 [参考: 第X条]\n"+
            "3. 如需引用未提供的法条,明确说明'需进一步检索'",
        ),
        schema.UserMessage("【可用法条】\n{legal_context}\n\n【问题】\n{question}"),
    )

    // 组装完整 RAG Chain
    return compose.NewChain[string, *schema.Message]().
        AppendRetriever(vectorStore).        // 第一层: 语义检索
        AppendLambda(expandReferences).      // 第二层: 引用扩展
        AppendLambda(enrichContext).         // 第三层: 上下文补全
        AppendChatTemplate(strictPrompt).    // 严格 Prompt
        AppendChatModel(chatModel).          // LLM 生成
        Compile(ctx)
}
```

### 10.4 引用一致性验证

**文件**: `lex-veritas-backend/validation/citation_validator.go`

```go
package validation

import (
    "regexp"
    "github.com/cloudwego/eino/schema"
)

type CitationValidator struct {
    pattern *regexp.Regexp
}

func NewCitationValidator() *CitationValidator {
    return &CitationValidator{
        pattern: regexp.MustCompile("第[零一二三四五六七八九十百千]+条"),
    }
}

func (cv *CitationValidator) Verify(
    answer string,
    retrievedDocs []*schema.Document,
) ValidationResult {
    // 提取答案中的法条引用
    mentioned := cv.pattern.FindAllString(answer, -1)

    // 构建已检索法条集合
    retrieved := make(map[string]bool)
    for _, doc := range retrievedDocs {
        if articleNum, ok := doc.MetaData["article_number"].(string); ok {
            retrieved[articleNum] = true
        }
    }

    // 检测未验证的引用
    unverified := []string{}
    for _, article := range mentioned {
        if !retrieved[article] {
            unverified = append(unverified, article)
        }
    }

    if len(unverified) > 0 {
        return ValidationResult{
            Status:         "NEEDS_REVIEW",
            UnverifiedRefs: unverified,
            Action:         "补充检索或拒绝回答",
        }
    }

    return ValidationResult{Status: "VERIFIED"}
}
```

### 10.5 渐进式实施路径

**Phase 1:基础架构(优先)**

- [ ] Markdown 转换管道 (PDF/DOCX → Markdown)
- [ ] 基于 eino 的法律文档分块器
- [ ] 元数据提取与注入
- [ ] 向量化与存储 (Milvus)

**Phase 2:引用增强(短期)**

- [ ] 三层检索链实现
- [ ] 引用关系图谱化存储
- [ ] 引用一致性校验

**Phase 3:生产级特性(中长期)**

- [ ] 法条版本溯源系统
- [ ] 法律冲突检测
- [ ] 图数据库集成 (Neo4j)
- [ ] 区块链存证集成

### 10.6 技术栈建议

| 功能模块     | 推荐技术            | 理由                     |
| :----------- | :------------------ | :----------------------- |
| **文档处理** | **eino + eino-ext** | **Golang 原生,性能优异** |
| 格式转换     | Pandoc              | 支持多格式无损转换       |
| 向量检索     | Milvus / Weaviate   | 高性能语义搜索           |
| 引用关系存储 | Neo4j / ArangoDB    | 原生支持图结构查询       |
| 版本管理     | PostgreSQL (JSONB)  | 支持复杂查询和事务       |
| Merkle Tree  | Go crypto/sha256    | 标准库,区块链友好        |

### 10.3 预期效果

通过实施上述完整方案，LexVeritas 系统将实现：

1.  **零截断**：保证所有检索到的法条都是完整的，无法律歧义。
2.  **零幻觉**：通过引用一致性校验，杜绝 LLM 捏造法条。
3.  **100% 引用完整性**：三层检索确保所有引用链都被加载。
4.  **时间点准确性**：版本溯源保证历史查询的准确性。
5.  **冲突解决能力**：自动识别并按法律原则排序冲突法条。
6.  **可解释性**：每一条回答都能精准定位到具体的法条编号、版本和区块链存证。

这一方案是构建专业级、生产级法律 AI 助手的完整技术基石。
