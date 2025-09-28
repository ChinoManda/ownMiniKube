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
	id string
)
// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(id)
		cli := core.InitClient()
		ctx := core.InitNamespace()
    pod, err := core.GetPodByID(ctx, cli, id)
    if err != nil {
      fmt.Printf("error cargando contenedor %s: %w", id, err)
    }
		killCode, killErr := pod.Kill()
		fmt.Println(killCode, killErr)
	  delErr := pod.Delete()
		fmt.Println(delErr)
  
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&id, "id", "d", "", "Container ID")
	deleteCmd.MarkFlagRequired("id")
}
