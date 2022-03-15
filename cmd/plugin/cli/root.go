package cli

import (
	"fmt"
	"os"

	"github.com/NimbleArchitect/kubectl-ice/pkg/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// auto updated version via gorelaser
var version = "0.0.0"

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubectl-ice",
		Short: "view container settings",
		Long: `ice lets you view configuration settings of containers inside pods.
you can run ice through kubectl with: kubectl ice [command]`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cobra.OnInitialize(initConfig)
	plugin.InitSubCommands(cmd)

	return cmd
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}
