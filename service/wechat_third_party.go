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

// qrCodeDirectResponse 平台直接返回的二维码响应（不包装在 success/data 中）
type qrCodeDirectResponse struct {
	SessionId string `json:"sessionId"`
	QRCodeUrl string `json:"qrCodeUrl"`
	ExpiresAt string `json:"expiresAt"`
	ExpiresIn int    `json:"expiresIn"`
	// 错误情况
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
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

	// 平台直接返回数据对象，不包装在 success/data 中
	var directResp qrCodeDirectResponse
	if err := json.Unmarshal(respBody, &directResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查是否有错误
	if directResp.Error != "" || directResp.Message != "" {
		errMsg := directResp.Message
		if errMsg == "" {
			errMsg = directResp.Error
		}
		return nil, fmt.Errorf("生成二维码失败: %s", errMsg)
	}

	// 检查必要字段
	if directResp.SessionId == "" || directResp.QRCodeUrl == "" {
		return nil, fmt.Errorf("响应数据不完整")
	}

	// 解析过期时间
	expiresAt, err := time.Parse(time.RFC3339, directResp.ExpiresAt)
	if err != nil {
		// 如果解析失败，使用 ExpiresIn 计算
		expiresAt = time.Now().Add(time.Duration(directResp.ExpiresIn) * time.Second)
	}

	return &QRCodeResponse{
		SessionId: directResp.SessionId,
		QRCodeUrl: directResp.QRCodeUrl,
		ExpiresAt: expiresAt,
	}, nil
}

// statusDirectResponse 平台直接返回的状态响应
type statusDirectResponse struct {
	SessionId string `json:"sessionId"`
	Status    string `json:"status"` // pending, scanned, confirmed, expired
	State     string `json:"state,omitempty"`
	ExpiresAt string `json:"expiresAt,omitempty"`
	// 用户信息（confirmed 时返回）
	UserInfo *UserInfo `json:"userInfo,omitempty"`
	// 错误情况
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
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

	// 平台直接返回数据对象
	var directResp statusDirectResponse
	if err := json.Unmarshal(respBody, &directResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查是否有错误
	if directResp.Error != "" || directResp.Message != "" {
		errMsg := directResp.Message
		if errMsg == "" {
			errMsg = directResp.Error
		}
		return nil, fmt.Errorf("查询状态失败: %s", errMsg)
	}

	// 解析过期时间
	var expiresAt time.Time
	if directResp.ExpiresAt != "" {
		expiresAt, _ = time.Parse(time.RFC3339, directResp.ExpiresAt)
	}

	return &LoginStatusResponse{
		SessionId: directResp.SessionId,
		Status:    directResp.Status,
		UserInfo:  directResp.UserInfo,
		State:     directResp.State,
		ExpiresAt: expiresAt,
	}, nil
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


