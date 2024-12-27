package bot

import "github.com/spf13/cobra"

var (
	logConfig  string
	configPath string
)

const serviceName string = "bot"

var BotCmd = &cobra.Command{
	Use:     "bot",
	Short:   "Starts meepo bot service",
	Example: "meepo start bot -c config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		startBot()
	},
}

func init() {
	BotCmd.Flags().StringVarP(&configPath, "config-path", "c", "config/config.yaml", "path of the configuration YAML file")
}

func startBot() {
	// Start the bot service
}
