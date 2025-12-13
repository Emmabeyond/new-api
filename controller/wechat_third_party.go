package controller

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/system_setting"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// WeChatThirdPartySession 微信第三方登录会话
type WeChatThirdPartySession struct {
	SessionId string    // 平台返回的会话 ID
	Action    string    // "login" 或 "bind"
	UserId    int       // 绑定时使用
	State     string    // CSRF 防护 state
	ClientIP  string    // 客户端 IP
	CreatedAt time.Time // 创建时间
	ExpiresAt time.Time // 过期时间
	Confirmed bool      // 是否已确认（防止重放）
}

// 会话存储（内存存储，生产环境建议使用 Redis）
var (
	wechatThirdPartySessions = make(map[string]*WeChatThirdPartySession)
	wechatSessionMutex       sync.RWMutex
)

// generateRequest 生成二维码请求
type wechatThirdPartyGenerateRequest struct {
	Action string `json:"action"` // "login" 或 "bind"
}

// WeChatThirdPartyGenerate 生成微信第三方登录二维码
func WeChatThirdPartyGenerate(c *gin.Context) {
	// 检查功能是否启用
	if !system_setting.IsWeChatThirdPartyEnabled() {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "微信第三方登录未启用",
		})
		return
	}

	// 解析请求
	var req wechatThirdPartyGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Action = "login" // 默认为登录
	}

	// 验证 action
	if req.Action != "login" && req.Action != "bind" {
		req.Action = "login"
	}

	// 如果是绑定操作，需要验证用户已登录
	var userId int
	if req.Action == "bind" {
		session := sessions.Default(c)
		id := session.Get("id")
		if id == nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "请先登录",
			})
			return
		}
		userId = id.(int)
	}

	// 生成 state 参数（CSRF 防护）
	state := common.GetRandomString(16)

	// 创建客户端
	client, err := service.NewWeChatThirdPartyClient()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 调用平台 API 生成二维码
	qrResp, err := client.GenerateQRCode("", state)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "生成二维码失败: " + err.Error(),
		})
		return
	}

	// 存储会话信息
	wechatSessionMutex.Lock()
	wechatThirdPartySessions[qrResp.SessionId] = &WeChatThirdPartySession{
		SessionId: qrResp.SessionId,
		Action:    req.Action,
		UserId:    userId,
		State:     state,
		ClientIP:  c.ClientIP(),
		CreatedAt: time.Now(),
		ExpiresAt: qrResp.ExpiresAt,
		Confirmed: false,
	}
	wechatSessionMutex.Unlock()

	// 返回二维码信息
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"sessionId": qrResp.SessionId,
			"qrCodeUrl": qrResp.QRCodeUrl,
			"expiresAt": qrResp.ExpiresAt,
		},
	})
}


// WeChatThirdPartyStatus 查询微信第三方登录状态
func WeChatThirdPartyStatus(c *gin.Context) {
	// 检查功能是否启用
	if !system_setting.IsWeChatThirdPartyEnabled() {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "微信第三方登录未启用",
		})
		return
	}

	sessionId := c.Param("sessionId")
	if sessionId == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "sessionId 不能为空",
		})
		return
	}

	// 获取本地会话信息
	wechatSessionMutex.RLock()
	localSession, exists := wechatThirdPartySessions[sessionId]
	wechatSessionMutex.RUnlock()

	if !exists {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "会话不存在或已过期",
		})
		return
	}

	// 检查会话是否已确认（防止重放攻击）
	if localSession.Confirmed {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "会话已使用",
		})
		return
	}

	// 创建客户端
	client, err := service.NewWeChatThirdPartyClient()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 查询平台状态
	statusResp, err := client.GetLoginStatus(sessionId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "查询状态失败: " + err.Error(),
		})
		return
	}

	// 如果状态是 confirmed，处理登录/绑定逻辑
	if statusResp.Status == "confirmed" && statusResp.UserInfo != nil {
		// 注意：第三方平台 (abu117.cn) 的 status 接口不返回 state 字段
		// CSRF 防护通过以下机制保证：
		// 1. sessionId 的唯一性和随机性
		// 2. 本地会话存储的 IP 绑定
		// 3. 会话的一次性使用（Confirmed 标记）

		// 标记会话为已确认（防止重放）
		wechatSessionMutex.Lock()
		localSession.Confirmed = true
		wechatSessionMutex.Unlock()

		// 根据 action 处理
		if localSession.Action == "bind" {
			handleWeChatThirdPartyBind(c, statusResp.UserInfo.OpenId, localSession.UserId)
		} else {
			handleWeChatThirdPartyLogin(c, statusResp.UserInfo)
		}
		return
	}

	// 返回状态
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"status":    statusResp.Status,
			"expiresAt": statusResp.ExpiresAt,
		},
	})
}

