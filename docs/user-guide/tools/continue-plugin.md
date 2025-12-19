# Continue 插件配置指南

<div align="center">

**开源 AI 编程助手插件**

*支持 VS Code / JetBrains IDE*

</div>

---

## 简介

Continue 是一款开源的 AI 编程助手插件，支持 VS Code 和 JetBrains 系列 IDE。它可以连接各种 AI 模型，提供代码补全、对话、重构等功能。

---

## 安装

### VS Code

1. 打开 VS Code
2. 按 `Cmd+Shift+X` (macOS) 或 `Ctrl+Shift+X` (Windows/Linux)
3. 搜索 "Continue"
4. 点击安装

或使用命令行：

```bash
code --install-extension Continue.continue
```

### JetBrains IDE

1. 打开 IDE（IntelliJ IDEA、PyCharm、WebStorm 等）
2. 进入 `Settings/Preferences` → `Plugins`
3. 搜索 "Continue"
4. 点击安装并重启

---

## 配置 BigAI Pro API

### 步骤一：打开配置

#### VS Code

按 `Cmd+Shift+P` → 输入 "Continue: Open Config"

#### JetBrains

`Settings` → `Tools` → `Continue`

### 步骤二：编辑配置文件

配置文件位置：`~/.continue/config.json`

```json
{
  "models": [
    {
      "title": "GPT-5.2 Instant",
      "provider": "openai",
      "model": "gpt-5.2-instant",
      "apiKey": "sk-xxxxxxxxxxxxxxxx",
      "apiBase": "https://api.bigaipro.com/v1"
    },
    {
      "title": "Claude Sonnet 4.5",
      "provider": "openai",
      "model": "claude-sonnet-4.5",
      "apiKey": "sk-xxxxxxxxxxxxxxxx",
      "apiBase": "https://api.bigaipro.com/v1"
    },
    {
      "title": "DeepSeek R1",
      "provider": "openai",
      "model": "deepseek-r1",
      "apiKey": "sk-xxxxxxxxxxxxxxxx",
      "apiBase": "https://api.bigaipro.com/v1"
    }
  ],
  "tabAutocompleteModel": {
    "title": "GPT-4.1 Mini",
    "provider": "openai",
    "model": "gpt-4.1-mini",
    "apiKey": "sk-xxxxxxxxxxxxxxxx",
    "apiBase": "https://api.bigaipro.com/v1"
  }
}
```

---

## 完整配置示例

```json
{
  "models": [
    {
      "title": "GPT-5.2 Pro",
      "provider": "openai",
      "model": "gpt-5.2-pro",
      "apiKey": "sk-xxxxxxxxxxxxxxxx",
      "apiBase": "https://api.bigaipro.com/v1",
      "contextLength": 1000000
    },
    {
      "title": "GPT-5.2 Instant",
      "provider": "openai",
      "model": "gpt-5.2-instant",
      "apiKey": "sk-xxxxxxxxxxxxxxxx",
      "apiBase": "https://api.bigaipro.com/v1",
      "contextLength": 1000000
    },
    {
      "title": "Claude Sonnet 4.5",
      "provider": "openai",
      "model": "claude-sonnet-4.5",
      "apiKey": "sk-xxxxxxxxxxxxxxxx",
      "apiBase": "https://api.bigaipro.com/v1",
      "contextLength": 200000
    },
    {
      "title": "o3",
      "provider": "openai",
      "model": "o3",
      "apiKey": "sk-xxxxxxxxxxxxxxxx",
      "apiBase": "https://api.bigaipro.com/v1",
      "contextLength": 200000
    },
    {
      "title": "DeepSeek R1",
      "provider": "openai",
      "model": "deepseek-r1",
      "apiKey": "sk-xxxxxxxxxxxxxxxx",
      "apiBase": "https://api.bigaipro.com/v1",
      "contextLength": 128000
    },
    {
      "title": "Qwen3 235B",
      "provider": "openai",
      "model": "qwen3-235b",
      "apiKey": "sk-xxxxxxxxxxxxxxxx",
      "apiBase": "https://api.bigaipro.com/v1",
      "contextLength": 131000
    }
  ],
  "tabAutocompleteModel": {
    "title": "GPT-4.1 Mini",
    "provider": "openai",
    "model": "gpt-4.1-mini",
    "apiKey": "sk-xxxxxxxxxxxxxxxx",
    "apiBase": "https://api.bigaipro.com/v1"
  },
  "embeddingsProvider": {
    "provider": "openai",
    "model": "text-embedding-3-large",
    "apiKey": "sk-xxxxxxxxxxxxxxxx",
    "apiBase": "https://api.bigaipro.com/v1"
  },
  "customCommands": [
    {
      "name": "test",
      "prompt": "为选中的代码编写单元测试",
      "description": "生成单元测试"
    },
    {
      "name": "review",
      "prompt": "审查选中的代码，指出潜在问题和改进建议",
      "description": "代码审查"
    },
    {
      "name": "doc",
      "prompt": "为选中的代码生成详细的文档注释",
      "description": "生成文档"
    }
  ],
  "contextProviders": [
    {
      "name": "code",
      "params": {}
    },
    {
      "name": "docs",
      "params": {}
    },
    {
      "name": "diff",
      "params": {}
    },
    {
      "name": "terminal",
      "params": {}
    },
    {
      "name": "problems",
      "params": {}
    }
  ],
  "slashCommands": [
    {
      "name": "edit",
      "description": "编辑选中的代码"
    },
    {
      "name": "comment",
      "description": "添加注释"
    },
    {
      "name": "share",
      "description": "分享对话"
    }
  ]
}
```

