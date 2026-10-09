package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/xueaaaa/go-uptime/internal/site/service"
	"github.com/xueaaaa/go-uptime/internal/util"
)

func NewSiteRemoveCmd(siteSvc service.SiteService) *cobra.Command {
	remove := &cobra.Command{
		Use:   "remove <ip|url>",
		Short: "Removes the specified site from the list of sites to be checked",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			arg := args[0]

			h, err := util.CleanAddr(arg, false /* TODO: config */)

			if err != nil {
				return err
			}

			got, err := siteSvc.GetByUrl(c.Context(), h)
			if err != nil {
				return err
			}

			if err = siteSvc.Delete(c.Context(), got.ID); err != nil {
				return err
			}

			pterm.Success.Printfln("Site (%s) successfully deleted", arg)
			return nil
		},
	}

	return remove
}
