package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(timerCmd)
}

var timerCmd = &cobra.Command{
	Use:   "timer",
	Short: "Sets the interval between notifications",
	Long: `Use this command to set how often you want to get price update (in minutes), the lowest you can set is
  one minute`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
