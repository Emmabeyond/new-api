package controller

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/QuantumNous/new-api/setting/system_setting"

	"github.com/Calcium-Ion/go-epay/epay"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

const (
	PaymentMethodLinuxDo = "linuxdo"
)

// LinuxDoPayRequest LINUX DO Credit 支付请求
type LinuxDoPayRequest struct {
	Amount int64 `json:"amount"`
}

// GetLinuxDoClient 获取 LINUX DO Credit 易支付客户端
func GetLinuxDoClient() *epay.Client {
	if setting.LinuxDoPayAddress == "" || setting.LinuxDoClientId == "" || setting.LinuxDoClientSecret == "" {
		return nil
	}
	withUrl, err := epay.NewClient(&epay.Config{
		PartnerID: setting.LinuxDoClientId,
		Key:       setting.LinuxDoClientSecret,
	}, setting.LinuxDoPayAddress)
	if err != nil {
		return nil
	}
	return withUrl
}

// GetLinuxDoTopUpInfo 获取 LINUX DO Credit 充值信息
func GetLinuxDoTopUpInfo(c *gin.Context) {
	// 验证配置
	isConfigured := setting.LinuxDoClientId != "" && setting.LinuxDoClientSecret != ""
	
	data := gin.H{
		"enable_linuxdo_topup": isConfigured,
		"linuxdo_min_topup":    setting.LinuxDoMinTopUp,
	}
	common.ApiSuccess(c, data)
}

// maskSecret 掩码敏感信息（用于日志输出）
func maskSecret(secret string) string {
	if len(secret) == 0 {
		return "(empty)"
	}
	if len(secret) <= 4 {
		return "****"
	}
	return secret[:2] + "****" + secret[len(secret)-2:]
}

// getLinuxDoMinTopup 获取 LINUX DO Credit 最小充值金额
func getLinuxDoMinTopup() int64 {
	minTopup := setting.LinuxDoMinTopUp
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		dMinTopup := decimal.NewFromInt(int64(minTopup))
		dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
		minTopup = int(dMinTopup.Mul(dQuotaPerUnit).IntPart())
	}
	return int64(minTopup)
}

// getLinuxDoPayMoney 计算 LINUX DO Credit 支付金额（积分）
func getLinuxDoPayMoney(amount int64, group string) float64 {
	dAmount := decimal.NewFromInt(amount)
	// 充值金额以"展示类型"为准：
	// - USD/CNY: 前端传 amount 为金额单位；TOKENS: 前端传 tokens，需要换成 USD 金额
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
		dAmount = dAmount.Div(dQuotaPerUnit)
	}

	topupGroupRatio := common.GetTopupGroupRatio(group)
	if topupGroupRatio == 0 {
		topupGroupRatio = 1
	}

	dTopupGroupRatio := decimal.NewFromFloat(topupGroupRatio)

	// 使用 LinuxDoUnitPrice 作为充值比例
	// LinuxDoUnitPrice 表示：1美元 = 多少积分
	// 例如：LinuxDoUnitPrice = 7.3，则充值 1 美元需要支付 7.3 积分
	unitPrice := setting.LinuxDoUnitPrice
	if unitPrice <= 0 {
		unitPrice = 1.0 // 默认 1:1
	}
	dUnitPrice := decimal.NewFromFloat(unitPrice)

	// payMoney = amount * topupGroupRatio * unitPrice
	payMoney := dAmount.Mul(dTopupGroupRatio).Mul(dUnitPrice)

	return payMoney.InexactFloat64()
}

// RequestLinuxDoPay 请求 LINUX DO Credit 支付
func RequestLinuxDoPay(c *gin.Context) {
	var req LinuxDoPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "参数错误"})
		return
	}

	if req.Amount < getLinuxDoMinTopup() {
		c.JSON(200, gin.H{"message": "error", "data": fmt.Sprintf("充值数量不能小于 %d", getLinuxDoMinTopup())})
		return
	}

	id := c.GetInt("id")
	group, err := model.GetUserGroup(id, true)
	if err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "获取用户分组失败"})
		return
	}

	payMoney := getLinuxDoPayMoney(req.Amount, group)
	if payMoney < 0.01 {
		c.JSON(200, gin.H{"message": "error", "data": "充值金额过低"})
		return
	}

	callBackAddress := service.GetCallbackAddress()
	returnUrl, _ := url.Parse(system_setting.ServerAddress + "/console/log")
	notifyUrl, _ := url.Parse(callBackAddress + "/api/user/linuxdo/notify")
	tradeNo := fmt.Sprintf("%s%d", common.GetRandomString(6), time.Now().Unix())
	tradeNo = fmt.Sprintf("LDO%dNO%s", id, tradeNo)

	client := GetLinuxDoClient()
	if client == nil {
		c.JSON(200, gin.H{"message": "error", "data": "LINUX DO Credit 未配置"})
		return
	}

	uri, params, err := client.Purchase(&epay.PurchaseArgs{
		Type:           "epay", // LINUX DO Credit 使用 epay 类型
		ServiceTradeNo: tradeNo,
		Name:           fmt.Sprintf("TUC%d", req.Amount),
		Money:          strconv.FormatFloat(payMoney, 'f', 2, 64),
		Device:         epay.PC,
		NotifyUrl:      notifyUrl,
		ReturnUrl:      returnUrl,
	})
	if err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "拉起支付失败"})
		return
	}

	// 计算实际充值金额
	amount := req.Amount
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		dAmount := decimal.NewFromInt(int64(amount))
		dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
		amount = dAmount.Div(dQuotaPerUnit).IntPart()
	}

	topUp := &model.TopUp{
		UserId:        id,
		Amount:        amount,
		Money:         payMoney,
		TradeNo:       tradeNo,
		PaymentMethod: PaymentMethodLinuxDo,
		CreateTime:    time.Now().Unix(),
		Status:        "pending",
	}
	if err := topUp.Insert(); err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "创建订单失败"})
		return
	}

	c.JSON(200, gin.H{"message": "success", "data": params, "url": uri})
}

