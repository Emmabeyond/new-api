package model

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var commonGroupCol string
var commonKeyCol string
var commonTrueVal string
var commonFalseVal string

var logKeyCol string
var logGroupCol string

func initCol() {
	// init common column names
	if common.UsingPostgreSQL {
		commonGroupCol = `"group"`
		commonKeyCol = `"key"`
		commonTrueVal = "true"
		commonFalseVal = "false"
	} else {
		commonGroupCol = "`group`"
		commonKeyCol = "`key`"
		commonTrueVal = "1"
		commonFalseVal = "0"
	}
	if os.Getenv("LOG_SQL_DSN") != "" {
		switch common.LogSqlType {
		case common.DatabaseTypePostgreSQL:
			logGroupCol = `"group"`
			logKeyCol = `"key"`
		default:
			logGroupCol = commonGroupCol
			logKeyCol = commonKeyCol
		}
	} else {
		// LOG_SQL_DSN 为空时，日志数据库与主数据库相同
		if common.UsingPostgreSQL {
			logGroupCol = `"group"`
			logKeyCol = `"key"`
		} else {
			logGroupCol = commonGroupCol
			logKeyCol = commonKeyCol
		}
	}
	// log sql type and database type
	//common.SysLog("Using Log SQL Type: " + common.LogSqlType)
}

var DB *gorm.DB

var LOG_DB *gorm.DB

func createRootAccountIfNeed() error {
	var user User
	//if user.Status != common.UserStatusEnabled {
	if err := DB.First(&user).Error; err != nil {
		common.SysLog("no user exists, create a root user for you: username is root, password is 123456")
		hashedPassword, err := common.Password2Hash("123456")
		if err != nil {
			return err
		}
		rootUser := User{
			Username:    "root",
			Password:    hashedPassword,
			Role:        common.RoleRootUser,
			Status:      common.UserStatusEnabled,
			DisplayName: "Root User",
			AccessToken: nil,
			Quota:       100000000,
		}
		DB.Create(&rootUser)
	}
	return nil
}

func CheckSetup() {
	setup := GetSetup()
	if setup == nil {
		// No setup record exists, check if we have a root user
		if RootUserExists() {
			common.SysLog("system is not initialized, but root user exists")
			// Create setup record
			newSetup := Setup{
				Version:       common.Version,
				InitializedAt: time.Now().Unix(),
			}
			err := DB.Create(&newSetup).Error
			if err != nil {
				common.SysLog("failed to create setup record: " + err.Error())
			}
			constant.Setup = true
		} else {
			common.SysLog("system is not initialized and no root user exists")
			constant.Setup = false
		}
	} else {
		// Setup record exists, system is initialized
		common.SysLog("system is already initialized at: " + time.Unix(setup.InitializedAt, 0).String())
		constant.Setup = true
	}
}

func chooseDB(envName string, isLog bool) (*gorm.DB, error) {
	defer func() {
		initCol()
	}()
	dsn := os.Getenv(envName)
	if dsn != "" {
		if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
			// Use PostgreSQL
			common.SysLog("using PostgreSQL as database")
			if !isLog {
				common.UsingPostgreSQL = true
			} else {
				common.LogSqlType = common.DatabaseTypePostgreSQL
			}
			return gorm.Open(postgres.New(postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true, // disables implicit prepared statement usage
			}), &gorm.Config{
				PrepareStmt: true, // precompile SQL
			})
		}
		if strings.HasPrefix(dsn, "local") {
			common.SysLog("SQL_DSN not set, using SQLite as database")
			if !isLog {
				common.UsingSQLite = true
			} else {
				common.LogSqlType = common.DatabaseTypeSQLite
			}
			return gorm.Open(sqlite.Open(common.SQLitePath), &gorm.Config{
				PrepareStmt: true, // precompile SQL
			})
		}
		// Use MySQL
		common.SysLog("using MySQL as database")
		// check parseTime
		if !strings.Contains(dsn, "parseTime") {
			if strings.Contains(dsn, "?") {
				dsn += "&parseTime=true"
			} else {
				dsn += "?parseTime=true"
			}
		}
		if !isLog {
			common.UsingMySQL = true
		} else {
			common.LogSqlType = common.DatabaseTypeMySQL
		}
		return gorm.Open(mysql.Open(dsn), &gorm.Config{
			PrepareStmt: true, // precompile SQL
		})
	}
	// Use SQLite
	common.SysLog("SQL_DSN not set, using SQLite as database")
	common.UsingSQLite = true
	return gorm.Open(sqlite.Open(common.SQLitePath), &gorm.Config{
		PrepareStmt: true, // precompile SQL
	})
}

