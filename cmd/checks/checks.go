package checks

import (
	"slices"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/xueaaaa/go-uptime/internal/check/model"
	"github.com/xueaaaa/go-uptime/internal/check/service"
	"github.com/xueaaaa/go-uptime/internal/output"
	service2 "github.com/xueaaaa/go-uptime/internal/site/service"
)

func NewChecksCmd(checkSvc service.CheckService, siteSvc service2.SiteService) *cobra.Command {
	checks := &cobra.Command{
		Use:   "checks",
		Short: "Displays all site checks (limit - 100 latest checks)",
		RunE: func(c *cobra.Command, args []string) error {
			output.Header.Println("Table of checks")

			checks, err := checkSvc.GetAll(c.Context())
			if err != nil {
				return err
			}

			slices.SortFunc(checks, func(a, b model.Check) int {
				return b.CheckedAt.Compare(a.CheckedAt)
			})

			if len(checks) > 100 {
				checks = checks[:100]
			}

			data := pterm.TableData{
				{"Site", "Status", "Status code", "Latency", "Error", "Checked at"},
			}

			for _, check := range checks {
				site, err := siteSvc.Get(c.Context(), check.SiteID)
				if err != nil {
					return err
				}

				checkData := []string{
					site.URL,
					output.StatusColor(check.Status),
					output.StatusCodeColor(check.StatusCode),
					output.PingColor(check.Latency),
					check.Error,
					check.CheckedAt.Format("2006 Jan 02 15:04:05.000"),
				}

				data = append(data, checkData)
			}

			if err = output.Table.WithData(data).Render(); err != nil {
				return err
			}

			return nil
		},
	}

	return checks
}
