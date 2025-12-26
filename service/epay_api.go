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

	"github.com/QuantumNous/new-api/setting/operation_setting"
)

// EpayAPIResponse 易支付 API 响应结构
type EpayAPIResponse struct {
	Code      int    `json:"code"`       // 1=成功，其他=失败
	Msg       string `json:"msg"`        // 错误信息
	TradeNo   string `json:"trade_no"`   // 订单号
	PayURL    string `json:"payurl"`     // 支付跳转URL
	QRCode    string `json:"qrcode"`     // 二维码链接
	URLScheme string `json:"urlscheme"`  // 小程序跳转URL
}

// EpayAPIRequest 易支付 API 请求参数
type EpayAPIRequest struct {
	PID        string // 商户ID
	Type       string // 支付方式: alipay/wxpay
	OutTradeNo string // 商户订单号
	NotifyURL  string // 异步通知地址
	ReturnURL  string // 跳转通知地址
	Name       string // 商品名称
	Money      string // 商品金额
	ClientIP   string // 用户IP地址
	Device     string // 设备类型: pc/mobile
	Param      string // 业务扩展参数
}

// CalculateEpaySign 计算易支付签名
// 签名规则：将参数按key排序，拼接成 key=value& 格式，去掉最后的&，加上商户密钥，MD5加密
func CalculateEpaySign(params map[string]string, key string) string {
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

// CallEpayAPI 调用易支付 mapi.php 接口获取支付二维码
func CallEpayAPI(req *EpayAPIRequest) (*EpayAPIResponse, error) {
	// 检查配置
	if operation_setting.PayAddress == "" || operation_setting.EpayId == "" || operation_setting.EpayKey == "" {
		return nil, fmt.Errorf("易支付配置不完整")
	}

	// 构建请求参数
	params := map[string]string{
		"pid":          operation_setting.EpayId,
		"type":         req.Type,
		"out_trade_no": req.OutTradeNo,
		"notify_url":   req.NotifyURL,
		"name":         req.Name,
		"money":        req.Money,
		"clientip":     req.ClientIP,
		"device":       req.Device,
	}

	// 可选参数
	if req.ReturnURL != "" {
		params["return_url"] = req.ReturnURL
	}
	if req.Param != "" {
		params["param"] = req.Param
	}

	// 计算签名
	sign := CalculateEpaySign(params, operation_setting.EpayKey)
	params["sign"] = sign
	params["sign_type"] = "MD5"

	// 构建 API URL
	apiURL := strings.TrimSuffix(operation_setting.PayAddress, "/") + "/mapi.php"

	// 构建 POST 数据
	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}

	// 发送请求
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.PostForm(apiURL, formData)
	if err != nil {
		return nil, fmt.Errorf("请求易支付API失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析响应
	var apiResp EpayAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v, body: %s", err, string(body))
	}

	// 检查返回状态
	if apiResp.Code != 1 {
		return &apiResp, fmt.Errorf("易支付返回错误: %s", apiResp.Msg)
	}

	return &apiResp, nil
}

// QueryEpayOrder 查询易支付订单状态
func QueryEpayOrder(tradeNo string) (bool, error) {
	// 检查配置
	if operation_setting.PayAddress == "" || operation_setting.EpayId == "" || operation_setting.EpayKey == "" {
		return false, fmt.Errorf("易支付配置不完整")
	}

	// 构建查询 URL
	apiURL := fmt.Sprintf("%s/api.php?act=order&pid=%s&key=%s&trade_no=%s",
		strings.TrimSuffix(operation_setting.PayAddress, "/"),
		operation_setting.EpayId,
		operation_setting.EpayKey,
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
