# Claude Code 使用教程

<div align="center">

**AI 驱动的终端编程助手**

*支持 macOS / Windows / Linux*

</div>

---

## 简介

Claude Code 是 Anthropic 推出的 AI 编程助手，直接在终端中运行，可以理解你的代码库、执行命令、编辑文件，帮助你更高效地编程。

通过配置 BigAI Pro API，你可以使用 Claude Code 的全部功能。

---

## 安装指南

### macOS

#### 方式一：Homebrew（推荐）

```bash
# 安装 Homebrew（如果没有）
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 安装 Claude Code
brew install claude-code
```

#### 方式二：NPM

```bash
# 确保已安装 Node.js 18+
node --version

# 全局安装
npm install -g @anthropic-ai/claude-code
```

#### 方式三：直接下载

```bash
# 下载最新版本
curl -fsSL https://claude.ai/code/install.sh | bash
```

### Windows

#### 方式一：Winget（推荐）

```powershell
# 使用 Windows 包管理器
winget install Anthropic.ClaudeCode
```

#### 方式二：NPM

```powershell
# 确保已安装 Node.js 18+
node --version

# 全局安装
npm install -g @anthropic-ai/claude-code
```

#### 方式三：Scoop

```powershell
# 添加 bucket
scoop bucket add extras

# 安装
scoop install claude-code
```

#### 方式四：手动安装

