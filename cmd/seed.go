/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/raspiantoro/goms/internal/app"
	"github.com/spf13/cobra"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed [seed name]",
	Short: "Add new seed file",
	Long: `Create and add a new seed file to your Golang project. 
The seed name will be prefixed with the current datetime in the format 'YYYYMMDDhhmm' as a seed version/identifier in the generated file name.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Print("Error: add seed requires a name for the seed.\n\n")
			cmd.Help()
			return
		} else if len(args) > 1 {
			fmt.Print("Error: add seed only accept one arguments [seed name]\n\n")
			cmd.Help()
			return
		}

		cfg := app.Config{Name: args[0]}

		if flagDir != "" {
			cfg.Path = flagDir
		}

		app.CreateSeed(cfg)
	},
}

func init() {
	addCmd.AddCommand(seedCmd)

	// Here you will define your flags and configuration settings.
	seedCmd.Flags().StringVarP(&flagDir, "dir", "d", "", "gormseed files directory")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
