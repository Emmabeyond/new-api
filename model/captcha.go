package model

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
)

// CaptchaChallenge 验证码挑战数据
type CaptchaChallenge struct {
	SessionID string    `json:"session_id"`
	TargetX   int       `json:"target_x"`
	PuzzleY   int       `json:"puzzle_y"`
	ImageIdx  int       `json:"image_idx"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Verified  bool      `json:"verified"`
}

// CaptchaToken 验证成功后生成的令牌
type CaptchaToken struct {
	Token     string    `json:"token"`
	UserIP    string    `json:"user_ip"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
}

// CaptchaIPRecord IP 限流记录
type CaptchaIPRecord struct {
	FailedAttempts int       `json:"failed_attempts"`
	FirstFailTime  time.Time `json:"first_fail_time"`
	BlockedUntil   time.Time `json:"blocked_until"`
}

// 内存存储（当 Redis 不可用时使用）
var (
	captchaChallengeStore = make(map[string]*CaptchaChallenge)
	captchaTokenStore     = make(map[string]*CaptchaToken)
	captchaIPStore        = make(map[string]*CaptchaIPRecord)
	captchaStoreMutex     sync.RWMutex
)

// 常量定义
const (
	CaptchaChallengeExpiration = 5 * time.Minute
	CaptchaTokenExpiration     = 10 * time.Minute
	CaptchaIPBlockDuration     = 5 * time.Minute
	CaptchaIPWindowDuration    = 1 * time.Minute
	CaptchaMaxFailedAttempts   = 5
)

// GenerateSessionID 生成唯一的会话 ID
func GenerateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateCaptchaToken 生成验证令牌
func GenerateCaptchaToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// StoreCaptchaChallenge 存储验证码挑战
func StoreCaptchaChallenge(challenge *CaptchaChallenge) error {
	if common.RedisEnabled {
		return storeCaptchaChallengeRedis(challenge)
	}
	return storeCaptchaChallengeMemory(challenge)
}

// GetCaptchaChallenge 获取验证码挑战
func GetCaptchaChallenge(sessionID string) (*CaptchaChallenge, error) {
	if common.RedisEnabled {
		return getCaptchaChallengeRedis(sessionID)
	}
	return getCaptchaChallengeMemory(sessionID)
}

// DeleteCaptchaChallenge 删除验证码挑战
func DeleteCaptchaChallenge(sessionID string) error {
	if common.RedisEnabled {
		return deleteCaptchaChallengeRedis(sessionID)
	}
	return deleteCaptchaChallengeMemory(sessionID)
}

// StoreCaptchaToken 存储验证令牌
func StoreCaptchaToken(token *CaptchaToken) error {
	if common.RedisEnabled {
		return storeCaptchaTokenRedis(token)
	}
	return storeCaptchaTokenMemory(token)
}

// GetCaptchaToken 获取验证令牌
func GetCaptchaToken(token string) (*CaptchaToken, error) {
	if common.RedisEnabled {
		return getCaptchaTokenRedis(token)
	}
	return getCaptchaTokenMemory(token)
}

// MarkCaptchaTokenUsed 标记令牌已使用
func MarkCaptchaTokenUsed(token string) error {
	if common.RedisEnabled {
		return markCaptchaTokenUsedRedis(token)
	}
	return markCaptchaTokenUsedMemory(token)
}

// GetCaptchaIPRecord 获取 IP 限流记录
func GetCaptchaIPRecord(ip string) (*CaptchaIPRecord, error) {
	if common.RedisEnabled {
		return getCaptchaIPRecordRedis(ip)
	}
	return getCaptchaIPRecordMemory(ip)
}

// IncrementCaptchaIPFailure 增加 IP 失败次数
func IncrementCaptchaIPFailure(ip string) error {
	if common.RedisEnabled {
		return incrementCaptchaIPFailureRedis(ip)
	}
	return incrementCaptchaIPFailureMemory(ip)
}

// ResetCaptchaIPRecord 重置 IP 记录
func ResetCaptchaIPRecord(ip string) error {
	if common.RedisEnabled {
		return resetCaptchaIPRecordRedis(ip)
	}
	return resetCaptchaIPRecordMemory(ip)
}

// 内存存储实现
func storeCaptchaChallengeMemory(challenge *CaptchaChallenge) error {
	captchaStoreMutex.Lock()
	defer captchaStoreMutex.Unlock()
	captchaChallengeStore[challenge.SessionID] = challenge
	return nil
}

func getCaptchaChallengeMemory(sessionID string) (*CaptchaChallenge, error) {
	captchaStoreMutex.RLock()
	defer captchaStoreMutex.RUnlock()
	challenge, exists := captchaChallengeStore[sessionID]
	if !exists {
		return nil, errors.New("session not found")
	}
	if time.Now().After(challenge.ExpiresAt) {
		return nil, errors.New("session expired")
	}
	return challenge, nil
}

func deleteCaptchaChallengeMemory(sessionID string) error {
	captchaStoreMutex.Lock()
	defer captchaStoreMutex.Unlock()
	delete(captchaChallengeStore, sessionID)
	return nil
}

