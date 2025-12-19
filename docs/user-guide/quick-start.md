# 快速入门指南

<div align="center">

**5 分钟上手 BigAI Pro API**

</div>

---

## 概述

BigAI Pro 提供与 OpenAI 完全兼容的 API 接口，支持 GPT、Claude、Gemini 等主流大语言模型。只需简单配置，即可无缝接入。

## 第一步：获取 API Key

1. 访问 [BigAI Pro 控制台](https://api.bigaipro.com/console)
2. 注册/登录账户
3. 进入「令牌管理」页面
4. 点击「创建新令牌」
5. 复制生成的 API Key（格式：`sk-xxxxxxxxxxxxxxxx`）

> ⚠️ **安全提示**: API Key 仅显示一次，请妥善保存

## 第二步：配置 API 端点

| 配置项 | 值 |
|--------|-----|
| **API Base URL** | `https://api.bigaipro.com/` |
| **API Key** | `sk-xxxxxxxxxxxxxxxx` |

## 第三步：发起第一个请求

### 使用 cURL

```bash
curl https://api.bigaipro.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-xxxxxxxxxxxxxxxx" \
  -d '{
    "model": "gpt-5.2-instant",
    "messages": [
      {"role": "user", "content": "你好，请介绍一下你自己"}
    ]
  }'
```

### 使用 Python

```python
from openai import OpenAI

client = OpenAI(
    api_key="sk-xxxxxxxxxxxxxxxx",
    base_url="https://api.bigaipro.com/v1"
)

response = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[
        {"role": "user", "content": "你好，请介绍一下你自己"}
    ]
)

print(response.choices[0].message.content)
```

### 使用 Node.js

```javascript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: 'sk-xxxxxxxxxxxxxxxx',
  baseURL: 'https://api.bigaipro.com/v1'
});

const response = await client.chat.completions.create({
  model: 'gpt-5.2-instant',
  messages: [
    { role: 'user', content: '你好，请介绍一下你自己' }
  ]
});

console.log(response.choices[0].message.content);
```

## 支持的模型

BigAI Pro 支持以下主流模型（2025 年 12 月更新）：

### OpenAI 系列
- `gpt-5.2-pro` - 最强旗舰模型
- `gpt-5.2-thinking` - 深度推理版
- `gpt-5.2-instant` - 快速响应版
- `gpt-4.1` - 100 万上下文
- `o3` / `o4-mini` - 推理模型

### Claude 系列
- `claude-sonnet-4.5` - 全球最强编码模型
- `claude-opus-4.5` - 最强综合能力
- `claude-haiku-4.5` - 快速响应版

### Gemini 系列
- `gemini-3.0-pro` - 最新旗舰
- `gemini-2.5-pro` - 100 万上下文
- `gemini-2.5-flash` - 高性价比

### 国产模型
- `qwen3-235b` - 通义千问3 旗舰
- `deepseek-r1` - 推理能力媲美 o1
- `deepseek-chat` - 性价比之王

> 📖 完整模型列表请查看 [模型总览](./models/overview.md)

## 环境变量配置

推荐使用环境变量管理 API 配置：

```bash
# Linux / macOS
export OPENAI_API_KEY="sk-xxxxxxxxxxxxxxxx"
export OPENAI_API_BASE="https://api.bigaipro.com/v1"

# Windows PowerShell
$env:OPENAI_API_KEY="sk-xxxxxxxxxxxxxxxx"
$env:OPENAI_API_BASE="https://api.bigaipro.com/v1"

# Windows CMD
set OPENAI_API_KEY=sk-xxxxxxxxxxxxxxxx
set OPENAI_API_BASE=https://api.bigaipro.com/v1
```

## 下一步

- 📖 [API Key 管理](./api-key-management.md) - 了解如何管理多个密钥
- 💰 [计费与额度](./billing-and-quota.md) - 了解计费规则
- 🤖 [模型使用指南](./models/overview.md) - 深入了解各模型特性
- 🛠️ [开发工具集成](./tools/claude-code.md) - 配置你喜欢的开发工具

---

<div align="center">

**遇到问题？** 查看 [错误处理指南](./advanced/error-handling.md) 或联系客服

</div>
