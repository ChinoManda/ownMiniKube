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
 Image string
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(Image)
		cli := core.InitClient()
		ctx := core.InitNamespace()
		containers := core.ListPods(cli, ctx)
		running := core.ListRunningPods(containers, ctx)
	 if len(running) > 0 {
    err := core.RollingUpdate(running, ctx, Image, cli)
		if err == nil{
			fmt.Println("bien")
			return  err
		} 
		fmt.Println("mal")
		return  err
	 }
	 if len(running) == 0 {
		 img, _ := core.PullImage(cli, ctx, Image) 
		 pod, err := core.NewPod(cli, ctx, img, "AutoDeployed")
		 if err != nil {
			 return err
		 }
		 Rerr := pod.Run()
		 return Rerr
	 }
	 return nil
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
  deployCmd.Flags().StringVarP(&Image, "image", "i", "", "Image name without tag")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
