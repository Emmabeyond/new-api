# Gemini CLI 使用教程

<div align="center">

**Google 官方命令行 AI 编程助手**

*支持 macOS / Windows / Linux*

</div>

---

## 简介

Gemini CLI 是 Google 推出的开源命令行 AI 编程工具，可以直接在终端中与 Gemini 模型对话，执行代码任务、文件操作和系统命令。

它能够理解你的整个代码库，帮助你编写代码、调试问题、执行 Shell 命令，而无需离开终端。

---

## 安装指南

### 前置要求

- Node.js 18 或更高版本
- Google 账户（用于认证）

### macOS

```bash
# 使用 NPM（推荐）
npm install -g @google/gemini-cli

# 或使用 npx 直接运行（无需安装）
npx @google/gemini-cli
```

### Windows

```powershell
# 使用 NPM
npm install -g @google/gemini-cli
```

### Linux

```bash
# 使用 NPM
npm install -g @google/gemini-cli

# Ubuntu/Debian 先安装 Node.js
sudo apt update && sudo apt install nodejs npm
npm install -g @google/gemini-cli
```

### 验证安装

```bash
gemini --version
```

---

## 认证配置

Gemini CLI 支持多种认证方式：

### 方式一：Google 账户登录（推荐）

```bash
# 启动 Gemini CLI
gemini

# 选择 "Login with Google"
# 浏览器会自动打开，完成登录授权
```

免费层级额度：
- 60 次请求/分钟
- 1000 次请求/天

### 方式二：API Key（Google AI Studio）

