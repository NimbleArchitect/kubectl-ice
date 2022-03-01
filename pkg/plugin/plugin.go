package plugin

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func InitSubCommands(rootCmd *cobra.Command) {
	KubernetesConfigFlags := genericclioptions.NewConfigFlags(false)

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
}

func addCommonFlags(cmdObj *cobra.Command) {
	cmdObj.Flags().BoolP("all-namespaces", "A", false, "list containers form pods in all namespaces")
}
