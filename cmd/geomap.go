package cmd

import (
	"github.com/spf13/cobra"

	"github.com/QOLPlus/weather-map/cmd/geomap"
)

var geomapCmd = &cobra.Command{
	Use:   "geomap",
	Short: "Generate Geographic Map",
	Long: "Generate Geographic Map",
	Run: geomap.Process,
}

func init() {
	rootCmd.AddCommand(geomapCmd)
}
