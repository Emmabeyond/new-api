-- Channel Query Optimization - MySQL Index Script
-- 渠道查询优化 - MySQL 索引脚本
-- 
-- 执行前请确认：
-- 1. 已备份数据库
-- 2. 在低峰期执行（索引创建可能锁表）
-- 3. 根据数据量大小，可能需要几分钟到几十分钟

-- ============================================
-- 单列索引（如果 GORM 自动迁移未创建）
-- ============================================

-- 检查并创建 status 索引
SET @exist := (SELECT COUNT(*) FROM information_schema.statistics 
               WHERE table_schema = DATABASE() 
               AND table_name = 'channels' 
               AND index_name = 'idx_channels_status');
SET @sqlstmt := IF(@exist = 0, 
    'CREATE INDEX idx_channels_status ON channels(status)', 
    'SELECT ''Index idx_channels_status already exists''');
PREPARE stmt FROM @sqlstmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并创建 type 索引
SET @exist := (SELECT COUNT(*) FROM information_schema.statistics 
               WHERE table_schema = DATABASE() 
               AND table_name = 'channels' 
               AND index_name = 'idx_channels_type');
SET @sqlstmt := IF(@exist = 0, 
    'CREATE INDEX idx_channels_type ON channels(type)', 
    'SELECT ''Index idx_channels_type already exists''');
PREPARE stmt FROM @sqlstmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并创建 priority 索引
SET @exist := (SELECT COUNT(*) FROM information_schema.statistics 
               WHERE table_schema = DATABASE() 
               AND table_name = 'channels' 
               AND index_name = 'idx_channels_priority');
SET @sqlstmt := IF(@exist = 0, 
    'CREATE INDEX idx_channels_priority ON channels(priority DESC)', 
    'SELECT ''Index idx_channels_priority already exists''');
PREPARE stmt FROM @sqlstmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 检查并创建 group 索引
SET @exist := (SELECT COUNT(*) FROM information_schema.statistics 
               WHERE table_schema = DATABASE() 
               AND table_name = 'channels' 
               AND index_name = 'idx_channels_group');
SET @sqlstmt := IF(@exist = 0, 
    'CREATE INDEX idx_channels_group ON channels(`group`)', 
    'SELECT ''Index idx_channels_group already exists''');
PREPARE stmt FROM @sqlstmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ============================================
-- 复合索引（核心优化）
-- ============================================

-- 复合索引：status + type + priority（用于常见筛选组合）
SET @exist := (SELECT COUNT(*) FROM information_schema.statistics 
               WHERE table_schema = DATABASE() 
               AND table_name = 'channels' 
               AND index_name = 'idx_channels_status_type_priority');
SET @sqlstmt := IF(@exist = 0, 
    'CREATE INDEX idx_channels_status_type_priority ON channels(status, type, priority DESC)', 
    'SELECT ''Index idx_channels_status_type_priority already exists''');
PREPARE stmt FROM @sqlstmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 复合索引：tag + status（用于标签模式查询）
SET @exist := (SELECT COUNT(*) FROM information_schema.statistics 
               WHERE table_schema = DATABASE() 
               AND table_name = 'channels' 
               AND index_name = 'idx_channels_tag_status');
SET @sqlstmt := IF(@exist = 0, 
    'CREATE INDEX idx_channels_tag_status ON channels(tag, status)', 
    'SELECT ''Index idx_channels_tag_status already exists''');
PREPARE stmt FROM @sqlstmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ============================================
-- 验证索引创建结果
-- ============================================
SELECT 
    index_name,
    column_name,
    seq_in_index
FROM information_schema.statistics 
WHERE table_schema = DATABASE() 
AND table_name = 'channels'
ORDER BY index_name, seq_in_index;
