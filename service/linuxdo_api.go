package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/setting"
)

// LinuxDoAPIResponse LINUX DO Credit API 响应结构
type LinuxDoAPIResponse struct {
	Code      int    `json:"code"`       // 1=成功，其他=失败（查询接口使用，支付接口无此字段）
	Msg       string `json:"msg"`        // 错误信息（查询接口使用）
	TradeNo   string `json:"trade_no"`   // 订单号（查询接口使用）
	OutTradeNo string `json:"out_trade_no"` // 业务单号（查询接口使用）
	PayURL    string `json:"payurl"`     // 支付跳转URL（兼容性字段）
	QRCode    string `json:"qrcode"`     // 二维码链接（兼容性字段）
	URLScheme string `json:"urlscheme"`  // 小程序跳转URL（兼容性字段）
	ErrorMsg  string `json:"error_msg"`  // 错误信息（支付接口使用）
	Status    int    `json:"status"`     // 支付状态 1=成功 0=处理中（查询接口使用）
}

// LinuxDoAPIRequest LINUX DO Credit API 请求参数
// 参考: https://credit.linux.do/docs/api
type LinuxDoAPIRequest struct {
	PID        string // 商户ID (Client ID)
	Type       string // 支付方式: epay (固定值)
	OutTradeNo string // 商户订单号
	NotifyURL  string // 异步通知地址 (仅参与签名)
	ReturnURL  string // 跳转通知地址 (仅参与签名)
	Name       string // 商品名称，最多64字符
	Money      string // 积分数量，最多2位小数
	Device     string // 终端标识，可选
}

