package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var versionsShort = "Display container versions and mount points"

var helpTemplate = `
{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}
More information at: https://nimblearchitect.github.io/kubectl-ice/

`

func Version(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	// 1234567890123456789012345678901234567890123456789012345678901234567890123456789
	fmt.Printf(`kubectl-ice kubernetes container viewer

version %s

the latest version can be found at: 
	https://nimblearchitect.github.io/kubectl-ice/downloads/

to view the documentation:
	https://nimblearchitect.github.io/kubectl-ice

or to raise issues: 
   https://github.com/NimbleArchitect/kubectl-ice

if you find this program useful please consider saying thanks I can be reached
 on twitter @nimblearchitect or you can buy me a coffee:
	https://nimblearchitect.github.io/kubectl-ice/donations/


if your just after the version string use: kubectl-ice -v

`, cmd.Parent().Version)
	return nil
}
