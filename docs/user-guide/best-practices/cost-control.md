# 成本控制指南

<div align="center">

**优化使用成本**

</div>

---

## 成本构成

```
总成本 = (输入 Token × 输入价格) + (输出 Token × 输出价格)
```

## 模型选择策略

### 按任务选择模型

| 任务类型 | 推荐模型 | 成本 |
|----------|----------|------|
| 简单问答 | gpt-4.1-mini | 💰 |
| 日常对话 | gpt-5.2-instant | 💰💰 |
| 代码生成 | claude-sonnet-4.5 | 💰💰 |
| 复杂分析 | gpt-5.2-pro | 💰💰💰 |
| 深度推理 | o3 | 💰💰💰💰 |

### 性价比排行

| 模型 | 输入 ($/1M) | 输出 ($/1M) | 性价比 |
|------|------------|------------|--------|
| deepseek-chat | $0.14 | $0.28 | ⭐⭐⭐⭐⭐ |
| qwen3-235b | $0.60 | $1.20 | ⭐⭐⭐⭐⭐ |
| gpt-4.1-mini | $0.15 | $0.60 | ⭐⭐⭐⭐ |
| gemini-2.5-flash | $0.075 | $0.30 | ⭐⭐⭐⭐⭐ |

## 减少 Token 消耗

### 1. 精简 System Prompt

```python
# ❌ 冗长 (~200 tokens)
system = """你是一个非常专业、知识渊博的AI助手。你的任务是帮助用户解决各种问题。
在回答问题时，你应该：
1. 首先理解用户的问题
2. 然后提供准确、详细的回答
3. 如果不确定，要诚实地说明
4. 回答要有条理，使用列表和标题
5. 语言要专业但易懂
..."""

# ✅ 精简 (~20 tokens)
system = "你是专业助手。回答准确简洁，不确定时说明。"
```

### 2. 控制输出长度

```python
response = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[
        {"role": "user", "content": "用一句话解释量子计算"}
    ],
    max_tokens=100  # 限制输出
)
```

### 3. 管理对话历史

```python
def trim_history(messages, max_tokens=4000):
    """保持历史在 token 限制内"""
    total = 0
    result = [messages[0]]  # 保留 system
    
    for msg in reversed(messages[1:]):
        tokens = len(msg["content"]) // 4  # 粗略估算
        if total + tokens > max_tokens:
            break
        result.insert(1, msg)
        total += tokens
    
    return result
```

## 智能路由

### 根据任务复杂度选择模型

```python
def smart_route(task):
    # 简单任务用便宜模型
    simple_keywords = ["翻译", "总结", "格式化"]
    if any(kw in task for kw in simple_keywords):
        return "gpt-4.1-mini"
    
    # 代码任务用 Claude
    code_keywords = ["代码", "编程", "函数", "bug"]
    if any(kw in task for kw in code_keywords):
        return "claude-sonnet-4.5"
    
    # 复杂任务用强模型
    complex_keywords = ["分析", "设计", "架构"]
    if any(kw in task for kw in complex_keywords):
        return "gpt-5.2-pro"
    
    # 默认
    return "gpt-5.2-instant"
```

### 分级处理

```python
async def tiered_processing(query):
    # 第一级：快速模型
    response = await call_model("gpt-4.1-mini", query)
    
    # 检查是否需要更强模型
    if "不确定" in response or "无法" in response:
        # 第二级：强模型
        response = await call_model("gpt-5.2-pro", query)
    
    return response
```

## 缓存策略

### 结果缓存

```python
import hashlib
import redis

r = redis.Redis()

def cached_completion(messages, ttl=3600):
    # 生成缓存键
    key = hashlib.md5(str(messages).encode()).hexdigest()
    
    # 检查缓存
    cached = r.get(key)
    if cached:
        return cached.decode()
    
    # 调用 API
    response = client.chat.completions.create(
        model="gpt-5.2-instant",
        messages=messages
    )
    result = response.choices[0].message.content
    
    # 存入缓存
    r.setex(key, ttl, result)
    
    return result
```

### 语义缓存

```python
from sentence_transformers import SentenceTransformer
import numpy as np

model = SentenceTransformer('all-MiniLM-L6-v2')
cache = {}  # {embedding: response}

def semantic_cache(query, threshold=0.9):
    query_embedding = model.encode(query)
    
    for cached_embedding, response in cache.items():
        similarity = np.dot(query_embedding, cached_embedding)
        if similarity > threshold:
            return response
    
    # 缓存未命中，调用 API
    response = call_api(query)
    cache[tuple(query_embedding)] = response
    
    return response
```

## 批量处理

### 使用 Batch API

```python
# 批量请求可以获得折扣
batch_requests = [
    {"model": "gpt-4.1-mini", "messages": [{"role": "user", "content": f"问题{i}"}]}
    for i in range(100)
]

# 批量处理通常有 50% 折扣
```

## 监控和预算

### 设置预算告警

在控制台设置：
- 日预算限制
- 月预算限制
- 余额预警阈值

### 追踪使用情况

```python
def track_usage(response):
    usage = response.usage
    
    # 记录使用情况
    log_usage({
        "model": response.model,
        "input_tokens": usage.prompt_tokens,
        "output_tokens": usage.completion_tokens,
        "total_tokens": usage.total_tokens,
        "timestamp": datetime.now()
    })
```

## 成本优化检查清单

- [ ] 选择合适的模型
- [ ] 精简 System Prompt
- [ ] 控制输出长度
- [ ] 管理对话历史
- [ ] 实现结果缓存
- [ ] 使用智能路由
- [ ] 设置预算告警
- [ ] 定期审查使用情况

## 预估成本计算器

```python
def estimate_cost(model, input_tokens, output_tokens):
    prices = {
        "gpt-5.2-instant": (2.50, 10.00),
        "gpt-4.1-mini": (0.15, 0.60),
        "claude-sonnet-4.5": (3.00, 15.00),
        "deepseek-chat": (0.14, 0.28),
    }
    
    input_price, output_price = prices.get(model, (1.0, 3.0))
    
    cost = (input_tokens * input_price + output_tokens * output_price) / 1_000_000
    
    return f"${cost:.4f}"

# 使用
print(estimate_cost("gpt-5.2-instant", 1000, 500))
```
