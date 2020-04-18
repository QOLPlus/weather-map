package cmd

import (
	"github.com/spf13/cobra"

	"github.com/QOLPlus/weather-map/cmd/geomap"
	"github.com/QOLPlus/weather-map/cmd/index"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Run index after running geomap",
	Long: "Run index after running geomap",

	Run: func(cmd *cobra.Command, args []string) {
		geomap.Process(cmd, args)
		index.Process(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
}