// LinuxDoNotify LINUX DO Credit 支付回调
// 回调方式: HTTP GET
// 参数: trade_no, out_trade_no, type, name, money, trade_status, sign, sign_type
func LinuxDoNotify(c *gin.Context) {
	// 从 URL Query 获取回调参数
	params := lo.Reduce(lo.Keys(c.Request.URL.Query()), func(r map[string]string, t string, i int) map[string]string {
		r[t] = c.Request.URL.Query().Get(t)
		return r
	}, map[string]string{})

	// 检查配置
	if setting.LinuxDoClientSecret == "" {
		log.Println("LINUX DO Credit 回调失败 未找到配置信息")
		c.Writer.Write([]byte("fail"))
		return
	}

	// 使用自定义签名验证（与 LINUX DO Credit API 文档一致）
	if !service.VerifyLinuxDoSign(params, setting.LinuxDoClientSecret) {
		log.Printf("LINUX DO Credit 回调签名验证失败 - params: %v", params)
		c.Writer.Write([]byte("fail"))
		return
	}

	// 签名验证通过，先返回 success
	c.Writer.Write([]byte("success"))

	// 获取回调参数
	tradeStatus := params["trade_status"]
	outTradeNo := params["out_trade_no"] // 商户订单号

	if tradeStatus == "TRADE_SUCCESS" {
		LockOrder(outTradeNo)
		defer UnlockOrder(outTradeNo)

		topUp := model.GetTopUpByTradeNo(outTradeNo)
		if topUp == nil {
			log.Printf("LINUX DO Credit 回调未找到订单: %s", outTradeNo)
			return
		}

		if topUp.Status == "pending" {
			topUp.Status = "success"
			if err := topUp.Update(); err != nil {
				log.Printf("LINUX DO Credit 回调更新订单失败: %v", topUp)
				return
			}

			dAmount := decimal.NewFromInt(int64(topUp.Amount))
			dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
			quotaToAdd := int(dAmount.Mul(dQuotaPerUnit).IntPart())

			if err := model.IncreaseUserQuota(topUp.UserId, quotaToAdd, true); err != nil {
				log.Printf("LINUX DO Credit 回调更新用户失败: %v", topUp)
				return
			}

			log.Printf("LINUX DO Credit 回调更新用户成功 %v", topUp)
			model.RecordLog(topUp.UserId, model.LogTypeTopup, fmt.Sprintf("使用 LINUX DO Credit 充值成功，充值金额: %v，支付积分：%f", logger.LogQuota(quotaToAdd), topUp.Money))
		}
	} else {
		log.Printf("LINUX DO Credit 异常回调状态: %s, params: %v", tradeStatus, params)
	}
}

