package worker

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/xueaaaa/go-uptime/internal/check/model"
	"github.com/xueaaaa/go-uptime/internal/job"
	model2 "github.com/xueaaaa/go-uptime/internal/site/model"
)

type Pool struct {
	jobs    <-chan job.Job
	results chan<- model.Check
}

func NewPool(
	jobs <-chan job.Job,
	results chan<- model.Check,
) *Pool {
	return &Pool{
		jobs:    jobs,
		results: results,
	}
}

func (p *Pool) Run(ctx context.Context) error {
	workers := 50 /* TODO: config */
	wg := sync.WaitGroup{}

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.work(ctx)
		}()
	}

	wg.Wait()
	close(p.results)
	return nil
}

func (p *Pool) work(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case j, ok := <-p.jobs:
			if !ok {
				return
			}

			res := p.check(ctx, j)
			select {
			case p.results <- res:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (p *Pool) check(ctx context.Context, job job.Job) model.Check {
	base := model.Check{
		SiteID:    job.SiteID,
		CheckedAt: time.Now(),
	}

	req, err := http.NewRequest(http.MethodGet, job.URL, nil)
	if err != nil {
		base.Status = model2.Unavailable
		base.Error = err.Error()
		return base
	}

	timeout := 10 * time.Second /* TODO: config */
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	req = req.WithContext(timeoutCtx)

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	base.Latency = int(time.Since(start).Milliseconds())
	if err != nil {
		base.Status = model2.Unavailable
		base.Error = err.Error()
		return base
	}
	defer resp.Body.Close()

	_, _ = io.Copy(io.Discard, resp.Body)

	switch {
	case http.StatusOK <= resp.StatusCode && resp.StatusCode <= http.StatusAlreadyReported:
		base.Status = model2.Available
	default:
		base.Status = model2.Unavailable
		base.Error = fmt.Sprintf("unexpected status code: %d", resp.StatusCode)
	}

	return base
}
