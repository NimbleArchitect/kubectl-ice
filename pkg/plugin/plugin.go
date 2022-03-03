package plugin

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type commonFlags struct {
	allNamespaces bool
	labels        string
}

func InitSubCommands(rootCmd *cobra.Command) {
	KubernetesConfigFlags := genericclioptions.NewConfigFlags(false)

	//commands
	var cmdCommands = &cobra.Command{
		Use:     "command",
		Short:   "retrieves the command line and any arguments specified at the container level",
		Long:    "",
		Aliases: []string{"cmd", "exec", "args"},
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Commands(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdCommands.Flags())
	addCommonFlags(cmdCommands)
	rootCmd.AddCommand(cmdCommands)

	//cpu
	var cmdCPU = &cobra.Command{
		Use:   "cpu",
		Short: "return cpu requests size, limits and usage of each container",
		Long:  "",
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Resources(cmd, KubernetesConfigFlags, args, "cpu"); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdCPU.Flags())
	cmdCPU.Flags().BoolP("raw", "r", false, "show raw values")
	addCommonFlags(cmdCPU)
	rootCmd.AddCommand(cmdCPU)

	//ip
	var cmdIP = &cobra.Command{
		Use:   "ip",
		Short: "list ip addresses of all pods in the namespace listed",
		Long:  "",
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := IP(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdIP.Flags())
	addCommonFlags(cmdIP)
	rootCmd.AddCommand(cmdIP)

	//image
	var cmdImage = &cobra.Command{
		Use:     "image",
		Short:   "list the image name and pull status for each container",
		Long:    "",
		Aliases: []string{"im"},
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Image(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdImage.Flags())
	addCommonFlags(cmdImage)
	rootCmd.AddCommand(cmdImage)

	//memory
	var cmdMemory = &cobra.Command{
		Use:     "memory",
		Short:   "return memory requests size, limits and usage of each container",
		Long:    "",
		Aliases: []string{"mem"},
		// SuggestFor: []string{""},
		// Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Resources(cmd, KubernetesConfigFlags, args, "memory"); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdMemory.Flags())
	cmdMemory.Flags().BoolP("raw", "r", false, "show raw values")
	addCommonFlags(cmdMemory)
	rootCmd.AddCommand(cmdMemory)

	//ports
	var cmdPorts = &cobra.Command{
		Use:     "ports",
		Short:   "shows ports exposed by the containers in a pod",
		Long:    "",
		Aliases: []string{"port", "po"},
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Ports(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdPorts.Flags())
	addCommonFlags(cmdPorts)
	rootCmd.AddCommand(cmdPorts)

	//probes
	var cmdProbes = &cobra.Command{
		Use:     "probes",
		Short:   "shows details of configured startup, readiness and liveness probes of each container",
		Long:    "",
		Aliases: []string{"probe"},
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Probes(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdProbes.Flags())
	addCommonFlags(cmdProbes)
	rootCmd.AddCommand(cmdProbes)

	//restarts
	var cmdRestart = &cobra.Command{
		Use:     "restarts",
		Short:   "show restart counts for each container in a named pod",
		Long:    "",
		Aliases: []string{"restart"},
		// SuggestFor: []string{""},
		// Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Restarts(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdRestart.Flags())
	addCommonFlags(cmdRestart)
	rootCmd.AddCommand(cmdRestart)

	//status
	var cmdStatus = &cobra.Command{
		Use:     "status",
		Short:   "list status of each container in a pod",
		Long:    "",
		Aliases: []string{"st"},
		// SuggestFor: []string{""},
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Status(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdStatus.Flags())
	cmdStatus.Flags().BoolP("previous", "p", false, "show previous state")
	addCommonFlags(cmdStatus)
	rootCmd.AddCommand(cmdStatus)

	//volumes
	var cmdVolume = &cobra.Command{
		Use:     "volumes",
		Short:   "list all container volumes with mount points",
		Long:    "",
		Aliases: []string{"volume", "vol"},
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Volumes(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdVolume.Flags())
	addCommonFlags(cmdVolume)
	rootCmd.AddCommand(cmdVolume)

}

// adds common flags to the passed command
func addCommonFlags(cmdObj *cobra.Command) {
	cmdObj.Flags().BoolP("all-namespaces", "A", false, "list containers form pods in all namespaces")
	cmdObj.Flags().StringP("selector", "l", "", `: Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2`)

}

func processCommonFlags(cmd *cobra.Command) commonFlags {
	f := commonFlags{}

	if cmd.Flag("all-namespaces").Value.String() == "true" {
		f.allNamespaces = true
	}

	//fmt.Println(cmd.Flag("selector"))
	if cmd.Flag("selector") != nil {
		if len(cmd.Flag("selector").Value.String()) > 0 {
			f.labels = cmd.Flag("selector").Value.String()
		}
	}

	return f
}
