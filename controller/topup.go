package controller

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"sync"
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

func GetTopUpInfo(c *gin.Context) {
	// 获取支付方式
	payMethods := operation_setting.PayMethods

	// 如果启用了 Stripe 支付，添加到支付方法列表
	if setting.StripeApiSecret != "" && setting.StripeWebhookSecret != "" && setting.StripePriceId != "" {
		// 检查是否已经包含 Stripe
		hasStripe := false
		for _, method := range payMethods {
			if method["type"] == "stripe" {
				hasStripe = true
				break
			}
		}

		if !hasStripe {
			stripeMethod := map[string]string{
				"name":      "Stripe",
				"type":      "stripe",
				"color":     "rgba(var(--semi-purple-5), 1)",
				"min_topup": strconv.Itoa(setting.StripeMinTopUp),
			}
			payMethods = append(payMethods, stripeMethod)
		}
	}

	data := gin.H{
		"enable_online_topup":  operation_setting.PayAddress != "" && operation_setting.EpayId != "" && operation_setting.EpayKey != "",
		"enable_stripe_topup":  setting.StripeApiSecret != "" && setting.StripeWebhookSecret != "" && setting.StripePriceId != "",
		"enable_creem_topup":   setting.CreemApiKey != "" && setting.CreemProducts != "[]",
		"enable_linuxdo_topup": setting.LinuxDoClientId != "" && setting.LinuxDoClientSecret != "",
		"creem_products":       setting.CreemProducts,
		"pay_methods":          payMethods,
		"min_topup":            operation_setting.MinTopUp,
		"stripe_min_topup":     setting.StripeMinTopUp,
		"linuxdo_min_topup":    setting.LinuxDoMinTopUp,
		"amount_options":       operation_setting.GetPaymentSetting().AmountOptions,
		"discount":             operation_setting.GetPaymentSetting().AmountDiscount,
	}
	common.ApiSuccess(c, data)
}

