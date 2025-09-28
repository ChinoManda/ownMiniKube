/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
  "ownkube/core"
	"github.com/spf13/cobra"
)

// listpodsCmd represents the listpods command
var listpodsCmd = &cobra.Command{
	Use:   "listpods",
	Short: "List all alive pods",
	Run: func(cmd *cobra.Command, args []string) {
		cli := core.InitClient()
		ctx := core.InitNamespace()
		containers := core.ListPods(cli, ctx)
		for _, c:= range containers {
    fmt.Println("ID:", c.ID())// ID del contenedor
		img, _ := c.Image(ctx)
    fmt.Println("Image:", img)         // La imagen asociada al contenedor
		}
	},
}

func init() {
	rootCmd.AddCommand(listpodsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listpodsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listpodsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