从 [Google AI Studio](https://aistudio.google.com/apikey) 获取 API Key：

#### macOS / Linux

```bash
# 编辑 ~/.bashrc 或 ~/.zshrc
export GEMINI_API_KEY="your-api-key-here"

# 使配置生效
source ~/.bashrc
```

#### Windows PowerShell

```powershell
# 设置环境变量
$env:GEMINI_API_KEY = "your-api-key-here"

# 永久设置
[Environment]::SetEnvironmentVariable("GEMINI_API_KEY", "your-api-key-here", "User")
```

### 方式三：使用 .env 文件

在项目目录或 `~/.gemini/` 创建 `.env` 文件：

```bash
# ~/.gemini/.env
GEMINI_API_KEY=your-api-key-here
```

---

## 配置 BigAI Pro API

> ⚠️ **注意**: Gemini CLI 目前官方仅支持 Google 官方 API 端点。自定义 API Base URL 功能正在开发中（[GitHub PR #2899](https://github.com/google-gemini/gemini-cli/pull/2899)）。

### 临时解决方案

如果你需要使用 BigAI Pro 的 Gemini 模型，可以通过以下方式：

#### 方案一：使用 BigAI Pro 的 OpenAI 兼容接口

BigAI Pro 提供 OpenAI 兼容的 Gemini 模型接口，可以配合 [Codex CLI](./codex-cli.md) 使用：

```bash
export OPENAI_API_KEY="sk-xxxxxxxxxxxxxxxx"
export OPENAI_API_BASE="https://api.bigaipro.com/v1"

# 使用 Codex CLI 调用 Gemini 模型
codex --model gemini-3.0-pro "你的问题"
```

#### 方案二：等待官方支持

Google 正在开发 `GEMINI_BASE_URL` 环境变量支持，届时可以这样配置：

```bash
# 未来版本支持（开发中）
export GEMINI_API_KEY="sk-xxxxxxxxxxxxxxxx"
export GEMINI_BASE_URL="https://api.bigaipro.com"
```

#### 方案三：直接使用 Google 官方 API

目前 Gemini CLI 可以直接使用 Google 官方 API，免费层级提供：
- 60 次请求/分钟
- 1000 次请求/天

---

## 基础使用

### 启动交互模式

```bash
# 在当前目录启动
gemini

# 在指定目录启动
gemini --cwd /path/to/project
```

### 单次命令

```bash
# 直接执行任务
gemini "创建一个 Python 脚本，读取 CSV 文件并生成报告"
```

### 常用内置命令

| 命令 | 说明 |
|------|------|
| `/help` | 显示帮助信息 |
| `/tools` | 查看可用工具 |
| `/settings` | 打开设置 |
| `/mcp` | 查看 MCP 服务器状态 |
| `/quit` | 退出 CLI |
| `!<command>` | 执行 Shell 命令 |

---

## 常用场景

### 代码生成

```bash
# 生成函数
> 写一个 Go 函数，实现 JWT 令牌验证

# 生成 API
> 用 Gin 框架创建一个用户注册接口
```

### 代码解释

```bash
# 解释代码
> 解释 main.go 中的 handleRequest 函数

# 解释错误
> 这个错误是什么意思？如何修复？
```

### 代码重构

```bash
# 重构函数
> 重构 utils.go 中的 formatDate 函数，提高可读性

# 优化性能
> 优化这个数据库查询的性能
```

### 文件操作

```bash
# 创建文件
> 创建一个 .gitignore 文件，包含 Go 项目常用的忽略规则

# 修改文件
> 在 go.mod 中添加 gin-gonic/gin 依赖
```

### Git 操作

```bash
# 提交更改
> 提交当前更改，生成有意义的 commit message

# 查看历史
> 显示最近一周的提交历史
```

---

## 配置文件

### 配置文件位置

| 类型 | 位置 |
|------|------|
| 用户配置 | `~/.gemini/settings.json` |
| 项目配置 | `<project>/.gemini/settings.json` |

### 配置示例

`~/.gemini/settings.json`：

```json
{
  "general": {
    "vimMode": false,
    "preferredEditor": "code"
  },
  "ui": {
    "theme": "GitHub",
    "hideBanner": false,
    "hideTips": false
  },
  "model": {
    "name": "gemini-3-flash-preview"
  },
  "context": {
    "fileName": ["GEMINI.md", "CONTEXT.md"]
  },
  "privacy": {
    "usageStatisticsEnabled": false
  }
}
```

### 主要配置项

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `general.vimMode` | 启用 Vim 模式 | false |
| `ui.theme` | 界面主题 | Default |
| `model.name` | 默认模型 | gemini-2.5-flash |
| `context.fileName` | 上下文文件名 | ["GEMINI.md"] |

---

## 项目上下文配置

在项目根目录创建 `GEMINI.md` 文件，为 AI 提供项目特定的指令和上下文：

```markdown
# 项目说明

这是一个 Go 语言的 API 网关项目。

## 代码规范

- 使用 Go 1.21+
- 遵循 Go 官方代码风格
- 错误处理使用 errors.Wrap
- 日志使用 zap 库

## 项目结构

- `/controller` - HTTP 处理器
- `/service` - 业务逻辑
- `/model` - 数据模型
- `/common` - 公共工具

## 注意事项

- 所有 API 需要添加认证中间件
- 数据库操作使用 GORM
```

---

## MCP 服务器集成

Gemini CLI 支持 Model Context Protocol (MCP)，可以扩展 AI 的能力。

### 配置 MCP 服务器

在 `.gemini/settings.json` 中添加：

```json
{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"],
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "$GITHUB_TOKEN"
      }
    },
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/allowed/dir"]
    }
  }
}
```

### 查看 MCP 工具

```bash
# 在 Gemini CLI 中
/mcp
```

---

## 沙箱模式

Gemini CLI 支持在 Docker 沙箱中执行命令，保护系统安全：

```bash
# 启用 Docker 沙箱
gemini --sandbox docker
```

或在配置文件中设置：

```json
{
  "tools": {
    "sandbox": "docker"
  }
}
```

---

## 常见问题

### Q: 如何使用代理？

```bash
# 设置代理环境变量
export HTTPS_PROXY="http://127.0.0.1:7890"
export HTTP_PROXY="http://127.0.0.1:7890"
```

### Q: 如何切换模型？

```bash
# 在配置文件中设置
{
  "model": {
    "name": "gemini-3-pro-preview"
  }
}
```

可用模型（2025 年 12 月更新）：

**Gemini 3 系列（最新）**
- `gemini-3-pro-preview` - 最强推理能力，100 万上下文
- `gemini-3-flash-preview` - Pro 级智能，Flash 级速度
- `gemini-3-pro-image-preview` - 最高质量图像生成

**Gemini 2.5 系列**
- `gemini-2.5-pro` - 稳定版旗舰
- `gemini-2.5-flash` - 高性价比（默认）
- `gemini-2.5-flash-lite` - 轻量快速

### Q: 如何查看历史记录？

```bash
# 历史记录存储在
~/.gemini/history/<project-hash>/
```

### Q: 认证失败怎么办？

```bash
# 清除缓存的认证信息
rm -rf ~/.gemini/auth/

# 重新登录
gemini
```

---

## 最佳实践

### 1. 创建项目上下文文件

在每个项目中创建 `GEMINI.md`，提供项目背景和代码规范。

### 2. 明确任务描述

```bash
# ✅ 好的描述
> 在 controller/user.go 中添加一个 UpdateProfile 函数，接收 JSON 请求体，更新用户的 nickname 和 avatar 字段

# ❌ 不好的描述
> 写个更新用户的函数
```

### 3. 分步执行复杂任务

```bash
> 第一步：分析当前项目结构
> 第二步：创建数据库模型
> 第三步：实现 CRUD 接口
```

### 4. 使用沙箱保护系统

对于不熟悉的项目或危险操作，启用 Docker 沙箱模式。

---

## 下一步

- 📖 [Claude Code 教程](./claude-code.md) - Anthropic 的编程助手
- 📖 [Codex CLI 教程](./codex-cli.md) - OpenAI 的命令行工具
- 📖 [Gemini 模型详解](../models/gemini-models.md) - 了解 Gemini 模型特性

---

<div align="center">

**Gemini CLI - 终端中的 AI 编程伙伴**

</div>