func InitDB() (err error) {
	db, err := chooseDB("SQL_DSN", false)
	if err == nil {
		if common.DebugEnabled {
			db = db.Debug()
		}
		DB = db
		// MySQL charset/collation startup check: ensure Chinese-capable charset
		if common.UsingMySQL {
			if err := checkMySQLChineseSupport(DB); err != nil {
				panic(err)
			}
		}
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		sqlDB.SetMaxIdleConns(common.GetEnvOrDefault("SQL_MAX_IDLE_CONNS", 50))
		sqlDB.SetMaxOpenConns(common.GetEnvOrDefault("SQL_MAX_OPEN_CONNS", 200))
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(common.GetEnvOrDefault("SQL_MAX_LIFETIME", 300)))
		sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(common.GetEnvOrDefault("SQL_MAX_IDLE_TIME", 60)))

		if !common.IsMasterNode {
			return nil
		}
		if common.UsingMySQL {
			//_, _ = sqlDB.Exec("ALTER TABLE channels MODIFY model_mapping TEXT;") // TODO: delete this line when most users have upgraded
		}
		common.SysLog("database migration started")
		err = migrateDB()
		if err != nil {
			return err
		}

		// 初始化 Token 多级缓存服务
		InitTokenCacheService()

		return nil
	} else {
		common.FatalLog(err)
	}
	return err
}

func InitLogDB() (err error) {
	if os.Getenv("LOG_SQL_DSN") == "" {
		LOG_DB = DB
		return
	}
	db, err := chooseDB("LOG_SQL_DSN", true)
	if err == nil {
		if common.DebugEnabled {
			db = db.Debug()
		}
		LOG_DB = db
		// If log DB is MySQL, also ensure Chinese-capable charset
		if common.LogSqlType == common.DatabaseTypeMySQL {
			if err := checkMySQLChineseSupport(LOG_DB); err != nil {
				panic(err)
			}
		}
		sqlDB, err := LOG_DB.DB()
		if err != nil {
			return err
		}
		sqlDB.SetMaxIdleConns(common.GetEnvOrDefault("LOG_SQL_MAX_IDLE_CONNS", 50))
		sqlDB.SetMaxOpenConns(common.GetEnvOrDefault("LOG_SQL_MAX_OPEN_CONNS", 200))
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(common.GetEnvOrDefault("LOG_SQL_MAX_LIFETIME", 300)))
		sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(common.GetEnvOrDefault("LOG_SQL_MAX_IDLE_TIME", 60)))

		if !common.IsMasterNode {
			return nil
		}
		common.SysLog("database migration started")
		err = migrateLOGDB()
		return err
	} else {
		common.FatalLog(err)
	}
	return err
}

func migrateDB() error {
	err := DB.AutoMigrate(
		&Channel{},
		&Token{},
		&User{},
		&PasskeyCredential{},
		&Option{},
		&Redemption{},
		&Ability{},
		&Log{},
		&Midjourney{},
		&TopUp{},
		&QuotaData{},
		&Task{},
		&Model{},
		&Vendor{},
		&PrefillGroup{},
		&Setup{},
		&TwoFA{},
		&TwoFABackupCode{},
		&Checkin{},
		&CheckinRecord{},
		&CheckinAudit{},
		&TokenPenalty{},
		&HelpCategory{},
		&HelpDocument{},
		&LevelConfig{},
		&LevelChangeLog{},
		&common.SecurityEvent{}, // 安全事件表
		&GuestbookMessage{},     // 留言板
	)
	if err != nil {
		return err
	}
	// Create composite indexes for channel queries optimization
	createChannelIndexes()
	// Create composite indexes for guestbook queries optimization
	createGuestbookIndexes()
	// Create composite indexes for log queries optimization (when LOG_SQL_DSN is not set, use DB directly)
	if os.Getenv("LOG_SQL_DSN") == "" {
		createLogIndexesWithDB(DB)
		// Also migrate LogStatsAggregation when using same DB
		if err := DB.AutoMigrate(&LogStatsAggregation{}); err != nil {
			common.SysLog("Warning: failed to migrate LogStatsAggregation: " + err.Error())
		}
	}
	// Create indexes for security_events table
	createSecurityEventIndexes()
	// Initialize security logger with database connection
	common.SetSecurityDB(DB)
	// Initialize default level configs
	if err := InitDefaultLevelConfigs(); err != nil {
		common.SysLog("Warning: failed to initialize default level configs: " + err.Error())
	}
	// 迁移：为空等级用户设置默认等级 tier_1
	migrateEmptyUserLevels()
	// 初始化默认留言板数据
	if err := InitDefaultGuestbookMessages(); err != nil {
		common.SysLog("Warning: failed to initialize default guestbook messages: " + err.Error())
	}
	// 初始化默认 AI 工具帮助文档
	if err := InitDefaultAIToolsHelpDocuments(); err != nil {
		common.SysLog("Warning: failed to initialize default AI tools help documents: " + err.Error())
	}
	// 初始化默认 FAQ
	if err := InitDefaultFAQ(); err != nil {
		common.SysLog("Warning: failed to initialize default FAQ: " + err.Error())
	}
	return nil
}

