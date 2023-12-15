package cmd

import (
	"fmt"
	"github.com/nextrevision/traci/providers"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "detect CI environment and print config",
	Run:   doDetect,
	Args:  cobra.MinimumNArgs(0),
}

func doDetect(cmd *cobra.Command, args []string) {
	config := getConfig()
	configStr, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("unable to marshal config to YAML: %v", err)
	}

	provider := providers.DetectProvider()
	if err != nil {
		log.Fatalf("error detecting CI provider: %s", err.Error())
	}

	fmt.Println("Config:")
	fmt.Printf("%s\n", configStr)

	fmt.Println("CI Settings")
	fmt.Printf("  provider: %s\n", provider.GetCIName())
	fmt.Printf("  service name: %s\n", provider.GetServiceName())
	fmt.Printf("  span name: %s\n", provider.GetSpanName())
	fmt.Printf("  trace value: %s\n", provider.GetTraceVal())
	fmt.Printf("  span value: %s\n", provider.GetSpanVal())
	fmt.Println("  attributes:")
	for k, v := range provider.GetAttributes() {
		fmt.Printf("    %s: %s\n", k, v)
	}
}

func init() {
	rootCmd.AddCommand(detectCmd)
}
