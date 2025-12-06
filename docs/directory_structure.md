# LexVeritas 项目目录结构说明

本文档详细说明了 LexVeritas 项目的目录结构及其各部分的作用。

## 根目录

- **`docker-compose.yaml`**: Docker 编排文件，用于在本地开发环境中启动 Milvus、Etcd、MinIO、后端和前端服务。
- **`project_setup_guide.md`**: 项目搭建指南，包含架构设计、技术栈选型和详细的实现步骤。

## 1. 后端服务 (`lex-veritas-backend/`)

基于 Go + Gin + Eino 框架实现的业务逻辑层，负责处理用户请求、RAG 流程编排和数据验证。

- **`cmd/server/`**
  - **`main.go`**: 程序的入口文件，负责初始化配置、数据库连接、区块链客户端，并启动 HTTP 服务器。
- **`internal/`**: 内部代码包，不对外暴露。
  - **`config/`**
    - **`config.go`**: 配置管理模块，负责加载和解析 `config.yaml` 配置文件。
  - **`handler/`**
    - **`chat_handler.go`**: HTTP 请求处理器，定义了 `/chat` 等 API 接口的处理逻辑。
  - **`graph/`**: Eino Graph 编排相关代码。
    - **`builder.go`**: 负责构建和编译 Eino Graph，将各个节点连接起来。
    - **`nodes/`**: 定义 Graph 中的各个节点。
      - **`retriever.go`**: 检索节点，负责调用 Milvus 检索相关 Chunks。
      - **`verifier.go`**: **核心验证节点**，负责计算 Chunk 哈希并与链上 Merkle Root 进行比对验证。
      - **`prompt.go`**: 提示词节点，负责构建包含验证通过信息的 Prompt。
      - **`llm.go`**: 大模型节点，负责调用 LLM API 生成回答。
  - **`client/`**: 外部服务客户端封装。
    - **`milvus_client.go`**: 封装 Milvus SDK，提供向量检索和数据查询功能。
    - **`blockchain_client.go`**: 封装区块链交互逻辑，用于获取 Merkle Root。
    - **`llm_client.go`**: 封装 LLM API (如 OpenAI) 的调用。
  - **`model/`**: 数据模型定义。
    - **`types.go`**: 定义系统中使用的数据结构，如 `Chunk`, `Citation` 等。
- **`go.mod`**: Go 模块定义文件，管理依赖。
- **`config.yaml`**: 后端服务的配置文件（端口、数据库地址、API Key 等）。

## 2. 区块链层 (`lex-veritas-blockchain/`)

基于 Hardhat 开发的智能合约项目，作为系统的信任锚定层。

- **`contracts/`**
  - **`LexKnowledgeBase.sol`**: 核心智能合约，用于存储法律知识库的版本信息和 Merkle Root。
- **`scripts/`**
  - **`deploy.ts`**: 合约部署脚本，用于将合约部署到 Polygon Amoy 测试网。
- **`hardhat.config.ts`**: Hardhat 配置文件，配置网络、编译器版本和插件。

## 4. 数据处理管道 (`lex-veritas-ingestion/`)

基于 Python 的数据处理脚本，负责知识库的构建和存证。

- **`ingestion.py`**: 核心脚本，执行以下流程：
  1.  加载 PDF 文档。
  2.  文本分块与清洗。
  3.  计算 Chunk 哈希。
  4.  构建 Merkle Tree。
  5.  调用智能合约发布 Merkle Root。
  6.  调用 OpenAI 生成 Embeddings。
  7.  将数据（含 Merkle Proof）存入 Milvus。
- **`legal_docs/`**: 存放原始法律文档（PDF/HTML）的目录。
