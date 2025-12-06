# 推送 LexVeritas 项目到远程仓库

本文档提供将 LexVeritas 项目(包括所有子模块)推送到远程仓库的详细步骤。

## 前提条件

- ✅ 已完成本地 Git 主项目和子模块配置
- ✅ 在 GitHub/GitLab/Gitee 上创建了以下空仓库:
  - `LexVeritas`
  - `lex-veritas-backend`
  - `lex-veritas-blockchain`
  - `lex-veritas-frontend`
  - `lex-veritas-ingestion`

> [!WARNING]
> 创建远程仓库时,**不要**初始化 README、.gitignore 或 LICENSE,保持仓库为空。

---

## 推送步骤

### 第一步: 推送所有子模块

> [!IMPORTANT] > **必须先推送子模块,再推送主项目**
>
> 原因: 主项目引用了子模块的特定提交哈希。如果子模块的提交还没有推送到远程仓库,其他人克隆主项目时会找不到对应的子模块提交,导致克隆失败。

#### 1.1 推送 backend 子模块

```bash
cd /tmp/lexveritas-submodules/lex-veritas-backend

# 添加远程仓库
git remote add origin https://github.com/nsxzhou/lex-veritas-backend.git

# 确保在 main 分支
git branch -M main

# 推送到远程仓库
git push -u origin main
```

#### 1.2 推送 blockchain 子模块

```bash
cd /tmp/lexveritas-submodules/lex-veritas-blockchain

git remote add origin https://github.com/nsxzhou/lex-veritas-blockchain.git
git branch -M main
git push -u origin main
```

#### 1.3 推送 frontend 子模块

```bash
cd /tmp/lexveritas-submodules/lex-veritas-frontend

git remote add origin https://github.com/nsxzhou/lex-veritas-frontend.git
git branch -M main
git push -u origin main
```

#### 1.4 推送 ingestion 子模块

```bash
cd /tmp/lexveritas-submodules/lex-veritas-ingestion

git remote add origin https://github.com/nsxzhou/lex-veritas-ingestion.git
git branch -M main
git push -u origin main
```

#### 1.5 验证子模块推送

访问 GitHub 确认所有子模块仓库都已成功推送:

- `https://github.com/nsxzhou/lex-veritas-backend`
- `https://github.com/nsxzhou/lex-veritas-blockchain`
- `https://github.com/nsxzhou/lex-veritas-frontend`
- `https://github.com/nsxzhou/lex-veritas-ingestion`

---

### 第二步: 更新 .gitmodules 文件

将子模块的 URL 从本地路径改为远程仓库地址。

```bash
cd /Users/zhouzirui/code/AI/LexVeritas

# 编辑 .gitmodules 文件
cat > .gitmodules << 'EOF'
[submodule "lex-veritas-backend"]
	path = lex-veritas-backend
	url = https://github.com/nsxzhou/lex-veritas-backend.git
[submodule "lex-veritas-blockchain"]
	path = lex-veritas-blockchain
	url = https://github.com/nsxzhou/lex-veritas-blockchain.git
[submodule "lex-veritas-frontend"]
	path = lex-veritas-frontend
	url = https://github.com/nsxzhou/lex-veritas-frontend.git
[submodule "lex-veritas-ingestion"]
	path = lex-veritas-ingestion
	url = https://github.com/nsxzhou/lex-veritas-ingestion.git
EOF
```

> [!TIP] > **使用 SSH 还是 HTTPS?**
>
> - **HTTPS**: `https://github.com/nsxzhou/repo.git` (推荐新手)
> - **SSH**: `git@github.com:nsxzhou/repo.git` (需要配置 SSH 密钥)

---

### 第三步: 同步子模块配置

```bash
# 同步 .gitmodules 的配置到 .git/config
git submodule sync

# 查看同步后的配置
git submodule status
```

---

### 第四步: 提交 .gitmodules 更改

```bash
# 添加 .gitmodules 文件
git add .gitmodules

# 提交更改
git commit -m "Update submodule URLs to remote repositories"

# 查看提交历史
git log --oneline -3
```

---

### 第五步: 推送主项目

```bash
# 添加主项目远程仓库
git remote add origin https://github.com/nsxzhou/LexVeritas.git

# 确保在 main 分支
git branch -M main

# 推送主项目
git push -u origin main
```

---

### 第六步: 验证推送结果

#### 6.1 检查远程仓库

```bash
# 查看远程仓库配置
git remote -v

# 应该显示:
# origin  https://github.com/nsxzhou/LexVeritas.git (fetch)
# origin  https://github.com/nsxzhou/LexVeritas.git (push)
```

