package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/bnadim/csf/csf"
)

var rootCmd = &cobra.Command{
	Use:   "csf [input file] [output file]",
	Short: "Concatenate swagger file based on $ref",
	Long: ``,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		outputPath := args[1]
		return csf.Convert(inputPath, outputPath)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}