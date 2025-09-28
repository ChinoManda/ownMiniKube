/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
  "ownkube/core"
	"github.com/spf13/cobra"
)

// listRunningCmd represents the listRunning command
var listRunningCmd = &cobra.Command{
	Use:   "listRunning",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		cli := core.InitClient()
		ctx := core.InitNamespace()
		containers := core.ListPods(cli, ctx)
		running := core.ListRunningPods(containers, ctx)
		for _,c := range running{
			  task, _:= c.Task(ctx, nil)
				id := c.ID()
        image, _ := c.Image(ctx)
				PID := task.Pid()
				fmt.Printf("ID: %s  IMAGE: %s  PID %d \n", id, image.Name(), PID)
		}
	},
}

func init() {
	rootCmd.AddCommand(listRunningCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listRunningCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listRunningCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
