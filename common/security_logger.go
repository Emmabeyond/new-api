package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

// SecurityEvent 安全事件
type SecurityEvent struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Timestamp   int64  `gorm:"index" json:"timestamp"`
	UserId      int    `gorm:"index" json:"user_id"`
	EventType   string `gorm:"index;type:varchar(50)" json:"event_type"` // "xss_attempt", "csrf_attempt", "invalid_origin"
	Details     string `gorm:"type:text" json:"details"`
	IPAddress   string `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent   string `gorm:"type:text" json:"user_agent"`
	RequestPath string `gorm:"type:varchar(255)" json:"request_path"`
}

// TableName 指定表名
func (SecurityEvent) TableName() string {
	return "security_events"
}

// ArchivedSecurityEvent 归档的安全事件
type ArchivedSecurityEvent struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	OriginalID  int64  `gorm:"index" json:"original_id"`
	Timestamp   int64  `gorm:"index" json:"timestamp"`
	UserId      int    `gorm:"index" json:"user_id"`
	EventType   string `gorm:"index;type:varchar(50)" json:"event_type"`
	Details     string `gorm:"type:text" json:"details"`
	IPAddress   string `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent   string `gorm:"type:text" json:"user_agent"`
	RequestPath string `gorm:"type:varchar(255)" json:"request_path"`
	ArchivedAt  int64  `gorm:"index" json:"archived_at"`
}

// TableName 指定归档表名
func (ArchivedSecurityEvent) TableName() string {
	return "archived_security_events"
}

