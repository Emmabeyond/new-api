# API Key 管理指南

<div align="center">

**安全高效地管理你的 API 密钥**

</div>

---

## 概述

API Key 是访问 BigAI Pro 服务的凭证。合理管理 API Key 可以提升安全性、便于追踪使用情况、控制访问权限。

## 创建 API Key

### 步骤

1. 登录 [BigAI Pro 控制台](https://api.bigaipro.com/console)
2. 进入「令牌管理」页面
3. 点击「创建新令牌」
4. 配置令牌属性：

| 属性 | 说明 |
|------|------|
| **名称** | 便于识别的描述性名称，如 `production-server`、`dev-local` |
| **额度限制** | 设置该令牌的最大使用额度（可选） |
| **过期时间** | 设置令牌有效期（可选） |
| **IP 白名单** | 限制可使用该令牌的 IP 地址（可选） |
| **模型权限** | 限制可访问的模型列表（可选） |

5. 点击「确认创建」
6. **立即复制并保存** 生成的 API Key

> ⚠️ **重要**: API Key 仅在创建时显示一次，无法再次查看。如果丢失，需要重新创建。

## API Key 格式

BigAI Pro 的 API Key 格式：

```
sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

- 前缀 `sk-` 表示 Secret Key
- 后续为 32 位随机字符串

## 最佳实践

### 1. 环境隔离

为不同环境创建独立的 API Key：

| 环境 | 命名建议 | 用途 |
|------|----------|------|
| 开发 | `dev-local` | 本地开发调试 |
| 测试 | `staging-server` | 测试环境 |
| 生产 | `prod-server-01` | 生产服务器 |

### 2. 权限最小化

- 为每个应用/服务创建独立的 Key
- 设置合理的额度限制
- 仅授予必要的模型访问权限

### 3. 定期轮换

建议每 90 天轮换一次 API Key：

1. 创建新的 API Key
2. 更新应用配置使用新 Key
3. 验证新 Key 工作正常
4. 删除旧的 API Key

### 4. 安全存储

**✅ 推荐做法**

```bash
# 使用环境变量
export OPENAI_API_KEY="sk-xxxxxxxx"

# 使用 .env 文件（添加到 .gitignore）
echo "OPENAI_API_KEY=sk-xxxxxxxx" >> .env

# 使用密钥管理服务
# AWS Secrets Manager / HashiCorp Vault / Azure Key Vault
```

**❌ 避免做法**

```python
# 不要硬编码在代码中
api_key = "sk-xxxxxxxx"  # 危险！

# 不要提交到版本控制
# .env 文件应该在 .gitignore 中
```

## 查看使用情况

在控制台可以查看每个 API Key 的：

- **调用次数**: 总请求数量
- **Token 消耗**: 输入/输出 Token 统计
- **费用统计**: 累计使用费用
- **最近活动**: 最后使用时间和 IP

## 管理操作

### 禁用 Key

临时禁用而不删除：

1. 进入「令牌管理」
2. 找到目标令牌
3. 点击「禁用」开关

### 删除 Key

永久删除令牌：

1. 进入「令牌管理」
2. 找到目标令牌
3. 点击「删除」按钮
4. 确认删除操作

> ⚠️ 删除操作不可恢复，请确保该 Key 不再使用

### 修改配置

可修改的属性：

- 名称
- 额度限制
- 过期时间
- IP 白名单
- 模型权限

## 故障排查

### 常见错误

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| `401` | API Key 无效 | 检查 Key 是否正确、是否已删除 |
| `403` | 权限不足 | 检查模型权限、IP 白名单设置 |
| `429` | 超出限制 | 检查额度限制、请求频率 |

### 验证 Key 有效性

```bash
curl https://api.bigaipro.com/v1/models \
  -H "Authorization: Bearer sk-xxxxxxxx"
```

成功响应表示 Key 有效。

## 下一步

- 💰 [计费与额度](./billing-and-quota.md) - 了解费用计算
- 🔒 [安全指南](./best-practices/security.md) - 深入安全最佳实践
- 🚀 [快速入门](./quick-start.md) - 开始使用 API

---

<div align="center">

**安全第一** - 妥善保管你的 API Key

</div>
