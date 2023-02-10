package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"c"},
	Short:   "creates a base cilia project",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println(url.ResolvePath("./go.mod", "/rotes/hi"))
		// ssr.Compile()

	},
}

func init() {
	// createCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "add addon globally")
	rootCmd.AddCommand(createCmd)
}
