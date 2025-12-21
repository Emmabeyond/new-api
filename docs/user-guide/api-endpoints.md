# API 端点说明

<div align="center">

**多端点部署，灵活选择**

</div>

---

## 概述

BigAI Pro 提供多个 API 端点，适应不同的使用场景和网络环境。所有端点功能完全一致，可根据需求灵活选择。

## 可用端点

### 主站点（推荐）

| 配置项 | 值 |
|--------|-----|
| **地址** | `https://api.bigaipro.com` |
| **API Base** | `https://api.bigaipro.com/v1` |
| **特点** | 功能完整，支持全部服务 |
| **防护** | Cloudflare CDN 加速 |

适用场景：
- 海外访问
- 需要完整控制台功能
- 浏览器直接访问

### 国内直连端点（新增）

| 配置项 | 值 |
|--------|-----|
| **地址** | `https://api.abu117.cn` |
| **API Base** | `https://api.abu117.cn/v1` |
| **特点** | 国内直连，无 CDN，低延迟 |
| **限制** | 仅支持 API 调用 |

适用场景：
- 国内服务器部署
- 生产环境 API 调用
- 需要更低延迟的场景
- 避免 Cloudflare 检测的自动化脚本

> **注意**: 该端点仅支持 API 调用，直接访问 `https://api.abu117.cn/` 会进入令牌使用记录查询页面

## 端点功能对比

| 功能 | api.bigaipro.com | api.abu117.cn |
|------|------------------|---------------|
| Chat Completions | ✅ | ✅ |
| Embeddings | ✅ | ✅ |
| Images | ✅ | ✅ |
| Audio | ✅ | ✅ |
| 控制台 | ✅ | ❌ |
| 网页登录 | ✅ | ❌ |
| CDN 加速 | ✅ (Cloudflare) | ❌ |
| 国内直连 | ❌ | ✅ |

## 使用示例

### 使用主站点

```python
from openai import OpenAI

client = OpenAI(
    api_key="sk-xxxxxxxxxxxxxxxx",
    base_url="https://api.bigaipro.com/v1"
)
```

### 使用国内直连端点

```python
from openai import OpenAI

client = OpenAI(
    api_key="sk-xxxxxxxxxxxxxxxx",
    base_url="https://api.abu117.cn/v1"  # 国内直连
)
```

### cURL 示例

```bash
# 使用主站点
curl https://api.bigaipro.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-xxxxxxxxxxxxxxxx" \
  -d '{"model": "gpt-4o", "messages": [{"role": "user", "content": "Hello"}]}'

# 使用国内直连
curl https://api.abu117.cn/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-xxxxxxxxxxxxxxxx" \
  -d '{"model": "gpt-4o", "messages": [{"role": "user", "content": "Hello"}]}'
```

### 环境变量配置

```bash
# 主站点
export OPENAI_API_BASE="https://api.bigaipro.com/v1"

# 国内直连
export OPENAI_API_BASE="https://api.abu117.cn/v1"
```

## 令牌使用记录查询

直接访问 `https://api.abu117.cn/` 可以查询令牌使用记录：

**功能说明**:
- 查询指定令牌的使用记录
- 支持时间范围筛选
- 支持导出 CSV 文件

**访问地址**: [api.abu117.cn](https://api.abu117.cn/)

### 导出 CSV

1. 访问 `https://api.abu117.cn/`
2. 输入您的令牌进行查询
3. 选择需要导出的时间范围
4. 点击「导出 CSV」下载记录文件

## 端点选择建议

| 场景 | 推荐端点 | 原因 |
|------|----------|------|
| 日常开发测试 | api.bigaipro.com | 功能完整 |
| 国内生产部署 | api.abu117.cn | 低延迟，直连 |
| Claude Code / Cursor | 两者皆可 | 根据网络环境选择 |
| 自动化脚本 | api.abu117.cn | 避免 CF 验证 |
| 需要控制台 | api.bigaipro.com | 唯一选择 |

## 故障转移

建议在代码中实现端点故障转移：

```python
from openai import OpenAI

ENDPOINTS = [
    "https://api.abu117.cn/v1",      # 优先国内直连
    "https://api.bigaipro.com/v1",   # 备用主站点
]

def create_client():
    for base_url in ENDPOINTS:
        try:
            client = OpenAI(
                api_key="sk-xxxxxxxxxxxxxxxx",
                base_url=base_url,
                timeout=10
            )
            # 测试连接
            client.models.list()
            return client
        except Exception:
            continue
    raise Exception("All endpoints unavailable")

client = create_client()
```

## 常见问题

### Q: 两个端点的 API Key 通用吗？

A: 是的，API Key 完全通用，可以在任意端点使用。

### Q: 国内直连端点稳定吗？

A: 稳定性与主站点一致，且对国内用户延迟更低。

### Q: 为什么访问 api.abu117.cn 显示的是查询页面？

A: 这是设计如此。该地址直接访问会进入令牌使用记录查询功能，API 调用需要使用 `/v1` 路径。

### Q: 数据是同步的吗？

A: 是的，两个端点连接同一后端，数据完全同步。

## 下一步

- 🚀 [快速入门](./quick-start.md) - 开始使用 API
- ⭐ [用户等级](./user-levels.md) - 了解等级权益
- 💰 [计费与额度](./billing-and-quota.md) - 了解计费规则

---

<div align="center">

**选择合适的端点，获得最佳体验**

</div>
