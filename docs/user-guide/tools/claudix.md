# Claudix - VS Code 扩展

<div align="center">

**在 VS Code 中使用 Claude Code 的图形界面**

</div>

---

## 简介

Claudix 是一款 VS Code 扩展，将 Claude Code 的能力直接集成到编辑器中，提供图形化的对话界面、历史记录管理和工具集成功能。

### 核心功能

- ✅ **图形化对话界面** - 在侧边栏直接与 Claude 对话
- ✅ **对话历史管理** - 保存和恢复历史对话
- ✅ **工具集成** - 支持代码编辑、文件操作等工具
- ✅ **代码理解** - 智能理解项目上下文

---

## 安装

### 从 VSIX 安装

1. 访问 [Claudix GitHub](https://github.com/Haleclipse/Claudix/releases)
2. 下载最新的 `.vsix` 文件
3. 在 VS Code 中：
   - 打开命令面板 `Cmd+Shift+P`
   - 输入 "Extensions: Install from VSIX"
   - 选择下载的 `.vsix` 文件

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/Haleclipse/Claudix.git
cd Claudix

# 安装依赖
npm install

# 构建
npm run build

# 打包
npm run package
```

---

## 配置 BigAI Pro API

Claudix 使用 Claude Code 的配置，需要先配置环境变量。

### macOS / Linux

```bash
# 编辑 ~/.zshrc 或 ~/.bashrc
export ANTHROPIC_API_KEY="sk-xxxxxxxxxxxxxxxx"
export ANTHROPIC_BASE_URL="https://api.bigaipro.com"

# 使配置生效
source ~/.zshrc
```

### Windows PowerShell

```powershell
# 永久设置
[Environment]::SetEnvironmentVariable("ANTHROPIC_API_KEY", "sk-xxxxxxxxxxxxxxxx", "User")
[Environment]::SetEnvironmentVariable("ANTHROPIC_BASE_URL", "https://api.bigaipro.com", "User")
```

### VS Code 设置

也可以在 VS Code 设置中配置终端环境变量：

```json
{
  "terminal.integrated.env.osx": {
    "ANTHROPIC_API_KEY": "sk-xxxxxxxxxxxxxxxx",
    "ANTHROPIC_BASE_URL": "https://api.bigaipro.com"
  },
  "terminal.integrated.env.windows": {
    "ANTHROPIC_API_KEY": "sk-xxxxxxxxxxxxxxxx",
    "ANTHROPIC_BASE_URL": "https://api.bigaipro.com"
  },
  "terminal.integrated.env.linux": {
    "ANTHROPIC_API_KEY": "sk-xxxxxxxxxxxxxxxx",
    "ANTHROPIC_BASE_URL": "https://api.bigaipro.com"
  }
}
```

---

## 使用方法

### 打开 Claudix 面板

1. 点击活动栏中的 Claudix 图标
2. 或使用命令面板 `Cmd+Shift+P` → "Claudix: Open Panel"

### 开始对话

1. 在输入框中输入问题
2. 按 Enter 发送
3. 等待 Claude 响应

### 常用操作

```
你：解释这个文件的架构
Claude：这个文件实现了...

你：为这个函数添加单元测试
Claude：我来为你创建测试文件...

你：重构这段代码，提高可读性
Claude：我建议以下修改...
```

### 工具操作

当 Claude 需要执行工具操作时（如编辑文件、运行命令），会提示你确认：

- ✅ **Approve** - 允许执行
- ❌ **Reject** - 拒绝执行
- 👁️ **Review** - 查看详情

---

## 功能特性

### 对话历史

- 自动保存对话历史
- 可以恢复之前的对话
- 支持多个对话会话

### 代码上下文

Claudix 会自动获取：
- 当前打开的文件
- 选中的代码
- 项目结构信息

### 工具集成

支持的工具操作：
- 文件读取和编辑
- 代码搜索
- 终端命令执行
- Git 操作

---

## 快捷键

| 操作 | 快捷键 |
|------|--------|
| 打开面板 | `Cmd+Shift+C` |
| 新建对话 | `Cmd+N` |
| 发送消息 | `Enter` |
| 换行 | `Shift+Enter` |

---

## 与 Claude Code CLI 对比

| 特性 | Claudix | Claude Code CLI |
|------|---------|-----------------|
| 界面 | 图形化 | 命令行 |
| 集成度 | VS Code 内置 | 独立终端 |
| 对话历史 | 可视化管理 | 文件存储 |
| 适用场景 | 日常开发 | 复杂任务 |

---

## 常见问题

### Q: 扩展无法启动？

1. 确保已安装 Claude Code CLI
2. 检查环境变量配置
3. 重启 VS Code

### Q: 对话响应慢？

1. 检查网络连接
2. 尝试使用更快的模型
3. 减少上下文大小

### Q: 工具操作失败？

1. 检查文件权限
2. 确认工作目录正确
3. 查看输出面板的错误信息

---

## 下一步

- 📖 [Claude Code 教程](./claude-code.md) - CLI 详细使用指南
- 📖 [VS Code CLI 集成](./vscode-cli.md) - 在 VS Code 中使用 CLI
- 📖 [CC-Switch](./cc-switch.md) - 配置管理工具

---

<div align="center">

**Claudix - VS Code 中的 Claude Code**

[GitHub](https://github.com/Haleclipse/Claudix)

</div>
