# Git 主项目和子模块配置指南

本文档详细说明如何从零开始配置 Git 主项目和子模块,以及如何将整个项目推送到远程仓库。

## 目录

1. [概述](#概述)
2. [前置准备](#前置准备)
3. [配置步骤](#配置步骤)
4. [推送到远程仓库](#推送到远程仓库)
5. [团队协作指南](#团队协作指南)
6. [常见问题](#常见问题)

---

## 概述

### 什么是 Git 子模块?

Git 子模块允许你将一个 Git 仓库作为另一个 Git 仓库的子目录。这使得你可以:

- 将大型项目拆分为多个独立的仓库
- 每个子项目有独立的版本控制
- 在主项目中精确控制每个子项目的版本

### LexVeritas 项目结构

```
LexVeritas/                    # 主项目
├── lex-veritas-backend/       # 子模块: Go 后端服务
├── lex-veritas-blockchain/    # 子模块: 智能合约
├── lex-veritas-frontend/      # 子模块: 前端应用
├── lex-veritas-ingestion/     # 子模块: 数据处理管道
├── docs/                      # 主项目文档
└── docker-compose.yaml        # 主项目配置文件
```

---

## 前置准备

### 1. 安装 Git

确保已安装 Git (版本 2.13 或更高):

```bash
git --version
```

### 2. 配置 Git 用户信息

```bash
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

### 3. 准备远程仓库

在 GitHub/GitLab/Gitee 上创建以下仓库:

- `LexVeritas` (主项目)
- `lex-veritas-backend`
- `lex-veritas-blockchain`
- `lex-veritas-frontend`
- `lex-veritas-ingestion`

> [!TIP]
> 创建仓库时,**不要**初始化 README、.gitignore 或 LICENSE,保持仓库为空。

---

## 配置步骤

### 方案一: 从零开始配置(新项目)

如果你还没有初始化任何 Git 仓库,按照以下步骤操作:

#### 步骤 1: 初始化主项目

```bash
cd /path/to/LexVeritas

# 初始化主项目 Git 仓库
git init

# 创建 .gitignore 文件
cat > .gitignore << 'EOF'
# macOS
.DS_Store

# IDE
.vscode/
.idea/

# Environment
.env
.env.local

# Logs
*.log
EOF

# 创建 README.md
cat > README.md << 'EOF'
# LexVeritas

基于区块链的可验证法律知识库 RAG 系统

## 克隆项目

\`\`\`bash
git clone --recursive <repository-url>
\`\`\`

## 更新子模块

\`\`\`bash
git submodule update --remote --merge
\`\`\`
EOF
```

#### 步骤 2: 为每个子项目初始化 Git 仓库

```bash
# Backend
cd lex-veritas-backend
git init
cat > .gitignore << 'EOF'
# Go
*.exe
*.dll
*.so
*.dylib
*.test
*.out
vendor/
.env
*.log
.DS_Store
EOF
git add -A
git commit -m "Initial commit: backend service"
cd ..

# Blockchain
cd lex-veritas-blockchain
git init
cat > .gitignore << 'EOF'
node_modules
.env
cache
artifacts
typechain
typechain-types
*.log
.DS_Store
EOF
git add -A
git commit -m "Initial commit: blockchain contracts"
cd ..

# Frontend
cd lex-veritas-frontend
git init
cat > .gitignore << 'EOF'
node_modules
.env
.env.local
/build
/dist
/.next
*.log
.DS_Store
EOF
git add -A
git commit -m "Initial commit: frontend application"
cd ..

# Ingestion
cd lex-veritas-ingestion
git init
cat > .gitignore << 'EOF'
__pycache__/
*.py[cod]
venv/
.env
*.log
.DS_Store
legal_docs/*.pdf
EOF
git add -A
git commit -m "Initial commit: data ingestion pipeline"
cd ..
```

#### 步骤 3: 配置本地文件传输(临时)

```bash
# 允许 Git 使用本地文件协议
git config --global protocol.file.allow always
```

#### 步骤 4: 将子项目移到临时位置

```bash
# 创建临时目录
mkdir -p /tmp/lexveritas-submodules

# 移动子项目
mv lex-veritas-backend /tmp/lexveritas-submodules/
mv lex-veritas-blockchain /tmp/lexveritas-submodules/
mv lex-veritas-frontend /tmp/lexveritas-submodules/
mv lex-veritas-ingestion /tmp/lexveritas-submodules/
```

#### 步骤 5: 添加子模块(使用本地路径)

```bash
# 添加子模块
git submodule add file:///tmp/lexveritas-submodules/lex-veritas-backend lex-veritas-backend
git submodule add file:///tmp/lexveritas-submodules/lex-veritas-blockchain lex-veritas-blockchain
git submodule add file:///tmp/lexveritas-submodules/lex-veritas-frontend lex-veritas-frontend
git submodule add file:///tmp/lexveritas-submodules/lex-veritas-ingestion lex-veritas-ingestion
```

#### 步骤 6: 提交主项目

```bash
# 添加所有文件并提交
git add -A
git commit -m "Initial commit: add submodules and project structure"
```

#### 步骤 7: 验证配置

```bash
# 查看子模块状态
git submodule status

# 查看 .gitmodules 文件
cat .gitmodules

# 验证主项目状态
git status
```

---

### 方案二: 已有 Git 仓库的项目

如果你的项目已经初始化了 Git 仓库,但还没有配置子模块:

#### 步骤 1: 为子目录初始化 Git 仓库

```bash
cd /path/to/LexVeritas

# 为每个子目录初始化 Git 仓库
cd lex-veritas-backend && git init && git add -A && git commit -m "Initial commit" && cd ..
cd lex-veritas-blockchain && git init && git add -A && git commit -m "Initial commit" && cd ..
cd lex-veritas-frontend && git init && git add -A && git commit -m "Initial commit" && cd ..
cd lex-veritas-ingestion && git init && git add -A && git commit -m "Initial commit" && cd ..
```

#### 步骤 2: 从主项目中移除子目录

```bash
# 从 Git 跟踪中移除(但保留文件)
git rm --cached -r lex-veritas-backend
git rm --cached -r lex-veritas-blockchain
git rm --cached -r lex-veritas-frontend
git rm --cached -r lex-veritas-ingestion

# 提交更改
git commit -m "Remove subdirectories before converting to submodules"
```

#### 步骤 3: 移动到临时位置并添加为子模块

```bash
# 移动到临时位置
mkdir -p /tmp/lexveritas-submodules
mv lex-veritas-* /tmp/lexveritas-submodules/

# 配置本地文件传输
git config --global protocol.file.allow always

# 添加为子模块
git submodule add file:///tmp/lexveritas-submodules/lex-veritas-backend lex-veritas-backend
git submodule add file:///tmp/lexveritas-submodules/lex-veritas-blockchain lex-veritas-blockchain
git submodule add file:///tmp/lexveritas-submodules/lex-veritas-frontend lex-veritas-frontend
git submodule add file:///tmp/lexveritas-submodules/lex-veritas-ingestion lex-veritas-ingestion

# 提交
git add -A
git commit -m "Convert subdirectories to submodules"
```

---

## 推送到远程仓库

配置完成后,按照以下步骤将项目推送到远程仓库。

### 步骤 1: 推送子模块到远程仓库

> [!IMPORTANT] > **必须先推送子模块,再推送主项目**,否则主项目引用的子模块提交在远程仓库中不存在。

```bash
# 推送 backend
cd /tmp/lexveritas-submodules/lex-veritas-backend
git remote add origin https://github.com/YOUR_USERNAME/lex-veritas-backend.git
git branch -M main
git push -u origin main

# 推送 blockchain
cd /tmp/lexveritas-submodules/lex-veritas-blockchain
git remote add origin https://github.com/YOUR_USERNAME/lex-veritas-blockchain.git
git branch -M main
git push -u origin main

# 推送 frontend
cd /tmp/lexveritas-submodules/lex-veritas-frontend
git remote add origin https://github.com/YOUR_USERNAME/lex-veritas-frontend.git
git branch -M main
git push -u origin main

# 推送 ingestion
cd /tmp/lexveritas-submodules/lex-veritas-ingestion
git remote add origin https://github.com/YOUR_USERNAME/lex-veritas-ingestion.git
git branch -M main
git push -u origin main
```

### 步骤 2: 更新 .gitmodules 文件

将子模块的 URL 从本地路径改为远程仓库地址:

```bash
cd /path/to/LexVeritas

# 编辑 .gitmodules 文件
cat > .gitmodules << 'EOF'
[submodule "lex-veritas-backend"]
	path = lex-veritas-backend
	url = https://github.com/YOUR_USERNAME/lex-veritas-backend.git
[submodule "lex-veritas-blockchain"]
	path = lex-veritas-blockchain
	url = https://github.com/YOUR_USERNAME/lex-veritas-blockchain.git
[submodule "lex-veritas-frontend"]
	path = lex-veritas-frontend
	url = https://github.com/YOUR_USERNAME/lex-veritas-frontend.git
[submodule "lex-veritas-ingestion"]
	path = lex-veritas-ingestion
	url = https://github.com/YOUR_USERNAME/lex-veritas-ingestion.git
EOF
```

> [!TIP]
> 将 `YOUR_USERNAME` 替换为你的 GitHub 用户名或组织名。

### 步骤 3: 同步子模块配置

```bash
# 同步 .gitmodules 的配置到 .git/config
git submodule sync

# 提交 .gitmodules 的更改
git add .gitmodules
git commit -m "Update submodule URLs to remote repositories"
```

### 步骤 4: 推送主项目到远程仓库

```bash
# 添加主项目远程仓库
git remote add origin https://github.com/YOUR_USERNAME/LexVeritas.git

# 推送主项目
git branch -M main
git push -u origin main
```

### 步骤 5: 验证推送结果

```bash
# 查看远程仓库
git remote -v

# 查看子模块状态
git submodule status

# 查看最近的提交
git log --oneline -5
```

---

## 团队协作指南

### 克隆项目(新成员)

```bash
# 方法 1: 克隆时自动初始化子模块
git clone --recursive https://github.com/YOUR_USERNAME/LexVeritas.git

# 方法 2: 先克隆主项目,再初始化子模块
git clone https://github.com/YOUR_USERNAME/LexVeritas.git
cd LexVeritas
git submodule update --init --recursive
```

### 更新项目

```bash
# 更新主项目和所有子模块
git pull --recurse-submodules

# 或者分步更新
git pull                                    # 更新主项目
git submodule update --remote --merge       # 更新子模块
```

### 在子模块中开发

```bash
# 1. 进入子模块目录
cd lex-veritas-backend

# 2. 创建并切换到新分支
git checkout -b feature/new-feature

# 3. 进行开发并提交
git add .
git commit -m "Add new feature"

# 4. 推送子模块分支
git push origin feature/new-feature

# 5. 返回主项目
cd ..

# 6. 更新主项目中的子模块引用
git add lex-veritas-backend
git commit -m "Update backend submodule to include new feature"
git push
```

### 切换子模块分支

```bash
# 进入子模块
cd lex-veritas-backend

# 切换到指定分支
git checkout develop

# 返回主项目并更新引用
cd ..
git add lex-veritas-backend
git commit -m "Switch backend submodule to develop branch"
git push
```

---

## 常见问题

### Q1: 为什么要先推送子模块再推送主项目?

**A**: 主项目的 `.gitmodules` 文件记录了子模块的提交哈希。如果子模块的提交还没有推送到远程仓库,其他人克隆主项目时会找不到对应的子模块提交。

### Q2: 如何查看子模块的当前版本?

```bash
# 查看所有子模块的提交哈希
git submodule status

# 查看特定子模块的详细信息
cd lex-veritas-backend
git log --oneline -1
```

### Q3: 子模块处于 "detached HEAD" 状态怎么办?

这是正常的。子模块默认指向特定的提交,而不是分支。如果需要在子模块中开发:

```bash
cd lex-veritas-backend
git checkout main  # 切换到 main 分支
```

### Q4: 如何删除子模块?

```bash
# 1. 从 .gitmodules 中删除对应条目
git config -f .gitmodules --remove-section submodule.lex-veritas-backend

# 2. 从 .git/config 中删除对应条目
git config -f .git/config --remove-section submodule.lex-veritas-backend

# 3. 从暂存区删除
git rm --cached lex-veritas-backend

# 4. 删除目录
rm -rf lex-veritas-backend
rm -rf .git/modules/lex-veritas-backend

# 5. 提交更改
git commit -m "Remove backend submodule"
```

### Q5: 如何更新所有子模块到最新版本?

```bash
# 更新所有子模块到远程仓库的最新版本
git submodule update --remote --merge

# 提交更新
git add -A
git commit -m "Update all submodules to latest versions"
git push
```

### Q6: 克隆项目后子模块目录是空的?

忘记初始化子模块了:

```bash
git submodule update --init --recursive
```

### Q7: 如何在 CI/CD 中使用子模块?

在 CI/CD 配置中添加:

```yaml
# GitHub Actions 示例
- name: Checkout code
  uses: actions/checkout@v3
  with:
    submodules: recursive

# GitLab CI 示例
variables:
  GIT_SUBMODULE_STRATEGY: recursive
```

---

## 最佳实践

### 1. 提交顺序

✅ **正确顺序**:

1. 在子模块中提交并推送更改
2. 在主项目中更新子模块引用并推送

❌ **错误顺序**:

1. 在主项目中更新子模块引用并推送
2. 在子模块中提交并推送更改

### 2. 分支管理

- 主项目和子模块可以有不同的分支策略
- 建议在主项目中明确记录每个子模块应该使用的分支
- 使用标签(tags)来标记重要的版本组合

### 3. 文档维护

在主项目的 README 中说明:

- 每个子模块的用途
- 推荐的开发工作流
- 如何克隆和更新项目

### 4. 自动化

创建脚本简化常见操作:

```bash
# update-all.sh - 更新所有子模块
#!/bin/bash
git pull
git submodule update --remote --merge
git add -A
git commit -m "Update all submodules"
git push
```

---

## 快速参考

### 常用命令

```bash
# 克隆项目(包含子模块)
git clone --recursive <url>

# 初始化子模块
git submodule update --init --recursive

# 更新子模块
git submodule update --remote --merge

# 查看子模块状态
git submodule status

# 在所有子模块中执行命令
git submodule foreach 'git pull origin main'

# 同步子模块配置
git submodule sync
```

### 配置文件

**`.gitmodules`** - 子模块配置(提交到仓库)

```ini
[submodule "lex-veritas-backend"]
    path = lex-veritas-backend
    url = https://github.com/YOUR_USERNAME/lex-veritas-backend.git
```

**`.git/config`** - 本地 Git 配置(不提交)

```ini
[submodule "lex-veritas-backend"]
    url = https://github.com/YOUR_USERNAME/lex-veritas-backend.git
    active = true
```

---

## 总结

Git 子模块是管理多仓库项目的强大工具。通过本指南,你应该能够:

✅ 从零开始配置 Git 主项目和子模块  
✅ 将项目推送到远程仓库  
✅ 在团队中协作开发  
✅ 解决常见问题

如有任何问题,请参考 [Git 官方文档](https://git-scm.com/book/zh/v2/Git-%E5%B7%A5%E5%85%B7-%E5%AD%90%E6%A8%A1%E5%9D%97)。
