package model

import (
	"errors"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"

	"github.com/bytedance/gopkg/util/gopool"
	"gorm.io/gorm"
)

const (
	BatchUpdateTypeUserQuota = iota
	BatchUpdateTypeTokenQuota
	BatchUpdateTypeUsedQuota
	BatchUpdateTypeChannelUsedQuota
	BatchUpdateTypeRequestCount
	BatchUpdateTypeCount // if you add a new type, you need to add a new map and a new lock
)

var batchUpdateStores []map[int]int
var batchUpdateLocks []sync.Mutex

func init() {
	for i := 0; i < BatchUpdateTypeCount; i++ {
		batchUpdateStores = append(batchUpdateStores, make(map[int]int))
		batchUpdateLocks = append(batchUpdateLocks, sync.Mutex{})
	}
}

func InitBatchUpdater() {
	gopool.Go(func() {
		for {
			time.Sleep(time.Duration(common.BatchUpdateInterval) * time.Second)
			batchUpdate()
		}
	})
}

func addNewRecord(type_ int, id int, value int) {
	batchUpdateLocks[type_].Lock()
	defer batchUpdateLocks[type_].Unlock()
	if _, ok := batchUpdateStores[type_][id]; !ok {
		batchUpdateStores[type_][id] = value
	} else {
		batchUpdateStores[type_][id] += value
	}
}

func batchUpdate() {
	// 收集所有需要更新的数据
	stores := make([]map[int]int, BatchUpdateTypeCount)
	hasData := false
	for i := 0; i < BatchUpdateTypeCount; i++ {
		batchUpdateLocks[i].Lock()
		if len(batchUpdateStores[i]) > 0 {
			hasData = true
			stores[i] = batchUpdateStores[i]
			batchUpdateStores[i] = make(map[int]int)
		}
		batchUpdateLocks[i].Unlock()
	}

	if !hasData {
		return
	}

	common.SysLog("batch update started")
	
	// 使用WaitGroup并行处理不同类型的更新
	var wg sync.WaitGroup
	for i := 0; i < BatchUpdateTypeCount; i++ {
		store := stores[i]
		if store == nil || len(store) == 0 {
			continue
		}
		wg.Add(1)
		go func(updateType int, data map[int]int) {
			defer wg.Done()
			for key, value := range data {
				switch updateType {
				case BatchUpdateTypeUserQuota:
					err := increaseUserQuota(key, value)
					if err != nil {
						common.SysLog("failed to batch update user quota: " + err.Error())
					}
				case BatchUpdateTypeTokenQuota:
					err := increaseTokenQuota(key, value)
					if err != nil {
						common.SysLog("failed to batch update token quota: " + err.Error())
					}
				case BatchUpdateTypeUsedQuota:
					updateUserUsedQuota(key, value)
				case BatchUpdateTypeRequestCount:
					updateUserRequestCount(key, value)
				case BatchUpdateTypeChannelUsedQuota:
					updateChannelUsedQuota(key, value)
				}
			}
		}(i, store)
	}
	wg.Wait()
	common.SysLog("batch update finished")
}

func RecordExist(err error) (bool, error) {
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

func shouldUpdateRedis(fromDB bool, err error) bool {
	return common.RedisEnabled && fromDB && err == nil
}
