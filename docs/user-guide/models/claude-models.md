# Claude 系列模型使用指南

<div align="center">

**Anthropic Claude - 全球最强编码模型**

*更新于 2025 年 12 月*

</div>

---

## 模型概览

### Claude 4.5 系列 (最新)

| 模型 | 上下文 | 特点 | 适用场景 |
|------|--------|------|----------|
| claude-sonnet-4.5 | 200K | 最强编码，Agent 首选 | 代码开发、自动化任务 |
| claude-opus-4.5 | 200K | 最强综合能力 | 复杂研究、深度分析 |
| claude-haiku-4.5 | 200K | 快速响应 | 实时应用、高并发 |

### Claude 4 系列

| 模型 | 上下文 | 特点 | 适用场景 |
|------|--------|------|----------|
| claude-4-opus | 200K | 深度推理 | 学术研究 |
| claude-4-sonnet | 200K | 平衡版本 | 通用任务 |

---

## Claude Sonnet 4.5 ⭐

### 为什么选择 Sonnet 4.5？

**全球最强编码模型**：
- SWE-Bench Verified: **77.2%** (业界第一)
- OSWorld 系统使用: **61.4%**
- 最佳 Agent 能力

### 核心特性

- **编码能力顶尖**: 代码生成、理解、调试全面领先
- **Agent 首选**: 最适合自动化任务和工具调用
- **扩展思考**: 支持可见的逐步推理过程
- **200K 上下文**: 处理大型代码库

### 基础调用

```python
from openai import OpenAI

client = OpenAI(
    api_key="sk-xxxxxxxx",
    base_url="https://api.bigaipro.com/v1"
)

response = client.chat.completions.create(
    model="claude-sonnet-4.5",
    messages=[
        {"role": "system", "content": "你是一位资深全栈工程师"},
        {"role": "user", "content": "用 Go 实现一个高性能的分布式任务队列"}
    ],
    max_tokens=8192
)

print(response.choices[0].message.content)
```

### 代码生成示例

```python
response = client.chat.completions.create(
    model="claude-sonnet-4.5",
    messages=[
        {
            "role": "system",
            "content": """你是世界顶级的软件工程师。
生成代码时：
1. 添加详细的中文注释
2. 遵循最佳实践和设计模式
3. 考虑边界情况和错误处理
4. 提供完整可运行的代码
5. 包含单元测试示例"""
        },
        {
            "role": "user",
            "content": "用 Rust 实现一个线程安全的 LRU 缓存，支持 TTL 过期"
        }
    ],
    max_tokens=8192
)
```

### Agent 任务

Sonnet 4.5 是构建 AI Agent 的最佳选择：

```python
tools = [
    {
        "type": "function",
        "function": {
            "name": "execute_code",
            "description": "执行 Python 代码",
            "parameters": {
                "type": "object",
                "properties": {
                    "code": {"type": "string", "description": "要执行的代码"}
                },
                "required": ["code"]
            }
        }
    },
    {
        "type": "function",
        "function": {
            "name": "read_file",
            "description": "读取文件内容",
            "parameters": {
                "type": "object",
                "properties": {
                    "path": {"type": "string", "description": "文件路径"}
                },
                "required": ["path"]
            }
        }
    }
]

response = client.chat.completions.create(
    model="claude-sonnet-4.5",
    messages=[
        {"role": "user", "content": "分析项目中的所有 Python 文件，找出潜在的性能问题"}
    ],
    tools=tools,
    tool_choice="auto"
)
```

---

## Claude Opus 4.5

### 核心特性

- **最强综合能力**: Claude 系列中能力最强
- **深度分析**: 适合复杂研究任务
- **高质量输出**: 文本质量最高

### 适用场景

- 学术研究和论文分析
- 复杂数据分析
- 高质量内容创作
- 专业咨询

```python
response = client.chat.completions.create(
    model="claude-opus-4.5",
    messages=[
        {
            "role": "user",
            "content": """分析以下研究论文，提供：
1. 核心论点和创新点
2. 方法论评估
3. 实验设计的优缺点
4. 与现有研究的对比
5. 未来研究方向建议

论文内容：..."""
        }
    ],
    max_tokens=8192
)
```

---

## Claude Haiku 4.5

### 核心特性

- **极速响应**: 毫秒级延迟
- **成本最低**: 价格仅为 Sonnet 的 1/10
- **能力均衡**: 满足大多数日常需求

### 适用场景

- 实时聊天机器人
- 内容分类和标签
- 快速摘要
- 高并发应用

```python
response = client.chat.completions.create(
    model="claude-haiku-4.5",
    messages=[
        {"role": "user", "content": "用一句话总结：人工智能正在改变..."}
    ],
    max_tokens=100
)
```

---

## 图像理解

Claude 4.5 系列支持强大的图像理解：