#### 6.2 检查子模块状态

```bash
# 查看子模块状态
git submodule status

# 应该显示类似:
# 29b619a442dacb61a75c484d28f143c25bba387b lex-veritas-backend (heads/main)
# f99db77812701730de24334f25974a8ab2e5f3fc lex-veritas-blockchain (heads/main)
# c2ea11339255fe61db9300d83d1ea0b7afe3fe78 lex-veritas-frontend (heads/main)
# aae53dd35105ad0934cca5e4a841a6507fb19785 lex-veritas-ingestion (heads/main)
```

#### 6.3 在 GitHub 上验证

访问主项目仓库: `https://github.com/nsxzhou/LexVeritas`

确认:

- ✅ 所有文件都已推送
- ✅ 子模块显示为特殊的目录图标(带 @ 符号)
- ✅ 点击子模块可以跳转到对应的子模块仓库

---

## 测试克隆

在另一个目录测试克隆项目,确保配置正确:

```bash
# 切换到其他目录
cd /tmp

# 克隆项目(包含子模块)
git clone --recursive https://github.com/nsxzhou/LexVeritas.git

# 进入项目目录
cd LexVeritas

# 检查子模块
git submodule status

# 检查子模块内容
ls -la lex-veritas-backend/
```

如果克隆成功且子模块目录有内容,说明配置正确! ✅

---

## 后续维护

### 更新子模块并推送

当你在子模块中进行了更改:

```bash
# 1. 在子模块中提交并推送
cd lex-veritas-backend
git add .
git commit -m "Add new feature"
git push origin main

# 2. 返回主项目并更新子模块引用
cd ..
git add lex-veritas-backend
git commit -m "Update backend submodule"
git push origin main
```

### 更新主项目文件并推送

当你修改了主项目的文件(如 `docker-compose.yaml`):

```bash
# 添加更改
git add docker-compose.yaml

# 提交
git commit -m "Update docker-compose configuration"

# 推送
git push origin main
```

---

## 常见问题

### Q1: 推送时提示 "Permission denied"

**原因**: 没有仓库的写权限或 SSH 密钥配置问题。

**解决方案**:

- 确认你是仓库的所有者或协作者
- 如果使用 SSH,确保已配置 SSH 密钥
- 尝试使用 HTTPS 方式

### Q2: 推送时提示 "failed to push some refs"

**原因**: 远程仓库有本地没有的提交。

**解决方案**:

```bash
# 先拉取远程更改
git pull --rebase origin main

# 再推送
git push origin main
```

### Q3: 子模块推送失败

**原因**: 子模块的远程仓库地址配置错误。

**解决方案**:

```bash
cd /tmp/lexveritas-submodules/lex-veritas-backend

# 检查远程仓库
git remote -v

# 如果地址错误,删除并重新添加
git remote remove origin
git remote add origin https://github.com/nsxzhou/lex-veritas-backend.git
git push -u origin main
```

### Q4: 如何切换到 SSH 方式?

```bash
# 更新 .gitmodules
sed -i '' 's|https://github.com/|git@github.com:|g' .gitmodules

# 同步配置
git submodule sync

# 提交更改
git add .gitmodules
git commit -m "Switch to SSH URLs for submodules"
git push
```

---

## 快速命令参考

### 一键推送所有子模块(脚本)

创建脚本 `push-submodules.sh`:

```bash
#!/bin/bash

# 子模块列表
SUBMODULES=("backend" "blockchain" "frontend" "ingestion")
USERNAME="nsxzhou"  # 替换为你的 GitHub 用户名

for module in "${SUBMODULES[@]}"; do
    echo "Pushing lex-veritas-${module}..."
    cd "/tmp/lexveritas-submodules/lex-veritas-${module}"

    # 添加远程仓库(如果还没有)
    git remote add origin "https://github.com/${USERNAME}/lex-veritas-${module}.git" 2>/dev/null

    # 推送
    git branch -M main
    git push -u origin main

    echo "✅ lex-veritas-${module} pushed successfully"
    echo ""
done

echo "🎉 All submodules pushed!"
```

使用方法:

```bash
chmod +x push-submodules.sh
./push-submodules.sh
```

---

## 总结

推送流程总结:

1. ✅ 推送所有子模块到远程仓库
2. ✅ 更新 `.gitmodules` 文件为远程 URL
3. ✅ 同步子模块配置
4. ✅ 提交 `.gitmodules` 更改
5. ✅ 推送主项目到远程仓库
6. ✅ 验证推送结果
7. ✅ 测试克隆

完成这些步骤后,你的项目就可以被团队成员克隆和协作开发了! 🚀
