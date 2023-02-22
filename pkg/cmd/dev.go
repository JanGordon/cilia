package cmd

import (
	"fmt"
	"strconv"

	"github.com/JanGordon/cilia-framework/pkg/server"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:     "development",
	Aliases: []string{"dev"},
	Args:    cobra.ExactArgs(1),
	Short:   "Runs a dev server wiht hot reload",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println(url.ResolvePath("./go.mod", "/rotes/hi"))
		// ssr.Compile()
		portNum, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(fmt.Errorf("please enter a valid number for port. e.g. 8080"))
		}
		server.Dev(portNum)
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
