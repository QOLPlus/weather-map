package index

import (
	"fmt"
	"github.com/spf13/cobra"
)

func Process(cmd *cobra.Command, args []string) {
	fmt.Println("index called")
}
