# Claude Code GUI (JetBrains 插件)

<div align="center">

**在 JetBrains IDE 中使用 Claude Code 的图形界面**

*第三方插件 - 支持 IntelliJ IDEA 及全系列 JetBrains IDE*

</div>

---

## 简介

Claude Code GUI 是一款第三方 JetBrains 插件，为 Claude Code 提供直观的图形界面，让你无需离开 IDE 即可使用 Claude AI 进行代码生成、分析和重构。

### 核心功能

- ✅ **图形化界面** - 直观的 GUI 界面，无需命令行
- ✅ **AI 代码生成** - 智能代码生成和重构
- ✅ **上下文感知** - 理解项目上下文的 AI 助手
- ✅ **代码分析** - 选中代码直接发送给 Claude 分析
- ✅ **集成 Claude Code CLI** - 完整的 Claude Code 功能

---

## 安装

### 从 JetBrains 插件市场安装

1. 打开 JetBrains IDE（IntelliJ IDEA、PyCharm、WebStorm 等）
2. 进入 `Settings/Preferences` → `Plugins`
3. 点击 `Marketplace` 标签
4. 搜索 "Claude Code GUI"
5. 点击 `Install` 安装
6. 重启 IDE

### 直接访问

[JetBrains 插件市场 - Claude Code GUI](https://plugins.jetbrains.com/plugin/29342-claude-code-gui)

---

## 前置要求

### 安装 Claude Code CLI

插件需要 Claude Code CLI 支持：

```bash
npm install -g @anthropic-ai/claude-code
```

验证安装：

```bash
claude --version
```

---

## 配置 BigAI Pro API

### 环境变量配置

#### macOS / Linux

```bash
# 编辑 ~/.zshrc 或 ~/.bashrc
export ANTHROPIC_API_KEY="sk-xxxxxxxxxxxxxxxx"
export ANTHROPIC_BASE_URL="https://api.bigaipro.com"

# 使配置生效
source ~/.zshrc
```

#### Windows PowerShell

```powershell
# 永久设置
[Environment]::SetEnvironmentVariable("ANTHROPIC_API_KEY", "sk-xxxxxxxxxxxxxxxx", "User")
[Environment]::SetEnvironmentVariable("ANTHROPIC_BASE_URL", "https://api.bigaipro.com", "User")
```

#### Windows CMD

```cmd
setx ANTHROPIC_API_KEY "sk-xxxxxxxxxxxxxxxx"
setx ANTHROPIC_BASE_URL "https://api.bigaipro.com"
```

### 重启 IDE

配置环境变量后，需要重启 JetBrains IDE 使配置生效。

---

## 使用方法

### 打开工具窗口

1. 在 IDE 右侧边栏找到 "Claude Code GUI" 图标
2. 点击打开工具窗口

### 发送代码到 Claude

**方式一：快捷键**

| 系统 | 快捷键 |
|------|--------|
| Windows/Linux | `Ctrl+Alt+K` |
| macOS | `Cmd+Alt+K` |

**方式二：右键菜单**

1. 在编辑器中选中代码
2. 右键 → "Send to Claude"

### 对话交互

1. 在工具窗口的输入框中输入问题
2. 按 Enter 发送
3. 等待 Claude 响应

---

## 常用场景

### 代码生成

```
你：创建一个用户认证服务类，包含登录、注册、JWT 验证方法
Claude：我来为你创建这个服务类...
```

### 代码解释

```
你：[选中代码] 解释这段代码的作用
Claude：这段代码实现了...
```

### 代码重构

```
你：[选中代码] 重构这段代码，提高可读性和性能
Claude：我建议以下修改...
```

### Bug 修复

```
你：[选中代码] 这段代码有什么问题？如何修复？
Claude：我发现以下问题...
```

### 单元测试

```
你：[选中代码] 为这个方法生成单元测试
Claude：我来创建测试用例...
```

---

## 快捷键

| 操作 | Windows/Linux | macOS |
|------|---------------|-------|
| 发送选中代码 | `Ctrl+Alt+K` | `Cmd+Alt+K` |
| 打开工具窗口 | 点击侧边栏图标 | 点击侧边栏图标 |

---

## 与官方插件对比

| 特性 | Claude Code GUI | 官方 Claude Code 插件 |
|------|-----------------|---------------------|
| 来源 | 第三方 (CodeMossAI) | Anthropic 官方 |
| 界面 | 完整 GUI | 终端 + Diff 查看 |
| 安装 | 插件市场直接安装 | 需要 CLI + 插件 |
| 功能 | 对话式交互 | CLI 集成 |
| 适用场景 | 快速问答 | 复杂任务 |

### 推荐组合

- **快速问答**：Claude Code GUI（图形界面）
- **复杂任务**：官方 Claude Code 插件（CLI 集成）

---

## 常见问题

### Q: 插件无法连接到 Claude？

1. 确认 Claude Code CLI 已安装
2. 检查环境变量配置
3. 重启 IDE

### Q: 响应速度慢？

1. 检查网络连接
2. 确认 API 端点配置正确
3. 尝试减少发送的代码量

### Q: 快捷键冲突？

1. 打开 `Settings` → `Keymap`
2. 搜索 "Claude"
3. 修改为其他快捷键

### Q: 如何更新插件？

1. 打开 `Settings` → `Plugins`
2. 点击 `Updates` 标签
3. 找到 Claude Code GUI 并更新

---

## 相关资源

- [GitHub 源码](https://github.com/zhukunpenglinyutong/idea-claude-code-gui)
- [JetBrains 插件市场](https://plugins.jetbrains.com/plugin/29342-claude-code-gui)

---

## 下一步

- 📖 [Claude Code 教程](./claude-code.md) - CLI 详细使用指南
- 📖 [Claude Code JetBrains 官方插件](./claude-code-jetbrains.md) - 官方插件
- 📖 [JetBrains IDE 配置](./jetbrains-ide.md) - AI Assistant 配置

---

<div align="center">

**Claude Code GUI - JetBrains IDE 的图形化 AI 助手**

</div>
