# Cursor IDE 配置指南

<div align="center">

**AI 原生的代码编辑器**

</div>

---

## 简介

Cursor 是一款 AI 原生的代码编辑器，基于 VS Code 构建，内置强大的 AI 编程能力。通过配置 BigAI Pro API，你可以使用各种顶尖模型进行编程。

---

## 安装 Cursor

### 下载安装

1. 访问 [Cursor 官网](https://cursor.sh)
2. 下载对应系统的安装包
3. 安装并启动 Cursor

### 系统要求

- macOS 10.15+
- Windows 10+
- Linux (Ubuntu 18.04+, Fedora 32+)

---

## 配置 BigAI Pro API

### 步骤一：打开设置

1. 启动 Cursor
2. 按 `Cmd+,` (macOS) 或 `Ctrl+,` (Windows/Linux) 打开设置
3. 搜索 "OpenAI" 或 "API"

### 步骤二：配置 API

在设置中找到以下选项：

```
OpenAI API Key: sk-xxxxxxxxxxxxxxxx
OpenAI API Base URL: https://api.bigaipro.com/v1
```

### 步骤三：选择模型

在模型设置中选择：

- `gpt-5.2-instant` - 日常编程
- `gpt-5.2-pro` - 复杂任务
- `claude-sonnet-4.5` - 代码生成
- `o3` - 推理任务

### 配置文件方式

也可以直接编辑配置文件 `~/.cursor/settings.json`：

```json
{
  "openai.apiKey": "sk-xxxxxxxxxxxxxxxx",
  "openai.apiBaseUrl": "https://api.bigaipro.com/v1",
  "cursor.model": "gpt-5.2-instant",
  "cursor.chat.model": "claude-sonnet-4.5",
  "cursor.autocomplete.model": "gpt-4.1-mini"
}
```

---

## 核心功能

### 1. AI Chat（对话）

按 `Cmd+L` (macOS) 或 `Ctrl+L` (Windows/Linux) 打开 AI 对话：

```
你：解释这段代码的作用
AI：这段代码实现了...

你：如何优化这个函数的性能？
AI：可以从以下几个方面优化...
```

### 2. Inline Edit（行内编辑）

选中代码后按 `Cmd+K` (macOS) 或 `Ctrl+K` (Windows/Linux)：

```
选中代码 → Cmd+K → 输入指令 → AI 修改代码
```

示例指令：
- "添加错误处理"
- "转换为 async/await"
- "添加类型注解"
- "优化性能"

### 3. Autocomplete（自动补全）

Cursor 会自动提供 AI 驱动的代码补全：

```python
def calculate_fibonacci(n):
    # 输入注释，AI 自动补全实现
    |  # 光标位置，按 Tab 接受建议
```

### 4. Composer（作曲家模式）

按 `Cmd+I` (macOS) 或 `Ctrl+I` (Windows/Linux) 打开 Composer：

- 多文件编辑
- 项目级重构
- 功能实现

```
你：创建一个用户认证系统，包括注册、登录、JWT 验证
AI：我将创建以下文件...
```

---

## 快捷键

| 功能 | macOS | Windows/Linux |
|------|-------|---------------|
| AI 对话 | `Cmd+L` | `Ctrl+L` |
| 行内编辑 | `Cmd+K` | `Ctrl+K` |
| Composer | `Cmd+I` | `Ctrl+I` |
| 接受建议 | `Tab` | `Tab` |
| 拒绝建议 | `Esc` | `Esc` |
| 下一个建议 | `Alt+]` | `Alt+]` |
| 上一个建议 | `Alt+[` | `Alt+[` |

---

## 高级配置

### 模型配置

```json
{
  "cursor.chat.model": "claude-sonnet-4.5",
  "cursor.autocomplete.model": "gpt-4.1-mini",
  "cursor.composer.model": "gpt-5.2-pro",
  "cursor.inline.model": "gpt-5.2-instant"
}
```

### 上下文配置

```json
{
  "cursor.context.maxFiles": 20,
  "cursor.context.maxTokens": 100000,
  "cursor.context.includePatterns": ["**/*.ts", "**/*.tsx"],
  "cursor.context.excludePatterns": ["node_modules", "dist"]
}
```

### 自动补全配置

```json
{
  "cursor.autocomplete.enabled": true,
  "cursor.autocomplete.delay": 100,
  "cursor.autocomplete.maxSuggestions": 3,
  "cursor.autocomplete.triggerCharacters": [".", "(", "{", "["]
}
```

---

## 使用技巧

### 1. 高效对话

```
# 提供上下文
@file:src/main.ts 解释这个文件的架构

# 引用代码
@selection 优化选中的代码

# 引用错误
@error 修复这个错误
```

### 2. 项目级操作

```
# 在 Composer 中
创建一个完整的 REST API，包括：
- 用户模型
- CRUD 接口
- JWT 认证
- 输入验证
- 错误处理
```

### 3. 代码审查

```
# 选中代码后
Cmd+K → "审查这段代码，指出潜在问题"
```

### 4. 测试生成

```
# 选中函数后
Cmd+K → "为这个函数生成单元测试"
```

---

## 常见问题

### Q: API 连接失败？

1. 检查 API Key 是否正确
2. 检查 Base URL 是否正确
3. 检查网络连接

```bash
# 测试 API
curl https://api.bigaipro.com/v1/models \
  -H "Authorization: Bearer sk-xxxxxxxx"
```

### Q: 自动补全不工作？

1. 检查是否启用了自动补全
2. 检查模型配置
3. 重启 Cursor

### Q: 如何切换模型？

1. 打开设置
2. 搜索 "model"
3. 修改对应功能的模型

---

## 推荐配置

### 日常开发

```json
{
  "cursor.chat.model": "gpt-5.2-instant",
  "cursor.autocomplete.model": "gpt-4.1-mini",
  "cursor.composer.model": "claude-sonnet-4.5"
}
```

### 复杂项目

```json
{
  "cursor.chat.model": "claude-sonnet-4.5",
  "cursor.autocomplete.model": "gpt-4.1",
  "cursor.composer.model": "gpt-5.2-pro"
}
```

### 性价比优先

```json
{
  "cursor.chat.model": "deepseek-chat",
  "cursor.autocomplete.model": "gpt-4.1-nano",
  "cursor.composer.model": "qwen3-235b"
}
```

---

## 下一步

- 📖 [Continue 插件](./continue-plugin.md) - VS Code/JetBrains 插件
- 📖 [Kiro IDE 配置](./kiro-ide.md) - Kiro 智能开发环境
- 📖 [模型选择指南](../models/overview.md) - 选择合适的模型

---

<div align="center">

**Cursor - AI 时代的代码编辑器**

</div>