func storeCaptchaTokenMemory(token *CaptchaToken) error {
	captchaStoreMutex.Lock()
	defer captchaStoreMutex.Unlock()
	captchaTokenStore[token.Token] = token
	return nil
}

func getCaptchaTokenMemory(tokenStr string) (*CaptchaToken, error) {
	captchaStoreMutex.RLock()
	defer captchaStoreMutex.RUnlock()
	token, exists := captchaTokenStore[tokenStr]
	if !exists {
		return nil, errors.New("token not found")
	}
	if time.Now().After(token.ExpiresAt) {
		return nil, errors.New("token expired")
	}
	return token, nil
}

func markCaptchaTokenUsedMemory(tokenStr string) error {
	captchaStoreMutex.Lock()
	defer captchaStoreMutex.Unlock()
	token, exists := captchaTokenStore[tokenStr]
	if !exists {
		return errors.New("token not found")
	}
	token.Used = true
	return nil
}

func getCaptchaIPRecordMemory(ip string) (*CaptchaIPRecord, error) {
	captchaStoreMutex.RLock()
	defer captchaStoreMutex.RUnlock()
	record, exists := captchaIPStore[ip]
	if !exists {
		return nil, nil
	}
	return record, nil
}

func incrementCaptchaIPFailureMemory(ip string) error {
	captchaStoreMutex.Lock()
	defer captchaStoreMutex.Unlock()

	now := time.Now()
	record, exists := captchaIPStore[ip]
	if !exists {
		captchaIPStore[ip] = &CaptchaIPRecord{
			FailedAttempts: 1,
			FirstFailTime:  now,
		}
		return nil
	}

	// 检查是否超过时间窗口，重置计数
	if now.Sub(record.FirstFailTime) > CaptchaIPWindowDuration {
		record.FailedAttempts = 1
		record.FirstFailTime = now
		record.BlockedUntil = time.Time{}
		return nil
	}

	record.FailedAttempts++
	if record.FailedAttempts >= CaptchaMaxFailedAttempts {
		record.BlockedUntil = now.Add(CaptchaIPBlockDuration)
	}
	return nil
}

func resetCaptchaIPRecordMemory(ip string) error {
	captchaStoreMutex.Lock()
	defer captchaStoreMutex.Unlock()
	delete(captchaIPStore, ip)
	return nil
}

// Redis 存储实现（占位，后续实现）
func storeCaptchaChallengeRedis(challenge *CaptchaChallenge) error {
	// TODO: 实现 Redis 存储
	return storeCaptchaChallengeMemory(challenge)
}

func getCaptchaChallengeRedis(sessionID string) (*CaptchaChallenge, error) {
	// TODO: 实现 Redis 获取
	return getCaptchaChallengeMemory(sessionID)
}

func deleteCaptchaChallengeRedis(sessionID string) error {
	// TODO: 实现 Redis 删除
	return deleteCaptchaChallengeMemory(sessionID)
}

func storeCaptchaTokenRedis(token *CaptchaToken) error {
	// TODO: 实现 Redis 存储
	return storeCaptchaTokenMemory(token)
}

func getCaptchaTokenRedis(tokenStr string) (*CaptchaToken, error) {
	// TODO: 实现 Redis 获取
	return getCaptchaTokenMemory(tokenStr)
}

func markCaptchaTokenUsedRedis(tokenStr string) error {
	// TODO: 实现 Redis 更新
	return markCaptchaTokenUsedMemory(tokenStr)
}

func getCaptchaIPRecordRedis(ip string) (*CaptchaIPRecord, error) {
	// TODO: 实现 Redis 获取
	return getCaptchaIPRecordMemory(ip)
}

func incrementCaptchaIPFailureRedis(ip string) error {
	// TODO: 实现 Redis 更新
	return incrementCaptchaIPFailureMemory(ip)
}

func resetCaptchaIPRecordRedis(ip string) error {
	// TODO: 实现 Redis 删除
	return resetCaptchaIPRecordMemory(ip)
}

// CleanupExpiredCaptchaData 清理过期数据
func CleanupExpiredCaptchaData() {
	captchaStoreMutex.Lock()
	defer captchaStoreMutex.Unlock()

	now := time.Now()

	// 清理过期的挑战
	for sessionID, challenge := range captchaChallengeStore {
		if now.After(challenge.ExpiresAt) {
			delete(captchaChallengeStore, sessionID)
		}
	}

	// 清理过期的令牌
	for tokenStr, token := range captchaTokenStore {
		if now.After(token.ExpiresAt) {
			delete(captchaTokenStore, tokenStr)
		}
	}

	// 清理过期的 IP 记录
	for ip, record := range captchaIPStore {
		if !record.BlockedUntil.IsZero() && now.After(record.BlockedUntil) {
			delete(captchaIPStore, ip)
		} else if now.Sub(record.FirstFailTime) > CaptchaIPWindowDuration {
			delete(captchaIPStore, ip)
		}
	}
}
