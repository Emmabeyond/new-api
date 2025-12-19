# Cherry Studio 配置指南

<div align="center">

**全能 AI 桌面客户端**

*支持 macOS / Windows / Linux*

</div>

---

## 简介

Cherry Studio 是一款功能强大的 AI 桌面客户端，支持多种 AI 模型和服务商，提供对话、知识库、Agent 等丰富功能。

---

## 安装

### macOS

```bash
# Homebrew
brew install --cask cherry-studio

# 或下载 DMG
# https://cherry-ai.com/download
```

### Windows

1. 访问 [Cherry Studio 下载页面](https://cherry-ai.com/download)
2. 下载 Windows 安装包
3. 运行安装程序

### Linux

```bash
# AppImage
chmod +x Cherry-Studio.AppImage
./Cherry-Studio.AppImage

# 或使用 Flatpak
flatpak install flathub com.cherry.studio
```

---

## 配置 BigAI Pro API

### 步骤一：打开设置

1. 启动 Cherry Studio
2. 点击左下角设置图标
3. 选择「模型服务」

### 步骤二：添加服务商

1. 点击「添加服务商」
2. 选择「OpenAI 兼容」
3. 填写配置：

| 配置项 | 值 |
|--------|-----|
| 名称 | BigAI Pro |
| API 地址 | https://api.bigaipro.com |
| API Key | sk-xxxxxxxxxxxxxxxx |

### 步骤三：添加模型

点击「添加模型」，配置以下模型：

```
# OpenAI 系列
gpt-5.2-pro
gpt-5.2-thinking
gpt-5.2-instant
gpt-4.1
gpt-4.1-mini
o3
o4-mini

# Claude 系列
claude-sonnet-4.5
claude-opus-4.5
claude-haiku-4.5

# Gemini 系列
gemini-3.0-pro
gemini-2.5-pro
gemini-2.5-flash

# 国产模型
qwen3-235b
deepseek-r1
deepseek-chat
glm-4-plus
```

---

## 核心功能

### 1. 多模型对话

Cherry Studio 支持同时与多个模型对话：

- 创建多个对话窗口
- 每个窗口使用不同模型
- 对比不同模型的回答

### 2. 知识库

创建本地知识库，让 AI 基于你的文档回答：

1. 点击「知识库」
2. 创建新知识库
3. 上传文档（PDF、Word、TXT 等）
4. 在对话中引用知识库

### 3. Agent

创建自定义 Agent：

```json
{
  "name": "代码审查专家",
  "description": "专业的代码审查助手",
  "systemPrompt": "你是一位资深的代码审查专家...",
  "model": "claude-sonnet-4.5",
  "temperature": 0.3
}
```

### 4. 提示词库

保存和管理常用提示词：

- 代码生成模板
- 文档写作模板
- 翻译模板
- 分析模板

### 5. 多模态

支持图像输入：

- 拖拽图片到对话框
- 粘贴剪贴板图片
- 支持多图对话

---

## 高级配置

### 完整配置示例

```json
{
  "providers": [
    {
      "id": "bigai-pro",
      "name": "BigAI Pro",
      "type": "openai-compatible",
      "baseUrl": "https://api.bigaipro.com",
      "apiKey": "sk-xxxxxxxxxxxxxxxx",
      "models": [
        {
          "id": "gpt-5.2-pro",
          "name": "GPT-5.2 Pro",
          "contextLength": 1000000,
          "maxTokens": 16384
        },
        {
          "id": "claude-sonnet-4.5",
          "name": "Claude Sonnet 4.5",
          "contextLength": 200000,
          "maxTokens": 8192
        },
        {
          "id": "deepseek-r1",
          "name": "DeepSeek R1",
          "contextLength": 128000,
          "maxTokens": 8192
        }
      ]
    }
  ],
  "defaultModel": "gpt-5.2-instant",
  "theme": "dark",
  "language": "zh-CN",
  "proxy": {
    "enabled": false,
    "url": ""
  }
}
```

### 快捷键配置

| 功能 | macOS | Windows |
|------|-------|---------|
| 新对话 | `Cmd+N` | `Ctrl+N` |
| 发送消息 | `Cmd+Enter` | `Ctrl+Enter` |
| 清空对话 | `Cmd+Shift+D` | `Ctrl+Shift+D` |
| 切换模型 | `Cmd+M` | `Ctrl+M` |
| 打开设置 | `Cmd+,` | `Ctrl+,` |

---

## 使用技巧

### 1. 模型对比

同时向多个模型提问，对比回答质量：

1. 创建多个对话标签
2. 每个标签使用不同模型
3. 发送相同问题
4. 对比回答

### 2. 知识库增强

```
1. 上传项目文档到知识库
2. 在对话中引用：@知识库名称
3. AI 会基于文档内容回答
```

### 3. 批量处理

使用 Agent 批量处理任务：

```
1. 创建处理 Agent
2. 设置批量输入
3. 自动处理并导出结果
```

---

## 常见问题

### Q: 连接失败？

1. 检查 API 地址是否正确
2. 检查 API Key 是否有效
3. 检查网络连接

### Q: 如何导出对话？

点击对话右上角菜单 → 导出 → 选择格式（Markdown/JSON/TXT）

### Q: 如何同步数据？

Cherry Studio 支持云同步：

1. 设置 → 同步
2. 登录账号
3. 开启同步

---

## 推荐配置

### 日常使用

```json
{
  "defaultModel": "gpt-5.2-instant",
  "temperature": 0.7
}
```

### 代码开发

```json
{
  "defaultModel": "claude-sonnet-4.5",
  "temperature": 0.3
}
```

### 创意写作

```json
{
  "defaultModel": "gpt-5.2-pro",
  "temperature": 1.0
}
```

---

## 下一步

- 📖 [快速入门](../quick-start.md) - API 基础使用
- 📖 [模型选择指南](../models/overview.md) - 选择合适的模型
- 📖 [Claude Code 教程](./claude-code.md) - 终端编程助手

---

<div align="center">

**Cherry Studio - 你的 AI 桌面助手**

</div>