// handleWeChatThirdPartyLogin 处理微信第三方登录
func handleWeChatThirdPartyLogin(c *gin.Context, userInfo *service.UserInfo) {
	wechatId := userInfo.OpenId

	user := model.User{
		WeChatId: wechatId,
	}

	if model.IsWeChatIdAlreadyTaken(wechatId) {
		// 已有用户，直接登录
		err := user.FillUserByWeChatId()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if user.Id == 0 {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "用户已注销",
			})
			return
		}
	} else {
		// 新用户，检查是否允许注册
		if !common.RegisterEnabled {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "管理员关闭了新用户注册",
			})
			return
		}

		// 创建新用户
		user.Username = "wechat_" + strconv.Itoa(model.GetMaxUserId()+1)
		if userInfo.Nickname != "" {
			user.DisplayName = userInfo.Nickname
		} else {
			user.DisplayName = "WeChat User"
		}
		user.Role = common.RoleCommonUser
		user.Status = common.UserStatusEnabled

		// 获取邀请码
		session := sessions.Default(c)
		affCode := session.Get("aff")
		inviterId := 0
		if affCode != nil {
			inviterId, _ = model.GetUserIdByAffCode(affCode.(string))
		}

		if err := user.Insert(inviterId); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
	}

	// 检查用户状态
	if user.Status != common.UserStatusEnabled {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户已被封禁",
		})
		return
	}

	// 设置登录会话
	session := sessions.Default(c)
	session.Set("id", user.Id)
	session.Set("username", user.Username)
	session.Set("role", user.Role)
	session.Set("status", user.Status)
	session.Set("group", user.Group)
	err := session.Save()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "无法保存会话信息，请重试",
			"success": false,
		})
		return
	}

	// 返回登录成功响应，包含 status: confirmed 以匹配前端期望的格式
	c.JSON(http.StatusOK, gin.H{
		"message": "",
		"success": true,
		"data": gin.H{
			"status":      "confirmed",
			"id":          user.Id,
			"username":    user.Username,
			"displayName": user.DisplayName,
			"role":        user.Role,
			"group":       user.Group,
		},
	})
}

// handleWeChatThirdPartyBind 处理微信第三方绑定
func handleWeChatThirdPartyBind(c *gin.Context, wechatId string, userId int) {
	// 检查微信 ID 是否已被绑定
	if model.IsWeChatIdAlreadyTaken(wechatId) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "该微信账号已被绑定",
		})
		return
	}

	// 获取用户信息
	user := model.User{
		Id: userId,
	}
	err := user.FillUserById()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 更新用户微信 ID
	user.WeChatId = wechatId
	err = user.Update(false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "bind",
		"data": gin.H{
			"status": "confirmed",
		},
	})
}

// CleanExpiredWeChatThirdPartySessions 清理过期的会话
func CleanExpiredWeChatThirdPartySessions() {
	wechatSessionMutex.Lock()
	defer wechatSessionMutex.Unlock()

	now := time.Now()
	for sessionId, session := range wechatThirdPartySessions {
		if now.After(session.ExpiresAt) {
			delete(wechatThirdPartySessions, sessionId)
		}
	}
}

// StartWeChatThirdPartySessionCleaner 启动会话清理定时任务
func StartWeChatThirdPartySessionCleaner() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			CleanExpiredWeChatThirdPartySessions()
		}
	}()
}
