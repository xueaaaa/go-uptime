package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/xueaaaa/go-uptime/internal/output"
	"github.com/xueaaaa/go-uptime/internal/site/service"
)

func NewSitesCmd(siteSvc service.SiteService) *cobra.Command {
	sites := &cobra.Command{
		Use:   "sites",
		Short: "Displays a summary table of added sites",
		RunE: func(c *cobra.Command, args []string) error {
			output.Header.Println("Table of sites")

			sites, err := siteSvc.GetAll(c.Context())
			if err != nil {
				return err
			}

			data := pterm.TableData{
				{"URL", "Status", "Interval", "Last check at", "Next check at", "Created at"},
			}

			for _, site := range sites {
				var lastCheckAtStr string
				if site.LastCheckAt == nil {
					lastCheckAtStr = "N/A"
				} else {
					lastCheckAtStr = site.LastCheckAt.Format("2006 Jan 02 15:04:05.000")
				}

				siteData := []string{
					site.URL,
					output.StatusColor(site.Status),
					site.Interval.String(),
					lastCheckAtStr,
					site.NextCheckAt.Format("2006 Jan 02 15:04:05.000"),
					site.CreatedAt.Format("2006 Jan 02 15:04:05.000"),
				}
				data = append(data, siteData)
			}

			if err = output.Table.WithData(data).Render(); err != nil {
				return err
			}

			return nil
		},
	}

	sites.AddCommand(
		NewSiteAddCmd(siteSvc),
		NewSiteRemoveCmd(siteSvc),
		NewSiteIntervalCmd(siteSvc),
	)
	return sites
}
