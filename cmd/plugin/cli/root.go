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

var rootShort = "View pod information at the container level"

var rootDescription = ` Ice lets you view configuration and settings of the containers that run inside pods.

 Suggestions and improvements can be made by raising an issue here: 
    https://github.com/NimbleArchitect/kubectl-ice

`

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "kubectl-ice",
		Short:         rootShort,
		Long:          fmt.Sprintf("%s\n\n%s", rootShort, rootDescription),
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}
