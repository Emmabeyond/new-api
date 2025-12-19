# GPT 系列模型使用指南

<div align="center">

**OpenAI GPT 系列 - AI 行业标杆**

*更新于 2025 年 12 月*

</div>

---

## 模型概览

### GPT-5.2 系列 (最新)

| 模型 | 上下文 | 特点 | 适用场景 |
|------|--------|------|----------|
| gpt-5.2-pro | 1M | 最强综合能力 | 复杂任务、研究分析 |
| gpt-5.2-thinking | 1M | 深度推理 | 数学、科学、逻辑 |
| gpt-5.2-instant | 1M | 快速响应 | 日常对话、实时应用 |

### GPT-4.1 系列

| 模型 | 上下文 | 特点 | 适用场景 |
|------|--------|------|----------|
| gpt-4.1 | 1M | 全能旗舰 | 通用任务、长文档 |
| gpt-4.1-mini | 1M | 轻量高效 | 高性价比场景 |
| gpt-4.1-nano | 128K | 极速响应 | 实时应用 |

### o 系列推理模型

| 模型 | 上下文 | 特点 | 适用场景 |
|------|--------|------|----------|
| o3 | 200K | 最强推理 | 复杂数学、编程、科研 |
| o3-pro | 200K | 推理增强 | 需要最高可靠性的任务 |
| o4-mini | 200K | 推理轻量 | 快速推理、成本敏感 |

---

## GPT-5.2 系列

### 核心特性

- **PhD 级智能**: 在专业领域达到博士水平
- **100 万上下文**: 处理超长文档和代码库
- **原生多模态**: 文本、图像、音频、视频统一处理
- **知识截止**: 2025 年 8 月

### GPT-5.2 Pro

最强综合能力，适合复杂任务：

```python
from openai import OpenAI

client = OpenAI(
    api_key="sk-xxxxxxxx",
    base_url="https://api.bigaipro.com/v1"
)

response = client.chat.completions.create(
    model="gpt-5.2-pro",
    messages=[
        {"role": "system", "content": "你是一位资深的技术架构师"},
        {"role": "user", "content": "设计一个支持千万级用户的实时消息系统架构"}
    ],
    max_tokens=4096
)

print(response.choices[0].message.content)
```

### GPT-5.2 Thinking

深度推理模式，适合复杂逻辑任务：

```python
response = client.chat.completions.create(
    model="gpt-5.2-thinking",
    messages=[
        {"role": "user", "content": "证明：对于任意正整数 n，n⁵ - n 能被 30 整除"}
    ]
)
```

### GPT-5.2 Instant

快速响应，适合实时应用：

```python
response = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[
        {"role": "user", "content": "用一句话解释量子纠缠"}
    ],
    max_tokens=100
)
```

---

## GPT-4.1 系列

### 核心特性

- **100 万上下文**: 业界领先的长上下文能力
- **指令遵循增强**: 更精准地执行复杂指令
- **代码能力提升**: 编程任务表现大幅提升
- **成本优化**: 相比 GPT-4o 更具性价比

### 基础调用

```python
response = client.chat.completions.create(
    model="gpt-4.1",
    messages=[
        {"role": "system", "content": "你是一个专业的代码审查专家"},
        {"role": "user", "content": "审查以下代码的安全性和性能..."}
    ],
    max_tokens=4096
)
```

### 长文档处理

```python
# 处理整个代码库
with open("large_codebase.txt", "r") as f:
    codebase = f.read()  # 可达数十万行代码

response = client.chat.completions.create(
    model="gpt-4.1",
    messages=[
        {
            "role": "user",
            "content": f"""分析以下代码库：

{codebase}

请提供：
1. 架构概述
2. 潜在的安全漏洞
3. 性能优化建议"""
        }
    ]
)
```

---

## o3 系列推理模型

### 核心特性

- **深度推理**: 专为复杂推理任务设计
- **工具调用**: 支持网页浏览、代码执行
- **视觉推理**: 可以"思考图像"
- **SWE-Bench 69.1%**: 代码任务业界领先

### o3 使用

```python
# o3 支持完整的工具调用
response = client.chat.completions.create(
    model="o3",
    messages=[
        {"role": "user", "content": "分析这个算法的时间复杂度，并给出优化方案"}
    ]
)
```

### o4-mini 使用

更快速、更经济的推理模型：

