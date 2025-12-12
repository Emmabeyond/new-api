以下是 Windows、macOS 和 Linux 系统下设置 `ANTHROPIC_BASE_URL` 和 `ANTHROPIC_AUTH_TOKEN` 环境变量的详细方法：

------

### **Windows 系统**

#### 方法1：配置settings.json

- 创建 `~/.claude/settings.json` 文件，内容如下：

  bash

  ```bash
  {
    "env": {
      "ANTHROPIC_AUTH_TOKEN": "替换为您的API Key",
      "ANTHROPIC_BASE_URL": "https://api.bigaipro.com",
      "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": 1
    },
    "permissions": {
      "allow": [],
      "deny": []
    }
  }
  ```

  复制

- vscode中插件使用，创建文件 ~/.claude/config.json

  bash

  ```bash
  {
      "primaryApiKey": "fox"
  }
  ```

  复制

#### 方法2：临时设置（仅当前终端有效）

- 在 **PowerShell** 或 **CMD** 中执行：

  powershell

  ```powershell
  # PowerShell
  $env:ANTHROPIC_BASE_URL="https://api.bigaipro.com"
  $env:ANTHROPIC_AUTH_TOKEN="替换为您的API Key"

  # CMD
  set ANTHROPIC_BASE_URL=https://api.bigaipro.com
  set ANTHROPIC_AUTH_TOKEN=替换为您的API Key
  ```

  复制

#### 方法3：永久设置（全局生效）

1. **图形界面**：

   - 右键「此电脑」→「属性」→「高级系统设置」→「环境变量」
   - 在「用户变量」或「系统变量」中新建：
     - 变量名：`ANTHROPIC_BASE_URL`
     - 变量值：https://api.bigaipro.com
   - 同样方法添加 `ANTHROPIC_AUTH_TOKEN`

2. **PowerShell 永久设置**：

   powershell

   ```powershell
   [System.Environment]::SetEnvironmentVariable('ANTHROPIC_BASE_URL', 'https://api.bigaipro.com', 'User')
   [System.Environment]::SetEnvironmentVariable('ANTHROPIC_AUTH_TOKEN', '替换为您的API Key', 'User')
   ```

   复制

   - 重启终端后生效。

------

### **macOS 系统**

#### 方法1：配置settings.json

- 创建 `~/.claude/settings.json` 文件，内容如下：

  bash

  ```bash
  {
    "env": {
      "ANTHROPIC_AUTH_TOKEN": "替换为您的API Key",
      "ANTHROPIC_BASE_URL": "https://api.bigaipro.com",
      "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": 1
    },
    "permissions": {
      "allow": [],
      "deny": []
    }
  }
  ```

  复制

- vscode中插件使用，创建文件 ~/.claude/config.json

  bash

  ```bash
  {
      "primaryApiKey": "fox"
  }
  ```

  复制

#### 方法2：临时设置（仅当前终端有效）

- 在 **终端** 中执行：

  bash

  ```bash
  export ANTHROPIC_BASE_URL="https://api.bigaipro.com"
  export ANTHROPIC_AUTH_TOKEN="替换为您的API Key"
  ```

  复制

#### 方法3：永久设置

1. 编辑 shell 配置文件（根据使用的 shell 选择）：

   bash

   ```bash
   # 如果是 bash（默认）
   echo 'export ANTHROPIC_BASE_URL="https://api.bigaipro.com"' >> ~/.bash_profile
   echo 'export ANTHROPIC_AUTH_TOKEN="替换为您的API Key"' >> ~/.bash_profile

   # 如果是 zsh
   echo 'export ANTHROPIC_BASE_URL="https://api.bigaipro.com"' >> ~/.zshrc
   echo 'export ANTHROPIC_AUTH_TOKEN="替换为您的API Key"' >> ~/.zshrc
   ```

   复制

2. 立即生效：

   bash

   ```bash
   source ~/.bash_profile  # 或 source ~/.zshrc
   ```

   复制

------

### **Linux 系统**

#### 方法1：配置settings.json

- 创建 `~/.claude/settings.json` 文件，内容如下：

  bash

  ```bash
  {
    "env": {
      "ANTHROPIC_AUTH_TOKEN": "替换为您的API Key",
      "ANTHROPIC_BASE_URL": "https://api.bigaipro.com",
      "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": 1
    },
    "permissions": {
      "allow": [],
      "deny": []
    }
  }
  ```

  复制

#### 方法2：临时设置（仅当前终端有效）

- 在 **终端** 中执行：

  bash

  ```bash
  export ANTHROPIC_BASE_URL="https://api.bigaipro.com"
  export ANTHROPIC_AUTH_TOKEN="替换为您的API Key"
  ```

  复制

#### 方法3：永久设置

1. 编辑 shell 配置文件（根据使用的 shell 选择）：

   bash

   ```bash
   # 如果是 bash
   echo 'export ANTHROPIC_BASE_URL="https://api.bigaipro.com"' >> ~/.bashrc
   echo 'export ANTHROPIC_AUTH_TOKEN="替换为您的API Key"' >> ~/.bashrc

   # 如果是 zsh
   echo 'export ANTHROPIC_BASE_URL="https://api.bigaipro.com"' >> ~/.zshrc
   echo 'export ANTHROPIC_AUTH_TOKEN="替换为您的API Key"' >> ~/.zshrc
   ```

   复制

2. 立即生效：

   bash

   ```bash
   source ~/.bashrc  # 或 source ~/.zshrc
   ```

   复制

------

### **通用验证方法**

在所有系统中，可以通过以下命令验证是否设置成功：

bash

```bash
# macOS/Linux
echo $ANTHROPIC_BASE_URL
echo $ANTHROPIC_AUTH_TOKEN

# Windows PowerShell
echo $env:ANTHROPIC_BASE_URL
echo $env:ANTHROPIC_AUTH_TOKEN

# Windows CMD
echo %ANTHROPIC_BASE_URL%
echo %ANTHROPIC_AUTH_TOKEN%
```
