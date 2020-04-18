package cmd

import (
	"github.com/spf13/cobra"

	"github.com/QOLPlus/weather-map/cmd/index"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Generate searchable index using geomap",
	Long: "Generate searchable index usin ggeomap",
	Run: index.Process,
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
