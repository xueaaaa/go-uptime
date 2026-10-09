package cmd

import (
	"github.com/spf13/cobra"
	cmdChecks "github.com/xueaaaa/go-uptime/cmd/checks"
	cmdSites "github.com/xueaaaa/go-uptime/cmd/sites"
	service2 "github.com/xueaaaa/go-uptime/internal/check/service"
	"github.com/xueaaaa/go-uptime/internal/site/service"
)

func NewRootCmd(siteSvc service.SiteService, checkSvc service2.CheckService) *cobra.Command {
	root := &cobra.Command{
		Use:           "go-uptime",
		Short:         "Automatic website availability monitoring",
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	root.AddCommand(
		cmdSites.NewSitesCmd(siteSvc),
		cmdChecks.NewChecksCmd(checkSvc, siteSvc),
		NewServeCmd(siteSvc, checkSvc),
	)

	return root
}
