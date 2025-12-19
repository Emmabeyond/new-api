# 流式响应处理

<div align="center">

**SSE 流式输出**

</div>

---

## 概述

流式响应允许在模型生成内容时实时接收输出，提供更好的用户体验。

## Python 实现

```python
from openai import OpenAI

client = OpenAI(
    api_key="sk-xxxxxxxx",
    base_url="https://api.bigaipro.com/v1"
)

stream = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[{"role": "user", "content": "写一篇文章"}],
    stream=True
)

for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="", flush=True)
```

## Node.js 实现

```javascript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: 'sk-xxxxxxxx',
  baseURL: 'https://api.bigaipro.com/v1'
});

const stream = await client.chat.completions.create({
  model: 'gpt-5.2-instant',
  messages: [{ role: 'user', content: '写一篇文章' }],
  stream: true
});

for await (const chunk of stream) {
  const content = chunk.choices[0]?.delta?.content || '';
  process.stdout.write(content);
}
```

## cURL 实现

```bash
curl https://api.bigaipro.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-xxxxxxxx" \
  -d '{
    "model": "gpt-5.2-instant",
    "messages": [{"role": "user", "content": "Hello"}],
    "stream": true
  }' \
  --no-buffer
```

## 前端实现

```javascript
const response = await fetch('https://api.bigaipro.com/v1/chat/completions', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer sk-xxxxxxxx'
  },
  body: JSON.stringify({
    model: 'gpt-5.2-instant',
    messages: [{ role: 'user', content: 'Hello' }],
    stream: true
  })
});

const reader = response.body.getReader();
const decoder = new TextDecoder();

while (true) {
  const { done, value } = await reader.read();
  if (done) break;
  
  const chunk = decoder.decode(value);
  const lines = chunk.split('\n').filter(line => line.startsWith('data: '));
  
  for (const line of lines) {
    const data = line.slice(6);
    if (data === '[DONE]') continue;
    
    const json = JSON.parse(data);
    const content = json.choices[0]?.delta?.content || '';
    console.log(content);
  }
}
```

## SSE 数据格式

```
data: {"id":"chatcmpl-xxx","object":"chat.completion.chunk","created":1234567890,"model":"gpt-5.2-instant","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}

data: {"id":"chatcmpl-xxx","object":"chat.completion.chunk","created":1234567890,"model":"gpt-5.2-instant","choices":[{"index":0,"delta":{"content":" World"},"finish_reason":null}]}

data: {"id":"chatcmpl-xxx","object":"chat.completion.chunk","created":1234567890,"model":"gpt-5.2-instant","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}

data: [DONE]
```

## 错误处理

```python
try:
    stream = client.chat.completions.create(
        model="gpt-5.2-instant",
        messages=[{"role": "user", "content": "Hello"}],
        stream=True
    )
    
    for chunk in stream:
        if chunk.choices[0].delta.content:
            print(chunk.choices[0].delta.content, end="")
            
except Exception as e:
    print(f"Stream error: {e}")
```

## 最佳实践

1. **超时设置**: 流式请求建议设置较长的超时时间
2. **断线重连**: 实现断线重连机制
3. **缓冲处理**: 适当缓冲输出，避免频繁 DOM 更新
4. **取消请求**: 提供取消流式请求的能力