```python
import base64

# URL 方式
response = client.chat.completions.create(
    model="claude-sonnet-4.5",
    messages=[
        {
            "role": "user",
            "content": [
                {"type": "text", "text": "分析这张架构图，指出潜在问题"},
                {
                    "type": "image_url",
                    "image_url": {
                        "url": "https://example.com/architecture.png"
                    }
                }
            ]
        }
    ]
)

# Base64 方式
with open("diagram.png", "rb") as f:
    image_data = base64.b64encode(f.read()).decode()

response = client.chat.completions.create(
    model="claude-sonnet-4.5",
    messages=[
        {
            "role": "user",
            "content": [
                {"type": "text", "text": "这个 UML 图有什么问题？"},
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

---

## 流式输出

```python
stream = client.chat.completions.create(
    model="claude-sonnet-4.5",
    messages=[{"role": "user", "content": "实现一个完整的 REST API"}],
    stream=True,
    max_tokens=4096
)

for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="", flush=True)
```

---

## 扩展思考模式

Sonnet 4.5 支持可见的逐步推理：

```python
response = client.chat.completions.create(
    model="claude-sonnet-4.5",
    messages=[
        {
            "role": "user",
            "content": "设计一个高可用的微服务架构，请详细说明你的思考过程"
        }
    ],
    max_tokens=8192
)
```

---

## Claude 特有功能

### XML 标签结构化

Claude 对 XML 标签有很好的理解：

```python
response = client.chat.completions.create(
    model="claude-sonnet-4.5",
    messages=[
        {
            "role": "user",
            "content": """分析以下代码并按格式输出：

<code>
async function fetchData(url) {
    const response = await fetch(url);
    return response.json();
}
</code>

请按以下格式输出：
<analysis>
<security>安全性分析</security>
<performance>性能分析</performance>
<improvements>改进建议</improvements>
<refactored_code>重构后的代码</refactored_code>
</analysis>"""
        }
    ]
)
```

### 长代码库分析

```python
# 读取整个项目
import os

def read_project(directory):
    code = []
    for root, dirs, files in os.walk(directory):
        dirs[:] = [d for d in dirs if d not in ['node_modules', '.git', 'dist']]
        for file in files:
            if file.endswith(('.ts', '.tsx', '.js', '.jsx', '.py', '.go')):
                path = os.path.join(root, file)
                with open(path, 'r', encoding='utf-8') as f:
                    code.append(f"// {path}\n{f.read()}")
    return "\n\n".join(code)

project_code = read_project("./src")

response = client.chat.completions.create(
    model="claude-sonnet-4.5",
    messages=[
        {
            "role": "user",
            "content": f"""分析以下项目代码：

{project_code}

请提供：
1. 项目架构概述
2. 代码质量评估
3. 潜在的 bug 和安全问题
4. 重构建议
5. 测试覆盖建议"""
        }
    ],
    max_tokens=8192
)
```

---

## 参数说明

| 参数 | 说明 | Claude 特点 |
|------|------|-------------|
| max_tokens | 最大输出 | 必须指定，最大 8192 |
| temperature | 随机性 | 默认 1，范围 0-1 |
| top_p | 核采样 | 与 temperature 二选一 |
| stop | 停止序列 | 支持多个停止词 |

---

## 定价

| 模型 | 输入 ($/1M tokens) | 输出 ($/1M tokens) |
|------|-------------------|-------------------|
| claude-sonnet-4.5 | $3.00 | $15.00 |
| claude-opus-4.5 | $15.00 | $75.00 |
| claude-haiku-4.5 | $0.25 | $1.25 |

> 支持 Prompt Caching 最高节省 90%，Batch 处理节省 50%

---

## 模型选择建议

| 任务类型 | 推荐模型 | 原因 |
|----------|----------|------|
| 代码开发 | claude-sonnet-4.5 | 全球最强编码 |
| Agent 任务 | claude-sonnet-4.5 | 最佳 Agent 能力 |
| 学术研究 | claude-opus-4.5 | 最强综合能力 |
| 实时对话 | claude-haiku-4.5 | 响应最快 |
| 代码审查 | claude-sonnet-4.5 | 深度理解代码 |
| 内容创作 | claude-opus-4.5 | 质量最高 |

---

## 最佳实践

### 1. 明确指令

```python
messages = [
    {
        "role": "user",
        "content": """任务：重构以下代码

要求：
- 使用 TypeScript 严格模式
- 添加完整的类型定义
- 遵循 SOLID 原则
- 添加 JSDoc 注释
- 包含错误处理

代码：
...
"""
    }
]
```

### 2. 分步骤处理

```python
messages = [
    {
        "role": "user",
        "content": """请按以下步骤分析这个系统：

步骤 1：理解当前架构
步骤 2：识别性能瓶颈
步骤 3：提出优化方案
步骤 4：给出实施计划
步骤 5：评估风险和收益

系统描述：
...
"""
    }
]
```

### 3. 提供示例

```python
messages = [
    {
        "role": "user",
        "content": """将 API 响应转换为 TypeScript 接口。

示例输入：
{"id": 1, "name": "test", "active": true}

示例输出：
interface Response {
  id: number;
  name: string;
  active: boolean;
}

请转换：
{"users": [{"id": 1, "email": "a@b.com"}], "total": 100}
"""
    }
]
```

---

## 下一步

- 📖 [Gemini 系列](./gemini-models.md) - 超长上下文
- 📖 [Claude Code 教程](../tools/claude-code.md) - 使用 Claude 进行开发
- 📖 [多模态使用](../advanced/multimodal.md) - 图像处理详解

---

<div align="center">

**Claude Sonnet 4.5 - 编码之王**

</div>
