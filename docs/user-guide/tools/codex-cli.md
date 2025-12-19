# Codex CLI 使用教程

<div align="center">

**OpenAI 官方命令行编程助手**

*支持 macOS / Windows / Linux*

</div>

---

## 简介

Codex CLI 是 OpenAI 推出的命令行 AI 编程工具，可以直接在终端中与 AI 对话，执行代码任务、文件操作和系统命令。

通过配置 BigAI Pro API，你可以使用 Codex CLI 的全部功能。

---

## 安装指南

### macOS

```bash
# 使用 Homebrew
brew install openai-codex

# 或使用 NPM
npm install -g @openai/codex
```

### Windows

```powershell
# 使用 NPM
npm install -g @openai/codex

# 或使用 Winget
winget install OpenAI.Codex
```

### Linux

```bash
# 使用 NPM（推荐）
npm install -g @openai/codex

# Ubuntu/Debian
curl -fsSL https://openai.com/codex/install.sh | bash
```

### 验证安装

```bash
codex --version
```

---

## 配置 BigAI Pro API

### 方式一：环境变量（推荐）

#### macOS / Linux

```bash
# 编辑 ~/.bashrc 或 ~/.zshrc
export OPENAI_API_KEY="sk-xxxxxxxxxxxxxxxx"
export OPENAI_API_BASE="https://api.bigaipro.com/v1"

# 使配置生效
source ~/.bashrc
```

#### Windows PowerShell

```powershell
# 设置环境变量
$env:OPENAI_API_KEY = "sk-xxxxxxxxxxxxxxxx"
$env:OPENAI_API_BASE = "https://api.bigaipro.com/v1"

# 永久设置
[Environment]::SetEnvironmentVariable("OPENAI_API_KEY", "sk-xxxxxxxxxxxxxxxx", "User")
[Environment]::SetEnvironmentVariable("OPENAI_API_BASE", "https://api.bigaipro.com/v1", "User")
```

#### Windows CMD

```cmd
setx OPENAI_API_KEY "sk-xxxxxxxxxxxxxxxx"
setx OPENAI_API_BASE "https://api.bigaipro.com/v1"
```

### 方式二：配置文件

创建 `~/.codex/config.json`：

```json
{
  "apiKey": "sk-xxxxxxxxxxxxxxxx",
  "baseUrl": "https://api.bigaipro.com/v1",
  "model": "gpt-5.2-instant",
  "approvalMode": "suggest"
}
```

### 方式三：命令行参数

```bash
codex --api-key "sk-xxxxxxxx" --api-base "https://api.bigaipro.com/v1"
```

---

## 验证配置

```bash
# 检查配置
codex config show

# 测试连接
codex "说你好"
```

---

## 基础使用

### 交互模式

```bash
# 启动交互模式
codex

# 在指定目录启动
codex --cwd /path/to/project
```

### 单次命令

```bash
# 直接执行任务
codex "创建一个 Python 脚本，读取 CSV 文件并生成报告"

# 指定模型
codex --model gpt-5.2-pro "设计一个数据库架构"
```

### 运行模式

Codex CLI 有三种运行模式：

| 模式 | 说明 | 命令 |
|------|------|------|
| suggest | 建议模式，需要确认 | `codex --approval suggest` |
| auto | 自动执行安全操作 | `codex --approval auto` |
| full-auto | 完全自动执行 | `codex --approval full-auto` |

```bash
# 建议模式（默认，最安全）
codex --approval suggest "重构这个函数"

# 自动模式（自动执行读取操作）
codex --approval auto "分析项目结构"

# 完全自动模式（谨慎使用）
codex --approval full-auto "修复所有 lint 错误"
```

---

## 常用场景

### 代码生成

```bash
# 生成函数
codex "写一个 Python 函数，计算斐波那契数列"

# 生成类
codex "创建一个 TypeScript 类，实现 LRU 缓存"

# 生成 API
codex "用 FastAPI 创建一个用户注册接口"
```

### 代码解释

```bash
# 解释代码
codex "解释 src/main.py 中的 process_data 函数"

# 解释错误
codex "这个错误是什么意思？如何修复？"
```

### 代码重构

```bash
# 重构函数
codex "重构 utils.js 中的 formatDate 函数，使用现代 JavaScript"

# 优化性能
codex "优化这个数据库查询的性能"
```

### 文件操作

```bash
# 创建文件
codex "创建一个 .gitignore 文件，包含 Node.js 项目常用的忽略规则"

# 修改文件
codex "在 package.json 中添加 jest 作为开发依赖"

# 批量操作
codex "将所有 .js 文件重命名为 .ts"
```

### 系统命令

