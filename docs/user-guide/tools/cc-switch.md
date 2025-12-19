# CC-Switch 使用教程

<div align="center">

**Claude Code / Codex / Gemini CLI 一体化管理工具**

*支持 macOS / Windows / Linux*

</div>

---

## 简介

CC-Switch 是一款跨平台桌面应用，用于统一管理 Claude Code、Codex CLI 和 Gemini CLI 的配置。通过图形界面，你可以轻松切换 API 供应商、管理 MCP 服务器、配置 Skills 扩展和系统提示词。

### 核心功能

- ✅ **一键切换 API 配置** - 在多个 API 提供商之间快速切换
- ✅ **可视化配置管理** - 通过图形界面轻松管理所有配置
- ✅ **MCP 服务器管理** - 管理 Model Context Protocol 服务器
- ✅ **Skills 扩展管理** - 管理自定义技能扩展
- ✅ **系统托盘快捷操作** - 通过托盘菜单快速切换
- ✅ **云同步支持** - 支持 Dropbox、OneDrive、iCloud 等云同步

---

## 安装指南

### Windows

1. 访问 [GitHub Releases](https://github.com/farion1231/cc-switch/releases)
2. 下载 `CC-Switch-v{version}-Windows.msi` 安装包
3. 运行安装程序完成安装

或下载便携版：`CC-Switch-v{version}-Windows-Portable.zip`

### macOS

#### 方式一：Homebrew（推荐）

```bash
# 添加 tap 源
brew tap farion1231/ccswitch

# 安装 CC-Switch
brew install --cask cc-switch
```

更新：

```bash
brew upgrade --cask cc-switch
```

#### 方式二：手动下载

1. 从 [Releases](https://github.com/farion1231/cc-switch/releases) 下载 `CC-Switch-v{version}-macOS.zip`
2. 解压后拖入应用程序文件夹

> 首次启动可能提示"未知开发者"，前往「系统设置」→「隐私与安全性」→ 点击「仍要打开」

### Linux

#### Debian/Ubuntu

```bash
# 下载 .deb 包
wget https://github.com/farion1231/cc-switch/releases/latest/download/cc-switch_x.x.x_amd64.deb

# 安装
sudo dpkg -i cc-switch_x.x.x_amd64.deb
```

#### Arch Linux

```bash
paru -S cc-switch
```

#### AppImage

从 Releases 下载 `CC-Switch-v{version}-Linux.AppImage`

---

## 环境检查

在配置 CC-Switch 前，请确保已安装相关 CLI 工具。

### 检查 Node.js

```bash
node -v  # 需要 18+
```

### 检查 CLI 工具

```bash
# Claude Code
claude --version

# Codex CLI
codex --version

# Gemini CLI
gemini --version
```

如未安装，请参考对应教程：
- [Claude Code 安装](./claude-code.md)
- [Codex CLI 安装](./codex-cli.md)
- [Gemini CLI 安装](./gemini-cli.md)

---

## 配置 BigAI Pro

### Claude Code 配置

1. 打开 CC-Switch
2. 在分组条中选择 **Claude**
3. 点击 **添加供应商**
4. 填写配置：

| 配置项 | 值 |
|--------|-----|
| **名称** | BigAI Pro |
| **API Key** | `sk-xxxxxxxxxxxxxxxx` |
| **API Base URL** | `https://api.bigaipro.com` |

5. 在配置 JSON 中添加 `apiKeyHelper`（与 `env` 同级）：

```json
{
  "env": {
    "ANTHROPIC_API_KEY": "sk-xxxxxxxxxxxxxxxx",
    "ANTHROPIC_BASE_URL": "https://api.bigaipro.com"
  },
  "apiKeyHelper": "echo 'sk-xxxxxxxxxxxxxxxx'"
}
```

6. 点击 **添加** 保存
7. 在主界面点击 **启用** 按钮

### Codex CLI 配置

1. 在分组条中选择 **Codex**
2. 点击 **添加供应商**
3. 填写配置：

| 配置项 | 值 |
|--------|-----|
| **名称** | BigAI Pro |
| **API Key** | `sk-xxxxxxxxxxxxxxxx` |
| **API 请求地址** | `https://api.bigaipro.com/v1` |

4. 点击 **添加** 保存
5. 在主界面点击 **启用** 按钮

### Gemini CLI 配置

1. 在分组条中选择 **Gemini**
2. 点击 **添加供应商**
3. 填写配置：

| 配置项 | 值 |
|--------|-----|
| **名称** | BigAI Pro |
| **API Key** | `sk-xxxxxxxxxxxxxxxx` |

> 注意：Gemini CLI 目前官方不支持自定义 API Base URL，详见 [Gemini CLI 教程](./gemini-cli.md)

4. 点击 **添加** 保存
5. 在主界面点击 **启用** 按钮

---

## 验证配置

配置完成后，重启终端并测试：

```bash
# 测试 Claude Code
claude
> 你好

# 测试 Codex CLI
codex "你好"

# 测试 Gemini CLI
gemini
> 你好
```

---

## MCP 服务器管理

CC-Switch 提供可视化的 MCP 服务器管理功能。

### 添加 MCP 服务器

1. 点击 **MCP 管理** 标签
2. 点击 **添加服务器**
3. 填写配置：

```json
{
  "name": "filesystem",
  "command": "npx",
  "args": ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/dir"]
}
```

4. 点击保存

### 常用 MCP 服务器

| 服务器 | 用途 |
|--------|------|
| `@modelcontextprotocol/server-filesystem` | 文件系统访问 |
| `@modelcontextprotocol/server-github` | GitHub 操作 |
| `@modelcontextprotocol/server-postgres` | PostgreSQL 数据库 |
| `@anthropic-ai/mcp-server-fetch` | HTTP 请求 |

---

## Skills 扩展管理

Skills 是自定义的命令扩展，可以增强 CLI 工具的能力。

### 添加 Skill

1. 点击 **Skills 管理** 标签
2. 点击 **添加 Skill**
3. 配置 Skill 内容

### 示例 Skill

```json
{
  "name": "code-review",
  "description": "代码审查",
  "prompt": "请审查以下代码，指出潜在问题和改进建议"
}
```

---

## 系统提示词管理

管理不同场景的系统提示词。

### 添加提示词

1. 点击 **Prompts 管理** 标签
2. 点击 **添加提示词**
3. 填写名称和内容

### 示例提示词

```
你是一位资深的 Go 语言开发者，擅长：
- 高性能后端开发
- 微服务架构设计
- 数据库优化

请用简洁专业的方式回答问题。
```

---

## 云同步配置

支持跨设备同步配置。

### 设置云同步

1. 打开 **设置**
2. 找到 **自定义配置目录**
3. 选择云同步文件夹（Dropbox、OneDrive、iCloud 等）
4. 重启应用生效
5. 在其他设备上重复此步骤

---

## 配置文件位置

CC-Switch 管理的配置文件位置：

### Claude Code

| 系统 | 路径 |
|------|------|
| macOS | `~/.claude.json` |
| Windows | `%USERPROFILE%\.claude.json` |
| Linux | `~/.claude.json` |

### Codex CLI

| 系统 | 路径 |
|------|------|
| macOS | `~/.codex/config.json` |
| Windows | `%USERPROFILE%\.codex\config.json` |
| Linux | `~/.codex/config.json` |

### Gemini CLI

| 系统 | 路径 |
|------|------|
| macOS | `~/.gemini/settings.json` |
| Windows | `%USERPROFILE%\.gemini\settings.json` |
| Linux | `~/.gemini/settings.json` |

### CC-Switch 存储

| 系统 | 路径 |
|------|------|
| macOS | `~/Library/Application Support/cc-switch/` |
| Windows | `%APPDATA%\cc-switch\` |
| Linux | `~/.config/cc-switch/` |

---

## 快捷操作

### 系统托盘

CC-Switch 运行时会在系统托盘显示图标，支持：

- 快速切换供应商
- 查看当前配置
- 打开主界面
- 退出应用

### 快捷键

| 操作 | 快捷键 |
|------|--------|
| 打开主界面 | 点击托盘图标 |
| 快速切换 | 右键托盘图标 |

---

## 常见问题

### Q: 配置后 CLI 工具无响应？

1. 确保已点击"启用"按钮
2. 重启终端
3. 检查 API Key 是否正确

### Q: 如何恢复官方配置？

1. 选择对应的 CLI 分组
2. 选择"官方登录"预设
3. 点击启用
4. 重启 CLI 工具并重新登录

### Q: 配置丢失怎么办？

1. 检查云同步是否正常
2. 查看 CC-Switch 存储目录
3. 从备份恢复

### Q: 如何更新 CC-Switch？

```bash
# macOS (Homebrew)
brew upgrade --cask cc-switch

# 其他系统
# 从 GitHub Releases 下载最新版本
```

---

## 下一步

- 📖 [Claude Code 教程](./claude-code.md) - 详细使用指南
- 📖 [Codex CLI 教程](./codex-cli.md) - 详细使用指南
- 📖 [Gemini CLI 教程](./gemini-cli.md) - 详细使用指南

---

<div align="center">

**CC-Switch - 一站式 AI CLI 管理**

[GitHub](https://github.com/farion1231/cc-switch) · [Releases](https://github.com/farion1231/cc-switch/releases)

</div>
