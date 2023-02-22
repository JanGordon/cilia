package cmd

import (
	"fmt"
	"strconv"

	"github.com/JanGordon/cilia-framework/pkg/server"
	"github.com/spf13/cobra"
)

var prodCmd = &cobra.Command{
	Use:     "production",
	Aliases: []string{"prod"},
	Args:    cobra.ExactArgs(1),
	Short:   "runs the production server for the project",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println(url.ResolvePath("./go.mod", "/rotes/hi"))
		// ssr.Compile()
		portNum, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(fmt.Errorf("please enter a valid number for port. e.g. 8080"))
		}
		server.Prod(portNum)
	},
}

func init() {
	rootCmd.AddCommand(prodCmd)
}
