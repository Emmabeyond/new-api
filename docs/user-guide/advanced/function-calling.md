# 函数调用 (Function Calling)

<div align="center">

**让 AI 调用你的函数**

</div>

---

## 概述

Function Calling 允许模型根据用户输入决定调用哪些函数，并生成函数参数。

## 基础示例

```python
from openai import OpenAI
import json

client = OpenAI(
    api_key="sk-xxxxxxxx",
    base_url="https://api.bigaipro.com/v1"
)

# 定义工具
tools = [
    {
        "type": "function",
        "function": {
            "name": "get_weather",
            "description": "获取指定城市的天气信息",
            "parameters": {
                "type": "object",
                "properties": {
                    "city": {
                        "type": "string",
                        "description": "城市名称，如：北京、上海"
                    },
                    "unit": {
                        "type": "string",
                        "enum": ["celsius", "fahrenheit"],
                        "description": "温度单位"
                    }
                },
                "required": ["city"]
            }
        }
    }
]

# 发送请求
response = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[
        {"role": "user", "content": "北京今天天气怎么样？"}
    ],
    tools=tools,
    tool_choice="auto"
)

# 处理响应
message = response.choices[0].message

if message.tool_calls:
    for tool_call in message.tool_calls:
        function_name = tool_call.function.name
        arguments = json.loads(tool_call.function.arguments)
        print(f"调用函数: {function_name}")
        print(f"参数: {arguments}")
```

## 完整流程

```python
def get_weather(city: str, unit: str = "celsius") -> dict:
    """模拟获取天气"""
    return {
        "city": city,
        "temperature": 25,
        "unit": unit,
        "condition": "晴天"
    }

def run_conversation(user_message: str):
    messages = [{"role": "user", "content": user_message}]
    
    # 第一次调用
    response = client.chat.completions.create(
        model="gpt-5.2-instant",
        messages=messages,
        tools=tools,
        tool_choice="auto"
    )
    
    message = response.choices[0].message
    
    # 如果需要调用函数
    if message.tool_calls:
        messages.append(message)
        
        for tool_call in message.tool_calls:
            function_name = tool_call.function.name
            arguments = json.loads(tool_call.function.arguments)
            
            # 执行函数
            if function_name == "get_weather":
                result = get_weather(**arguments)
            
            # 添加函数结果
            messages.append({
                "role": "tool",
                "tool_call_id": tool_call.id,
                "content": json.dumps(result, ensure_ascii=False)
            })
        
        # 第二次调用，获取最终回复
        final_response = client.chat.completions.create(
            model="gpt-5.2-instant",
            messages=messages
        )
        
        return final_response.choices[0].message.content
    
    return message.content

# 使用
result = run_conversation("北京今天天气怎么样？")
print(result)
```

## 多函数示例

```python
tools = [
    {
        "type": "function",
        "function": {
            "name": "get_weather",
            "description": "获取天气信息",
            "parameters": {
                "type": "object",
                "properties": {
                    "city": {"type": "string", "description": "城市名称"}
                },
                "required": ["city"]
            }
        }
    },
    {
        "type": "function",
        "function": {
            "name": "search_flights",
            "description": "搜索航班",
            "parameters": {
                "type": "object",
                "properties": {
                    "from_city": {"type": "string", "description": "出发城市"},
                    "to_city": {"type": "string", "description": "到达城市"},
                    "date": {"type": "string", "description": "日期，格式：YYYY-MM-DD"}
                },
                "required": ["from_city", "to_city", "date"]
            }
        }
    },
    {
        "type": "function",
        "function": {
            "name": "book_hotel",
            "description": "预订酒店",
            "parameters": {
                "type": "object",
                "properties": {
                    "city": {"type": "string", "description": "城市"},
                    "check_in": {"type": "string", "description": "入住日期"},
                    "check_out": {"type": "string", "description": "退房日期"}
                },
                "required": ["city", "check_in", "check_out"]
            }
        }
    }
]
```

## 强制调用特定函数

```python
# 强制调用 get_weather
response = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[{"role": "user", "content": "你好"}],
    tools=tools,
    tool_choice={"type": "function", "function": {"name": "get_weather"}}
)

# 禁止调用函数
response = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[{"role": "user", "content": "北京天气"}],
    tools=tools,
    tool_choice="none"
)
```

## 并行函数调用

模型可能同时调用多个函数：

```python
response = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[
        {"role": "user", "content": "查询北京和上海的天气"}
    ],
    tools=tools
)

# 可能返回多个 tool_calls
for tool_call in response.choices[0].message.tool_calls:
    print(f"函数: {tool_call.function.name}")
    print(f"参数: {tool_call.function.arguments}")
```

## 支持的模型

| 模型 | Function Calling |
|------|-----------------|
| gpt-5.2 系列 | ✅ |
| gpt-4.1 系列 | ✅ |
| o3 / o4-mini | ✅ |
| claude-sonnet-4.5 | ✅ |
| gemini-3.0-pro | ✅ |

## 最佳实践

1. **清晰的描述**: 函数和参数描述要清晰准确
2. **参数验证**: 验证模型生成的参数
3. **错误处理**: 处理函数执行失败的情况
4. **安全考虑**: 不要让模型调用敏感函数