```bash
# 查找文件
codex "找到所有超过 1MB 的日志文件"

# 系统信息
codex "显示当前系统的内存使用情况"

# 进程管理
codex "找到占用 8080 端口的进程"
```

### Git 操作

```bash
# 提交更改
codex "提交当前更改，生成有意义的 commit message"

# 查看历史
codex "显示最近一周的提交历史"

# 分支管理
codex "创建一个新分支用于开发用户认证功能"
```

---

## 高级功能

### 上下文管理

```bash
# 添加文件到上下文
codex --context src/main.py "优化这个文件"

# 添加多个文件
codex --context "src/*.py" "分析这些文件的代码质量"

# 排除文件
codex --exclude "node_modules,dist" "搜索所有 TODO 注释"
```

### 输出控制

```bash
# 输出到文件
codex "生成 API 文档" --output docs/api.md

# JSON 格式输出
codex "列出所有函数" --format json

# 静默模式
codex "运行测试" --quiet
```

### 模型选择

```bash
# 使用 GPT-5.2 Pro（复杂任务）
codex --model gpt-5.2-pro "设计微服务架构"

# 使用 GPT-5.2 Instant（快速响应）
codex --model gpt-5.2-instant "格式化这段代码"

# 使用 o3（推理任务）
codex --model o3 "分析这个算法的时间复杂度"
```

---

## 配置选项

### 完整配置示例

`~/.codex/config.json`：

```json
{
  "apiKey": "sk-xxxxxxxxxxxxxxxx",
  "baseUrl": "https://api.bigaipro.com/v1",
  "model": "gpt-5.2-instant",
  "approvalMode": "suggest",
  "maxTokens": 4096,
  "temperature": 0.7,
  "shell": "/bin/zsh",
  "editor": "vim",
  "theme": "dark",
  "history": {
    "enabled": true,
    "maxSize": 1000
  },
  "safety": {
    "allowFileWrite": true,
    "allowFileDelete": false,
    "allowNetworkAccess": true,
    "allowSystemCommands": true,
    "blockedCommands": ["rm -rf /", "format"]
  },
  "context": {
    "autoInclude": ["README.md", "package.json"],
    "exclude": ["node_modules", ".git", "dist"]
  }
}
```

### 配置说明

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| apiKey | API 密钥 | 必填 |
| baseUrl | API 地址 | https://api.openai.com/v1 |
| model | 默认模型 | gpt-5.2-instant |
| approvalMode | 审批模式 | suggest |
| maxTokens | 最大输出 | 4096 |
| safety.allowFileDelete | 允许删除文件 | false |

---

## 安全设置

### 审批模式

```bash
# 建议模式：所有操作需要确认
codex --approval suggest

# 自动模式：只有写操作需要确认
codex --approval auto

# 完全自动：所有操作自动执行（危险）
codex --approval full-auto
```

### 命令黑名单

在配置文件中设置禁止执行的命令：

```json
{
  "safety": {
    "blockedCommands": [
      "rm -rf",
      "format",
      "dd if=",
      "mkfs",
      "> /dev/sda"
    ]
  }
}
```

---

## 常见问题

### Q: 连接失败？

```bash
# 检查网络
curl https://api.bigaipro.com/v1/models -H "Authorization: Bearer sk-xxx"

# 检查配置
codex config show

# 启用调试
codex --debug "测试"
```

### Q: 如何使用代理？

```bash
# 设置代理环境变量
export HTTPS_PROXY="http://127.0.0.1:7890"
export HTTP_PROXY="http://127.0.0.1:7890"
```

### Q: 如何清除历史？

```bash
codex history clear
```

---

## 最佳实践

### 1. 明确任务描述

```bash
# ✅ 好的描述
codex "在 src/utils/date.ts 中创建一个函数，将 ISO 日期字符串转换为中文格式（如：2025年12月19日）"

# ❌ 不好的描述
codex "写个日期函数"
```

### 2. 使用上下文

```bash
# 提供相关文件作为上下文
codex --context src/types.ts "基于这些类型定义创建验证函数"
```

### 3. 分步执行

```bash
# 复杂任务分步执行
codex "第一步：分析当前项目结构"
codex "第二步：创建数据库模型"
codex "第三步：实现 CRUD 接口"
```

---

## 下一步

- 📖 [Claude Code 教程](./claude-code.md) - Anthropic 的编程助手
- 📖 [Cursor IDE 配置](./cursor-ide.md) - AI 增强的代码编辑器
- 📖 [GPT 模型详解](../models/gpt-models.md) - 了解 GPT 模型特性

---

<div align="center">

**Codex CLI - 命令行中的 AI 助手**

</div>
