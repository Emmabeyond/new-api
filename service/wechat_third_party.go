package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/setting/system_setting"
)

// WeChatThirdPartyClient 微信第三方认证平台客户端
type WeChatThirdPartyClient struct {
	baseURL      string
	clientKey    string
	clientSecret string
	httpClient   *http.Client
}

// QRCodeResponse 生成二维码响应
type QRCodeResponse struct {
	SessionId string    `json:"sessionId"`
	QRCodeUrl string    `json:"qrCodeUrl"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// LoginStatusResponse 登录状态响应
type LoginStatusResponse struct {
	SessionId string    `json:"sessionId"`
	Status    string    `json:"status"` // pending, scanned, confirmed, expired
	UserInfo  *UserInfo `json:"userInfo,omitempty"`
	State     string    `json:"state,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}

// UserInfo 微信用户信息
type UserInfo struct {
	OpenId   string `json:"openid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// apiResponse 平台 API 通用响应结构
type apiResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Message string          `json:"message,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// generateRequest 生成二维码请求体
type generateRequest struct {
	CallbackUrl string `json:"callbackUrl,omitempty"`
	State       string `json:"state,omitempty"`
}

// NewWeChatThirdPartyClient 创建微信第三方认证客户端
func NewWeChatThirdPartyClient() (*WeChatThirdPartyClient, error) {
	settings := system_setting.GetWeChatThirdPartySettings()
	if !settings.Enabled {
		return nil, errors.New("微信第三方登录未启用")
	}

	if settings.ClientKey == "" || settings.ClientSecret == "" {
		return nil, errors.New("微信第三方登录配置不完整")
	}

	baseURL := settings.BaseURL
	if baseURL == "" {
		baseURL = "https://www.abu117.cn"
	}

	return &WeChatThirdPartyClient{
		baseURL:      baseURL,
		clientKey:    settings.ClientKey,
		clientSecret: settings.ClientSecret,
		httpClient:   GetHttpClient(),
	}, nil
}


// buildAuthToken 构建 Bearer Token
// 格式: Base64(clientKey:clientSecret)
func (c *WeChatThirdPartyClient) buildAuthToken() string {
	credentials := c.clientKey + ":" + c.clientSecret
	return base64.StdEncoding.EncodeToString([]byte(credentials))
}

// GenerateQRCode 生成登录二维码
// callbackUrl: 登录成功后的回调地址（可选）
// state: 自定义状态参数，会在回调时原样返回（可选）
func (c *WeChatThirdPartyClient) GenerateQRCode(callbackUrl, state string) (*QRCodeResponse, error) {
	apiUrl := c.baseURL + "/api/wechat-auth/v1/login/generate"

	// 获取服务器地址作为来源域名
	serverAddress := system_setting.ServerAddress
	var originDomain string
	if serverAddress != "" {
		if parsed, err := url.Parse(serverAddress); err == nil && parsed.Host != "" {
			originDomain = parsed.Scheme + "://" + parsed.Host
		} else {
			originDomain = serverAddress
		}
	}

	// 构建请求体
	reqBody := generateRequest{
		CallbackUrl: callbackUrl,
		State:       state,
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+c.buildAuthToken())
	req.Header.Set("Content-Type", "application/json")
	
	// 设置来源域名（abu117.cn 平台需要验证）
	if originDomain != "" {
		req.Header.Set("Origin", originDomain)
		req.Header.Set("Referer", originDomain+"/")
	}

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	var apiResp apiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查响应状态
	if !apiResp.Success {
		errMsg := apiResp.Message
		if errMsg == "" {
			errMsg = apiResp.Error
		}
		if errMsg == "" {
			errMsg = "未知错误"
		}
		return nil, fmt.Errorf("生成二维码失败: %s", errMsg)
	}

	// 解析数据
	var qrResp QRCodeResponse
	if err := json.Unmarshal(apiResp.Data, &qrResp); err != nil {
		return nil, fmt.Errorf("解析二维码数据失败: %w", err)
	}

	return &qrResp, nil
}

// GetLoginStatus 查询登录状态
// sessionId: 会话 ID
func (c *WeChatThirdPartyClient) GetLoginStatus(sessionId string) (*LoginStatusResponse, error) {
	apiUrl := c.baseURL + "/api/wechat-auth/v1/login/status/" + sessionId

	// 获取服务器地址作为来源域名
	serverAddress := system_setting.ServerAddress
	var originDomain string
	if serverAddress != "" {
		if parsed, err := url.Parse(serverAddress); err == nil && parsed.Host != "" {
			originDomain = parsed.Scheme + "://" + parsed.Host
		} else {
			originDomain = serverAddress
		}
	}

	// 创建请求
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+c.buildAuthToken())
	
	// 设置来源域名（abu117.cn 平台需要验证）
	if originDomain != "" {
		req.Header.Set("Origin", originDomain)
		req.Header.Set("Referer", originDomain+"/")
	}

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	var apiResp apiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查响应状态
	if !apiResp.Success {
		errMsg := apiResp.Message
		if errMsg == "" {
			errMsg = apiResp.Error
		}
		if errMsg == "" {
			errMsg = "未知错误"
		}
		return nil, fmt.Errorf("查询状态失败: %s", errMsg)
	}

	// 解析数据
	var statusResp LoginStatusResponse
	if err := json.Unmarshal(apiResp.Data, &statusResp); err != nil {
		return nil, fmt.Errorf("解析状态数据失败: %w", err)
	}

	return &statusResp, nil
}

// IsValidStatus 检查状态是否有效
func IsValidStatus(status string) bool {
	switch status {
	case "pending", "scanned", "confirmed", "expired":
		return true
	default:
		return false
	}
}
