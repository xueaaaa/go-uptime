package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/xueaaaa/go-uptime/internal/site/model"
	"github.com/xueaaaa/go-uptime/internal/site/service"
)

func NewAddCmd(siteSvc service.SiteService) *cobra.Command {
	var interval time.Duration

	add := &cobra.Command{
		Use:   "add <ip|url>",
		Short: "Adds the specified site to the list of sites to be checked",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			arg := args[0]

			if _, err := url.ParseRequestURI(arg); err != nil {
				return errors.New("specified argument is not a valid url")
			}

			if interval < 10*time.Second {
				return fmt.Errorf("interval must be at least 10s, got %s", interval)
			}

			site := model.Site{
				URL:       arg,
				Interval:  interval,
				CreatedAt: time.Now(),
			}

			_, err := siteSvc.Create(c.Context(), site)
			if err != nil {
				return err
			}

			pterm.Success.Printfln("Site (%s) successfully added", arg)
			return nil
		},
	}

	add.Flags().DurationVarP(&interval, "interval", "i", time.Minute,
		"how often to check the site's availability (e.g. 30s, 5m)")

	return add
}
