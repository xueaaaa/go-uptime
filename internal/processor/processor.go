package processor

import (
	"context"

	"github.com/xueaaaa/go-uptime/internal/check/model"
	service2 "github.com/xueaaaa/go-uptime/internal/check/service"
	"github.com/xueaaaa/go-uptime/internal/site/service"
)

type Processor struct {
	results      <-chan model.Check
	siteService  service.SiteService
	checkService service2.CheckService
}

func NewProcessor(
	results <-chan model.Check,
	siteService service.SiteService,
	checkService service2.CheckService,
) *Processor {
	return &Processor{
		results:      results,
		siteService:  siteService,
		checkService: checkService,
	}
}

func (p *Processor) Run(ctx context.Context) error {
	for {
		select {
		case res, ok := <-p.results:
			if !ok {
				return nil
			}

			if _, err := p.checkService.Create(ctx, res); err != nil {
				return err /* TODO: Log instead of return */
			}

			if err := p.siteService.UpdateByCheck(ctx, res); err != nil {
				return err /* TODO: Log instead of return */
			}
		}
	}
}
