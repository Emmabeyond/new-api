package common

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/bytedance/gopkg/util/gopool"
)

var relayGoPool gopool.Pool

// 默认goroutine池大小，可通过环境变量 GOPOOL_SIZE 配置
const defaultGoPoolSize = 10000

func init() {
	poolSize := defaultGoPoolSize
	if envSize := os.Getenv("GOPOOL_SIZE"); envSize != "" {
		if size, err := strconv.Atoi(envSize); err == nil && size > 0 {
			poolSize = size
		}
	}
	relayGoPool = gopool.NewPool("gopool.RelayPool", int32(poolSize), gopool.NewConfig())
	relayGoPool.SetPanicHandler(func(ctx context.Context, i interface{}) {
		if stopChan, ok := ctx.Value("stop_chan").(chan bool); ok {
			SafeSendBool(stopChan, true)
		}
		SysError(fmt.Sprintf("panic in gopool.RelayPool: %v", i))
	})
}

func RelayCtxGo(ctx context.Context, f func()) {
	relayGoPool.CtxGo(ctx, f)
}
