package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/go-redis/redis/v8"
)

// ModelRequest represents a single model request record
type ModelRequest struct {
	ModelName string    `json:"model_name"`
	Timestamp time.Time `json:"timestamp"`
}

// ModelSwitchTracker tracks model switching behavior for tokens
type ModelSwitchTracker interface {
	// RecordModelRequest records a model request for a token
	RecordModelRequest(tokenID int, modelName string) error

	// GetDistinctModelCount returns the number of distinct models requested within the time window
	GetDistinctModelCount(tokenID int, windowMinutes int) (int, error)

	// GetModelHistory returns the model request history within the time window
	GetModelHistory(tokenID int, windowMinutes int) ([]ModelRequest, error)
}

// redisModelSwitchTracker implements ModelSwitchTracker using Redis
type redisModelSwitchTracker struct{}

// memoryModelSwitchTracker implements ModelSwitchTracker using in-memory storage
type memoryModelSwitchTracker struct {
	mu      sync.RWMutex
	records map[int][]ModelRequest // tokenID -> list of requests
}

var (
	modelSwitchTracker     ModelSwitchTracker
	modelSwitchTrackerOnce sync.Once
)

// GetModelSwitchTracker returns the singleton ModelSwitchTracker instance
func GetModelSwitchTracker() ModelSwitchTracker {
	modelSwitchTrackerOnce.Do(func() {
		if common.RedisEnabled {
			modelSwitchTracker = &redisModelSwitchTracker{}
		} else {
			modelSwitchTracker = &memoryModelSwitchTracker{
				records: make(map[int][]ModelRequest),
			}
		}
	})
	return modelSwitchTracker
}


// Redis key format for model switch tracking
const modelSwitchKeyPrefix = "abuse:model_switch:"

func getModelSwitchKey(tokenID int) string {
	return fmt.Sprintf("%s%d", modelSwitchKeyPrefix, tokenID)
}

// Redis implementation

func (r *redisModelSwitchTracker) RecordModelRequest(tokenID int, modelName string) error {
	ctx := context.Background()
	key := getModelSwitchKey(tokenID)
	now := time.Now().Unix()

	// Use ZADD to add model request with timestamp as score
	// Member is "modelName:timestamp" to allow multiple requests for same model
	member := fmt.Sprintf("%s:%d", modelName, now)
	err := common.RDB.ZAdd(ctx, key, &redis.Z{
		Score:  float64(now),
		Member: member,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to record model request: %w", err)
	}

	// Set expiration to clean up old data (24 hours)
	common.RDB.Expire(ctx, key, 24*time.Hour)

	return nil
}

func (r *redisModelSwitchTracker) GetDistinctModelCount(tokenID int, windowMinutes int) (int, error) {
	history, err := r.GetModelHistory(tokenID, windowMinutes)
	if err != nil {
		return 0, err
	}

	// Count distinct models
	modelSet := make(map[string]struct{})
	for _, req := range history {
		modelSet[req.ModelName] = struct{}{}
	}

	return len(modelSet), nil
}

func (r *redisModelSwitchTracker) GetModelHistory(tokenID int, windowMinutes int) ([]ModelRequest, error) {
	ctx := context.Background()
	key := getModelSwitchKey(tokenID)

	// Calculate time window
	minScore := float64(time.Now().Add(-time.Duration(windowMinutes) * time.Minute).Unix())
	maxScore := float64(time.Now().Unix())

	// Get all members within the time window
	results, err := common.RDB.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min: fmt.Sprintf("%f", minScore),
		Max: fmt.Sprintf("%f", maxScore),
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get model history: %w", err)
	}

	history := make([]ModelRequest, 0, len(results))
	for _, z := range results {
		member := z.Member.(string)
		// Parse "modelName:timestamp" format
		modelName := parseModelNameFromMember(member)
		history = append(history, ModelRequest{
			ModelName: modelName,
			Timestamp: time.Unix(int64(z.Score), 0),
		})
	}

	return history, nil
}

// parseModelNameFromMember extracts model name from "modelName:timestamp" format
func parseModelNameFromMember(member string) string {
	// Find the last colon (timestamp separator)
	for i := len(member) - 1; i >= 0; i-- {
		if member[i] == ':' {
			return member[:i]
		}
	}
	return member
}


// Memory implementation

func (m *memoryModelSwitchTracker) RecordModelRequest(tokenID int, modelName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	req := ModelRequest{
		ModelName: modelName,
		Timestamp: time.Now(),
	}

	m.records[tokenID] = append(m.records[tokenID], req)

	// Clean up old records (older than 24 hours)
	m.cleanupOldRecords(tokenID, 24*60)

	return nil
}

func (m *memoryModelSwitchTracker) GetDistinctModelCount(tokenID int, windowMinutes int) (int, error) {
	history, err := m.GetModelHistory(tokenID, windowMinutes)
	if err != nil {
		return 0, err
	}

	// Count distinct models
	modelSet := make(map[string]struct{})
	for _, req := range history {
		modelSet[req.ModelName] = struct{}{}
	}

	return len(modelSet), nil
}

func (m *memoryModelSwitchTracker) GetModelHistory(tokenID int, windowMinutes int) ([]ModelRequest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	records, exists := m.records[tokenID]
	if !exists {
		return []ModelRequest{}, nil
	}

	cutoff := time.Now().Add(-time.Duration(windowMinutes) * time.Minute)
	history := make([]ModelRequest, 0)

	for _, req := range records {
		if req.Timestamp.After(cutoff) {
			history = append(history, req)
		}
	}

	return history, nil
}

// cleanupOldRecords removes records older than the specified window
func (m *memoryModelSwitchTracker) cleanupOldRecords(tokenID int, windowMinutes int) {
	records, exists := m.records[tokenID]
	if !exists {
		return
	}

	cutoff := time.Now().Add(-time.Duration(windowMinutes) * time.Minute)
	newRecords := make([]ModelRequest, 0)

	for _, req := range records {
		if req.Timestamp.After(cutoff) {
			newRecords = append(newRecords, req)
		}
	}

	m.records[tokenID] = newRecords
}
