# 多模态使用指南

<div align="center">

**图像、音频、视频处理**

</div>

---

## 概述

BigAI Pro 支持多种多模态模型，可以处理图像、音频和视频输入。

## 支持的模型

| 模型 | 图像 | 音频 | 视频 |
|------|------|------|------|
| gpt-5.2 | ✅ | ✅ | ✅ |
| gpt-4.1 | ✅ | ✅ | ✅ |
| claude-sonnet-4.5 | ✅ | ✅ | ❌ |
| gemini-3.0-pro | ✅ | ✅ | ✅ |

## 图像输入

### URL 方式

```python
from openai import OpenAI

client = OpenAI(
    api_key="sk-xxxxxxxx",
    base_url="https://api.bigaipro.com/v1"
)

response = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[
        {
            "role": "user",
            "content": [
                {"type": "text", "text": "描述这张图片"},
                {
                    "type": "image_url",
                    "image_url": {
                        "url": "https://example.com/image.jpg"
                    }
                }
            ]
        }
    ]
)
```

### Base64 方式

```python
import base64

with open("image.png", "rb") as f:
    image_data = base64.b64encode(f.read()).decode()

response = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[
        {
            "role": "user",
            "content": [
                {"type": "text", "text": "分析这张图片"},
                {
                    "type": "image_url",
                    "image_url": {
                        "url": f"data:image/png;base64,{image_data}"
                    }
                }
            ]
        }
    ]
)
```

### 多图输入

```python
response = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[
        {
            "role": "user",
            "content": [
                {"type": "text", "text": "比较这两张图片"},
                {"type": "image_url", "image_url": {"url": "https://example.com/1.jpg"}},
                {"type": "image_url", "image_url": {"url": "https://example.com/2.jpg"}}
            ]
        }
    ]
)
```

## 音频输入

```python
response = client.chat.completions.create(
    model="gpt-5.2-instant",
    messages=[
        {
            "role": "user",
            "content": [
                {"type": "text", "text": "转录这段音频"},
                {
                    "type": "audio_url",
                    "audio_url": {
                        "url": "https://example.com/audio.mp3"
                    }
                }
            ]
        }
    ]
)
```

## 视频输入

```python
response = client.chat.completions.create(
    model="gemini-3.0-pro",
    messages=[
        {
            "role": "user",
            "content": [
                {"type": "text", "text": "总结这个视频"},
                {
                    "type": "video_url",
                    "video_url": {
                        "url": "https://example.com/video.mp4"
                    }
                }
            ]
        }
    ]
)
```

## 图像生成

使用 DALL-E 3 生成图像：

```python
response = client.images.generate(
    model="dall-e-3",
    prompt="一只可爱的机器猫在星空下",
    size="1024x1024",
    quality="hd",
    n=1
)

image_url = response.data[0].url
```

## 最佳实践

1. **图像大小**: 建议不超过 20MB
2. **图像格式**: 支持 PNG、JPEG、GIF、WebP
3. **Base64**: 大图建议使用 URL 方式
4. **视频长度**: 建议不超过 10 分钟
