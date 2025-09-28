/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
  "ownkube/core"
	"github.com/spf13/cobra"
)

var (
 image string
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {	
		cli := core.InitClient()
		ctx := core.InitNamespace()
		containers := core.ListPods(cli, ctx)
		running := core.ListRunningPods(containers, ctx)
		
    for _, r:= range running {
    	r.Image(ctx)
    }
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
  deployCmd.Flags().StringVarP(&image, "image", "i", "" "Image name without tag")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
