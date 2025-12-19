# Gemini 系列模型使用指南

<div align="center">

**Google Gemini - 原生多模态 AI**

*更新于 2025 年 12 月*

</div>

---

## 模型概览

### Gemini 3.0 系列 (最新)

| 模型 | 上下文 | 特点 | 适用场景 |
|------|--------|------|----------|
| gemini-3.0-pro | 1M+ | 最新旗舰 | 复杂任务、多模态 |
| gemini-3.0-deep-think | 1M+ | 深度思考 | 复杂推理 |

### Gemini 2.5 系列

| 模型 | 上下文 | 特点 | 适用场景 |
|------|--------|------|----------|
| gemini-2.5-pro | 1M | 专业版 | 长文档、复杂任务 |
| gemini-2.5-flash | 1M | 快速版 | 高性价比、实时应用 |
| gemini-2.5-flash-lite | 1M | 极速版 | 最低成本、高吞吐 |

---

## Gemini 3.0 Pro

### 核心特性

- **稀疏 MoE 架构**: 更高效的计算
- **64K 输出**: 支持超长输出
- **原生多模态**: 文本、图像、音频、视频统一处理
- **深度思考**: 可控的推理深度

### 基础调用

```python
from openai import OpenAI

client = OpenAI(
    api_key="sk-xxxxxxxx",
    base_url="https://api.bigaipro.com/v1"
)

response = client.chat.completions.create(
    model="gemini-3.0-pro",
    messages=[
        {"role": "user", "content": "分析人工智能对未来社会的影响"}
    ],
    max_tokens=8192
)

print(response.choices[0].message.content)
```

### 深度思考模式

```python
response = client.chat.completions.create(
    model="gemini-3.0-deep-think",
    messages=[
        {"role": "user", "content": "证明黎曼猜想的一个特殊情况"}
    ]
)
```

---

## Gemini 2.5 系列

### 核心优势

#### 超长上下文

Gemini 支持高达 **100 万 Token** 的上下文窗口：

- 处理整本书籍
- 分析数小时的视频
- 理解大型代码库

#### 思考模型

Gemini 2.5 是"思考模型"，可以在响应前进行推理：

```python
# 可以控制思考预算
response = client.chat.completions.create(
    model="gemini-2.5-pro",
    messages=[
        {"role": "user", "content": "设计一个分布式数据库的一致性协议"}
    ]
)
```

#### 极致性价比

- Flash-Lite 版本价格极低
- 适合大规模应用

---

## Gemini 2.5 Pro

### 特点

- **100 万上下文**: 业界最长
- **思考能力**: 可控的推理深度
- **多模态支持**: 图像、音频、视频

### 长文档处理

```python
# 处理超长文档
with open("book.txt", "r") as f:
    book_content = f.read()  # 可以是整本书

response = client.chat.completions.create(
    model="gemini-2.5-pro",
    messages=[
        {
            "role": "user",
            "content": f"""请阅读以下书籍内容并回答问题：

{book_content}

问题：
1. 这本书的主要论点是什么？
2. 作者使用了哪些论证方法？
3. 你认为最有说服力的章节是哪一章？为什么？"""
        }
    ],
    max_tokens=8192
)
```

### 代码库分析

```python
import os

def read_codebase(directory):
    code = []
    for root, dirs, files in os.walk(directory):
        dirs[:] = [d for d in dirs if d not in ['node_modules', '.git']]
        for file in files:
            if file.endswith(('.py', '.js', '.ts', '.go', '.rs')):
                path = os.path.join(root, file)
                with open(path, 'r') as f:
                    code.append(f"// {path}\n{f.read()}")
    return "\n\n".join(code)

codebase = read_codebase("./src")

response = client.chat.completions.create(
    model="gemini-2.5-pro",
    messages=[
        {
            "role": "user",
            "content": f"""分析以下代码库：

{codebase}

请提供：
1. 架构概述
2. 主要模块功能
3. 潜在的改进建议"""
        }
    ]
)
```

---

## Gemini 2.5 Flash

### 特点

- **极速响应**: 毫秒级延迟
- **超低成本**: 最具性价比
- **思考能力**: 支持可控推理

### 适用场景

- 实时聊天
- 内容分类
- 快速摘要
- 高并发服务

```python
response = client.chat.completions.create(
    model="gemini-2.5-flash",
    messages=[
        {"role": "user", "content": "用三句话总结人工智能的发展历程"}
    ],
    max_tokens=200
)
```

---

## Gemini 2.5 Flash-Lite

### 特点

- **最低延迟**: 首 Token 时间最短
- **最低成本**: 价格最优
- **高吞吐**: 适合大规模处理
- **思考可选**: 默认关闭思考，可按需开启

### 适用场景

- 大规模分类任务
- 内容摘要
- 高吞吐应用

```python
response = client.chat.completions.create(
    model="gemini-2.5-flash-lite",
    messages=[
        {"role": "user", "content": "分类这条评论的情感：这个产品太棒了！"}
    ],
    max_tokens=50
)
```

---

## 多模态使用

### 图像理解

