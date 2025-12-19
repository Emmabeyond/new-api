# 国产模型使用指南

<div align="center">

**中文优化，本土首选**

*更新于 2025 年 12 月*

</div>

---

## 模型概览

| 厂商 | 模型系列 | 特点 | 适用场景 |
|------|----------|------|----------|
| 阿里云 | Qwen3 (通义千问3) | 开源旗舰，中文最强 | 通用任务、代码 |
| DeepSeek | DeepSeek-R1, V3 | 推理顶尖，性价比极高 | 推理、代码 |
| 智谱 | GLM-4 Plus | 推理能力强 | 复杂任务 |
| 月之暗面 | Moonshot | 超长上下文 | 文档处理 |
| 字节跳动 | 豆包 | 多模态 | 创意内容 |

---

## 通义千问 Qwen3 ⭐

### 模型列表

```
# Qwen3 系列 (2025年最新)
qwen3-235b          # 旗舰版 (235B 参数，22B 激活)
qwen3-32b           # 专业版
qwen3-14b           # 标准版
qwen3-8b            # 轻量版
qwen3-4b            # 极速版

# 特殊版本
qwen3-235b-instruct-2507  # 最新指令优化版
qwen3-thinking      # 思考版本
```

### 核心特性

- **开源旗舰**: Apache 2.0 许可，完全开源
- **中文最强**: 中文理解和生成能力顶尖
- **131K 上下文**: 处理长文档
- **MoE 架构**: 235B 参数，仅 22B 激活，高效推理
- **媲美顶尖**: 与 DeepSeek-R1、o1、Grok-3 竞争

### 基础调用

```python
from openai import OpenAI

client = OpenAI(
    api_key="sk-xxxxxxxx",
    base_url="https://api.bigaipro.com/v1"
)

response = client.chat.completions.create(
    model="qwen3-235b",
    messages=[
        {"role": "system", "content": "你是一个专业的中文助手"},
        {"role": "user", "content": "请用文言文写一首关于人工智能的诗"}
    ],
    max_tokens=2000
)

print(response.choices[0].message.content)
```

### 代码生成

```python
response = client.chat.completions.create(
    model="qwen3-235b",
    messages=[
        {
            "role": "system",
            "content": "你是一位资深的软件工程师，擅长编写高质量代码"
        },
        {
            "role": "user",
            "content": "用 Python 实现一个高性能的异步爬虫框架"
        }
    ],
    max_tokens=4096
)
```

### 定价

| 模型 | 输入 ($/1M tokens) | 输出 ($/1M tokens) |
|------|-------------------|-------------------|
| qwen3-235b | $0.60 | $1.20 |
| qwen3-32b | $0.30 | $0.60 |
| qwen3-14b | $0.15 | $0.30 |
| qwen3-8b | $0.08 | $0.16 |

> 💡 Qwen3-235B 价格不到同级闭源模型的 1/10

---

## DeepSeek ⭐

### 模型列表

```
# DeepSeek-R1 系列 (推理模型)
deepseek-r1         # 推理旗舰，媲美 o1
deepseek-r1-0528    # 最新版本，推理增强

# DeepSeek-V3 系列
deepseek-v3         # 对话模型 (671B 参数)
deepseek-chat       # 通用对话
deepseek-coder      # 代码专用
```

### 核心特性

- **推理顶尖**: DeepSeek-R1 媲美 OpenAI o1
- **性价比极高**: 价格仅为 GPT-4 的 1/100
- **开源友好**: 模型开源，社区活跃
- **代码能力强**: 编程任务表现优异

### DeepSeek-R1 使用

```python
response = client.chat.completions.create(
    model="deepseek-r1",
    messages=[
        {"role": "user", "content": "证明：对于任意正整数 n，n³ - n 能被 6 整除"}
    ]
)
```

### DeepSeek-R1-0528

最新版本，推理能力大幅提升：

```python
response = client.chat.completions.create(
    model="deepseek-r1-0528",
    messages=[
        {"role": "user", "content": "设计一个分布式一致性算法"}
    ]
)
```

### 代码生成

```python
response = client.chat.completions.create(
    model="deepseek-coder",
    messages=[
        {
            "role": "system",
            "content": "你是一个专业的程序员，擅长编写高质量代码"
        },
        {
            "role": "user",
            "content": "用 Rust 实现一个高性能的 HTTP 服务器"
        }
    ],
    max_tokens=4096
)
```

### 定价

| 模型 | 输入 ($/1M tokens) | 输出 ($/1M tokens) |
|------|-------------------|-------------------|
| deepseek-r1 | $0.55 | $2.19 |
| deepseek-v3 | $0.27 | $1.10 |
| deepseek-chat | $0.14 | $0.28 |
| deepseek-coder | $0.14 | $0.28 |