// createChannelIndexes creates composite indexes for optimizing channel queries
func createChannelIndexes() {
	var err error
	// Composite index for common filter combinations (status + type + priority)
	if common.UsingPostgreSQL {
		err = DB.Exec("CREATE INDEX IF NOT EXISTS idx_channels_status_type_priority ON channels(status, type, priority DESC)").Error
	} else if common.UsingMySQL {
		// MySQL doesn't support IF NOT EXISTS for indexes, check first
		var count int64
		DB.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'channels' AND index_name = 'idx_channels_status_type_priority'").Scan(&count)
		if count == 0 {
			err = DB.Exec("CREATE INDEX idx_channels_status_type_priority ON channels(status, type, priority DESC)").Error
		}
	} else {
		// SQLite
		err = DB.Exec("CREATE INDEX IF NOT EXISTS idx_channels_status_type_priority ON channels(status, type, priority DESC)").Error
	}
	if err != nil {
		common.SysLog(fmt.Sprintf("Warning: failed to create idx_channels_status_type_priority: %v", err))
	}

	// Composite index for tag mode queries (tag + status)
	if common.UsingPostgreSQL {
		err = DB.Exec("CREATE INDEX IF NOT EXISTS idx_channels_tag_status ON channels(tag, status)").Error
	} else if common.UsingMySQL {
		var count int64
		DB.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'channels' AND index_name = 'idx_channels_tag_status'").Scan(&count)
		if count == 0 {
			err = DB.Exec("CREATE INDEX idx_channels_tag_status ON channels(tag, status)").Error
		}
	} else {
		err = DB.Exec("CREATE INDEX IF NOT EXISTS idx_channels_tag_status ON channels(tag, status)").Error
	}
	if err != nil {
		common.SysLog(fmt.Sprintf("Warning: failed to create idx_channels_tag_status: %v", err))
	}

	common.SysLog("Channel composite indexes created/verified")
}

