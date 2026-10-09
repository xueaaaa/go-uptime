package output

import (
	"context"
	"slices"
	"time"

	"github.com/pterm/pterm"
	"github.com/xueaaaa/go-uptime/internal/check/model"
	"github.com/xueaaaa/go-uptime/internal/check/service"
	service2 "github.com/xueaaaa/go-uptime/internal/site/service"
)

func LiveChecks(ctx context.Context, checkSvc service.CheckService, siteSvc service2.SiteService) error {
	startTime := time.Now()
	area, err := pterm.DefaultArea.WithRemoveWhenDone(false).Start()
	if err != nil {
		return err
	}
	defer area.Stop()

	render := func() error {
		checks, err := checkSvc.GetAll(ctx)
		if err != nil {
			return err
		}

		filtered := slices.DeleteFunc(checks, func(check model.Check) bool {
			return check.CheckedAt.Before(startTime)
		})

		data := pterm.TableData{{"Site", "Status", "Status code", "Latency", "Error", "Checked at"}}
		for _, c := range filtered {
			site, err := siteSvc.Get(ctx, c.SiteID)
			if err != nil {
				return err
			}

			data = append(data, []string{
				site.URL,
				StatusColor(c.Status),
				StatusCodeColor(c.StatusCode),
				PingColor(c.Latency),
				c.Error,
				c.CheckedAt.Format("2006 Jan 02 15:04:05.000"),
			})
		}

		out, err := pterm.DefaultTable.WithHasHeader().WithData(data).Srender()
		if err != nil {
			return err
		}

		area.Update(out)
		return nil
	}

	if err := render(); err != nil {
		return err
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := render(); err != nil {
				return err
			}
		}
	}
}