> 💡 DeepSeek 是目前性价比最高的模型之一

---

## 智谱 GLM-4

### 模型列表

```
glm-4-plus          # 增强版
glm-4               # 标准版
glm-4-flash         # 快速版
glm-4v              # 视觉版
```

### 基础调用

```python
response = client.chat.completions.create(
    model="glm-4-plus",
    messages=[
        {"role": "user", "content": "分析中国经济发展的主要趋势"}
    ]
)
```

### 特点

- **推理能力强**: 逻辑推理和数学能力优秀
- **知识丰富**: 中文知识库完善
- **对话自然**: 中文表达流畅

---

## 月之暗面 Moonshot

### 模型列表

```
moonshot-v1-8k      # 8K 上下文
moonshot-v1-32k     # 32K 上下文
moonshot-v1-128k    # 128K 上下文
```

### 基础调用

```python
response = client.chat.completions.create(
    model="moonshot-v1-128k",
    messages=[
        {"role": "user", "content": "总结这篇长文档的要点..."}
    ]
)
```

### 特点

- **超长上下文**: 最高支持 128K
- **文档处理**: 适合长文档分析
- **中文优化**: 中文理解能力强

---

## 模型对比

### 能力对比

| 能力 | Qwen3-235B | DeepSeek-R1 | GLM-4 Plus |
|------|-----------|-------------|------------|
| 中文理解 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| 代码能力 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| 推理能力 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| 长上下文 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| 性价比 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

### 价格对比

| 模型 | 输入 ($/1M) | 输出 ($/1M) | 性价比 |
|------|------------|------------|--------|
| qwen3-235b | $0.60 | $1.20 | ⭐⭐⭐⭐⭐ |
| deepseek-r1 | $0.55 | $2.19 | ⭐⭐⭐⭐⭐ |
| deepseek-chat | $0.14 | $0.28 | ⭐⭐⭐⭐⭐ |
| glm-4-plus | $1.50 | $5.00 | ⭐⭐⭐⭐ |

---

## 选择建议

| 场景 | 推荐模型 | 原因 |
|------|----------|------|
| 中文写作 | qwen3-235b | 中文表达最自然 |
| 代码开发 | deepseek-coder | 代码能力强，价格低 |
| 复杂推理 | deepseek-r1 | 推理能力媲美 o1 |
| 长文档 | moonshot-v1-128k | 超长上下文 |
| 知识问答 | glm-4-plus | 知识库丰富 |
| 高性价比 | deepseek-chat | 价格最低 |
| 开源部署 | qwen3-235b | Apache 2.0 许可 |

---

## 最佳实践

### 1. 中文 Prompt 优化

```python
response = client.chat.completions.create(
    model="qwen3-235b",
    messages=[
        {
            "role": "system",
            "content": """你是一位资深的中文编辑。

写作要求：
1. 语言流畅自然，符合中文表达习惯
2. 适当使用成语和修辞手法
3. 段落结构清晰，逻辑严密
4. 避免翻译腔和生硬表达"""
        },
        {"role": "user", "content": "写一篇关于科技创新的文章"}
    ]
)
```

### 2. 专业领域应用

```python
# 法律文书
response = client.chat.completions.create(
    model="glm-4-plus",
    messages=[
        {
            "role": "system",
            "content": "你是一位专业的法律顾问，熟悉中国法律法规"
        },
        {"role": "user", "content": "分析这份合同的法律风险..."}
    ]
)

# 技术分析
response = client.chat.completions.create(
    model="deepseek-r1",
    messages=[
        {
            "role": "system",
            "content": "你是一位资深的技术架构师"
        },
        {"role": "user", "content": "设计一个高可用的微服务架构"}
    ]
)
```

### 3. 多模型协作

```python
def smart_route(task_type, content):
    """根据任务类型选择最优模型"""
    model_map = {
        "code": "deepseek-coder",
        "reasoning": "deepseek-r1",
        "writing": "qwen3-235b",
        "analysis": "glm-4-plus",
        "long_doc": "moonshot-v1-128k",
        "quick": "deepseek-chat"
    }
    
    return client.chat.completions.create(
        model=model_map.get(task_type, "qwen3-235b"),
        messages=[{"role": "user", "content": content}]
    )
```

---

## 下一步

- 📖 [模型总览](./overview.md) - 查看所有模型
- 📖 [快速入门](../quick-start.md) - 开始使用
- 📖 [成本控制](../best-practices/cost-control.md) - 优化使用成本

---

<div align="center">

**国产大模型，中文更懂你**

</div>
