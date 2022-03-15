package plugin

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type commonFlags struct {
	allNamespaces bool
	labels        string
	container     string
	sortList      []string
	outputAs      string
}

func InitSubCommands(rootCmd *cobra.Command) {
	KubernetesConfigFlags := genericclioptions.NewConfigFlags(false)

	//commands
	var cmdCommands = &cobra.Command{
		Use:     "command",
		Short:   "retrieves the command line and any arguments specified at the container level",
		Long:    "",
		Example: fmt.Sprintf(commandsExample, rootCmd.CommandPath()),
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
		Use:     "cpu",
		Short:   "return cpu requests size, limits and usage of each container",
		Long:    "",
		Example: fmt.Sprintf(cpuExample, rootCmd.CommandPath()),
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
		Use:     "ip",
		Short:   "list ip addresses of all pods in the namespace listed",
		Long:    "",
		Example: fmt.Sprintf(ipExample, rootCmd.CommandPath()),
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
		Example: fmt.Sprintf(imageExample, rootCmd.CommandPath()),
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
		Example: fmt.Sprintf(memoryExample, rootCmd.CommandPath()),
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
		Example: fmt.Sprintf(portsExample, rootCmd.CommandPath()),
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
		Example: fmt.Sprintf(probesExample, rootCmd.CommandPath()),
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
		Example: fmt.Sprintf(restartsExample, rootCmd.CommandPath()),
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
		Example: fmt.Sprintf(statusExample, rootCmd.CommandPath()),
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
		Short:   "Display container volumes and mount points",
		Long:    "",
		Example: fmt.Sprintf(volumesExample, rootCmd.CommandPath()),
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
	cmdObj.Flags().StringP("selector", "l", "", `Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2`)
	cmdObj.Flags().StringP("container", "c", "", `Container name. If omitted show all containers in the pod`)
	cmdObj.Flags().StringP("sort", "", "", `Sort by column`)
	cmdObj.Flags().StringP("output", "o", "", `Output format, only json is supported`)

}

func processCommonFlags(cmd *cobra.Command) (commonFlags, error) {
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

	if cmd.Flag("container") != nil {
		if len(cmd.Flag("container").Value.String()) > 0 {
			f.container = cmd.Flag("container").Value.String()
		}
	}

	if cmd.Flag("output") != nil {
		if len(cmd.Flag("output").Value.String()) > 0 {
			outAs := cmd.Flag("output").Value.String()
			// we use a switch to match -o flag so I can expand in future
			switch strings.ToLower(outAs) {
			case "json":
				f.outputAs = "json"

			default:
				return commonFlags{}, errors.New("unknown output format only json is supported")
			}
		}
	}

	if cmd.Flag("sort") != nil {
		// based on a whitelist approach sort just removes invalid chars,
		// we cant check header names as we dont know them at this point
		if len(cmd.Flag("sort").Value.String()) > 0 {
			rawSortString := cmd.Flag("sort").Value.String()
			rawSortList := strings.Split(rawSortString, ",")
			for i := 0; i < len(rawSortList); i++ {
				safeStr := ""
				rawItem := strings.TrimSpace(rawSortList[i])
				if len(rawItem) <= 0 {
					continue
				}

				// current used chars in headers are A-Z ! and % nothing else is needed
				// so pointless using regex
				rawUpper := strings.ToUpper(rawItem)
				for _, v := range strings.Split(rawUpper, "") {
					if strings.Contains("ABCDEFGHIJKLMNOPQRSTUVWXYZ!%-", v) {
						safeStr += v
					}
				}

				if len(safeStr) != len(rawItem) {
					return commonFlags{}, errors.New("invalid characters in column name")
				}
				f.sortList = append(f.sortList, safeStr)
			}

		}
	}
	return f, nil
}
