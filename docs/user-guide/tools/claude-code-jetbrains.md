# Claude Code JetBrains 插件

<div align="center">

**在 JetBrains IDE 中使用 Claude Code**

*支持 IntelliJ IDEA / PyCharm / WebStorm / GoLand 等*

</div>

---

## 简介

Claude Code 官方提供了 JetBrains IDE 插件，可以将 Claude Code 的能力集成到 JetBrains 系列 IDE 中，提供交互式 Diff 查看、选中代码上下文共享等功能。

### 核心功能

- ✅ **交互式 Diff 查看** - 在 IDE 内查看和应用代码变更
- ✅ **选中代码共享** - 将选中的代码发送给 Claude
- ✅ **IDE 集成** - 与 JetBrains IDE 深度集成
- ✅ **终端集成** - 在 IDE 终端中运行 Claude Code

---

## 支持的 IDE

| IDE | 版本要求 |
|-----|---------|
| IntelliJ IDEA | 2024.1+ |
| PyCharm | 2024.1+ |
| WebStorm | 2024.1+ |
| GoLand | 2024.1+ |
| Rider | 2024.1+ |
| CLion | 2024.1+ |
| PhpStorm | 2024.1+ |
| RubyMine | 2024.1+ |

---

## 安装

### 步骤一：安装 Claude Code CLI

首先确保已安装 Claude Code CLI：

```bash
npm install -g @anthropic-ai/claude-code
```

验证安装：

```bash
claude --version
```

### 步骤二：安装 JetBrains 插件

1. 打开 JetBrains IDE
2. 进入 `Settings/Preferences` → `Plugins`
3. 搜索 "Claude Code"
4. 点击 "Install" 安装
5. 重启 IDE

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

#### Windows

```powershell
# 永久设置
[Environment]::SetEnvironmentVariable("ANTHROPIC_API_KEY", "sk-xxxxxxxxxxxxxxxx", "User")
[Environment]::SetEnvironmentVariable("ANTHROPIC_BASE_URL", "https://api.bigaipro.com", "User")
```

### 插件设置

1. 打开 `Settings` → `Tools` → `Claude Code [Beta]`
2. 配置选项：

| 配置项 | 说明 |
|--------|------|
| Claude Command Path | Claude CLI 路径（通常自动检测） |
| Auto-edit Mode | 自动应用编辑 |
| Show Notifications | 显示通知 |

---

## 使用方法

### 从 IDE 终端使用

1. 打开 IDE 内置终端
2. 运行 `claude` 命令
3. 所有集成功能自动激活

```bash
# 在 IDE 终端中
claude

# 开始对话
> 解释这个项目的架构
> 为 UserService 添加单元测试
```

### 从外部终端连接

如果在外部终端运行 Claude Code，使用 `/ide` 命令连接到 IDE：

```bash
# 在外部终端
claude

# 连接到 IDE
/ide

# 现在可以使用 IDE 集成功能
```

### 配置 Diff 工具

1. 运行 `claude`
2. 输入 `/config` 命令
3. 设置 diff 工具为 `auto`（自动检测 IDE）

---

## 功能特性

### 交互式 Diff 查看

当 Claude 建议代码修改时：

1. 修改会在 IDE 的 Diff 查看器中显示
2. 可以逐行查看变更
3. 选择接受或拒绝修改

### 选中代码共享

1. 在编辑器中选中代码
2. 右键 → "Send to Claude Code"
3. 或使用快捷键

### 文件上下文

Claude Code 会自动获取：
- 当前打开的文件
- 项目结构
- 相关依赖

---

## 快捷键

| 操作 | 快捷键 |
|------|--------|
| 发送选中代码 | `Alt+C` |
| 打开 Claude 面板 | `Alt+Shift+C` |
| 中断操作 | `Esc` |

### ESC 键配置

如果 ESC 键无法中断 Claude Code 操作：

1. 打开 `Settings` → `Tools` → `Terminal`
2. 选择以下任一方式：
   - 禁用 "Override IDE shortcuts"
   - 从覆盖列表中移除 ESC 键
3. 应用更改

---

## 特殊配置

### 远程开发

使用 JetBrains 远程开发时：

- 插件必须安装在**远程主机**上
- 不是安装在本地客户端

### WSL 配置

Windows WSL 用户需要额外配置：

1. 确保 WSL 中已安装 Claude Code
2. 在插件设置中配置 WSL 命令格式
3. 可能需要配置 `wsl.exe` 路径

---

## 与 JetBrains AI Assistant 对比

| 特性 | Claude Code 插件 | JetBrains AI Assistant |
|------|-----------------|----------------------|
| 模型 | Claude 系列 | 多种模型 |
| 自定义 API | ✅ 支持 | ✅ 支持 BYOK |
| Agent 能力 | ✅ 强大 | ✅ Junie Agent |
| 代码补全 | ❌ | ✅ |
| 价格 | 按量付费 | 订阅制 |

### 推荐组合

- **日常开发**：JetBrains AI Assistant（代码补全）
- **复杂任务**：Claude Code 插件（Agent 能力）

---

## 常见问题

### Q: 插件无法检测到 Claude Code？

1. 确认 Claude Code 已安装：`claude --version`
2. 检查 PATH 环境变量
3. 在插件设置中手动配置路径

### Q: IDE 未被检测到？

1. 确保从 IDE 终端运行 Claude
2. 或使用 `/ide` 命令手动连接
3. 检查插件是否正确安装

### Q: 命令未找到？

1. 验证安装：`npm list -g @anthropic-ai/claude-code`
2. 在插件设置中配置 Claude 命令路径
3. WSL 用户使用 WSL 命令格式

### Q: Diff 查看器不工作？

1. 运行 `/config` 设置 diff 工具为 `auto`
2. 重启 IDE
3. 检查插件版本是否最新

---

## 安全注意事项

当 Claude Code 在 JetBrains IDE 中以自动编辑模式运行时：

- 可能修改 IDE 配置文件
- 这些文件可能被 IDE 自动执行
- 建议在自动编辑模式下谨慎操作

推荐做法：
- 使用建议模式而非自动模式
- 审查所有代码变更
- 定期备份重要配置

---

## 下一步

- 📖 [Claude Code 教程](./claude-code.md) - CLI 详细使用指南
- 📖 [JetBrains IDE 配置](./jetbrains-ide.md) - AI Assistant 配置
- 📖 [CC-Switch](./cc-switch.md) - 配置管理工具

---

<div align="center">

**Claude Code - JetBrains IDE 的 AI 编程伙伴**

[官方文档](https://code.claude.com/docs/en/jetbrains)

</div>
