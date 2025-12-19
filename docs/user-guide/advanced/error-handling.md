# 错误处理指南

<div align="center">

**错误码与重试策略**

</div>

---

## 常见错误码

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 400 | 请求参数错误 | 检查请求格式和参数 |
| 401 | API Key 无效 | 检查 API Key 是否正确 |
| 402 | 余额不足 | 充值账户 |
| 403 | 权限不足 | 检查 API Key 权限设置 |
| 404 | 资源不存在 | 检查模型名称是否正确 |
| 429 | 请求过于频繁 | 降低请求频率或等待 |
| 500 | 服务器错误 | 稍后重试 |
| 502 | 网关错误 | 稍后重试 |
| 503 | 服务不可用 | 稍后重试 |

## 错误响应格式

```json
{
  "error": {
    "message": "错误描述",
    "type": "invalid_request_error",
    "param": "model",
    "code": "model_not_found"
  }
}
```

## Python 错误处理

```python
from openai import OpenAI, APIError, RateLimitError, AuthenticationError
import time

client = OpenAI(
    api_key="sk-xxxxxxxx",
    base_url="https://api.bigaipro.com/v1"
)

def call_api_with_retry(messages, max_retries=3):
    for attempt in range(max_retries):
        try:
            response = client.chat.completions.create(
                model="gpt-5.2-instant",
                messages=messages
            )
            return response
            
        except AuthenticationError as e:
            # API Key 错误，不重试
            print(f"认证失败: {e}")
            raise
            
        except RateLimitError as e:
            # 频率限制，等待后重试
            wait_time = 2 ** attempt
            print(f"频率限制，等待 {wait_time} 秒...")
            time.sleep(wait_time)
            
        except APIError as e:
            # 其他 API 错误
            if e.status_code >= 500:
                # 服务器错误，重试
                wait_time = 2 ** attempt
                print(f"服务器错误，等待 {wait_time} 秒...")
                time.sleep(wait_time)
            else:
                # 客户端错误，不重试
                print(f"API 错误: {e}")
                raise
                
        except Exception as e:
            print(f"未知错误: {e}")
            raise
    
    raise Exception("重试次数已用完")
```

## Node.js 错误处理

```javascript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: 'sk-xxxxxxxx',
  baseURL: 'https://api.bigaipro.com/v1'
});

async function callApiWithRetry(messages, maxRetries = 3) {
  for (let attempt = 0; attempt < maxRetries; attempt++) {
    try {
      const response = await client.chat.completions.create({
        model: 'gpt-5.2-instant',
        messages
      });
      return response;
      
    } catch (error) {
      if (error.status === 401) {
        // 认证失败，不重试
        throw error;
      }
      
      if (error.status === 429 || error.status >= 500) {
        // 频率限制或服务器错误，重试
        const waitTime = Math.pow(2, attempt) * 1000;
        console.log(`等待 ${waitTime}ms 后重试...`);
        await new Promise(resolve => setTimeout(resolve, waitTime));
        continue;
      }
      
      throw error;
    }
  }
  
  throw new Error('重试次数已用完');
}
```

## 指数退避策略

```python
import time
import random

def exponential_backoff(attempt, base_delay=1, max_delay=60):
    """计算指数退避延迟时间"""
    delay = min(base_delay * (2 ** attempt), max_delay)
    # 添加随机抖动
    jitter = random.uniform(0, delay * 0.1)
    return delay + jitter

def call_with_backoff(func, max_retries=5):
    for attempt in range(max_retries):
        try:
            return func()
        except Exception as e:
            if attempt == max_retries - 1:
                raise
            delay = exponential_backoff(attempt)
            print(f"重试 {attempt + 1}/{max_retries}，等待 {delay:.2f}s")
            time.sleep(delay)
```

## 超时处理

```python
from openai import OpenAI

client = OpenAI(
    api_key="sk-xxxxxxxx",
    base_url="https://api.bigaipro.com/v1",
    timeout=60.0  # 60 秒超时
)

# 或者单次请求设置超时
response = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[{"role": "user", "content": "Hello"}],
    timeout=30.0
)
```

## 常见问题排查

### 401 Unauthorized

```bash
# 检查 API Key
curl https://api.bigaipro.com/v1/models \
  -H "Authorization: Bearer sk-xxxxxxxx"
```

### 429 Rate Limit

1. 检查请求频率
2. 检查 API Key 的频率限制设置
3. 实现请求队列和限流

### 500 Internal Server Error

1. 检查请求参数是否正确
2. 尝试其他模型
3. 稍后重试

## 最佳实践

1. **始终处理错误**: 不要忽略异常
2. **实现重试机制**: 使用指数退避
3. **设置超时**: 避免请求无限等待
4. **记录日志**: 记录错误信息便于排查
5. **优雅降级**: 提供备选方案
