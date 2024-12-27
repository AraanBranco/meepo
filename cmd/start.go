package cmd

import (
	"github.com/AraanBranco/meepo/cmd/bot"
	"github.com/AraanBranco/meepo/cmd/managementapi"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the provided meepo component service",
}

func init() {
	startCmd.AddCommand(managementapi.ManagementApiCmd)
	startCmd.AddCommand(bot.BotCmd)
}
