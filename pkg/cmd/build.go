package cmd

import (
	"github.com/JanGordon/cilia-framework/pkg/global"
	"github.com/JanGordon/cilia-framework/pkg/ssr"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "builds the project",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println(url.ResolvePath("./go.mod", "/rotes/hi"))
		// ssr.Compile()
		ssr.Compile(global.ProjectRoot, false, "")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
