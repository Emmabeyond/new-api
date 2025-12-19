# VS Code 中使用 CLI 工具

<div align="center">

**在 VS Code / Cursor / Kiro 中集成命令行 AI 工具**

</div>

---

## 简介

VS Code 及其衍生编辑器（Cursor、Kiro、Windsurf 等）都支持集成终端。你可以在编辑器内直接使用 Claude Code、Codex CLI、Gemini CLI 等命令行 AI 工具，实现无缝的 AI 编程体验。

---

## 环境配置

### 配置 BigAI Pro API

在使用 CLI 工具前，先配置环境变量。

#### macOS / Linux

编辑 `~/.zshrc` 或 `~/.bashrc`：

```bash
# BigAI Pro API 配置
export OPENAI_API_KEY="sk-xxxxxxxxxxxxxxxx"
export OPENAI_API_BASE="https://api.bigaipro.com/v1"

# Anthropic API 配置（用于 Claude Code）
export ANTHROPIC_API_KEY="sk-xxxxxxxxxxxxxxxx"
export ANTHROPIC_BASE_URL="https://api.bigaipro.com"
```

使配置生效：

```bash
source ~/.zshrc
```

#### Windows PowerShell

```powershell
# 永久设置环境变量
[Environment]::SetEnvironmentVariable("OPENAI_API_KEY", "sk-xxxxxxxxxxxxxxxx", "User")
[Environment]::SetEnvironmentVariable("OPENAI_API_BASE", "https://api.bigaipro.com/v1", "User")
[Environment]::SetEnvironmentVariable("ANTHROPIC_API_KEY", "sk-xxxxxxxxxxxxxxxx", "User")
[Environment]::SetEnvironmentVariable("ANTHROPIC_BASE_URL", "https://api.bigaipro.com", "User")
```

重启终端使配置生效。

---

## 在 VS Code 中使用

### 打开集成终端

| 操作 | macOS | Windows/Linux |
|------|-------|---------------|
| 打开终端 | `` Ctrl+` `` | `` Ctrl+` `` |
| 新建终端 | `Cmd+Shift+`` | `Ctrl+Shift+`` |
| 切换终端 | `Cmd+Shift+[/]` | `Ctrl+Shift+[/]` |

### 使用 Claude Code

```bash
# 安装
npm install -g @anthropic-ai/claude-code

# 在项目目录启动
claude

# 常用命令
> 解释 src/main.ts 的架构
> 为 UserService 添加单元测试
> 重构这个函数，提高可读性
```

### 使用 Codex CLI

```bash
# 安装
npm install -g @openai/codex

# 在项目目录启动
codex

# 常用命令
> 创建一个 REST API 接口
> 修复这个 bug
> 生成 API 文档
```

### 使用 Gemini CLI

```bash
# 安装
npm install -g @google/gemini-cli

# 在项目目录启动
gemini

# 常用命令
> 分析项目结构
> 优化数据库查询
> 生成测试用例
```

---

## 在 Cursor 中使用

Cursor 基于 VS Code，完全兼容上述配置。

### 推荐工作流

1. 使用 Cursor 内置 AI 进行快速编辑（`Cmd+K`）
2. 使用 CLI 工具进行复杂任务
3. 结合使用获得最佳体验

### 终端配置

Cursor 设置 → Terminal → 配置默认 Shell：

```json
{
  "terminal.integrated.defaultProfile.osx": "zsh",
  "terminal.integrated.env.osx": {
    "OPENAI_API_KEY": "sk-xxxxxxxxxxxxxxxx",
    "OPENAI_API_BASE": "https://api.bigaipro.com/v1"
  }
}
```

---

## 在 Kiro 中使用

Kiro 同样基于 VS Code，支持集成终端。

### 结合 Spec 使用

1. 使用 Kiro Spec 生成任务列表
2. 在终端中使用 CLI 工具执行具体任务
3. 使用 Kiro 的 Agent Hooks 自动化流程

### 终端配置

在 `.kiro/settings/` 中创建环境配置：

```json
{
  "terminal.integrated.env.osx": {
    "ANTHROPIC_API_KEY": "sk-xxxxxxxxxxxxxxxx",
    "ANTHROPIC_BASE_URL": "https://api.bigaipro.com"
  }
}
```

