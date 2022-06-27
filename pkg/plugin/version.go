package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var versionsShort = "Display container versions and mount points"

func Version(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	// 1234567890123456789012345678901234567890123456789012345678901234567890123456789
	fmt.Printf(`kubectl-ice kubernetes container viewer

version %s

the latest version can be found at: 
   https://github.com/NimbleArchitect/kubectl-ice/releases

to view the documentation or raise issues: 
   https://github.com/NimbleArchitect/kubectl-ice

if you find this program useful please consider saying thanks I can be reached
 on twitter @nimblearchitect or you can buy me a coffee:
   https://www.buymeacoffee.com/NimbleArchitect


for just the version string output use: kubectl-ice -v
`, cmd.Parent().Version)
	return nil
}