1. 访问 [Claude Code 下载页面](https://claude.ai/code/download)
2. 下载 Windows 安装包 (.msi)
3. 运行安装程序

### Linux

#### Ubuntu / Debian

```bash
# 添加 APT 源
curl -fsSL https://claude.ai/code/gpg | sudo gpg --dearmor -o /usr/share/keyrings/claude-code.gpg
echo "deb [signed-by=/usr/share/keyrings/claude-code.gpg] https://claude.ai/code/apt stable main" | sudo tee /etc/apt/sources.list.d/claude-code.list

# 安装
sudo apt update
sudo apt install claude-code
```

#### Fedora / RHEL / CentOS

```bash
# 添加 YUM 源
sudo tee /etc/yum.repos.d/claude-code.repo << 'EOF'
[claude-code]
name=Claude Code
baseurl=https://claude.ai/code/rpm
enabled=1
gpgcheck=1
gpgkey=https://claude.ai/code/gpg
EOF

# 安装
sudo dnf install claude-code
```

#### Arch Linux

```bash
# 使用 AUR
yay -S claude-code
```

#### 通用方式：NPM

```bash
# 确保已安装 Node.js 18+
node --version

# 全局安装
npm install -g @anthropic-ai/claude-code
```

---

## 配置 BigAI Pro API

### 方式一：环境变量（推荐）

#### macOS / Linux

编辑 `~/.bashrc` 或 `~/.zshrc`：

```bash
# BigAI Pro API 配置
export ANTHROPIC_API_KEY="sk-xxxxxxxxxxxxxxxx"
export ANTHROPIC_BASE_URL="https://api.bigaipro.com"
```

使配置生效：

```bash
source ~/.bashrc  # 或 source ~/.zshrc
```

#### Windows PowerShell

编辑 PowerShell 配置文件：

```powershell
# 打开配置文件
notepad $PROFILE

# 添加以下内容
$env:ANTHROPIC_API_KEY = "sk-xxxxxxxxxxxxxxxx"
$env:ANTHROPIC_BASE_URL = "https://api.bigaipro.com"
```

或设置系统环境变量：

```powershell
# 设置用户级环境变量
[Environment]::SetEnvironmentVariable("ANTHROPIC_API_KEY", "sk-xxxxxxxxxxxxxxxx", "User")
[Environment]::SetEnvironmentVariable("ANTHROPIC_BASE_URL", "https://api.bigaipro.com", "User")
```

#### Windows CMD

```cmd
# 设置用户级环境变量
setx ANTHROPIC_API_KEY "sk-xxxxxxxxxxxxxxxx"
setx ANTHROPIC_BASE_URL "https://api.bigaipro.com"
```

### 方式二：配置文件

创建配置文件 `~/.claude-code/config.json`：

```json
{
  "apiKey": "sk-xxxxxxxxxxxxxxxx",
  "baseUrl": "https://api.bigaipro.com",
  "model": "claude-sonnet-4.5",
  "maxTokens": 8192
}
```

### 方式三：命令行参数

```bash
claude-code --api-key "sk-xxxxxxxx" --base-url "https://api.bigaipro.com"
```

---

## 验证配置

```bash
# 检查版本
claude-code --version

# 验证 API 连接
claude-code --check

# 查看当前配置
claude-code config show
```

---

## 基础使用

### 启动 Claude Code

```bash
# 在当前目录启动
claude-code

# 在指定目录启动
claude-code /path/to/project

# 使用特定模型
claude-code --model claude-sonnet-4.5
```

### 常用命令

在 Claude Code 交互界面中：

```
# 询问代码问题
> 解释这个函数的作用

# 编辑文件
> 修改 src/main.py，添加错误处理

# 执行命令
> 运行测试

# 搜索代码
> 找到所有使用 deprecated API 的地方

# 生成代码
> 创建一个 REST API 端点处理用户注册

# 重构代码
> 将这个类重构为使用依赖注入

# 调试
> 这个错误是什么原因？如何修复？
```

### 快捷键

| 快捷键 | 功能 |
|--------|------|
| `Ctrl+C` | 取消当前操作 |
| `Ctrl+D` | 退出 Claude Code |
| `Ctrl+L` | 清屏 |
| `↑` / `↓` | 浏览历史命令 |
| `Tab` | 自动补全 |

---

## 高级功能

### Checkpoints（检查点）

Claude Code 2.0 新增的检查点功能，可以保存进度并回滚：

```bash
# 创建检查点
> /checkpoint save "添加用户认证前"

# 查看检查点
> /checkpoint list

# 回滚到检查点
> /checkpoint restore "添加用户认证前"
```

### 多文件操作

```bash
# 批量修改
> 将所有 .js 文件中的 var 替换为 const

# 跨文件重构
> 将 UserService 类拆分到单独的文件

# 项目范围搜索
> 找到所有未使用的导入
```

### 与 Git 集成

```bash
# 查看更改
> 显示我今天的所有更改

# 生成提交信息
> 为当前更改生成 commit message

# 代码审查
> 审查最近的 PR 更改
```

### 测试生成

```bash
# 生成单元测试
> 为 UserService 类生成单元测试

# 生成集成测试
> 为 API 端点生成集成测试

# 修复失败的测试
> 这个测试为什么失败？修复它
```

---

## 配置选项

### 完整配置示例

`~/.claude-code/config.json`：

```json
{
  "apiKey": "sk-xxxxxxxxxxxxxxxx",
  "baseUrl": "https://api.bigaipro.com",
  "model": "claude-sonnet-4.5",
  "maxTokens": 8192,
  "temperature": 0.7,
  "autoSave": true,
  "checkpointEnabled": true,
  "theme": "dark",
  "editor": "vim",
  "shell": "/bin/zsh",
  "ignorePatterns": [
    "node_modules",
    ".git",
    "dist",
    "*.log"
  ],
  "contextFiles": [
    "README.md",
    "package.json"
  ]
}
```

### 配置说明

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| apiKey | API 密钥 | 必填 |
| baseUrl | API 地址 | https://api.anthropic.com |
| model | 使用的模型 | claude-sonnet-4.5 |
| maxTokens | 最大输出长度 | 4096 |
| temperature | 随机性 | 0.7 |
| autoSave | 自动保存更改 | true |
| checkpointEnabled | 启用检查点 | true |
| ignorePatterns | 忽略的文件模式 | [] |

---

## 常见问题

### Q: 连接失败怎么办？

```bash
# 检查网络
curl -I https://api.bigaipro.com/v1/models

# 检查 API Key
claude-code --check

# 查看详细日志
claude-code --debug
```

### Q: 如何切换模型？

```bash
# 命令行指定
claude-code --model claude-opus-4.5

# 或在配置文件中修改
# ~/.claude-code/config.json
{
  "model": "claude-opus-4.5"
}
```

### Q: 如何处理大型项目？

```bash
# 使用 ignorePatterns 排除不需要的文件
# 在配置文件中添加：
{
  "ignorePatterns": [
    "node_modules",
    ".git",
    "dist",
    "build",
    "*.min.js",
    "*.map"
  ]
}
```

### Q: 如何提高响应速度？

1. 使用 `claude-haiku-4.5` 模型（更快但能力稍弱）
2. 减少上下文文件数量
3. 使用 `ignorePatterns` 排除大文件

---

## 最佳实践

### 1. 项目初始化

```bash
# 进入项目目录
cd /path/to/project

# 启动 Claude Code
claude-code

# 让 Claude 了解项目
> 分析这个项目的结构和技术栈
```

### 2. 高效提问

```bash
# ✅ 好的提问
> 在 src/services/user.ts 中，为 createUser 函数添加输入验证

# ❌ 不好的提问
> 帮我写代码
```

### 3. 迭代开发

```bash
# 创建检查点
> /checkpoint save "开始重构"

# 进行更改
> 重构 UserService 使用依赖注入

# 测试
> 运行测试

# 如果有问题，回滚
> /checkpoint restore "开始重构"
```

---

## 下一步

- 📖 [Codex CLI 教程](./codex-cli.md) - OpenAI 的命令行工具
- 📖 [Cursor IDE 配置](./cursor-ide.md) - AI 增强的代码编辑器
- 📖 [Claude 模型详解](../models/claude-models.md) - 了解 Claude 模型特性

---

<div align="center">

**Claude Code - 让编程更智能**

</div>
