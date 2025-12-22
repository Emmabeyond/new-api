-- Log Query Optimization - Index Script
-- 日志查询优化 - 索引脚本
-- 
-- 执行前请确认：
-- 1. 已备份数据库
-- 2. 在低峰期执行（索引创建可能锁表）
-- 3. 根据数据量大小，可能需要几分钟到几十分钟
--
-- 注意：以下索引已通过 GORM 自动迁移创建：
-- - idx_created_at_id (created_at, id) - 时间范围查询优化
-- - idx_created_at_type (type, created_at) - 类型+时间查询优化
-- - idx_user_created (user_id, created_at) - 用户+时间查询优化
-- - idx_username_created (username, created_at) - 用户名+时间查询优化
-- - idx_channel_created (channel_id, created_at) - 渠道+时间查询优化

-- ============================================
-- MySQL 版本
-- ============================================

-- 复合索引：user_id + type + created_at（用于用户按类型查询日志）
-- 适用场景：用户查询自己的消费/充值/系统日志
SET @exist := (SELECT COUNT(*) FROM information_schema.statistics 
               WHERE table_schema = DATABASE() 
               AND table_name = 'logs' 
               AND index_name = 'idx_user_type_created');
SET @sqlstmt := IF(@exist = 0, 
    'CREATE INDEX idx_user_type_created ON logs(user_id, type, created_at)', 
    'SELECT ''Index idx_user_type_created already exists''');
PREPARE stmt FROM @sqlstmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 复合索引：type + created_at（如果 GORM 未自动创建）
-- 适用场景：管理员按类型查询全局日志
SET @exist := (SELECT COUNT(*) FROM information_schema.statistics 
               WHERE table_schema = DATABASE() 
               AND table_name = 'logs' 
               AND index_name = 'idx_type_created');
SET @sqlstmt := IF(@exist = 0, 
    'CREATE INDEX idx_type_created ON logs(type, created_at)', 
    'SELECT ''Index idx_type_created already exists''');
PREPARE stmt FROM @sqlstmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ============================================
-- 验证索引创建结果 (MySQL)
-- ============================================
SELECT 
    index_name,
    column_name,
    seq_in_index
FROM information_schema.statistics 
WHERE table_schema = DATABASE() 
AND table_name = 'logs'
ORDER BY index_name, seq_in_index;

-- ============================================
-- PostgreSQL 版本（如果使用 PostgreSQL）
-- ============================================
-- 
-- -- 复合索引：user_id + type + created_at
-- CREATE INDEX IF NOT EXISTS idx_user_type_created ON logs(user_id, type, created_at);
-- 
-- -- 复合索引：type + created_at
-- CREATE INDEX IF NOT EXISTS idx_type_created ON logs(type, created_at);
-- 
-- -- 验证索引
-- SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'logs';

-- ============================================
-- SQLite 版本（如果使用 SQLite）
-- ============================================
-- 
-- -- 复合索引：user_id + type + created_at
-- CREATE INDEX IF NOT EXISTS idx_user_type_created ON logs(user_id, type, created_at);
-- 
-- -- 复合索引：type + created_at
-- CREATE INDEX IF NOT EXISTS idx_type_created ON logs(type, created_at);
-- 
-- -- 验证索引
-- SELECT name FROM sqlite_master WHERE type='index' AND tbl_name='logs';

-- ============================================
-- 索引使用说明
-- ============================================
-- 
-- 1. idx_user_type_created (user_id, type, created_at)
--    - 用途：用户按类型查询自己的日志
--    - 查询示例：SELECT * FROM logs WHERE user_id = ? AND type = ? AND created_at BETWEEN ? AND ?
--    - 优化效果：避免全表扫描，直接定位到特定用户的特定类型日志
--
-- 2. idx_type_created (type, created_at)
--    - 用途：管理员按类型查询全局日志
--    - 查询示例：SELECT * FROM logs WHERE type = ? AND created_at BETWEEN ? AND ?
--    - 优化效果：快速筛选特定类型的日志记录
--
-- 注意事项：
-- - 索引会占用额外存储空间
-- - 索引会略微降低写入性能
-- - 建议在数据量超过 100 万条时创建这些索引
-- - 可以使用 EXPLAIN 分析查询是否使用了索引
