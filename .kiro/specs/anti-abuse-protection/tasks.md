# Implementation Plan

## Phase 1: Security Setting Module

- [x] 1. Create Security Setting structure and configuration management
  - [x] 1.1 Create `setting/security_setting/security_setting.go` with SecuritySetting struct
    - Define all configuration fields (channel masking + anti-abuse)
    - Register with GlobalConfig manager
    - Implement GetSecuritySetting() function
    - _Requirements: 1.1, 1.2, 1.3_
  - [ ]* 1.2 Write property test for settings persistence round-trip
    - **Property 1: Settings Persistence Round-Trip**
    - **Validates: Requirements 1.2**
  - [x] 1.3 Migrate channel masking variables from `common/constants.go` to SecuritySetting
    - Update `common/channel_mask.go` to use SecuritySetting
    - Remove environment variable initialization from `common/init.go`
    - _Requirements: 1.4, 1.5_
  - [ ]* 1.4 Write property test for channel masking behavior
    - **Property 2: Channel Masking Behavior**
    - **Validates: Requirements 1.4, 1.5**

- [x] 2. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Phase 2: Core Detection Components

- [x] 3. Implement Model Switch Tracker
  - [x] 3.1 Create `service/model_switch_tracker.go` with ModelSwitchTracker interface and implementation
    - Implement Redis-based storage for model requests
    - Implement memory-based fallback storage
    - Add RecordModelRequest and GetDistinctModelCount methods
    - _Requirements: 2.1, 2.2_
  - [ ]* 3.2 Write property test for model request recording
    - **Property 3: Model Request Recording**
    - **Validates: Requirements 2.1**
  - [ ]* 3.3 Write property test for model switch detection
    - **Property 4: Model Switch Detection and Penalty**
    - **Validates: Requirements 2.2, 2.4**

- [x] 4. Implement Test Content Detector
  - [x] 4.1 Create `service/test_content_detector.go` with TestContentDetector interface and implementation
    - Implement pattern matching for test content
    - Implement content length check
    - Add RecordTestContent and GetTestContentCount methods
    - _Requirements: 3.1, 3.2, 3.4_
  - [ ]* 4.2 Write property test for test content detection
    - **Property 5: Test Content Detection and Counting**
    - **Validates: Requirements 3.1, 3.4**
  - [ ]* 4.3 Write property test for test content threshold flagging
    - **Property 6: Test Content Threshold Flagging**
    - **Validates: Requirements 3.2**

- [x] 5. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Phase 3: Abuse Scoring and Penalty Management

- [x] 6. Implement Abuse Score Calculator
  - [x] 6.1 Create `service/abuse_score.go` with AbuseScoreCalculator
    - Implement score calculation from model switch rate, test content frequency
    - Add configurable weights for each factor
    - _Requirements: 4.1_
  - [ ]* 6.2 Write property test for abuse score calculation and penalty
    - **Property 7: Abuse Score Calculation and Penalty**
    - **Validates: Requirements 4.1, 4.3**

- [x] 7. Implement Penalty Manager
  - [x] 7.1 Create `service/penalty_manager.go` with PenaltyManager interface and implementation
    - Implement ApplyPenalty, CheckPenalty, LiftPenalty methods
    - Add Redis-based penalty storage with TTL
    - _Requirements: 5.1, 5.2_
  - [x] 7.2 Create `model/token_penalty.go` for audit log database model
    - Define TokenPenalty struct with GORM tags
    - Add migration for token_penalties table
    - _Requirements: 5.3_
  - [ ]* 7.3 Write property test for temporary ban auto-expiration
    - **Property 8: Temporary Ban Auto-Expiration**
    - **Validates: Requirements 5.2**
  - [ ]* 7.4 Write property test for penalty audit logging
    - **Property 9: Penalty Audit Logging**
    - **Validates: Requirements 5.3**

- [x] 8. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Phase 4: Abuse Detector and Middleware

- [x] 9. Implement Abuse Detector
  - [x] 9.1 Create `service/abuse_detector.go` with AbuseDetector interface and implementation
    - Integrate ModelSwitchTracker, TestContentDetector, AbuseScoreCalculator
    - Implement CheckRequest method with whitelist check
    - Implement IsWhitelisted method
    - _Requirements: 2.2, 3.2, 4.3, 6.1_
  - [ ]* 9.2 Write property test for whitelist bypass with metrics recording
    - **Property 10: Whitelist Bypass with Metrics Recording**
    - **Validates: Requirements 6.1, 6.3**

- [x] 10. Create Anti-Abuse Middleware
  - [x] 10.1 Create `middleware/anti_abuse.go` with AntiAbuseCheck middleware
    - Check if anti-abuse is enabled
    - Get Token info from context
    - Call AbuseDetector.CheckRequest
    - Apply penalty response if needed
    - _Requirements: 2.4, 4.3_
  - [x] 10.2 Integrate middleware into relay router
    - Add AntiAbuseCheck to relayV1Router in `router/relay-router.go`
    - Position after TokenAuth and before ModelRequestRateLimit
    - _Requirements: 2.4_

- [x] 11. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Phase 5: Admin API and Frontend

- [x] 12. Create Security Setting API
  - [x] 12.1 Create `controller/security_setting.go` with API handlers
    - GET /api/security/settings - Get current settings
    - PUT /api/security/settings - Update settings
    - GET /api/security/penalties - Get active penalties
    - DELETE /api/security/penalties/:token_id - Lift penalty
    - _Requirements: 1.1, 5.4_
  - [x] 12.2 Add routes to `router/api-router.go`
    - Add security settings routes under admin group
    - Apply RootAuth middleware
    - _Requirements: 1.1_

- [x] 13. Create Token Abuse Info API
  - [x] 13.1 Add abuse info to token detail API in `controller/token.go`
    - Add GetTokenAbuseInfo handler
    - Return abuse score and contributing factors
    - _Requirements: 4.4_

- [x] 14. Final Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.
