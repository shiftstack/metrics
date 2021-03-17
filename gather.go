package metrics

import (
	"context"
	"sync"
	"time"
)

func refreshAll(services []service) {
	var wg sync.WaitGroup
	for _, s := range services {
		wg.Add(1)
		go func() {
			s.Refresh()
			wg.Done()
		}()
	}
	wg.Wait()
}

func keepFresh(ctx context.Context, ttl time.Duration, services []service) error {
	ticker := time.NewTicker(ttl)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return ctx.Err()
		case <-ticker.C:
			refreshAll(services)
		}
	}
}
