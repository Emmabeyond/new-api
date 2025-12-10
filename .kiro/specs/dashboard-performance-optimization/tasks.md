# Implementation Plan

- [x] 1. 添加时间范围验证器
  - [x] 1.1 创建 `common/time_validator.go` 文件
    - 实现 `ValidateTimeRange` 函数
    - 定义 `MaxQueryDays` 常量 (90 天)
    - 定义错误类型 `ErrInvalidTimeRange`, `ErrInvalidTimestamp`
    - _Requirements: 5.1, 5.2_
  - [ ]* 1.2 编写属性测试 - 时间范围限制
    - **Property 4: 时间范围限制**
    - **Validates: Requirements 5.1**
  - [ ]* 1.3 编写属性测试 - 无效时间范围处理
    - **Property 5: 无效时间范围处理**
    - **Validates: Requirements 5.2**

- [x] 2. 实现 QuotaData 缓存服务
  - [x] 2.1 创建 `model/usedata_cache.go` 文件
    - 实现 `QuotaDataCacheKey` 函数生成缓存 key
    - 实现 `GetQuotaDataWithCache` 带缓存查询函数
    - 实现 `InvalidateQuotaDataCache` 缓存失效函数
    - 定义缓存 TTL 常量 (5 分钟)
    - _Requirements: 3.1, 3.2, 3.3, 3.4_
  - [ ]* 2.2 编写属性测试 - 缓存一致性
    - **Property 2: 缓存一致性**
    - **Validates: Requirements 3.1, 3.2, 3.3**
  - [ ]* 2.3 编写属性测试 - 缓存失效正确性
    - **Property 3: 缓存失效正确性**
    - **Validates: Requirements 3.4**

- [x] 3. 添加数据库复合索引
  - [x] 3.1 修改 `model/usedata.go` 添加索引定义
    - 添加复合索引 `idx_quota_data_time_model` (created_at, model_name)
    - 添加复合索引 `idx_quota_data_user_time` (user_id, created_at)
    - _Requirements: 4.1, 4.2, 4.3_
  - [x] 3.2 创建索引迁移 SQL 文件
    - 创建 `docs/sql/quota_data_indexes.sql`
    - 包含手动执行索引的 SQL 语句
    - _Requirements: 4.1, 4.2, 4.3_

- [x] 4. 集成到 Controller 层
  - [x] 4.1 修改 `controller/usedata.go`
    - 在 `GetAllQuotaDates` 中集成时间范围验证
    - 在 `GetAllQuotaDates` 中集成缓存查询
    - 在 `GetUserQuotaDates` 中集成时间范围验证
    - 在 `GetUserQuotaDates` 中集成缓存查询
    - _Requirements: 1.1, 2.1, 3.1_
  - [ ]* 4.2 编写属性测试 - 用户数据隔离
    - **Property 1: 用户数据隔离**
    - **Validates: Requirements 2.2**

- [x] 5. 缓存失效集成
  - [x] 5.1 修改 `model/usedata.go` 的 `SaveQuotaDataCache` 函数
    - 在保存数据后调用缓存失效
    - _Requirements: 3.4_

- [x] 6. Checkpoint - 确保所有测试通过
  - Ensure all tests pass, ask the user if questions arise.