// CalculateLinuxDoSign 计算 LINUX DO Credit 签名
// 签名规则与易支付相同：将参数按key排序，拼接成 key=value& 格式，去掉最后的&，加上商户密钥，MD5加密
func CalculateLinuxDoSign(params map[string]string, key string) string {
	// 获取所有key并排序
	keys := make([]string, 0, len(params))
	for k := range params {
		// 排除 sign 和 sign_type，以及空值
		if k != "sign" && k != "sign_type" && params[k] != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 拼接参数
	var signStr strings.Builder
	for i, k := range keys {
		if i > 0 {
			signStr.WriteString("&")
		}
		signStr.WriteString(k)
		signStr.WriteString("=")
		signStr.WriteString(params[k])
	}

	// 加上商户密钥
	signStr.WriteString(key)

	// MD5加密
	hash := md5.Sum([]byte(signStr.String()))
	return hex.EncodeToString(hash[:])
}

// VerifyLinuxDoSign 验证 LINUX DO Credit 回调签名
func VerifyLinuxDoSign(params map[string]string, key string) bool {
	sign, ok := params["sign"]
	if !ok || sign == "" {
		return false
	}
	expectedSign := CalculateLinuxDoSign(params, key)
	return strings.EqualFold(sign, expectedSign)
}

// CallLinuxDoAPI 调用 LINUX DO Credit submit.php 接口获取支付链接
func CallLinuxDoAPI(req *LinuxDoAPIRequest) (*LinuxDoAPIResponse, error) {
	// 检查配置
	if setting.LinuxDoPayAddress == "" || setting.LinuxDoClientId == "" || setting.LinuxDoClientSecret == "" {
		return nil, fmt.Errorf("LINUX DO Credit 配置不完整")
	}

	// 构建请求参数 - 仅包含 API 文档中定义的字段
	// 参考: https://credit.linux.do/docs/api
	// 必需: pid, type, name, money, sign, sign_type
	// 可选: out_trade_no, notify_url, return_url, device
	// 注意: clientip 不在文档中，不应包含
	params := map[string]string{
		"pid":          setting.LinuxDoClientId,
		"type":         req.Type,
		"out_trade_no": req.OutTradeNo,
		"notify_url":   req.NotifyURL,
		"name":         req.Name,
		"money":        req.Money,
	}

	// 可选参数
	if req.ReturnURL != "" {
		params["return_url"] = req.ReturnURL
	}
	// device 是可选参数，仅在非空时添加
	if req.Device != "" {
		params["device"] = req.Device
	}

	// 计算签名
	sign := CalculateLinuxDoSign(params, setting.LinuxDoClientSecret)
	params["sign"] = sign
	params["sign_type"] = "MD5"

	// 构建 API URL
	// LINUX DO Credit 兼容易支付协议，支付接口为 /pay/submit.php
	apiURL := strings.TrimSuffix(setting.LinuxDoPayAddress, "/") + "/pay/submit.php"

	// 构建 POST 数据
	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}

	// 发送请求
	// 禁止自动跟随重定向，因为我们需要获取 Location 响应头中的支付链接
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 返回一个错误以禁止跟随重定向
			return http.ErrUseLastResponse
		},
	}

	// 创建 POST 请求而不是使用 PostForm，以便设置自定义 headers
	httpReq, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Printf("[DEBUG] LINUX DO Credit API 创建请求失败 - Error: %v\n", err)
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置必要的 HTTP 头，避免被 Cloudflare 拦截
	// Cloudflare 会检查多个请求头来验证请求是否来自真实浏览器
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	httpReq.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	httpReq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	httpReq.Header.Set("Accept-Encoding", "gzip, deflate, br")
	httpReq.Header.Set("Connection", "keep-alive")
	httpReq.Header.Set("Upgrade-Insecure-Requests", "1")
	httpReq.Header.Set("Sec-Fetch-Dest", "document")
	httpReq.Header.Set("Sec-Fetch-Mode", "navigate")
	httpReq.Header.Set("Sec-Fetch-Site", "same-origin")
	httpReq.Header.Set("Referer", strings.TrimSuffix(setting.LinuxDoPayAddress, "/") + "/")
	httpReq.Header.Set("DNT", "1")
	// 某些 Cloudflare 配置会检查这个头
	httpReq.Header.Set("X-Requested-With", "")

	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Printf("[DEBUG] LINUX DO Credit API 请求失败 - URL: %s, Error: %v\n", apiURL, err)
		return nil, fmt.Errorf("请求 LINUX DO Credit API 失败: %v", err)
	}
	defer resp.Body.Close()

	// 记录请求参数
	fmt.Printf("[DEBUG] LINUX DO Credit API 请求参数 - PID: %s, OutTradeNo: %s, Money: %s\n", 
		params["pid"], params["out_trade_no"], params["money"])

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[DEBUG] 读取 LINUX DO Credit 响应失败 - Error: %v\n", err)
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 成功响应：HTTP 302/303 重定向，Location 头中包含支付链接
	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusSeeOther {
		location := resp.Header.Get("Location")
		if location == "" {
			fmt.Printf("[DEBUG] LINUX DO Credit 重定向但无 Location 头 - StatusCode: %d\n", resp.StatusCode)
			return nil, fmt.Errorf("支付链接重定向失败：无 Location 响应头")
		}

		// 从 Location URL 中提取 order_no 参数作为 trade_no
		// 格式: https://credit.linux.do/paying?order_no=XXX
		parsedURL, err := url.Parse(location)
		if err == nil {
			orderNo := parsedURL.Query().Get("order_no")
			if orderNo != "" {
				fmt.Printf("[DEBUG] LINUX DO Credit API 响应成功 - PayURL: %s, OrderNo: %s\n", location, orderNo)
				return &LinuxDoAPIResponse{
					Code:   1,
					PayURL: location,
				}, nil
			}
		}

		fmt.Printf("[DEBUG] LINUX DO Credit API 响应成功 - PayURL: %s\n", location)
		return &LinuxDoAPIResponse{
			Code:   1,
			PayURL: location,
		}, nil
	}

	// 失败响应：HTTP 200 返回 JSON 错误信息
	if resp.StatusCode == http.StatusOK {
		var apiResp LinuxDoAPIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			fmt.Printf("[DEBUG] LINUX DO Credit 响应解析失败 - Body: %s, Error: %v\n", string(body), err)
			return nil, fmt.Errorf("解析响应失败: %v, body: %s", err, string(body))
		}

		// 检查是否为错误响应
		if apiResp.ErrorMsg != "" {
			fmt.Printf("[DEBUG] LINUX DO Credit 返回错误 - ErrorMsg: %s\n", apiResp.ErrorMsg)
			return &apiResp, fmt.Errorf("LINUX DO Credit 返回错误: %s", apiResp.ErrorMsg)
		}

		fmt.Printf("[DEBUG] LINUX DO Credit API 响应 - Code: %d, Msg: %s\n", apiResp.Code, apiResp.Msg)
		return &apiResp, nil
	}

	// 其他 HTTP 状态码
	fmt.Printf("[DEBUG] LINUX DO Credit API 返回异常状态码 - StatusCode: %d, Body: %s\n", resp.StatusCode, string(body))
	return nil, fmt.Errorf("LINUX DO Credit API 返回异常状态码: %d", resp.StatusCode)
}

// QueryLinuxDoOrder 查询 LINUX DO Credit 订单状态
func QueryLinuxDoOrder(tradeNo string) (bool, error) {
	// 检查配置
	if setting.LinuxDoPayAddress == "" || setting.LinuxDoClientId == "" || setting.LinuxDoClientSecret == "" {
		return false, fmt.Errorf("LINUX DO Credit 配置不完整")
	}

	// 构建查询 URL
	apiURL := fmt.Sprintf("%s/api.php?act=order&pid=%s&key=%s&trade_no=%s",
		strings.TrimSuffix(setting.LinuxDoPayAddress, "/"),
		setting.LinuxDoClientId,
		setting.LinuxDoClientSecret,
		tradeNo,
	)

	// 发送请求
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return false, fmt.Errorf("查询订单失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析响应
	var result struct {
		Code   int    `json:"code"`
		Status int    `json:"status"` // 1=支付成功，0=未支付
		Msg    string `json:"msg"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, fmt.Errorf("解析响应失败: %v", err)
	}

	if result.Code != 1 {
		return false, fmt.Errorf("查询失败: %s", result.Msg)
	}

	return result.Status == 1, nil
}