// RequestLinuxDoQRCode 获取 LINUX DO Credit 二维码
func RequestLinuxDoQRCode(c *gin.Context) {
	var req LinuxDoPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "参数错误"})
		return
	}

	if req.Amount < getLinuxDoMinTopup() {
		c.JSON(200, gin.H{"message": "error", "data": fmt.Sprintf("充值数量不能小于 %d", getLinuxDoMinTopup())})
		return
	}

	id := c.GetInt("id")
	group, err := model.GetUserGroup(id, true)
	if err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "获取用户分组失败"})
		return
	}

	payMoney := getLinuxDoPayMoney(req.Amount, group)
	if payMoney < 0.01 {
		c.JSON(200, gin.H{"message": "error", "data": "充值金额过低"})
		return
	}

	// 生成订单号
	tradeNo := fmt.Sprintf("%s%d", common.GetRandomString(6), time.Now().Unix())
	tradeNo = fmt.Sprintf("LDO%dNO%s", id, tradeNo)

	// 获取回调地址
	callBackAddress := service.GetCallbackAddress()
	notifyUrl := callBackAddress + "/api/user/linuxdo/notify"
	returnUrl := system_setting.ServerAddress + "/console/topup"

	// 调用 LINUX DO Credit API 获取支付链接
	apiReq := &service.LinuxDoAPIRequest{
		Type:       "epay",
		OutTradeNo: tradeNo,
		NotifyURL:  notifyUrl,
		ReturnURL:  returnUrl,
		Name:       fmt.Sprintf("TUC%d", req.Amount),
		Money:      strconv.FormatFloat(payMoney, 'f', 2, 64),
		Device:     "pc",
	}

	apiResp, err := service.CallLinuxDoAPI(apiReq)
	if err != nil {
		errMsg := "获取支付链接失败"
		if apiResp != nil && apiResp.Msg != "" {
			errMsg = apiResp.Msg
		}
		c.JSON(200, gin.H{"message": "error", "data": errMsg})
		return
	}

	// 计算实际充值金额
	amount := req.Amount
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		dAmount := decimal.NewFromInt(int64(amount))
		dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
		amount = dAmount.Div(dQuotaPerUnit).IntPart()
	}

	// 创建订单记录
	topUp := &model.TopUp{
		UserId:        id,
		Amount:        amount,
		Money:         payMoney,
		TradeNo:       tradeNo,
		PaymentMethod: PaymentMethodLinuxDo,
		CreateTime:    time.Now().Unix(),
		Status:        "pending",
	}
	if err := topUp.Insert(); err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "创建订单失败"})
		return
	}

	// 返回支付链接
	responseData := gin.H{
		"trade_no": tradeNo,
		"amount":   req.Amount,
		"money":    strconv.FormatFloat(payMoney, 'f', 2, 64),
	}

	// 优先返回支付链接
	if apiResp.PayURL != "" {
		responseData["payurl"] = apiResp.PayURL
	} else if apiResp.QRCode != "" {
		responseData["qrcode"] = apiResp.QRCode
	}

	c.JSON(200, gin.H{"message": "success", "data": responseData})
}

// RequestLinuxDoAmount 获取 LINUX DO Credit 支付金额
func RequestLinuxDoAmount(c *gin.Context) {
	var req LinuxDoPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "参数错误"})
		return
	}

	if req.Amount < getLinuxDoMinTopup() {
		c.JSON(200, gin.H{"message": "error", "data": fmt.Sprintf("充值数量不能小于 %d", getLinuxDoMinTopup())})
		return
	}

	id := c.GetInt("id")
	group, err := model.GetUserGroup(id, true)
	if err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "获取用户分组失败"})
		return
	}

	payMoney := getLinuxDoPayMoney(req.Amount, group)
	if payMoney <= 0.01 {
		c.JSON(200, gin.H{"message": "error", "data": "充值金额过低"})
		return
	}

	c.JSON(200, gin.H{"message": "success", "data": strconv.FormatFloat(payMoney, 'f', 2, 64)})
}

// GetLinuxDoOrderStatus 查询 LINUX DO Credit 订单状态
func GetLinuxDoOrderStatus(c *gin.Context) {
	tradeNo := c.Query("trade_no")
	if tradeNo == "" {
		c.JSON(200, gin.H{"message": "error", "data": "订单号不能为空"})
		return
	}

	// 从数据库查询订单
	topUp := model.GetTopUpByTradeNo(tradeNo)
	if topUp == nil {
		c.JSON(200, gin.H{"message": "error", "data": "订单不存在"})
		return
	}

	// 验证订单归属
	userId := c.GetInt("id")
	if topUp.UserId != userId {
		c.JSON(200, gin.H{"message": "error", "data": "无权查询此订单"})
		return
	}

	// 如果订单状态是 pending，尝试从 LINUX DO Credit 查询最新状态
	if topUp.Status == "pending" {
		isPaid, err := service.QueryLinuxDoOrder(tradeNo)
		if err == nil && isPaid {
			// 更新订单状态
			LockOrder(tradeNo)
			defer UnlockOrder(tradeNo)

			// 重新获取订单，防止并发问题
			topUp = model.GetTopUpByTradeNo(tradeNo)
			if topUp != nil && topUp.Status == "pending" {
				topUp.Status = "success"
				if err := topUp.Update(); err == nil {
					// 增加用户额度
					dAmount := decimal.NewFromInt(int64(topUp.Amount))
					dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
					quotaToAdd := int(dAmount.Mul(dQuotaPerUnit).IntPart())
					model.IncreaseUserQuota(topUp.UserId, quotaToAdd, true)
					model.RecordLog(topUp.UserId, model.LogTypeTopup, fmt.Sprintf("使用 LINUX DO Credit 充值成功，充值金额: %v，支付积分：%f", logger.LogQuota(quotaToAdd), topUp.Money))
				}
			}
		}
	}

	// 返回订单状态
	responseData := gin.H{
		"trade_no": topUp.TradeNo,
		"status":   topUp.Status,
	}

	if topUp.Status == "success" {
		responseData["amount"] = topUp.Amount
	}

	c.JSON(200, gin.H{"message": "success", "data": responseData})
}

