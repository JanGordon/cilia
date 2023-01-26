package cmd

import (
	"fmt"
	"strconv"

	"github.com/JanGordon/cilia-framework/pkg/server"
	"github.com/spf13/cobra"
)

var global bool

var addCmd = &cobra.Command{
	Use:   "add",
	Args:  cobra.ExactArgs(1),
	Short: "adds a component or addon",
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
	addCmd.PersistentFlags().BoolVarP(&global, "global", "g", false, "add addon globally")
	rootCmd.AddCommand(addCmd)
}
