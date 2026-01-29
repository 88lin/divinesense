//go:build perf

package queryengine

import (
	"context"
	"testing"
	"time"
)

// TestQueryRouter_TodayPerformance 性能测试：对比"今日"vs"今天".
// Run with: go test -tags=perf ./server/queryengine/...
func TestQueryRouter_TodayPerformance(t *testing.T) {
	router := NewQueryRouter()
	ctx := context.Background()

	queries := []string{
		"今天有什么安排",
		"今日日程",
		"明天的计划",
		"明日安排",
	}

	// 预热
	for _, query := range queries {
		router.Route(ctx, query, nil)
	}

	// 性能测试：1000次路由
	iterations := 1000
	start := time.Now()

	for i := 0; i < iterations; i++ {
		for _, query := range queries {
			router.Route(ctx, query, nil)
		}
	}

	duration := time.Since(start)
	avgDuration := duration / time.Duration(iterations*len(queries))

	t.Logf("Performance: %d routes in %v", iterations*len(queries), duration)
	t.Logf("Average time per route: %v", avgDuration)

	// 目标：平均路由时间 < 10μs
	if avgDuration > 10*time.Microsecond {
		t.Errorf("Route() too slow: %v, want < 10μs", avgDuration)
	}

	t.Logf("✓ Performance target met: %v < 10μs", avgDuration)
}