// FailedWebhookNotification 失败的 Webhook 通知记录
type FailedWebhookNotification struct {
	ID         int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	WebhookURL string `gorm:"type:varchar(500)" json:"webhook_url"`
	EventID    int64  `gorm:"index" json:"event_id"`
	Payload    string `gorm:"type:text" json:"payload"`
	Attempts   int    `json:"attempts"`
	LastError  string `gorm:"type:text" json:"last_error"`
	CreatedAt  int64  `gorm:"index" json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

// TableName 指定失败通知表名
func (FailedWebhookNotification) TableName() string {
	return "failed_webhook_notifications"
}

// SecurityArchiveConfig 归档配置
type SecurityArchiveConfig struct {
	Enabled           bool  // 是否启用归档
	RetentionDays     int   // 保留天数（超过此天数的日志将被归档）
	ArchiveThreshold  int64 // 归档阈值（当日志数量超过此值时触发归档）
	ArchiveIntervalHr int   // 归档检查间隔（小时）
}

// WebhookConfig Webhook 配置
type WebhookConfig struct {
	Enabled    bool   // 是否启用 Webhook 通知
	URL        string // Webhook URL
	MaxRetries int    // 最大重试次数
	TimeoutSec int    // 超时时间（秒）
}

// SecurityEventQueryParams 查询参数
type SecurityEventQueryParams struct {
	StartTime int64
	EndTime   int64
	UserId    int
	EventType string
	Page      int
	PageSize  int
}

var (
	securityLogFile     *os.File
	securityLogMutex    sync.Mutex
	securityLogDir      = "logs"
	securityLogFileName = "security_events.log"
	securityDB          *gorm.DB // 数据库连接，由外部设置

	// 归档配置
	archiveConfig = &SecurityArchiveConfig{
		Enabled:           false,
		RetentionDays:     30,
		ArchiveThreshold:  10000,
		ArchiveIntervalHr: 24,
	}

	// Webhook 配置
	webhookConfig = &WebhookConfig{
		Enabled:    false,
		URL:        "",
		MaxRetries: 3,
		TimeoutSec: 5,
	}

	// 归档任务停止通道
	archiveStopChan chan struct{}
)

// InitSecurityLogger 初始化安全日志记录器
func InitSecurityLogger() error {
	securityLogMutex.Lock()
	defer securityLogMutex.Unlock()

	// 确保日志目录存在
	if err := os.MkdirAll(securityLogDir, 0755); err != nil {
		return fmt.Errorf("failed to create security log directory: %v", err)
	}

	logPath := filepath.Join(securityLogDir, securityLogFileName)
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open security log file: %v", err)
	}

	securityLogFile = f
	return nil
}

// SetSecurityDB 设置安全日志数据库连接
func SetSecurityDB(db *gorm.DB) {
	securityLogMutex.Lock()
	defer securityLogMutex.Unlock()
	securityDB = db
}

// LogSecurityEvent 记录安全事件（异步写入）
func LogSecurityEvent(event *SecurityEvent) {
	if event == nil {
		return
	}

	// 设置时间戳
	if event.Timestamp == 0 {
		event.Timestamp = time.Now().Unix()
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				SysLog(fmt.Sprintf("Failed to log security event: %v", r))
			}
		}()

		// 优先写入数据库
		if securityDB != nil {
			if err := writeSecurityEventToDB(event); err != nil {
				SysLog(fmt.Sprintf("Failed to write security event to database: %v, falling back to file", err))
				// 数据库写入失败，回退到本地文件
				writeSecurityEventToFile(event)
			}
		} else {
			// 没有数据库连接，写入本地文件
			writeSecurityEventToFile(event)
		}

		// 发送 Webhook 通知
		SendSecurityWebhook(event)
	}()
}

// writeSecurityEventToDB 写入安全事件到数据库
func writeSecurityEventToDB(event *SecurityEvent) error {
	if securityDB == nil {
		return fmt.Errorf("security database not initialized")
	}
	return securityDB.Create(event).Error
}

// writeSecurityEventToFile 写入安全事件到本地文件
func writeSecurityEventToFile(event *SecurityEvent) {
	securityLogMutex.Lock()
	defer securityLogMutex.Unlock()

	if securityLogFile == nil {
		// 尝试初始化
		if err := InitSecurityLogger(); err != nil {
			SysLog(fmt.Sprintf("Failed to initialize security logger: %v", err))
			return
		}
	}

	// 序列化事件
	eventJSON, err := json.Marshal(event)
	if err != nil {
		SysLog(fmt.Sprintf("Failed to marshal security event: %v", err))
		return
	}

	// 写入文件
	logLine := fmt.Sprintf("[%s] %s\n",
		time.Unix(event.Timestamp, 0).Format(time.RFC3339),
		string(eventJSON))

	if _, err := securityLogFile.WriteString(logLine); err != nil {
		SysLog(fmt.Sprintf("Failed to write security event to file: %v", err))
	}
}

// CloseSecurityLogger 关闭安全日志记录器
func CloseSecurityLogger() {
	securityLogMutex.Lock()
	defer securityLogMutex.Unlock()

	if securityLogFile != nil {
		securityLogFile.Close()
		securityLogFile = nil
	}
}

// GetSecurityEvents 查询安全事件
// 支持按时间范围、用户ID、事件类型筛选
func GetSecurityEvents(params *SecurityEventQueryParams) ([]*SecurityEvent, int64, error) {
	if securityDB == nil {
		return nil, 0, fmt.Errorf("security database not initialized")
	}

	var events []*SecurityEvent
	var total int64

	query := securityDB.Model(&SecurityEvent{})

	// 时间范围筛选
	if params.StartTime > 0 {
		query = query.Where("timestamp >= ?", params.StartTime)
	}
	if params.EndTime > 0 {
		query = query.Where("timestamp <= ?", params.EndTime)
	}

	// 用户ID筛选
	if params.UserId > 0 {
		query = query.Where("user_id = ?", params.UserId)
	}

	// 事件类型筛选
	if params.EventType != "" {
		query = query.Where("event_type = ?", params.EventType)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count security events: %v", err)
	}

	// 分页
	page := params.Page
	pageSize := params.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// 查询数据，按时间倒序
	if err := query.Order("timestamp DESC").Offset(offset).Limit(pageSize).Find(&events).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to query security events: %v", err)
	}

	return events, total, nil
}

// GetSecurityEventByID 根据ID查询单个安全事件
func GetSecurityEventByID(id int64) (*SecurityEvent, error) {
	if securityDB == nil {
		return nil, fmt.Errorf("security database not initialized")
	}

	var event SecurityEvent
	if err := securityDB.First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// GetSecurityEventStats 获取安全事件统计
func GetSecurityEventStats(startTime, endTime int64) (map[string]int64, error) {
	if securityDB == nil {
		return nil, fmt.Errorf("security database not initialized")
	}

	type Result struct {
		EventType string
		Count     int64
	}

	var results []Result
	query := securityDB.Model(&SecurityEvent{}).
		Select("event_type, COUNT(*) as count").
		Group("event_type")

	if startTime > 0 {
		query = query.Where("timestamp >= ?", startTime)
	}
	if endTime > 0 {
		query = query.Where("timestamp <= ?", endTime)
	}

	if err := query.Find(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get security event stats: %v", err)
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.EventType] = r.Count
	}

	return stats, nil
}

// DeleteOldSecurityEvents 删除旧的安全事件（用于归档后清理）
func DeleteOldSecurityEvents(beforeTimestamp int64) (int64, error) {
	if securityDB == nil {
		return 0, fmt.Errorf("security database not initialized")
	}

	result := securityDB.Where("timestamp < ?", beforeTimestamp).Delete(&SecurityEvent{})
	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete old security events: %v", result.Error)
	}

	return result.RowsAffected, nil
}

// ==================== 归档功能 ====================

// LoadArchiveConfigFromEnv 从环境变量加载归档配置
// Requirements: 7.1, 7.2, 7.3, 7.5 - 从环境变量读取配置，验证配置，使用默认值回退
func LoadArchiveConfigFromEnv() *SecurityArchiveConfig {
	config := &SecurityArchiveConfig{
		Enabled:           GetEnvOrDefaultBool("SECURITY_ARCHIVE_ENABLED", false),
		RetentionDays:     GetEnvOrDefault("SECURITY_ARCHIVE_RETENTION_DAYS", 30),
		ArchiveThreshold:  int64(GetEnvOrDefault("SECURITY_ARCHIVE_THRESHOLD", 10000)),
		ArchiveIntervalHr: GetEnvOrDefault("SECURITY_ARCHIVE_INTERVAL_HOURS", 24),
	}

	// 配置验证
	if config.RetentionDays < 1 || config.RetentionDays > 365 {
		SysError(fmt.Sprintf("Invalid SECURITY_ARCHIVE_RETENTION_DAYS: %d, using default: 30", config.RetentionDays))
		config.RetentionDays = 30
	}

	if config.ArchiveThreshold < 100 || config.ArchiveThreshold > 1000000 {
		SysError(fmt.Sprintf("Invalid SECURITY_ARCHIVE_THRESHOLD: %d, using default: 10000", config.ArchiveThreshold))
		config.ArchiveThreshold = 10000
	}

	if config.ArchiveIntervalHr < 1 || config.ArchiveIntervalHr > 168 { // 最多 7 天
		SysError(fmt.Sprintf("Invalid SECURITY_ARCHIVE_INTERVAL_HOURS: %d, using default: 24", config.ArchiveIntervalHr))
		config.ArchiveIntervalHr = 24
	}

	return config
}

// SetArchiveConfig 设置归档配置
func SetArchiveConfig(config *SecurityArchiveConfig) {
	securityLogMutex.Lock()
	defer securityLogMutex.Unlock()
	archiveConfig = config
}

// GetArchiveConfig 获取当前归档配置
func GetArchiveConfig() *SecurityArchiveConfig {
	securityLogMutex.Lock()
	defer securityLogMutex.Unlock()
	return archiveConfig
}

// StartSecurityLogArchiver 启动安全日志归档任务
func StartSecurityLogArchiver() {
	config := LoadArchiveConfigFromEnv()
	SetArchiveConfig(config)

	if !config.Enabled {
		SysLog("Security log archiver is disabled")
		return
	}

	archiveStopChan = make(chan struct{})

	go func() {
		interval := time.Duration(config.ArchiveIntervalHr) * time.Hour
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		SysLog(fmt.Sprintf("Security log archiver started, interval: %d hours, retention: %d days",
			config.ArchiveIntervalHr, config.RetentionDays))

		// 启动时先执行一次检查
		checkAndArchive()

		for {
			select {
			case <-ticker.C:
				checkAndArchive()
			case <-archiveStopChan:
				SysLog("Security log archiver stopped")
				return
			}
		}
	}()
}

// StopSecurityLogArchiver 停止安全日志归档任务
func StopSecurityLogArchiver() {
	if archiveStopChan != nil {
		close(archiveStopChan)
	}
}

// checkAndArchive 检查并执行归档
func checkAndArchive() {
	if securityDB == nil {
		return
	}

	config := GetArchiveConfig()
	if !config.Enabled {
		return
	}

	// 检查日志数量是否超过阈值
	var count int64
	if err := securityDB.Model(&SecurityEvent{}).Count(&count).Error; err != nil {
		SysLog(fmt.Sprintf("Failed to count security events for archive check: %v", err))
		return
	}

	if count < config.ArchiveThreshold {
		return
	}

	// 计算归档截止时间
	cutoffTime := time.Now().AddDate(0, 0, -config.RetentionDays).Unix()

	// 执行归档
	archived, err := ArchiveSecurityEvents(cutoffTime)
	if err != nil {
		SysLog(fmt.Sprintf("Failed to archive security events: %v", err))
		return
	}

	if archived > 0 {
		SysLog(fmt.Sprintf("Archived %d security events older than %d days", archived, config.RetentionDays))
	}
}

// ArchiveSecurityEvents 归档指定时间之前的安全事件
func ArchiveSecurityEvents(beforeTimestamp int64) (int64, error) {
	if securityDB == nil {
		return 0, fmt.Errorf("security database not initialized")
	}

	// 确保归档表存在
	if err := securityDB.AutoMigrate(&ArchivedSecurityEvent{}); err != nil {
		return 0, fmt.Errorf("failed to create archive table: %v", err)
	}

	// 开始事务
	tx := securityDB.Begin()
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	// 查询需要归档的事件
	var events []*SecurityEvent
	if err := tx.Where("timestamp < ?", beforeTimestamp).Find(&events).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to query events for archive: %v", err)
	}

	if len(events) == 0 {
		tx.Commit()
		return 0, nil
	}

	// 转换为归档事件
	archivedAt := time.Now().Unix()
	archivedEvents := make([]*ArchivedSecurityEvent, len(events))
	for i, event := range events {
		archivedEvents[i] = &ArchivedSecurityEvent{
			OriginalID:  event.ID,
			Timestamp:   event.Timestamp,
			UserId:      event.UserId,
			EventType:   event.EventType,
			Details:     event.Details,
			IPAddress:   event.IPAddress,
			UserAgent:   event.UserAgent,
			RequestPath: event.RequestPath,
			ArchivedAt:  archivedAt,
		}
	}

	// 批量插入归档表
	if err := tx.CreateInBatches(archivedEvents, 100).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert archived events: %v", err)
	}

	// 删除原表中的已归档事件
	if err := tx.Where("timestamp < ?", beforeTimestamp).Delete(&SecurityEvent{}).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to delete archived events from original table: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("failed to commit archive transaction: %v", err)
	}

	return int64(len(events)), nil
}

// GetArchivedSecurityEvents 查询归档的安全事件
func GetArchivedSecurityEvents(params *SecurityEventQueryParams) ([]*ArchivedSecurityEvent, int64, error) {
	if securityDB == nil {
		return nil, 0, fmt.Errorf("security database not initialized")
	}

	var events []*ArchivedSecurityEvent
	var total int64

	query := securityDB.Model(&ArchivedSecurityEvent{})

	// 时间范围筛选
	if params.StartTime > 0 {
		query = query.Where("timestamp >= ?", params.StartTime)
	}
	if params.EndTime > 0 {
		query = query.Where("timestamp <= ?", params.EndTime)
	}

	// 用户ID筛选
	if params.UserId > 0 {
		query = query.Where("user_id = ?", params.UserId)
	}

	// 事件类型筛选
	if params.EventType != "" {
		query = query.Where("event_type = ?", params.EventType)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count archived security events: %v", err)
	}

	// 分页
	page := params.Page
	pageSize := params.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// 查询数据，按时间倒序
	if err := query.Order("timestamp DESC").Offset(offset).Limit(pageSize).Find(&events).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to query archived security events: %v", err)
	}

	return events, total, nil
}


// ==================== Webhook 通知功能 ====================

// LoadWebhookConfigFromEnv 从环境变量加载 Webhook 配置
// Requirements: 7.1, 7.2, 7.3, 7.5 - 从环境变量读取配置，验证配置，使用默认值回退
func LoadWebhookConfigFromEnv() *WebhookConfig {
	config := &WebhookConfig{
		Enabled:    GetEnvOrDefaultBool("SECURITY_WEBHOOK_ENABLED", false),
		URL:        GetEnvOrDefaultString("SECURITY_WEBHOOK_URL", ""),
		MaxRetries: GetEnvOrDefault("SECURITY_WEBHOOK_MAX_RETRIES", 3),
		TimeoutSec: GetEnvOrDefault("SECURITY_WEBHOOK_TIMEOUT", 5),
	}

	// 配置验证
	if config.Enabled {
		if config.URL == "" {
			SysError("SECURITY_WEBHOOK_ENABLED is true but SECURITY_WEBHOOK_URL is empty, disabling webhook")
			config.Enabled = false
		} else if !isValidWebhookURL(config.URL) {
			SysError(fmt.Sprintf("Invalid SECURITY_WEBHOOK_URL: %s, disabling webhook", config.URL))
			config.Enabled = false
		}
	}

	// 验证重试次数范围
	if config.MaxRetries < 0 || config.MaxRetries > 10 {
		SysError(fmt.Sprintf("Invalid SECURITY_WEBHOOK_MAX_RETRIES: %d, using default: 3", config.MaxRetries))
		config.MaxRetries = 3
	}

	// 验证超时时间范围
	if config.TimeoutSec < 1 || config.TimeoutSec > 60 {
		SysError(fmt.Sprintf("Invalid SECURITY_WEBHOOK_TIMEOUT: %d, using default: 5", config.TimeoutSec))
		config.TimeoutSec = 5
	}

	return config
}

// isValidWebhookURL 验证 Webhook URL 格式
func isValidWebhookURL(urlStr string) bool {
	if urlStr == "" {
		return false
	}
	// 必须是 http 或 https
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return false
	}
	return true
}

// SetWebhookConfig 设置 Webhook 配置
func SetWebhookConfig(config *WebhookConfig) {
	securityLogMutex.Lock()
	defer securityLogMutex.Unlock()
	webhookConfig = config
}

// GetWebhookConfig 获取当前 Webhook 配置
func GetWebhookConfig() *WebhookConfig {
	securityLogMutex.Lock()
	defer securityLogMutex.Unlock()
	return webhookConfig
}

// InitWebhookNotifier 初始化 Webhook 通知器
func InitWebhookNotifier() {
	config := LoadWebhookConfigFromEnv()
	SetWebhookConfig(config)

	if config.Enabled && config.URL != "" {
		SysLog(fmt.Sprintf("Security webhook notifier enabled, URL: %s", config.URL))
	} else {
		SysLog("Security webhook notifier is disabled")
	}
}

// WebhookPayload Webhook 请求体
type WebhookPayload struct {
	EventType   string `json:"event_type"`
	Timestamp   int64  `json:"timestamp"`
	UserId      int    `json:"user_id"`
	Details     string `json:"details"`
	IPAddress   string `json:"ip_address"`
	RequestPath string `json:"request_path"`
	Severity    string `json:"severity"`
}

// SendSecurityWebhook 发送安全事件 Webhook 通知
func SendSecurityWebhook(event *SecurityEvent) {
	config := GetWebhookConfig()
	if !config.Enabled || config.URL == "" {
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				SysLog(fmt.Sprintf("Panic in webhook notification: %v", r))
			}
		}()

		err := sendWebhookNotification(config, event)
		if err != nil {
			SysLog(fmt.Sprintf("Webhook notification failed: %v", err))
			// 记录失败的通知
			recordFailedWebhook(config.URL, event, err.Error())
		}
	}()
}

// sendWebhookNotification 发送 Webhook 通知（带重试）
func sendWebhookNotification(config *WebhookConfig, event *SecurityEvent) error {
	// 构建请求体
	payload := &WebhookPayload{
		EventType:   event.EventType,
		Timestamp:   event.Timestamp,
		UserId:      event.UserId,
		Details:     event.Details,
		IPAddress:   event.IPAddress,
		RequestPath: event.RequestPath,
		Severity:    getSeverityLevel(event.EventType),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %v", err)
	}

	// 指数退避重试
	var lastErr error
	for i := 0; i < config.MaxRetries; i++ {
		if i > 0 {
			// 指数退避：1s, 2s, 4s...
			delay := time.Duration(1<<uint(i-1)) * time.Second
			time.Sleep(delay)
		}

		err := doWebhookRequest(config.URL, payloadBytes, config.TimeoutSec)
		if err == nil {
			return nil
		}

		lastErr = err
		SysLog(fmt.Sprintf("Webhook notification attempt %d/%d failed: %v", i+1, config.MaxRetries, err))
	}

	return fmt.Errorf("webhook notification failed after %d retries: %v", config.MaxRetries, lastErr)
}

// doWebhookRequest 执行单次 Webhook 请求
func doWebhookRequest(url string, payload []byte, timeoutSec int) error {
	client := &http.Client{
		Timeout: time.Duration(timeoutSec) * time.Second,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "NewAPI-Security-Webhook/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// getSeverityLevel 根据事件类型获取严重级别
func getSeverityLevel(eventType string) string {
	switch eventType {
	case "xss_attempt", "sql_injection", "csrf_attempt":
		return "high"
	case "invalid_origin", "rate_limit_exceeded":
		return "medium"
	default:
		return "low"
	}
}

// recordFailedWebhook 记录失败的 Webhook 通知
func recordFailedWebhook(webhookURL string, event *SecurityEvent, errMsg string) {
	if securityDB == nil {
		return
	}

	// 确保失败通知表存在
	if err := securityDB.AutoMigrate(&FailedWebhookNotification{}); err != nil {
		SysLog(fmt.Sprintf("Failed to create failed webhook table: %v", err))
		return
	}

	payload, _ := json.Marshal(event)
	failedNotification := &FailedWebhookNotification{
		WebhookURL: webhookURL,
		EventID:    event.ID,
		Payload:    string(payload),
		Attempts:   GetWebhookConfig().MaxRetries,
		LastError:  errMsg,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	if err := securityDB.Create(failedNotification).Error; err != nil {
		SysLog(fmt.Sprintf("Failed to record failed webhook notification: %v", err))
	}
}

// GetFailedWebhookNotifications 获取失败的 Webhook 通知列表
func GetFailedWebhookNotifications(page, pageSize int) ([]*FailedWebhookNotification, int64, error) {
	if securityDB == nil {
		return nil, 0, fmt.Errorf("security database not initialized")
	}

	var notifications []*FailedWebhookNotification
	var total int64

	if err := securityDB.Model(&FailedWebhookNotification{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count failed notifications: %v", err)
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	if err := securityDB.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to query failed notifications: %v", err)
	}

	return notifications, total, nil
}

// RetryFailedWebhook 重试失败的 Webhook 通知
func RetryFailedWebhook(notificationID int64) error {
	if securityDB == nil {
		return fmt.Errorf("security database not initialized")
	}

	var notification FailedWebhookNotification
	if err := securityDB.First(&notification, notificationID).Error; err != nil {
		return fmt.Errorf("failed to find notification: %v", err)
	}

	config := GetWebhookConfig()
	if !config.Enabled || config.URL == "" {
		return fmt.Errorf("webhook is not enabled")
	}

	// 解析原始事件
	var event SecurityEvent
	if err := json.Unmarshal([]byte(notification.Payload), &event); err != nil {
		return fmt.Errorf("failed to unmarshal event payload: %v", err)
	}

	// 重试发送
	err := sendWebhookNotification(config, &event)
	if err != nil {
		// 更新失败记录
		notification.Attempts++
		notification.LastError = err.Error()
		notification.UpdatedAt = time.Now().Unix()
		securityDB.Save(&notification)
		return err
	}

	// 发送成功，删除失败记录
	securityDB.Delete(&notification)
	return nil
}

// DeleteFailedWebhookNotification 删除失败的 Webhook 通知记录
func DeleteFailedWebhookNotification(notificationID int64) error {
	if securityDB == nil {
		return fmt.Errorf("security database not initialized")
	}

	result := securityDB.Delete(&FailedWebhookNotification{}, notificationID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete notification: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// CleanupOldFailedWebhooks 清理旧的失败 Webhook 记录
func CleanupOldFailedWebhooks(beforeTimestamp int64) (int64, error) {
	if securityDB == nil {
		return 0, fmt.Errorf("security database not initialized")
	}

	result := securityDB.Where("created_at < ?", beforeTimestamp).Delete(&FailedWebhookNotification{})
	if result.Error != nil {
		return 0, fmt.Errorf("failed to cleanup old failed webhooks: %v", result.Error)
	}

	return result.RowsAffected, nil
}
