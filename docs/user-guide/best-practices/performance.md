# 性能优化指南

<div align="center">

**提升 API 调用效率**

</div>

---

## 减少延迟

### 1. 选择合适的模型

| 场景 | 推荐模型 | 延迟 |
|------|----------|------|
| 实时对话 | gpt-5.2-instant | 低 |
| 简单任务 | gpt-4.1-mini | 低 |
| 复杂任务 | gpt-5.2-pro | 中 |
| 推理任务 | o3 | 高 |

### 2. 使用流式输出

```python
# 流式输出可以更快看到结果
stream = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[{"role": "user", "content": "Hello"}],
    stream=True
)

for chunk in stream:
    print(chunk.choices[0].delta.content or "", end="")
```

### 3. 控制输出长度

```python
response = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[{"role": "user", "content": "简短回答：什么是 AI？"}],
    max_tokens=100  # 限制输出长度
)
```

## 减少 Token 消耗

### 1. 精简 System Prompt

```python
# ❌ 冗长的 prompt
system = """你是一个非常专业的助手，你需要帮助用户解决各种问题，
你的回答应该准确、详细、有帮助，你应该考虑用户的需求..."""

# ✅ 精简的 prompt
system = "你是专业助手，回答准确简洁。"
```

### 2. 管理对话历史

```python
def trim_messages(messages, max_messages=10):
    """保留最近的消息"""
    if len(messages) <= max_messages:
        return messages
    # 保留 system message 和最近的消息
    return [messages[0]] + messages[-(max_messages-1):]
```

### 3. 使用摘要

```python
def summarize_history(messages):
    """将长对话历史压缩为摘要"""
    if len(messages) < 20:
        return messages
    
    # 获取历史摘要
    history_text = "\n".join([m["content"] for m in messages[1:-5]])
    summary = client.chat.completions.create(
        model="gpt-4.1-mini",
        messages=[
            {"role": "user", "content": f"总结以下对话：\n{history_text}"}
        ],
        max_tokens=200
    )
    
    return [
        messages[0],  # system
        {"role": "assistant", "content": f"[历史摘要] {summary.choices[0].message.content}"},
        *messages[-5:]  # 最近的消息
    ]
```

## 并发处理

### 异步请求

```python
import asyncio
from openai import AsyncOpenAI

client = AsyncOpenAI(
    api_key="sk-xxxxxxxx",
    base_url="https://api.bigaipro.com/v1"
)

async def process_item(item):
    response = await client.chat.completions.create(
        model="gpt-4.1-mini",
        messages=[{"role": "user", "content": item}]
    )
    return response.choices[0].message.content

async def process_batch(items):
    tasks = [process_item(item) for item in items]
    return await asyncio.gather(*tasks)

# 使用
results = asyncio.run(process_batch(["问题1", "问题2", "问题3"]))
```

### 请求池

```python
import asyncio
from asyncio import Semaphore

async def process_with_limit(items, max_concurrent=5):
    semaphore = Semaphore(max_concurrent)
    
    async def limited_process(item):
        async with semaphore:
            return await process_item(item)
    
    tasks = [limited_process(item) for item in items]
    return await asyncio.gather(*tasks)
```

## 缓存策略

### 简单缓存

```python
from functools import lru_cache
import hashlib

@lru_cache(maxsize=1000)
def cached_completion(prompt_hash):
    # 实际调用 API
    pass

def get_completion(prompt):
    prompt_hash = hashlib.md5(prompt.encode()).hexdigest()
    return cached_completion(prompt_hash)
```

### Redis 缓存

```python
import redis
import json
import hashlib

r = redis.Redis(host='localhost', port=6379, db=0)

def get_cached_completion(messages, ttl=3600):
    # 生成缓存键
    key = hashlib.md5(json.dumps(messages).encode()).hexdigest()
    
    # 检查缓存
    cached = r.get(key)
    if cached:
        return json.loads(cached)
    
    # 调用 API
    response = client.chat.completions.create(
        model="gpt-5.2-instant",
        messages=messages
    )
    result = response.choices[0].message.content
    
    # 存入缓存
    r.setex(key, ttl, json.dumps(result))
    
    return result
```

## 批量处理

```python
# 使用 Batch API（如果支持）
batch_requests = [
    {"model": "gpt-4.1-mini", "messages": [{"role": "user", "content": f"问题{i}"}]}
    for i in range(100)
]

# 或者自己实现批量处理
async def batch_process(requests, batch_size=10):
    results = []
    for i in range(0, len(requests), batch_size):
        batch = requests[i:i+batch_size]
        batch_results = await asyncio.gather(*[
            process_request(req) for req in batch
        ])
        results.extend(batch_results)
    return results
```

## 监控和优化

### 记录性能指标

```python
import time

def timed_completion(messages):
    start = time.time()
    
    response = client.chat.completions.create(
        model="gpt-5.2-instant",
        messages=messages
    )
    
    elapsed = time.time() - start
    tokens = response.usage.total_tokens
    
    print(f"耗时: {elapsed:.2f}s, Tokens: {tokens}")
    
    return response
```

## 最佳实践总结

1. **选择合适的模型**: 简单任务用小模型
2. **控制上下文长度**: 精简 prompt 和历史
3. **使用流式输出**: 提升用户体验
4. **实现缓存**: 避免重复请求
5. **并发处理**: 提高吞吐量
6. **监控性能**: 持续优化
