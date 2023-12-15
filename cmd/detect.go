package cmd

import (
	"fmt"
	"github.com/nextrevision/traci/providers"
	"github.com/spf13/cobra"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "detect CI environment and print config",
	Run:   doDetect,
	Args:  cobra.MinimumNArgs(0),
}

func doDetect(cmd *cobra.Command, args []string) {
	provider := providers.DetectProvider()

	fmt.Println("CI Settings")
	fmt.Printf("  provider: %s\n", provider.GetCIName())
	fmt.Printf("  service name: %s\n", provider.GetServiceName())
	fmt.Printf("  span name: %s\n", provider.GetSpanName())
	fmt.Printf("  trace unique string: %s\n", provider.GetTraceVal())
	fmt.Println("  attributes:")
	for k, v := range provider.GetAttributes() {
		fmt.Printf("    %s: %s\n", k, v)
	}
}

func init() {
	rootCmd.AddCommand(detectCmd)
}
