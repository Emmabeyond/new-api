-- Dashboard 性能优化索引
-- 用于优化 quota_data 表的查询性能
-- 执行前请备份数据库

-- 1. 复合索引: (created_at, model_name) 用于 GROUP BY 聚合查询优化
-- 适用于管理员查询全局统计数据的场景
CREATE INDEX idx_qdt_time_model ON quota_data (created_at, model_name);

-- 2. 复合索引: (user_id, created_at) 用于用户查询优化
-- 适用于普通用户查询自己的统计数据的场景
CREATE INDEX idx_qdt_user_time ON quota_data (user_id, created_at);

-- MySQL 版本 (如果上述语法不兼容)
-- CREATE INDEX idx_qdt_time_model ON quota_data (created_at, model_name);
-- CREATE INDEX idx_qdt_user_time ON quota_data (user_id, created_at);

-- PostgreSQL 版本
-- CREATE INDEX idx_qdt_time_model ON quota_data (created_at, model_name);
-- CREATE INDEX idx_qdt_user_time ON quota_data (user_id, created_at);

-- 验证索引是否生效 (MySQL)
-- EXPLAIN SELECT model_name, sum(count) as count, sum(quota) as quota, sum(token_used) as token_used, created_at 
-- FROM quota_data 
-- WHERE created_at >= UNIX_TIMESTAMP() - 86400 AND created_at <= UNIX_TIMESTAMP() 
-- GROUP BY model_name, created_at;

-- 验证索引是否生效 (PostgreSQL)
-- EXPLAIN ANALYZE SELECT model_name, sum(count) as count, sum(quota) as quota, sum(token_used) as token_used, created_at 
-- FROM quota_data 
-- WHERE created_at >= EXTRACT(EPOCH FROM NOW()) - 86400 AND created_at <= EXTRACT(EPOCH FROM NOW()) 
-- GROUP BY model_name, created_at;