type EpayRequest struct {
	Amount        int64  `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	TopUpCode     string `json:"top_up_code"`
}

type AmountRequest struct {
	Amount    int64  `json:"amount"`
	TopUpCode string `json:"top_up_code"`
}

func GetEpayClient() *epay.Client {
	if operation_setting.PayAddress == "" || operation_setting.EpayId == "" || operation_setting.EpayKey == "" {
		return nil
	}
	withUrl, err := epay.NewClient(&epay.Config{
		PartnerID: operation_setting.EpayId,
		Key:       operation_setting.EpayKey,
	}, operation_setting.PayAddress)
	if err != nil {
		return nil
	}
	return withUrl
}

func getPayMoney(amount int64, group string) float64 {
	dAmount := decimal.NewFromInt(amount)
	// 充值金额以“展示类型”为准：
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
	dPrice := decimal.NewFromFloat(operation_setting.Price)
	// apply optional preset discount by the original request amount (if configured), default 1.0
	discount := 1.0
	if ds, ok := operation_setting.GetPaymentSetting().AmountDiscount[int(amount)]; ok {
		if ds > 0 {
			discount = ds
		}
	}
	dDiscount := decimal.NewFromFloat(discount)

	payMoney := dAmount.Mul(dPrice).Mul(dTopupGroupRatio).Mul(dDiscount)

	return payMoney.InexactFloat64()
}

func getMinTopup() int64 {
	minTopup := operation_setting.MinTopUp
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		dMinTopup := decimal.NewFromInt(int64(minTopup))
		dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
		minTopup = int(dMinTopup.Mul(dQuotaPerUnit).IntPart())
	}
	return int64(minTopup)
}

func RequestEpay(c *gin.Context) {
	var req EpayRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "参数错误"})
		return
	}
	if req.Amount < getMinTopup() {
		c.JSON(200, gin.H{"message": "error", "data": fmt.Sprintf("充值数量不能小于 %d", getMinTopup())})
		return
	}

	id := c.GetInt("id")
	group, err := model.GetUserGroup(id, true)
	if err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "获取用户分组失败"})
		return
	}
	payMoney := getPayMoney(req.Amount, group)
	if payMoney < 0.01 {
		c.JSON(200, gin.H{"message": "error", "data": "充值金额过低"})
		return
	}

	if !operation_setting.ContainsPayMethod(req.PaymentMethod) {
		c.JSON(200, gin.H{"message": "error", "data": "支付方式不存在"})
		return
	}

	callBackAddress := service.GetCallbackAddress()
	returnUrl, _ := url.Parse(system_setting.ServerAddress + "/console/log")
	notifyUrl, _ := url.Parse(callBackAddress + "/api/user/epay/notify")
	tradeNo := fmt.Sprintf("%s%d", common.GetRandomString(6), time.Now().Unix())
	tradeNo = fmt.Sprintf("USR%dNO%s", id, tradeNo)
	client := GetEpayClient()
	if client == nil {
		c.JSON(200, gin.H{"message": "error", "data": "当前管理员未配置支付信息"})
		return
	}
	uri, params, err := client.Purchase(&epay.PurchaseArgs{
		Type:           req.PaymentMethod,
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
		PaymentMethod: req.PaymentMethod,
		CreateTime:    time.Now().Unix(),
		Status:        "pending",
	}
	err = topUp.Insert()
	if err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "创建订单失败"})
		return
	}
	c.JSON(200, gin.H{"message": "success", "data": params, "url": uri})
}

// tradeNo lock
var orderLocks sync.Map
var createLock sync.Mutex

// LockOrder 尝试对给定订单号加锁
func LockOrder(tradeNo string) {
	lock, ok := orderLocks.Load(tradeNo)
	if !ok {
		createLock.Lock()
		defer createLock.Unlock()
		lock, ok = orderLocks.Load(tradeNo)
		if !ok {
			lock = new(sync.Mutex)
			orderLocks.Store(tradeNo, lock)
		}
	}
	lock.(*sync.Mutex).Lock()
}

// UnlockOrder 释放给定订单号的锁
func UnlockOrder(tradeNo string) {
	lock, ok := orderLocks.Load(tradeNo)
	if ok {
		lock.(*sync.Mutex).Unlock()
	}
}

func EpayNotify(c *gin.Context) {
	params := lo.Reduce(lo.Keys(c.Request.URL.Query()), func(r map[string]string, t string, i int) map[string]string {
		r[t] = c.Request.URL.Query().Get(t)
		return r
	}, map[string]string{})
	client := GetEpayClient()
	if client == nil {
		log.Println("易支付回调失败 未找到配置信息")
		_, err := c.Writer.Write([]byte("fail"))
		if err != nil {
			log.Println("易支付回调写入失败")
		}
		return
	}
	verifyInfo, err := client.Verify(params)
	if err == nil && verifyInfo.VerifyStatus {
		_, err := c.Writer.Write([]byte("success"))
		if err != nil {
			log.Println("易支付回调写入失败")
		}
	} else {
		_, err := c.Writer.Write([]byte("fail"))
		if err != nil {
			log.Println("易支付回调写入失败")
		}
		log.Println("易支付回调签名验证失败")
		return
	}

	if verifyInfo.TradeStatus == epay.StatusTradeSuccess {
		log.Println(verifyInfo)
		LockOrder(verifyInfo.ServiceTradeNo)
		defer UnlockOrder(verifyInfo.ServiceTradeNo)
		topUp := model.GetTopUpByTradeNo(verifyInfo.ServiceTradeNo)
		if topUp == nil {
			log.Printf("易支付回调未找到订单: %v", verifyInfo)
			return
		}
		if topUp.Status == "pending" {
			topUp.Status = "success"
			err := topUp.Update()
			if err != nil {
				log.Printf("易支付回调更新订单失败: %v", topUp)
				return
			}
			//user, _ := model.GetUserById(topUp.UserId, false)
			//user.Quota += topUp.Amount * 500000
			dAmount := decimal.NewFromInt(int64(topUp.Amount))
			dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
			quotaToAdd := int(dAmount.Mul(dQuotaPerUnit).IntPart())
			err = model.IncreaseUserQuota(topUp.UserId, quotaToAdd, true)
			if err != nil {
				log.Printf("易支付回调更新用户失败: %v", topUp)
				return
			}
			log.Printf("易支付回调更新用户成功 %v", topUp)
			model.RecordLog(topUp.UserId, model.LogTypeTopup, fmt.Sprintf("使用在线充值成功，充值金额: %v，支付金额：%f", logger.LogQuota(quotaToAdd), topUp.Money))
		}
	} else {
		log.Printf("易支付异常回调: %v", verifyInfo)
	}
}

func RequestAmount(c *gin.Context) {
	var req AmountRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "参数错误"})
		return
	}

	if req.Amount < getMinTopup() {
		c.JSON(200, gin.H{"message": "error", "data": fmt.Sprintf("充值数量不能小于 %d", getMinTopup())})
		return
	}
	id := c.GetInt("id")
	group, err := model.GetUserGroup(id, true)
	if err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "获取用户分组失败"})
		return
	}
	payMoney := getPayMoney(req.Amount, group)
	if payMoney <= 0.01 {
		c.JSON(200, gin.H{"message": "error", "data": "充值金额过低"})
		return
	}
	c.JSON(200, gin.H{"message": "success", "data": strconv.FormatFloat(payMoney, 'f', 2, 64)})
}

func GetUserTopUps(c *gin.Context) {
	userId := c.GetInt("id")
	pageInfo := common.GetPageQuery(c)
	keyword := c.Query("keyword")

	var (
		topups []*model.TopUp
		total  int64
		err    error
	)
	if keyword != "" {
		topups, total, err = model.SearchUserTopUps(userId, keyword, pageInfo)
	} else {
		topups, total, err = model.GetUserTopUps(userId, pageInfo)
	}
	if err != nil {
		common.ApiError(c, err)
		return
	}

	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(topups)
	common.ApiSuccess(c, pageInfo)
}

// GetAllTopUps 管理员获取全平台充值记录（支持多条件筛选）
func GetAllTopUps(c *gin.Context) {
	pageInfo := common.GetPageQuery(c)

	// 解析筛选参数
	filter := &model.TopUpFilter{}
	filter.Keyword = c.Query("keyword")
	filter.Status = c.Query("status")
	filter.PaymentMethod = c.Query("payment_method")
	if userIdStr := c.Query("user_id"); userIdStr != "" {
		if userId, err := strconv.Atoi(userIdStr); err == nil {
			filter.UserId = userId
		}
	}
	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		if startTime, err := strconv.ParseInt(startTimeStr, 10, 64); err == nil {
			filter.StartTime = startTime
		}
	}
	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		if endTime, err := strconv.ParseInt(endTimeStr, 10, 64); err == nil {
			filter.EndTime = endTime
		}
	}

	topups, total, err := model.GetAllTopUpsWithFilter(filter, pageInfo)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(topups)
	common.ApiSuccess(c, pageInfo)
}

// GetTopUpDetail 管理员获取充值订单详情
func GetTopUpDetail(c *gin.Context) {
	idStr := c.Param("topup_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ApiErrorMsg(c, "无效的订单ID")
		return
	}

	detail, err := model.GetTopUpDetail(id)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	common.ApiSuccess(c, detail)
}

// GetTopUpStats 管理员获取充值统计数据
func GetTopUpStats(c *gin.Context) {
	// 解析筛选参数
	filter := &model.TopUpFilter{}
	filter.Keyword = c.Query("keyword")
	filter.Status = c.Query("status")
	filter.PaymentMethod = c.Query("payment_method")
	if userIdStr := c.Query("user_id"); userIdStr != "" {
		if userId, err := strconv.Atoi(userIdStr); err == nil {
			filter.UserId = userId
		}
	}
	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		if startTime, err := strconv.ParseInt(startTimeStr, 10, 64); err == nil {
			filter.StartTime = startTime
		}
	}
	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		if endTime, err := strconv.ParseInt(endTimeStr, 10, 64); err == nil {
			filter.EndTime = endTime
		}
	}

	stats, err := model.GetTopUpStats(filter)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	common.ApiSuccess(c, stats)
}

// ExportTopUps 管理员导出充值记录
func ExportTopUps(c *gin.Context) {
	// 解析筛选参数
	filter := &model.TopUpFilter{}
	filter.Keyword = c.Query("keyword")
	filter.Status = c.Query("status")
	filter.PaymentMethod = c.Query("payment_method")
	if userIdStr := c.Query("user_id"); userIdStr != "" {
		if userId, err := strconv.Atoi(userIdStr); err == nil {
			filter.UserId = userId
		}
	}
	if startTimeStr := c.Query("start_time"); startTimeStr != "" {
		if startTime, err := strconv.ParseInt(startTimeStr, 10, 64); err == nil {
			filter.StartTime = startTime
		}
	}
	if endTimeStr := c.Query("end_time"); endTimeStr != "" {
		if endTime, err := strconv.ParseInt(endTimeStr, 10, 64); err == nil {
			filter.EndTime = endTime
		}
	}

	topups, err := model.GetAllTopUpsForExport(filter)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	// 设置响应头
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=topup_export.csv")

	// 写入 BOM 以支持 Excel 正确识别 UTF-8
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	// 写入 CSV 头
	c.Writer.WriteString("订单ID,用户ID,充值金额,支付金额,订单号,支付方式,创建时间,完成时间,状态\n")

	// 写入数据行
	for _, topup := range topups {
		createTime := time.Unix(topup.CreateTime, 0).Format("2006-01-02 15:04:05")
		completeTime := ""
		if topup.CompleteTime > 0 {
			completeTime = time.Unix(topup.CompleteTime, 0).Format("2006-01-02 15:04:05")
		}

		line := fmt.Sprintf("%d,%d,%d,%.2f,%s,%s,%s,%s,%s\n",
			topup.Id,
			topup.UserId,
			topup.Amount,
			topup.Money,
			topup.TradeNo,
			topup.PaymentMethod,
			createTime,
			completeTime,
			topup.Status,
		)
		c.Writer.WriteString(line)
	}
}

type AdminCompleteTopupRequest struct {
	TradeNo string `json:"trade_no"`
}

// AdminCompleteTopUp 管理员补单接口
func AdminCompleteTopUp(c *gin.Context) {
	var req AdminCompleteTopupRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.TradeNo == "" {
		common.ApiErrorMsg(c, "参数错误")
		return
	}

	// 订单级互斥，防止并发补单
	LockOrder(req.TradeNo)
	defer UnlockOrder(req.TradeNo)

	if err := model.ManualCompleteTopUp(req.TradeNo); err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, nil)
}


// QRCodeRequest 二维码支付请求
type QRCodeRequest struct {
	Amount        int64  `json:"amount"`
	PaymentMethod string `json:"payment_method"`
}

// RequestEpayQRCode 获取易支付二维码（也支持 LINUX DO Credit）
// POST /api/user/pay/qrcode
func RequestEpayQRCode(c *gin.Context) {
	var req QRCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "参数错误"})
		return
	}

	// 判断是否为 LINUX DO Credit 支付
	isLinuxDo := req.PaymentMethod == PaymentMethodLinuxDo

	// 获取最小充值金额
	var minTopupVal int64
	if isLinuxDo {
		minTopupVal = getLinuxDoMinTopup()
	} else {
		minTopupVal = getMinTopup()
	}

	if req.Amount < minTopupVal {
		c.JSON(200, gin.H{"message": "error", "data": fmt.Sprintf("充值数量不能小于 %d", minTopupVal)})
		return
	}

	id := c.GetInt("id")
	group, err := model.GetUserGroup(id, true)
	if err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "获取用户分组失败"})
		return
	}

	// 根据支付方式计算支付金额
	var payMoney float64
	if isLinuxDo {
		payMoney = getLinuxDoPayMoney(req.Amount, group)
	} else {
		payMoney = getPayMoney(req.Amount, group)
	}

	if payMoney < 0.01 {
		c.JSON(200, gin.H{"message": "error", "data": "充值金额过低"})
		return
	}

	// 非 LINUX DO Credit 需要检查支付方式是否存在
	if !isLinuxDo && !operation_setting.ContainsPayMethod(req.PaymentMethod) {
		c.JSON(200, gin.H{"message": "error", "data": "支付方式不存在"})
		return
	}

	// 生成订单号
	tradeNo := fmt.Sprintf("%s%d", common.GetRandomString(6), time.Now().Unix())
	if isLinuxDo {
		tradeNo = fmt.Sprintf("LDO%dNO%s", id, tradeNo)
	} else {
		tradeNo = fmt.Sprintf("USR%dNO%s", id, tradeNo)
	}

	// 获取回调地址
	callBackAddress := service.GetCallbackAddress()
	returnUrl := system_setting.ServerAddress + "/console/topup"

	var notifyUrl string
	if isLinuxDo {
		notifyUrl = callBackAddress + "/api/user/linuxdo/notify"
	} else {
		notifyUrl = callBackAddress + "/api/user/epay/notify"
	}

	var qrCode, payUrl, urlScheme string

	if isLinuxDo {
		// 调用 LINUX DO Credit API
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
			errMsg := "获取支付二维码失败"
			if apiResp != nil && apiResp.Msg != "" {
				errMsg = apiResp.Msg
			}
			c.JSON(200, gin.H{"message": "error", "data": errMsg})
			return
		}
		
		qrCode = apiResp.QRCode
		payUrl = apiResp.PayURL
		urlScheme = apiResp.URLScheme
	} else {
		// 获取用户 IP（易支付需要）
		clientIP := c.ClientIP()
		
		// 调用易支付 API
		apiReq := &service.EpayAPIRequest{
			Type:       req.PaymentMethod,
			OutTradeNo: tradeNo,
			NotifyURL:  notifyUrl,
			ReturnURL:  returnUrl,
			Name:       fmt.Sprintf("TUC%d", req.Amount),
			Money:      strconv.FormatFloat(payMoney, 'f', 2, 64),
			ClientIP:   clientIP,
			Device:     "pc",
		}

		apiResp, err := service.CallEpayAPI(apiReq)
		if err != nil {
			errMsg := "获取支付二维码失败"
			if apiResp != nil && apiResp.Msg != "" {
				errMsg = apiResp.Msg
			}
			c.JSON(200, gin.H{"message": "error", "data": errMsg})
			return
		}
		qrCode = apiResp.QRCode
		payUrl = apiResp.PayURL
		urlScheme = apiResp.URLScheme
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
		PaymentMethod: req.PaymentMethod,
		CreateTime:    time.Now().Unix(),
		Status:        "pending",
	}
	if err := topUp.Insert(); err != nil {
		c.JSON(200, gin.H{"message": "error", "data": "创建订单失败"})
		return
	}

	// 返回支付信息
	responseData := gin.H{
		"trade_no": tradeNo,
		"amount":   req.Amount,
		"money":    strconv.FormatFloat(payMoney, 'f', 2, 64),
	}

	// LINUX DO Credit: 优先返回支付链接（二维码需要登录支持，直接跳转支付更可靠）
	// 易支付: 优先返回二维码，其次返回支付链接
	if isLinuxDo {
		// LINUX DO Credit 直接返回支付链接，不返回二维码
		if payUrl != "" {
			responseData["payurl"] = payUrl
		}
	} else {
		// 易支付：优先返回二维码，其次返回支付链接
		if qrCode != "" {
			responseData["qrcode"] = qrCode
		} else if payUrl != "" {
			responseData["payurl"] = payUrl
		} else if urlScheme != "" {
			responseData["urlscheme"] = urlScheme
		}
	}

	c.JSON(200, gin.H{"message": "success", "data": responseData})
}

// GetOrderStatus 查询订单支付状态
// GET /api/user/topup/status?trade_no=xxx
func GetOrderStatus(c *gin.Context) {
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

	// 如果订单状态是 pending，根据支付方式查询最新状态
	if topUp.Status == "pending" {
		var isPaid bool
		var err error

		// 根据支付方式调用对应的查询接口
		if topUp.PaymentMethod == PaymentMethodLinuxDo {
			isPaid, err = service.QueryLinuxDoOrder(tradeNo)
		} else {
			isPaid, err = service.QueryEpayOrder(tradeNo)
		}

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

					logMsg := "使用在线充值成功"
					if topUp.PaymentMethod == PaymentMethodLinuxDo {
						logMsg = "使用 LINUX DO Credit 充值成功"
					}
					model.RecordLog(topUp.UserId, model.LogTypeTopup, fmt.Sprintf("%s，充值金额: %v，支付金额：%f", logMsg, logger.LogQuota(quotaToAdd), topUp.Money))
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