// createSecurityEventIndexes creates indexes for security_events table
func createSecurityEventIndexes() {
	var err error

	// Index definitions for security_events table
	indexes := []struct {
		name    string
		columns string
	}{
		{"idx_security_events_timestamp", "timestamp DESC"},
		{"idx_security_events_user_id", "user_id"},
		{"idx_security_events_event_type", "event_type"},
		{"idx_security_events_user_type_time", "user_id, event_type, timestamp DESC"},
	}

	for _, idx := range indexes {
		if common.UsingPostgreSQL {
			err = DB.Exec(fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON security_events(%s)", idx.name, idx.columns)).Error
		} else if common.UsingMySQL {
			var count int64
			DB.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'security_events' AND index_name = ?", idx.name).Scan(&count)
			if count == 0 {
				err = DB.Exec(fmt.Sprintf("CREATE INDEX %s ON security_events(%s)", idx.name, idx.columns)).Error
			}
		} else {
			// SQLite
			err = DB.Exec(fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON security_events(%s)", idx.name, idx.columns)).Error
		}
		if err != nil {
			common.SysLog(fmt.Sprintf("Warning: failed to create %s: %v", idx.name, err))
		}
	}

	common.SysLog("Security event indexes created/verified")
}

func migrateDBFast() error {

	var wg sync.WaitGroup

	migrations := []struct {
		model interface{}
		name  string
	}{
		{&Channel{}, "Channel"},
		{&Token{}, "Token"},
		{&User{}, "User"},
		{&PasskeyCredential{}, "PasskeyCredential"},
		{&Option{}, "Option"},
		{&Redemption{}, "Redemption"},
		{&Ability{}, "Ability"},
		{&Log{}, "Log"},
		{&Midjourney{}, "Midjourney"},
		{&TopUp{}, "TopUp"},
		{&QuotaData{}, "QuotaData"},
		{&Task{}, "Task"},
		{&Model{}, "Model"},
		{&Vendor{}, "Vendor"},
		{&PrefillGroup{}, "PrefillGroup"},
		{&Setup{}, "Setup"},
		{&TwoFA{}, "TwoFA"},
		{&TwoFABackupCode{}, "TwoFABackupCode"},
		{&Checkin{}, "Checkin"},
		{&CheckinRecord{}, "CheckinRecord"},
		{&CheckinAudit{}, "CheckinAudit"},
		{&TokenPenalty{}, "TokenPenalty"},
		{&LevelConfig{}, "LevelConfig"},
		{&LevelChangeLog{}, "LevelChangeLog"},
	}
	// 动态计算migration数量，确保errChan缓冲区足够大
	errChan := make(chan error, len(migrations))

	for _, m := range migrations {
		wg.Add(1)
		go func(model interface{}, name string) {
			defer wg.Done()
			if err := DB.AutoMigrate(model); err != nil {
				errChan <- fmt.Errorf("failed to migrate %s: %v", name, err)
			}
		}(m.model, m.name)
	}

	// Wait for all migrations to complete
	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	common.SysLog("database migrated")
	return nil
}

func migrateLOGDB() error {
	var err error
	if err = LOG_DB.AutoMigrate(&Log{}); err != nil {
		return err
	}
	// Migrate log stats aggregation table
	if err = LOG_DB.AutoMigrate(&LogStatsAggregation{}); err != nil {
		return err
	}
	// Create composite indexes for log queries optimization
	createLogIndexesWithDB(LOG_DB)
	return nil
}

// createLogIndexesWithDB creates composite indexes for optimizing log queries using the provided db instance
// These indexes support common query patterns: user-based, username-based, channel-based, and type-based filtering
func createLogIndexesWithDB(db *gorm.DB) {
	if db == nil {
		common.SysLog("Warning: cannot create log indexes, db is nil")
		return
	}

	var err error

	// Determine which database type is being used
	isPostgreSQL := common.LogSqlType == common.DatabaseTypePostgreSQL
	isMySQL := common.LogSqlType == common.DatabaseTypeMySQL

	// If LOG_SQL_DSN is not set, use main DB type flags
	if os.Getenv("LOG_SQL_DSN") == "" {
		isPostgreSQL = common.UsingPostgreSQL
		isMySQL = common.UsingMySQL
	}

	// Index definitions: (index_name, columns)
	indexes := []struct {
		name    string
		columns string
	}{
		{"idx_logs_type_created", "type, created_at DESC"},
		{"idx_logs_user_created", "user_id, created_at DESC"},
		{"idx_logs_username_created", "username, created_at DESC"},
		{"idx_logs_channel_created", "channel_id, created_at DESC"},
	}

	for _, idx := range indexes {
		if isPostgreSQL {
			err = db.Exec(fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON logs(%s)", idx.name, idx.columns)).Error
		} else if isMySQL {
			// MySQL doesn't support IF NOT EXISTS for indexes, check first
			var count int64
			db.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'logs' AND index_name = ?", idx.name).Scan(&count)
			if count == 0 {
				err = db.Exec(fmt.Sprintf("CREATE INDEX %s ON logs(%s)", idx.name, idx.columns)).Error
			}
		} else {
			// SQLite
			err = db.Exec(fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON logs(%s)", idx.name, idx.columns)).Error
		}
		if err != nil {
			common.SysLog(fmt.Sprintf("Warning: failed to create %s: %v", idx.name, err))
		}
	}

	common.SysLog("Log composite indexes created/verified")
}

func closeDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	return err
}

func CloseDB() error {
	if LOG_DB != DB {
		err := closeDB(LOG_DB)
		if err != nil {
			return err
		}
	}
	return closeDB(DB)
}

// checkMySQLChineseSupport ensures the MySQL connection and current schema
// default charset/collation can store Chinese characters. It allows common
// Chinese-capable charsets (utf8mb4, utf8, gbk, big5, gb18030) and panics otherwise.
func checkMySQLChineseSupport(db *gorm.DB) error {
	// 仅检测：当前库默认字符集/排序规则 + 各表的排序规则（隐含字符集）

	// Read current schema defaults
	var schemaCharset, schemaCollation string
	err := db.Raw("SELECT DEFAULT_CHARACTER_SET_NAME, DEFAULT_COLLATION_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = DATABASE()").Row().Scan(&schemaCharset, &schemaCollation)
	if err != nil {
		return fmt.Errorf("读取当前库默认字符集/排序规则失败 / Failed to read schema default charset/collation: %v", err)
	}

	toLower := func(s string) string { return strings.ToLower(s) }
	// Allowed charsets that can store Chinese text
	allowedCharsets := map[string]string{
		"utf8mb4": "utf8mb4_",
		"utf8":    "utf8_",
		"gbk":     "gbk_",
		"big5":    "big5_",
		"gb18030": "gb18030_",
	}
	isChineseCapable := func(cs, cl string) bool {
		csLower := toLower(cs)
		clLower := toLower(cl)
		if prefix, ok := allowedCharsets[csLower]; ok {
			if clLower == "" {
				return true
			}
			return strings.HasPrefix(clLower, prefix)
		}
		// 如果仅提供了排序规则，尝试按排序规则前缀判断
		for _, prefix := range allowedCharsets {
			if strings.HasPrefix(clLower, prefix) {
				return true
			}
		}
		return false
	}

	// 1) 当前库默认值必须支持中文
	if !isChineseCapable(schemaCharset, schemaCollation) {
		return fmt.Errorf("当前库默认字符集/排序规则不支持中文：schema(%s/%s)。请将库设置为 utf8mb4/utf8/gbk/big5/gb18030 / Schema default charset/collation is not Chinese-capable: schema(%s/%s). Please set to utf8mb4/utf8/gbk/big5/gb18030",
			schemaCharset, schemaCollation, schemaCharset, schemaCollation)
	}

	// 2) 所有物理表的排序规则（隐含字符集）必须支持中文
	type tableInfo struct {
		Name      string
		Collation *string
	}
	var tables []tableInfo
	if err := db.Raw("SELECT TABLE_NAME, TABLE_COLLATION FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_TYPE = 'BASE TABLE'").Scan(&tables).Error; err != nil {
		return fmt.Errorf("读取表排序规则失败 / Failed to read table collations: %v", err)
	}

	var badTables []string
	for _, t := range tables {
		// NULL 或空表示继承库默认设置，已在上面校验库默认，视为通过
		if t.Collation == nil || *t.Collation == "" {
			continue
		}
		cl := *t.Collation
		// 仅凭排序规则判断是否中文可用
		ok := false
		lower := strings.ToLower(cl)
		for _, prefix := range allowedCharsets {
			if strings.HasPrefix(lower, prefix) {
				ok = true
				break
			}
		}
		if !ok {
			badTables = append(badTables, fmt.Sprintf("%s(%s)", t.Name, cl))
		}
	}

	if len(badTables) > 0 {
		// 限制输出数量以避免日志过长
		maxShow := 20
		shown := badTables
		if len(shown) > maxShow {
			shown = shown[:maxShow]
		}
		return fmt.Errorf(
			"存在不支持中文的表，请修复其排序规则/字符集。示例（最多展示 %d 项）：%v / Found tables not Chinese-capable. Please fix their collation/charset. Examples (showing up to %d): %v",
			maxShow, shown, maxShow, shown,
		)
	}
	return nil
}

var (
	lastPingTime time.Time
	pingMutex    sync.Mutex
)

func PingDB() error {
	pingMutex.Lock()
	defer pingMutex.Unlock()

	if time.Since(lastPingTime) < time.Second*10 {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error getting sql.DB from GORM: %v", err)
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Printf("Error pinging DB: %v", err)
		return err
	}

	lastPingTime = time.Now()
	common.SysLog("Database pinged successfully")
	return nil
}


// migrateEmptyUserLevels 为空等级用户设置默认等级 tier_1
// 同时修复历史数据中等级名称与ID不一致的问题
func migrateEmptyUserLevels() {
	// 1. 修复空等级
	result := DB.Model(&User{}).Where("level = '' OR level IS NULL").Update("level", "tier_1")
	if result.Error != nil {
		common.SysLog("Warning: failed to migrate empty user levels: " + result.Error.Error())
	} else if result.RowsAffected > 0 {
		common.SysLog(fmt.Sprintf("Migrated %d users with empty level to tier_1", result.RowsAffected))
	}

	// 2. 修复等级名称与ID不一致的问题（历史数据可能存的是名称而非ID）
	levelNameToId := map[string]string{
		"Tier 1": "tier_1",
		"Tier 2": "tier_2",
		"Tier 3": "tier_3",
	}
	for name, id := range levelNameToId {
		result := DB.Model(&User{}).Where("level = ?", name).Update("level", id)
		if result.Error != nil {
			common.SysLog(fmt.Sprintf("Warning: failed to migrate level '%s' to '%s': %s", name, id, result.Error.Error()))
		} else if result.RowsAffected > 0 {
			common.SysLog(fmt.Sprintf("Migrated %d users from level '%s' to '%s'", result.RowsAffected, name, id))
		}
	}
}
