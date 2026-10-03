package cmd

import (
	"github.com/spf13/cobra"
	cmd "github.com/xueaaaa/go-uptime/cmd/sites"
	"github.com/xueaaaa/go-uptime/internal/site/service"
)

func NewRootCmd(siteSvc service.SiteService) *cobra.Command {
	root := &cobra.Command{
		Use:           "uptime",
		Short:         "Automatic website availability monitoring",
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	root.AddCommand(
		cmd.NewSitesCmd(siteSvc),
	)

	return root
}