```python
import base64

# URL 方式
response = client.chat.completions.create(
    model="gemini-2.5-pro",
    messages=[
        {
            "role": "user",
            "content": [
                {"type": "text", "text": "详细描述这张图片"},
                {
                    "type": "image_url",
                    "image_url": {
                        "url": "https://example.com/image.jpg"
                    }
                }
            ]
        }
    ]
)

# Base64 方式
with open("chart.png", "rb") as f:
    image_data = base64.b64encode(f.read()).decode()

response = client.chat.completions.create(
    model="gemini-2.5-pro",
    messages=[
        {
            "role": "user",
            "content": [
                {"type": "text", "text": "分析这个图表的数据趋势"},
                {
                    "type": "image_url",
                    "image_url": {
                        "url": f"data:image/png;base64,{image_data}"
                    }
                }
            ]
        }
    ]
)
```

### 多图对比

```python
response = client.chat.completions.create(
    model="gemini-2.5-pro",
    messages=[
        {
            "role": "user",
            "content": [
                {"type": "text", "text": "比较这两张设计图的优缺点"},
                {
                    "type": "image_url",
                    "image_url": {"url": "https://example.com/design1.jpg"}
                },
                {
                    "type": "image_url",
                    "image_url": {"url": "https://example.com/design2.jpg"}
                }
            ]
        }
    ]
)
```

### 视频理解

Gemini 原生支持视频理解：

```python
response = client.chat.completions.create(
    model="gemini-2.5-pro",
    messages=[
        {
            "role": "user",
            "content": [
                {"type": "text", "text": "总结这个视频的主要内容"},
                {
                    "type": "video_url",
                    "video_url": {
                        "url": "https://example.com/video.mp4"
                    }
                }
            ]
        }
    ]
)
```

### 音频处理

```python
response = client.chat.completions.create(
    model="gemini-2.5-pro",
    messages=[
        {
            "role": "user",
            "content": [
                {"type": "text", "text": "转录并总结这段音频"},
                {
                    "type": "audio_url",
                    "audio_url": {
                        "url": "https://example.com/audio.mp3"
                    }
                }
            ]
        }
    ]
)
```

---

## 流式输出

```python
stream = client.chat.completions.create(
    model="gemini-2.5-pro",
    messages=[{"role": "user", "content": "写一篇科技文章"}],
    stream=True
)

for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="", flush=True)
```

---

## 参数说明

| 参数 | 说明 | Gemini 特点 |
|------|------|-------------|
| max_tokens | 最大输出 | 3.0 Pro 支持 64K |
| temperature | 随机性 | 范围 0-2 |
| top_p | 核采样 | 范围 0-1 |
| top_k | Top-K 采样 | Gemini 特有 |
| stop | 停止序列 | 支持多个 |

---

## 定价

| 模型 | 输入 ($/1M tokens) | 输出 ($/1M tokens) |
|------|-------------------|-------------------|
| gemini-3.0-pro | $1.25 | $5.00 |
| gemini-2.5-pro | $1.25 | $5.00 |
| gemini-2.5-flash | $0.075 | $0.30 |
| gemini-2.5-flash-lite | $0.02 | $0.08 |

---

## 模型选择建议

| 任务类型 | 推荐模型 | 原因 |
|----------|----------|------|
| 长文档分析 | gemini-2.5-pro | 100 万上下文 |
| 视频理解 | gemini-2.5-pro | 原生多模态 |
| 实时对话 | gemini-2.5-flash | 速度快、成本低 |
| 大规模处理 | gemini-2.5-flash-lite | 极致性价比 |
| 最新功能 | gemini-3.0-pro | 最强能力 |
| 复杂推理 | gemini-3.0-deep-think | 深度思考 |

---

## 最佳实践

### 1. 利用长上下文

```python
response = client.chat.completions.create(
    model="gemini-2.5-pro",
    messages=[
        {
            "role": "user",
            "content": f"""以下是过去一年的所有会议记录：

{all_meeting_notes}

请：
1. 总结每个季度的主要决策
2. 识别反复出现的问题
3. 提取所有待办事项及其状态"""
        }
    ]
)
```

### 2. 多模态组合

```python
response = client.chat.completions.create(
    model="gemini-2.5-pro",
    messages=[
        {
            "role": "user",
            "content": [
                {"type": "text", "text": "根据这份产品设计图和需求文档，评估实现难度"},
                {"type": "image_url", "image_url": {"url": design_image_url}},
                {"type": "text", "text": f"需求文档：\n{requirements_doc}"}
            ]
        }
    ]
)
```

### 3. 批量处理

```python
import asyncio

async def process_batch(items):
    tasks = []
    for item in items:
        task = client.chat.completions.create(
            model="gemini-2.5-flash-lite",
            messages=[{"role": "user", "content": f"分类：{item}"}]
        )
        tasks.append(task)
    return await asyncio.gather(*tasks)
```

---

## 下一步

- 📖 [国产模型](./chinese-models.md) - 了解国产大模型
- 📖 [多模态使用](../advanced/multimodal.md) - 深入多模态处理
- 📖 [流式响应](../advanced/streaming.md) - 流式输出详解

---

<div align="center">

**Gemini - 多模态 AI 的未来**

</div>
