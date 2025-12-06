# 区块链层实现指南

> 智能合约部署与后端集成

---

## 📋 智能合约代码

### LexKnowledgeBase.sol

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";

contract LexKnowledgeBase is Ownable {
    struct KnowledgeVersion {
        bytes32 merkleRoot;        // Merkle Root 哈希
        string description;        // 版本描述
        uint256 timestamp;         // 发布时间
        uint256 chunkCount;        // Chunk 总数
    }

    KnowledgeVersion[] public versions;

    event VersionPublished(
        uint256 indexed versionId,
        bytes32 merkleRoot,
        string description,
        uint256 chunkCount
    );

    constructor() Ownable(msg.sender) {}

    function publishVersion(
        bytes32 _root,
        string memory _desc,
        uint256 _count
    ) public onlyOwner {
        versions.push(KnowledgeVersion({
            merkleRoot: _root,
            description: _desc,
            timestamp: block.timestamp,
            chunkCount: _count
        }));

        emit VersionPublished(versions.length - 1, _root, _desc, _count);
    }

    function getCurrentVersion() public view returns (KnowledgeVersion memory) {
        require(versions.length > 0, "No version published");
        return versions[versions.length - 1];
    }

    function verifyChunk(
        bytes32 _leaf,
        bytes32[] memory _proof,
        uint256 _index
    ) public view returns (bool) {
        require(versions.length > 0, "No version published");
        bytes32 root = versions[versions.length - 1].merkleRoot;
        return _verifyMerkleProof(_leaf, _proof, root, _index);
    }

    function _verifyMerkleProof(
        bytes32 leaf,
        bytes32[] memory proof,
        bytes32 root,
        uint256 index
    ) internal pure returns (bool) {
        bytes32 computedHash = leaf;

        for (uint256 i = 0; i < proof.length; i++) {
            bytes32 proofElement = proof[i];

            if (index % 2 == 0) {
                computedHash = keccak256(abi.encodePacked(computedHash, proofElement));
            } else {
                computedHash = keccak256(abi.encodePacked(proofElement, computedHash));
            }

            index = index / 2;
        }

        return computedHash == root;
    }
}
```

---

## 🚀 部署步骤

### 1. 初始化 Hardhat 项目

```bash
mkdir blockchain && cd blockchain
npx hardhat init  # 选择 TypeScript
npm install @openzeppelin/contracts
```

### 2. 配置网络

编辑 `hardhat.config.ts`:

```typescript
import { HardhatUserConfig } from "hardhat/config";

const config: HardhatUserConfig = {
  solidity: "0.8.20",
  networks: {
    amoy: {
      url: process.env.AMOY_RPC_URL || "https://rpc-amoy.polygon.technology",
      accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [],
    },
  },
};

export default config;
```

### 3. 部署脚本

`scripts/deploy.ts`:

```typescript
import { ethers } from "hardhat";

async function main() {
  const LexKnowledgeBase = await ethers.getContractFactory("LexKnowledgeBase");
  const contract = await LexKnowledgeBase.deploy();
  await contract.waitForDeployment();

  const address = await contract.getAddress();
  console.log(`LexKnowledgeBase deployed to: ${address}`);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
```

### 4. 执行部署

```bash
export AMOY_RPC_URL="https://rpc-amoy.polygon.technology"
export PRIVATE_KEY="your_private_key"

npx hardhat run scripts/deploy.ts --network amoy
```

---

## 🔗 后端 Go 集成

### 1. 生成 Go Bindings

```bash
# 编译合约
npx hardhat compile

# 安装 abigen
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# 生成 Go 代码
abigen --abi=artifacts/contracts/LexKnowledgeBase.sol/LexKnowledgeBase.json \
       --pkg=contracts \
       --out=backend/internal/contracts/lex_knowledge_base.go
```

### 2. Go 客户端代码

```go
package blockchain

import (
    "context"
    "math/big"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
    "your-org/backend/internal/contracts"
)

type BlockchainClient struct {
    client   *ethclient.Client
    contract *contracts.LexKnowledgeBase
}

func NewBlockchainClient(rpcURL, contractAddr string) (*BlockchainClient, error) {
    client, err := ethclient.Dial(rpcURL)
    if err != nil {
        return nil, err
    }

    contract, err := contracts.NewLexKnowledgeBase(
        common.HexToAddress(contractAddr),
        client,
    )
    if err != nil {
        return nil, err
    }

    return &BlockchainClient{
        client:   client,
        contract: contract,
    }, nil
}

func (bc *BlockchainClient) GetCurrentRoot(ctx context.Context) ([]byte, error) {
    version, err := bc.contract.GetCurrentVersion(nil)
    if err != nil {
        return nil, err
    }

    return version.MerkleRoot[:], nil
}

func (bc *BlockchainClient) VerifyChunk(
    ctx context.Context,
    leaf [32]byte,
    proof [][32]byte,
    index *big.Int,
) (bool, error) {
    return bc.contract.VerifyChunk(nil, leaf, proof, index)
}
```

---

## 🐍 Python 集成 (数据处理管道)

### 发布新版本

```python
from web3 import Web3
import json
import os

def publish_version_to_blockchain(merkle_root: bytes, description: str, chunk_count: int):
    """发布新版本到区块链"""
    w3 = Web3(Web3.HTTPProvider(os.getenv('AMOY_RPC_URL')))

    # 加载合约
    with open('contract-abi.json') as f:
        abi = json.load(f)

    contract = w3.eth.contract(
        address=os.getenv('CONTRACT_ADDRESS'),
        abi=abi
    )

    # 构建交易
    tx = contract.functions.publishVersion(
        merkle_root,  # bytes32
        description,
        chunk_count
    ).build_transaction({
        'from': os.getenv('ADMIN_ADDRESS'),
        'nonce': w3.eth.get_transaction_count(os.getenv('ADMIN_ADDRESS')),
        'gas': 300000,
        'maxFeePerGas': w3.to_wei('50', 'gwei'),
        'maxPriorityFeePerGas': w3.to_wei('2', 'gwei')
    })

    # 签名并发送
    signed_tx = w3.eth.account.sign_transaction(tx, os.getenv('PRIVATE_KEY'))
    tx_hash = w3.eth.send_raw_transaction(signed_tx.rawTransaction)

    # 等待确认
    receipt = w3.eth.wait_for_transaction_receipt(tx_hash, timeout=120)

    return {
        'tx_hash': receipt['transactionHash'].hex(),
        'block_number': receipt['blockNumber'],
        'status': receipt['status']
    }
```

---

## ❓ 常见问题

**Q: Gas 费用过高?**

- 使用 Polygon Amoy 测试网(免费)
- 批量发布版本而非单个 Chunk

**Q: 合约验证失败?**

```bash
npx hardhat verify --network amoy <CONTRACT_ADDRESS>
```

**Q: RPC 连接超时?**

- 使用备用 RPC: https://polygon-amoy.g.alchemy.com/v2/YOUR-API-KEY
- 增加超时时间: `timeout=300`

---

## 🔗 相关文档

- [← 快速开始](./01_quick_start.md)
- [向量数据库配置 →](./03_vector_database.md)
- [数据处理管道](./04_data_pipeline.md)