---

## 核心功能

### 1. AI 对话

按 `Cmd+L` (macOS) 或 `Ctrl+L` (Windows/Linux) 打开对话面板：

```
你：解释这段代码
AI：这段代码实现了...

你：@file:src/main.ts 分析这个文件
AI：这个文件的主要功能是...
```

### 2. 代码补全

Continue 提供 AI 驱动的代码补全：

- 自动触发补全建议
- 按 `Tab` 接受建议
- 按 `Esc` 拒绝建议

### 3. 行内编辑

选中代码后按 `Cmd+I` (macOS) 或 `Ctrl+I` (Windows/Linux)：

```
选中代码 → Cmd+I → 输入指令 → AI 修改
```

### 4. 斜杠命令

在对话中使用斜杠命令：

```
/edit 重构这个函数使用 async/await
/comment 添加详细注释
/test 生成单元测试
```

### 5. 上下文引用

使用 `@` 引用上下文：

```
@file:src/utils.ts 解释这个文件
@folder:src/components 分析这个目录的组件
@code 解释选中的代码
@docs 搜索文档
@terminal 查看终端输出
@problems 查看当前问题
```

---

## 快捷键

### VS Code

| 功能 | macOS | Windows/Linux |
|------|-------|---------------|
| 打开对话 | `Cmd+L` | `Ctrl+L` |
| 行内编辑 | `Cmd+I` | `Ctrl+I` |
| 接受补全 | `Tab` | `Tab` |
| 拒绝补全 | `Esc` | `Esc` |
| 添加到对话 | `Cmd+Shift+L` | `Ctrl+Shift+L` |

### JetBrains

| 功能 | macOS | Windows/Linux |
|------|-------|---------------|
| 打开对话 | `Cmd+J` | `Ctrl+J` |
| 行内编辑 | `Cmd+I` | `Ctrl+I` |
| 接受补全 | `Tab` | `Tab` |

---

## 自定义命令

在配置文件中添加自定义命令：

```json
{
  "customCommands": [
    {
      "name": "optimize",
      "prompt": "优化选中代码的性能，保持功能不变",
      "description": "性能优化"
    },
    {
      "name": "security",
      "prompt": "检查选中代码的安全漏洞",
      "description": "安全检查"
    },
    {
      "name": "refactor",
      "prompt": "重构选中代码，提高可读性和可维护性",
      "description": "代码重构"
    },
    {
      "name": "explain",
      "prompt": "用简单的语言解释选中代码的作用",
      "description": "代码解释"
    }
  ]
}
```

使用方式：

```
/optimize
/security
/refactor
/explain
```

---

## 模型切换

### 在对话中切换

点击对话面板顶部的模型名称，选择其他模型。

### 快捷切换

```
@model:claude-sonnet-4.5 使用 Claude 回答这个问题
@model:o3 分析这个算法的复杂度
```

---

## 常见问题

### Q: 连接失败？

1. 检查 API Key 和 Base URL
2. 检查网络连接
3. 查看 Continue 输出日志

```bash
# 测试 API
curl https://api.bigaipro.com/v1/models \
  -H "Authorization: Bearer sk-xxxxxxxx"
```

### Q: 补全不工作？

1. 检查 `tabAutocompleteModel` 配置
2. 确保模型支持补全
3. 重启 IDE

### Q: 如何查看日志？

#### VS Code

`View` → `Output` → 选择 "Continue"

#### JetBrains

`Help` → `Show Log in Finder/Explorer`

---

## 推荐配置

### 日常开发

```json
{
  "models": [
    {
      "title": "GPT-5.2 Instant",
      "model": "gpt-5.2-instant"
    }
  ],
  "tabAutocompleteModel": {
    "model": "gpt-4.1-mini"
  }
}
```

### 代码密集型

```json
{
  "models": [
    {
      "title": "Claude Sonnet 4.5",
      "model": "claude-sonnet-4.5"
    }
  ],
  "tabAutocompleteModel": {
    "model": "gpt-4.1"
  }
}
```

### 性价比优先

```json
{
  "models": [
    {
      "title": "DeepSeek Chat",
      "model": "deepseek-chat"
    }
  ],
  "tabAutocompleteModel": {
    "model": "gpt-4.1-nano"
  }
}
```

---

## 下一步

- 📖 [Kiro IDE 配置](./kiro-ide.md) - Kiro 智能开发环境
- 📖 [Cherry Studio](./cherry-studio.md) - 桌面客户端
- 📖 [模型选择指南](../models/overview.md) - 选择合适的模型

---

<div align="center">

**Continue - 开源 AI 编程助手**

</div>