```python
response = client.chat.completions.create(
    model="o4-mini",
    messages=[
        {"role": "user", "content": "解这道数学题：..."}
    ]
)
```

### o3 vs o4-mini

| 特性 | o3 | o4-mini |
|------|-----|---------|
| 推理深度 | 最深 | 适中 |
| 响应速度 | 较慢 | 较快 |
| 价格 | 较高 | 较低 |
| 工具调用 | ✅ | ✅ |
| 视觉推理 | ✅ | ✅ |
| 适用场景 | 复杂科研/数学 | 日常推理/编程 |

---

## 多模态能力

### 图像理解

```python
response = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[
        {
            "role": "user",
            "content": [
                {"type": "text", "text": "详细分析这张图片"},
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
```

### Base64 图像

```python
import base64

with open("image.png", "rb") as f:
    image_data = base64.b64encode(f.read()).decode()

response = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[
        {
            "role": "user",
            "content": [
                {"type": "text", "text": "这张图表显示了什么趋势？"},
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

### 视频理解

GPT-5.2 支持视频输入：

```python
response = client.chat.completions.create(
    model="gpt-5.2-pro",
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

---

## 流式输出

```python
stream = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[{"role": "user", "content": "写一篇关于 AI 未来的文章"}],
    stream=True
)

for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="", flush=True)
```

---

## Function Calling

```python
tools = [
    {
        "type": "function",
        "function": {
            "name": "get_weather",
            "description": "获取指定城市的天气",
            "parameters": {
                "type": "object",
                "properties": {
                    "city": {"type": "string", "description": "城市名称"},
                    "unit": {"type": "string", "enum": ["celsius", "fahrenheit"]}
                },
                "required": ["city"]
            }
        }
    }
]

response = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[{"role": "user", "content": "北京今天天气怎么样？"}],
    tools=tools,
    tool_choice="auto"
)
```

---

## 参数说明

| 参数 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| model | string | 模型名称 | 必填 |
| messages | array | 对话消息 | 必填 |
| temperature | float | 随机性 (0-2) | 1 |
| max_tokens | int | 最大输出长度 | 模型默认 |
| top_p | float | 核采样 (0-1) | 1 |
| frequency_penalty | float | 频率惩罚 (-2 to 2) | 0 |
| presence_penalty | float | 存在惩罚 (-2 to 2) | 0 |
| stop | array | 停止序列 | null |
| stream | bool | 流式输出 | false |

---

## 模型选择建议

| 任务类型 | 推荐模型 | 原因 |
|----------|----------|------|
| 日常对话 | gpt-5.2-instant | 快速智能 |
| 复杂分析 | gpt-5.2-pro | 最强能力 |
| 数学推理 | gpt-5.2-thinking / o3 | 深度推理 |
| 代码开发 | gpt-4.1 / o3 | 代码能力强 |
| 长文档 | gpt-4.1 | 100万上下文 |
| 高性价比 | gpt-4.1-mini | 能力强价格低 |
| 实时应用 | gpt-4.1-nano | 极速响应 |

---

## 最佳实践

### 1. System Prompt 优化

```python
messages = [
    {
        "role": "system",
        "content": """你是一位资深的全栈工程师。

规则：
1. 代码要有详细注释
2. 遵循最佳实践和设计模式
3. 考虑安全性和性能
4. 给出完整可运行的代码"""
    },
    {"role": "user", "content": "..."}
]
```

### 2. 上下文管理

```python
def manage_context(messages, max_messages=20):
    """保持对话历史，但控制长度"""
    if len(messages) > max_messages:
        return [messages[0]] + messages[-(max_messages-1):]
    return messages
```

### 3. 错误重试

```python
import time
from openai import RateLimitError, APIError

def call_with_retry(func, max_retries=3):
    for i in range(max_retries):
        try:
            return func()
        except RateLimitError:
            time.sleep(2 ** i)
        except APIError as e:
            if i == max_retries - 1:
                raise
            time.sleep(1)
    raise Exception("Max retries exceeded")
```

---

## 下一步

- 📖 [Claude 系列](./claude-models.md) - 最强编码模型
- 📖 [多模态使用](../advanced/multimodal.md) - 深入图像/视频处理
- 📖 [函数调用](../advanced/function-calling.md) - Function Calling 详解

---

<div align="center">

**GPT - 定义 AI 的未来**

</div>
