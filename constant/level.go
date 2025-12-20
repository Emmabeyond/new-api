package constant

// 等级相关错误码

const (
	// 等级配置错误 (4001xx)
	ErrLevelNotFound        = 400101 // 等级不存在
	ErrLevelIdDuplicate     = 400102 // 等级ID重复
	ErrLevelPriorityDup     = 400103 // 等级优先级重复
	ErrLevelHasUsers        = 400104 // 等级下有用户，无法删除
	ErrLevelConfigInvalid   = 400105 // 等级配置无效
	ErrDefaultLevelRequired = 400106 // 必须存在默认等级
	ErrDefaultLevelDelete   = 400107 // 不能删除默认等级

	// 等级权限错误 (4002xx)
	ErrGroupNotAllowed   = 400201 // 用户等级无权访问该分组
	ErrRateLimitExceeded = 400202 // 超过速率限制

	// 等级升级错误 (4003xx)
	ErrLevelDowngrade     = 400301 // 不允许降级（手动调整除外）
	ErrSyncRechargeFailed = 400302 // 同步累计充值失败
)

// 等级相关错误消息
var LevelErrorMessages = map[int]string{
	ErrLevelNotFound:        "等级不存在",
	ErrLevelIdDuplicate:     "等级ID已存在",
	ErrLevelPriorityDup:     "等级优先级已存在",
	ErrLevelHasUsers:        "该等级下有用户，无法删除",
	ErrLevelConfigInvalid:   "等级配置无效",
	ErrDefaultLevelRequired: "必须存在默认等级",
	ErrDefaultLevelDelete:   "不能删除默认等级",
	ErrGroupNotAllowed:      "您的等级无权访问该渠道分组",
	ErrRateLimitExceeded:    "请求过于频繁，请稍后重试",
	ErrLevelDowngrade:       "不允许自动降级",
	ErrSyncRechargeFailed:   "同步累计充值失败",
}

// GetLevelErrorMessage 获取等级错误消息
func GetLevelErrorMessage(code int) string {
	if msg, ok := LevelErrorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
