/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"ownkube/core"
)
var (
	image string
	number int
)
// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",

	RunE: func(cmd *cobra.Command, args []string) error{
		cli := core.InitClient()
		ctx := core.InitNamespace()
		img, err := core.PullImage(cli, ctx, image)
		if err != nil {
			return fmt.Errorf("error obteniendo imagen: %w", err)
		}
		for i := 0; i < number; i++ {
			pod, err := core.NewPod(cli, ctx, img, fmt.Sprintf("Pod-%d", i))
			if err != nil {
				return err
			}
			if err := pod.Run(); err != nil {
				return err
			}
			fmt.Printf("Pod %s creado y corriendo\n", pod.Id)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&image, "image", "i", "", "Container Image (ej: docker.io/library/nginx:latest)")
	createCmd.Flags().IntVarP(&number, "number", "n", 1, "Pod number")
	createCmd.MarkFlagRequired("image")
}
