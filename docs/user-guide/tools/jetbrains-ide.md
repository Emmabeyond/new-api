# JetBrains IDE 配置指南

<div align="center">

**IntelliJ IDEA / PyCharm / WebStorm / GoLand 等**

</div>

---

## 简介

JetBrains 系列 IDE（IntelliJ IDEA、PyCharm、WebStorm、GoLand、Rider 等）通过 AI Assistant 插件支持 AI 编程功能。2025 年 12 月起，JetBrains 正式支持 BYOK（Bring Your Own Key），可以使用自定义 API 端点。

通过配置 BigAI Pro API，你可以在 JetBrains IDE 中使用 GPT、Claude、Gemini 等顶尖模型。

---

## 安装 AI Assistant 插件

### 步骤一：打开插件市场

1. 启动 JetBrains IDE
2. 打开 `Settings/Preferences` → `Plugins`
3. 搜索 "AI Assistant"

### 步骤二：安装插件

1. 点击 "Install" 安装 JetBrains AI Assistant
2. 重启 IDE

### 支持的 IDE

| IDE | 版本要求 |
|-----|---------|
| IntelliJ IDEA | 2024.3+ |
| PyCharm | 2024.3+ |
| WebStorm | 2024.3+ |
| GoLand | 2024.3+ |
| Rider | 2024.3+ |
| CLion | 2024.3+ |
| PhpStorm | 2024.3+ |
| RubyMine | 2024.3+ |

---

## 配置 BigAI Pro API

### 方式一：BYOK 配置（推荐）

JetBrains 支持 OpenAI 兼容的 API 端点。

#### 步骤一：打开设置

1. 打开 `Settings/Preferences`
2. 导航到 `Tools` → `AI Assistant` → `Models`

#### 步骤二：配置 OpenAI 兼容端点

1. 在 "Third-party AI providers" 部分
2. 选择 Provider: `OpenAI API Compatible`
3. 填写配置：

| 配置项 | 值 |
|--------|-----|
| **URL** | `https://api.bigaipro.com/v1` |
| **API Key** | `sk-xxxxxxxxxxxxxxxx` |
| **Tool calling** | ✅ 启用 |

4. 点击 "Test Connection" 测试连接
5. 点击 "Apply" 保存

#### 步骤三：选择模型

连接成功后，在模型选择器中可以看到可用模型：

- `gpt-5.2-pro` - 复杂任务
- `gpt-5.2-instant` - 日常编程
- `claude-sonnet-4.5` - 代码生成
- `deepseek-chat` - 性价比之选

### 方式二：使用 Anthropic API

如果你主要使用 Claude 模型：

1. 选择 Provider: `Anthropic`
2. 填写配置：

| 配置项 | 值 |
|--------|-----|
| **API Key** | `sk-xxxxxxxxxxxxxxxx` |

> 注意：Anthropic 原生端点不支持自定义 URL，建议使用 OpenAI 兼容方式。

---

## 模型分配

JetBrains AI Assistant 允许为不同功能分配不同模型。

### 打开模型分配设置

`Settings` → `Tools` → `AI Assistant` → `Models Assignment`

### 推荐配置

| 功能 | 推荐模型 | 说明 |
|------|---------|------|
| Core Features | `gpt-5.2-instant` | AI Chat、代码解释 |
| Lightweight Features | `gpt-4.1-mini` | 快速补全、简单任务 |
| Code Completion | `gpt-4.1-mini` | 代码自动补全 |

### 配置示例

```
Core Model: gpt-5.2-instant
Lightweight Model: gpt-4.1-mini
Code Completion Model: gpt-4.1-mini
Context Window Size: 128000
```

---

## 核心功能

### 1. AI Chat

按 `Alt+Enter` 或点击工具栏 AI 图标打开对话：

```
你：解释这段代码的作用
AI：这段代码实现了...

你：如何优化这个函数？
AI：可以从以下几个方面优化...
```

### 2. 代码补全

在编辑器中输入代码时，AI 会自动提供补全建议：

```java
public class UserService {
    // 输入注释，AI 自动补全实现
    // 创建用户
    |  // 按 Tab 接受建议
}
```

### 3. 代码生成

选中代码或在编辑器中右键：

- `AI Actions` → `Generate Code`
- `AI Actions` → `Generate Tests`
- `AI Actions` → `Generate Documentation`

### 4. 代码解释

选中代码后：

- 右键 → `AI Actions` → `Explain Code`
- 或使用快捷键 `Alt+Enter` → `Explain with AI`

### 5. 重构建议

选中代码后：

- `AI Actions` → `Suggest Refactoring`
- AI 会分析代码并提供重构建议

### 6. Commit Message 生成

在 Git 提交窗口：

- 点击 AI 图标
- 自动生成有意义的 commit message

---

## 快捷键

| 功能 | 快捷键 |
|------|--------|
| 打开 AI Chat | `Alt+\` |
| AI 操作菜单 | `Alt+Enter` |
| 接受补全 | `Tab` |
| 拒绝补全 | `Esc` |
| 下一个建议 | `Alt+]` |
| 上一个建议 | `Alt+[` |

---

## 高级配置

### 上下文配置

`Settings` → `Tools` → `AI Assistant` → `Context`

```
Max Context Size: 128000
Include Project Files: ✅
Include Open Files: ✅
Include Git History: ✅
```

### 隐私设置

`Settings` → `Tools` → `AI Assistant` → `Privacy`

```
Send Code to AI: ✅
Send File Names: ✅
Telemetry: ❌
```

---

## 使用 Junie Agent

Junie 是 JetBrains 的 AI Agent，支持自主完成复杂任务。

### 启用 Junie

1. 确保已配置 BYOK
2. 在 AI Chat 中输入复杂任务
3. Junie 会自动规划和执行

### 示例任务

```
创建一个用户认证模块，包括：
- 用户注册接口
- 登录接口
- JWT 令牌验证
- 密码加密
```

Junie 会：
1. 分析需求
2. 创建文件结构
3. 生成代码
4. 添加测试

---

## 常见问题

### Q: 连接测试失败？

1. 检查 API Key 是否正确
2. 检查 URL 是否正确（需要 `/v1` 后缀）
3. 检查网络连接

```bash
# 测试 API
curl https://api.bigaipro.com/v1/models \
  -H "Authorization: Bearer sk-xxxxxxxx"
```

### Q: 模型列表为空？

1. 确保连接测试成功
2. 点击刷新按钮
3. 重启 IDE

### Q: 代码补全不工作？

1. 检查是否启用了代码补全
2. 检查模型分配设置
3. 确保文件类型受支持

### Q: 如何使用代理？

在 IDE 设置中配置代理：

`Settings` → `Appearance & Behavior` → `System Settings` → `HTTP Proxy`

---

## 推荐配置

### 日常开发

```
Core Model: gpt-5.2-instant
Lightweight Model: gpt-4.1-mini
Code Completion: gpt-4.1-mini
```

### 复杂项目

```
Core Model: claude-sonnet-4.5
Lightweight Model: gpt-4.1
Code Completion: gpt-4.1
```

### 性价比优先

```
Core Model: deepseek-chat
Lightweight Model: deepseek-chat
Code Completion: gpt-4.1-nano
```

---

## 下一步

- 📖 [Cursor IDE 配置](./cursor-ide.md) - AI 原生编辑器
- 📖 [Continue 插件](./continue-plugin.md) - 开源 AI 插件
- 📖 [模型选择指南](../models/overview.md) - 选择合适的模型

---

<div align="center">

**JetBrains IDE - 专业开发者的 AI 助手**

</div>
