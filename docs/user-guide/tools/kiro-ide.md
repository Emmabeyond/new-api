# Kiro IDE 配置指南

<div align="center">

**AWS 出品的智能开发环境**

</div>

---

## 简介

Kiro 是 AWS 推出的 AI 原生 IDE，基于 VS Code 构建，提供 Spec 驱动开发、Agent Hooks、Steering 等独特功能。通过配置 BigAI Pro API，可以使用各种顶尖模型。

---

## 安装 Kiro

### 下载安装

1. 访问 [Kiro 官网](https://kiro.dev)
2. 下载对应系统的安装包
3. 安装并启动 Kiro

### 系统要求

- macOS 10.15+
- Windows 10+
- Linux (Ubuntu 18.04+)

---

## 配置 BigAI Pro API

### 方式一：MCP 配置（推荐）

Kiro 支持通过 MCP (Model Context Protocol) 配置自定义 API。

#### 工作区配置

创建 `.kiro/settings/mcp.json`：

```json
{
  "mcpServers": {
    "bigai-pro": {
      "command": "npx",
      "args": ["-y", "@anthropic-ai/claude-code-mcp"],
      "env": {
        "ANTHROPIC_API_KEY": "sk-xxxxxxxxxxxxxxxx",
        "ANTHROPIC_BASE_URL": "https://api.bigaipro.com"
      }
    }
  }
}
```

#### 用户级配置

编辑 `~/.kiro/settings/mcp.json`：

```json
{
  "mcpServers": {
    "openai-compatible": {
      "command": "node",
      "args": ["path/to/openai-mcp-server.js"],
      "env": {
        "OPENAI_API_KEY": "sk-xxxxxxxxxxxxxxxx",
        "OPENAI_API_BASE": "https://api.bigaipro.com/v1"
      }
    }
  }
}
```

### 方式二：环境变量

```bash
# macOS / Linux
export OPENAI_API_KEY="sk-xxxxxxxxxxxxxxxx"
export OPENAI_API_BASE="https://api.bigaipro.com/v1"

# Windows PowerShell
$env:OPENAI_API_KEY = "sk-xxxxxxxxxxxxxxxx"
$env:OPENAI_API_BASE = "https://api.bigaipro.com/v1"
```

---

## Kiro 核心功能

### 1. Spec 驱动开发

Kiro 的 Spec 功能可以将想法转化为结构化的需求、设计和任务：

```
用户想法 → 需求文档 → 设计文档 → 任务列表 → 代码实现
```

#### 创建 Spec

1. 打开命令面板 `Cmd+Shift+P`
2. 输入 "Kiro: Create Spec"
3. 描述你的功能想法
4. Kiro 会生成：
   - `requirements.md` - 需求文档
   - `design.md` - 设计文档
   - `tasks.md` - 任务列表

#### Spec 文件结构

```
.kiro/specs/
└── feature-name/
    ├── requirements.md
    ├── design.md
    └── tasks.md
```

### 2. Agent Hooks

Agent Hooks 允许在特定事件触发时自动执行 AI 任务：

#### 创建 Hook

1. 打开 Explorer 视图
2. 找到 "Agent Hooks" 部分
3. 点击创建新 Hook

#### Hook 配置示例

```json
{
  "name": "auto-test",
  "trigger": "onFileSave",
  "pattern": "**/*.ts",
  "action": {
    "type": "message",
    "content": "为修改的文件更新测试"
  }
}
```

#### 常用 Hook 场景

| 触发器 | 场景 |
|--------|------|
| onFileSave | 保存时自动运行测试 |
| onSessionCreate | 新会话时加载项目上下文 |
| onAgentComplete | Agent 完成后执行后续任务 |
| manual | 手动触发的任务 |

### 3. Steering（引导）

Steering 文件为 AI 提供项目特定的上下文和指令：

#### 创建 Steering 文件

在 `.kiro/steering/` 目录下创建 `.md` 文件：

```markdown
---
inclusion: always
---

# 项目规范

## 代码风格
- 使用 TypeScript 严格模式
- 遵循 ESLint 规则
- 使用 Prettier 格式化

## 架构原则
- 使用依赖注入
- 遵循 SOLID 原则
- 分层架构：Controller → Service → Repository

## 测试要求
- 单元测试覆盖率 > 80%
- 使用 Jest 测试框架
```

#### Steering 包含模式

| 模式 | 说明 |
|------|------|
| `always` | 始终包含 |
| `fileMatch` | 匹配特定文件时包含 |
| `manual` | 手动引用时包含 |

#### 条件包含示例

```markdown
---
inclusion: fileMatch
fileMatchPattern: "**/*.test.ts"
---

# 测试规范

- 使用 describe/it 结构
- 每个测试只测试一个功能
- 使用 mock 隔离依赖
```

---

## 使用技巧

### 1. 上下文引用

在对话中使用 `#` 引用上下文：

```
#File:src/main.ts 解释这个文件
#Folder:src/components 分析这个目录
#Problems 查看当前问题
#Terminal 查看终端输出
#Git Diff 查看代码变更
#Codebase 搜索整个代码库
```

### 2. 图像支持

拖拽图像到对话框，或点击图像图标上传：

- 支持截图分析
- 支持设计稿理解
- 支持错误截图诊断

### 3. 自主模式

Kiro 支持两种工作模式：

| 模式 | 说明 |
|------|------|
| Autopilot | 自动执行文件修改 |
| Supervised | 修改后需要确认 |

---

## 快捷键

| 功能 | macOS | Windows/Linux |
|------|-------|---------------|
| 打开对话 | `Cmd+L` | `Ctrl+L` |
| 命令面板 | `Cmd+Shift+P` | `Ctrl+Shift+P` |
| 创建 Spec | `Cmd+Shift+S` | `Ctrl+Shift+S` |
| 执行任务 | `Cmd+Enter` | `Ctrl+Enter` |

---

## 推荐工作流

### 1. 新功能开发

```
1. 创建 Spec，描述功能想法
2. 审查生成的需求文档
3. 审查设计文档
4. 逐个执行任务
5. 代码审查和测试
```

### 2. Bug 修复

```
1. 使用 #Problems 查看错误
2. 让 AI 分析错误原因
3. 生成修复方案
4. 应用修复
5. 验证修复
```

### 3. 代码重构

```
1. 创建 Steering 文件定义重构规范
2. 创建 Spec 描述重构目标
3. 逐步执行重构任务
4. 运行测试验证
```

---

## 常见问题

### Q: 如何切换模型？

在对话中使用 `@model` 或在设置中配置默认模型。

### Q: Spec 文件在哪里？

所有 Spec 文件存储在 `.kiro/specs/` 目录下。

### Q: 如何查看 Agent 日志？

打开输出面板，选择 "Kiro Agent" 通道。

---

## 下一步

- 📖 [Cherry Studio](./cherry-studio.md) - 桌面客户端
- 📖 [Claude Code 教程](./claude-code.md) - 终端编程助手
- 📖 [模型选择指南](../models/overview.md) - 选择合适的模型

---

<div align="center">

**Kiro - Spec 驱动的智能开发**

</div>
