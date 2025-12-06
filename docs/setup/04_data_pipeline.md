# 数据处理管道

> 法律文档分块、向量化与入库

---

## 📄 文档格式转换

### Pandoc 转换

```bash
# PDF → Markdown
pandoc input.pdf -o output.md --extract-media=./media --wrap=none

# DOCX → Markdown
pandoc input.docx -o output.md -t markdown --wrap=none
```

---

## 🔪 结构化分块

### 核心实现

```python
import re
from typing import Dict, List

def extract_hierarchy(markdown_text: str) -> Dict:
    """提取法律层级结构"""
    hierarchy = {'parts': [], 'chapters': [], 'sections': []}

    lines = markdown_text.split('\n')
    for line in lines:
        if match := re.match(r'^#\s+第[零一二三四五六七八九十百千]+编\s+(.+)', line):
            hierarchy['parts'].append(match.group(0))
        elif match := re.match(r'^##\s+第[零一二三四五六七八九十百千]+章\s+(.+)', line):
            hierarchy['chapters'].append(match.group(0))

    return hierarchy

def split_by_articles(markdown_text: str) -> List[Dict]:
    """按法条分块"""
    articles = []
    current_hierarchy = {'part': None, 'chapter': None, 'section': None}
    current_article = None
    current_lines = []

    for line in markdown_text.split('\n'):
        # 更新层级
        if re.match(r'^#\s+第.+编', line):
            current_hierarchy['part'] = line.strip('#').strip()
        elif re.match(r'^##\s+第.+章', line):
            current_hierarchy['chapter'] = line.strip('#').strip()

        # 检测法条
        if match := re.match(r'^(第[零一二三四五六七八九十百千]+条)\s+(.+)', line):
            if current_article:
                articles.append({
                    'article_number': current_article,
                    'content': '\n'.join(current_lines),
                    'hierarchy': current_hierarchy.copy()
                })
            current_article = match.group(1)
            current_lines = [line]
        elif current_article:
            current_lines.append(line)

    return articles

def extract_references(content: str) -> Dict:
    """提取引用关系"""
    pattern = re.compile(r'(本法|依照)?第[零一二三四五六七八九十百千]+条')
    refs = pattern.findall(content)

    return {
        'direct_refs': list(set(refs)),
        'referenced_by': [],
        'related': []
    }

def chunk_legal_document(markdown_text: str, law_name: str, source_file: str) -> List[Dict]:
    """完整分块流程"""
    articles = split_by_articles(markdown_text)

    chunks = []
    for i, article in enumerate(articles):
        refs = extract_references(article['content'])

        hierarchy_path = ' > '.join(filter(None, [
            article['hierarchy']['part'],
            article['hierarchy']['chapter'],
            article['article_number']
        ]))

        chunks.append({
            'chunk_id': f"{source_file}_{i:04d}",
            'chunk_text': article['content'],
            'chunk_order': i,
            'source_file': source_file,
            'law_name': law_name,
            'law_hierarchy': hierarchy_path,
            'article_number': article['article_number'],
            'references': refs,
        })

    return chunks
```

---

## 🌲 Merkle Tree 构建

```python
from merkletools import MerkleTools
import hashlib

def build_merkle_tree(chunks):
    """构建 Merkle Tree"""
    mt = MerkleTools(hash_type="sha256")

    # 计算叶子节点哈希
    for chunk in chunks:
        chunk_hash = hashlib.sha256(chunk['chunk_text'].encode('utf-8')).hexdigest()
        chunk['chunk_hash'] = chunk_hash
        mt.add_leaf(chunk_hash, do_hash=False)

    # 构建树
    mt.make_tree()
    merkle_root = mt.get_merkle_root()

    # 生成 Merkle Proof
    for i, chunk in enumerate(chunks):
        proof = mt.get_proof(i)
        chunk['merkle_index'] = i
        chunk['merkle_proof'] = {
            'proof': [p['right'] if p['right'] else p['left'] for p in proof],
            'index': i
        }

    return merkle_root, chunks
```

---

## 📊 向量化

```python
from openai import OpenAI
import os

client = OpenAI(api_key=os.getenv('OPENAI_API_KEY'))

def generate_embeddings(chunks):
    """批量生成 Embeddings"""
    texts = [chunk['chunk_text'] for chunk in chunks]

    response = client.embeddings.create(
        model="text-embedding-3-small",
        input=texts
    )

    for i, data in enumerate(response.data):
        chunks[i]['embedding'] = data.embedding

    return chunks
```

---

## 🔄 完整流程

```python
import json
import os

def process_legal_document(pdf_path: str, law_name: str):
    """端到端处理流程"""

    # 1. 转换为 Markdown
    os.system(f"pandoc {pdf_path} -o temp.md --wrap=none")
    with open('temp.md', 'r', encoding='utf-8') as f:
        markdown_text = f.read()

    # 2. 结构化分块
    chunks = chunk_legal_document(
        markdown_text=markdown_text,
        law_name=law_name,
        source_file=pdf_path.replace('.pdf', '')
    )

    # 3. 构建 Merkle Tree
    merkle_root, chunks = build_merkle_tree(chunks)

    # 4. 向量化
    chunks = generate_embeddings(chunks)

    # 5. 上传到区块链
    from blockchain_client import publish_version_to_blockchain
    tx_receipt = publish_version_to_blockchain(
        merkle_root=bytes.fromhex(merkle_root),
        description=f"{law_name} v1.0",
        chunk_count=len(chunks)
    )

    # 6. 存入 Milvus
    from milvus_client import ingest_to_milvus
    ingest_to_milvus(chunks, version_id=tx_receipt['version_id'])

    print(f"✅ 处理完成: {len(chunks)} 个法条")
    print(f"📦 Merkle Root: {merkle_root}")
    print(f"🔗 TX Hash: {tx_receipt['tx_hash']}")

# 使用示例
if __name__ == "__main__":
    process_legal_document(
        pdf_path="./data/民法典.pdf",
        law_name="中华人民共和国民法典"
    )
```

---

## 🔗 相关文档

- [← 向量数据库配置](./03_vector_database.md)
- [后端开发指南 →](./05_backend_development.md)