---

## 工作流示例

### 示例一：代码审查

```bash
# 1. 打开终端
Ctrl+`

# 2. 启动 Claude Code
claude

# 3. 审查代码
> 审查 src/services/ 目录下的所有文件，指出潜在问题

# 4. 根据建议修改代码
> 修复你发现的问题
```

### 示例二：功能开发

```bash
# 1. 使用 Codex CLI 生成代码框架
codex "创建一个用户认证模块的基础结构"

# 2. 使用编辑器 AI 完善细节
# Cmd+K → "添加输入验证"

# 3. 使用 CLI 生成测试
codex "为 UserAuth 模块生成单元测试"
```

### 示例三：项目分析

```bash
# 1. 使用 Gemini CLI 分析项目
gemini

> 分析这个项目的架构，生成架构图

# 2. 获取优化建议
> 这个项目有哪些可以优化的地方？

# 3. 执行优化
> 优化 database.go 中的连接池配置
```

---

## 快捷键配置

### VS Code / Cursor

在 `keybindings.json` 中添加：

```json
[
  {
    "key": "cmd+shift+c",
    "command": "workbench.action.terminal.sendSequence",
    "args": { "text": "claude\n" }
  },
  {
    "key": "cmd+shift+x",
    "command": "workbench.action.terminal.sendSequence",
    "args": { "text": "codex\n" }
  },
  {
    "key": "cmd+shift+g",
    "command": "workbench.action.terminal.sendSequence",
    "args": { "text": "gemini\n" }
  }
]
```

---

## 终端分屏

### 并行使用多个 CLI

1. 打开终端
2. 点击分屏按钮或 `Cmd+\`
3. 在不同面板运行不同 CLI

```
┌─────────────────┬─────────────────┐
│   Claude Code   │   Codex CLI     │
│                 │                 │
│  > 分析代码     │  > 生成测试     │
│                 │                 │
└─────────────────┴─────────────────┘
```

---

## 任务自动化

### VS Code Tasks

创建 `.vscode/tasks.json`：

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "AI: 代码审查",
      "type": "shell",
      "command": "claude",
      "args": ["审查当前项目的代码质量"],
      "problemMatcher": []
    },
    {
      "label": "AI: 生成测试",
      "type": "shell",
      "command": "codex",
      "args": ["为 ${file} 生成单元测试"],
      "problemMatcher": []
    },
    {
      "label": "AI: 生成文档",
      "type": "shell",
      "command": "codex",
      "args": ["为 ${file} 生成 API 文档"],
      "problemMatcher": []
    }
  ]
}
```

运行任务：`Cmd+Shift+P` → "Tasks: Run Task"

---

## 常见问题

### Q: 终端中环境变量不生效？

1. 确保已重启终端
2. 检查 Shell 配置文件是否正确
3. 在 VS Code 设置中配置终端环境变量

### Q: CLI 工具安装失败？

```bash
# 检查 Node.js 版本
node -v  # 需要 18+

# 使用 npx 直接运行（无需安装）
npx @anthropic-ai/claude-code
npx @openai/codex
npx @google/gemini-cli
```

### Q: 如何在 WSL 中使用？

```bash
# 在 WSL 终端中配置
echo 'export OPENAI_API_KEY="sk-xxx"' >> ~/.bashrc
echo 'export OPENAI_API_BASE="https://api.bigaipro.com/v1"' >> ~/.bashrc
source ~/.bashrc
```

---

## 工具对比

| 工具 | 优势 | 适用场景 |
|------|------|---------|
| Claude Code | 代码理解强、上下文长 | 复杂重构、代码审查 |
| Codex CLI | 系统命令执行、文件操作 | 自动化任务、脚本生成 |
| Gemini CLI | 多模态、免费额度 | 项目分析、文档生成 |

---

## 下一步

- 📖 [Claude Code 教程](./claude-code.md) - 详细使用指南
- 📖 [Codex CLI 教程](./codex-cli.md) - 详细使用指南
- 📖 [Gemini CLI 教程](./gemini-cli.md) - 详细使用指南

---

<div align="center">

**CLI + IDE = 高效 AI 编程**

</div>
