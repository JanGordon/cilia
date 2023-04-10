package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "shows version",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println(url.ResolvePath("./go.mod", "/rotes/hi"))
		// ssr.Compile()
		fmt.Println("v0.0.2")

	},
}

func init() {
	// createCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "add addon globally")
	rootCmd.AddCommand(versionCmd)
}
