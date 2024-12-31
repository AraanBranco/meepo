package cmd

import (
	"github.com/AraanBranco/meepow/cmd/bot"
	"github.com/AraanBranco/meepow/cmd/managementapi"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the provided meepow component service",
}

func init() {
	startCmd.AddCommand(managementapi.ManagementApiCmd)
	startCmd.AddCommand(bot.BotCmd)
}
