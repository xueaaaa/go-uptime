package cmd

import (
	"fmt"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/xueaaaa/go-uptime/internal/site/service"
	"github.com/xueaaaa/go-uptime/internal/util"
)

func NewSiteIntervalCmd(siteSvc service.SiteService) *cobra.Command {
	interval := &cobra.Command{
		Use:   "interval <ip|url> <interval>",
		Short: "Changes the period for checking the website availability (minimum - 10 seconds)",
		Args:  cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			siteArg := args[0]

			h, err := util.CleanAddr(siteArg, false /* TODO: config */)
			if err != nil {
				return err
			}

			intervalArg := args[1]
			interval, err := time.ParseDuration(intervalArg)
			if err != nil {
				return fmt.Errorf("invalid interval: %s", err.Error())
			}

			if interval <= 10*time.Second {
				fmt.Errorf("interval must be at least 10s, got %s", interval)
			}

			site, err := siteSvc.GetByUrl(c.Context(), h)
			if err != nil {
				return err
			}

			site.Interval = interval
			if err = siteSvc.Update(c.Context(), site); err != nil {
				return err
			}

			pterm.Success.Printfln("New heck interval (%s) for site (%s) has been successfully set", interval, h)
			return nil
		},
	}

	return interval
}
