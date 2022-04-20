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
	allNamespaces bool     // should we search all namespaces
	container     string   // name of the container to search for
	filterList    []string // used to filter out rows form the table during Print function
	labels        string   // k8s pod labels
	showOddities  bool     // this isnt really common but it does sho up across 3+ commands and im lazy
	outputAs      string   // how to output the table, currently only accepts json
	sortList      []string //column names to sort on when table.Print() is called
}

func InitSubCommands(rootCmd *cobra.Command) {
	var odditiesShort string = "show only the outlier rows that dont fall within the computed range"
	KubernetesConfigFlags := genericclioptions.NewConfigFlags(false)

	//commands
	var cmdCommands = &cobra.Command{
		Use:     "command",
		Short:   commandsShort,
		Long:    fmt.Sprintf("%s\n\n%s", commandsShort, commandsDescription),
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
		Short:   resourceShort("cpu"),
		Long:    fmt.Sprintf("%s\n\n%s", resourceShort("cpu"), resourceDescription("cpu")),
		Example: fmt.Sprintf(resourceExample("cpu"), rootCmd.CommandPath()),
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
	cmdCPU.Flags().BoolP("oddities", "", false, odditiesShort)
	addCommonFlags(cmdCPU)
	rootCmd.AddCommand(cmdCPU)

	//environment
	var cmdEnvironment = &cobra.Command{
		Use:     "environment",
		Short:   environmentShort,
		Long:    fmt.Sprintf("%s\n\n%s", environmentShort, environmentDescription),
		Example: fmt.Sprintf(environmentExample, rootCmd.CommandPath()),
		Aliases: []string{"env", "vars"},
		// SuggestFor: []string{""},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := environment(cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdEnvironment.Flags())
	cmdEnvironment.Flags().BoolP("translate", "t", false, "read the configmap show its values")
	addCommonFlags(cmdEnvironment)
	rootCmd.AddCommand(cmdEnvironment)

	//ip
	var cmdIP = &cobra.Command{
		Use:     "ip",
		Short:   ipShort,
		Long:    fmt.Sprintf("%s\n\n%s", ipShort, ipDescription),
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
		Short:   imageShort,
		Long:    fmt.Sprintf("%s\n\n%s", imageShort, imageDescription),
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
		Short:   resourceShort("memory"),
		Long:    fmt.Sprintf("%s\n\n%s", resourceShort("memory"), resourceDescription("memory")),
		Example: fmt.Sprintf(resourceExample("memory"), rootCmd.CommandPath()),
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
	cmdMemory.Flags().BoolP("oddities", "", false, odditiesShort)
	addCommonFlags(cmdMemory)
	rootCmd.AddCommand(cmdMemory)

	//ports
	var cmdPorts = &cobra.Command{
		Use:     "ports",
		Short:   portsShort,
		Long:    fmt.Sprintf("%s\n\n%s", portsShort, portsDescription),
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
		Short:   probesShort,
		Long:    fmt.Sprintf("%s\n\n%s", probesShort, probesDescription),
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
		Short:   restartsShort,
		Long:    fmt.Sprintf("%s\n\n%s", restartsShort, restartsDescription),
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
	cmdRestart.Flags().BoolP("oddities", "", false, odditiesShort)
	addCommonFlags(cmdRestart)
	rootCmd.AddCommand(cmdRestart)

	//status
	var cmdStatus = &cobra.Command{
		Use:     "status",
		Short:   statusShort,
		Long:    fmt.Sprintf("%s\n\n%s", statusShort, statusDescription),
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
	cmdStatus.Flags().BoolP("oddities", "", false, odditiesShort)
	addCommonFlags(cmdStatus)
	rootCmd.AddCommand(cmdStatus)

	//volumes
	var cmdVolume = &cobra.Command{
		Use:     "volumes",
		Short:   volumesShort,
		Long:    fmt.Sprintf("%s\n\n%s", volumesShort, volumesDescription),
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
	cmdObj.Flags().StringP("match", "", "", `Filters out results, comma seperated list of COLUMN OP VALUE, where OP can be one of ==,<,>,<=,>= and != `)
}

func processCommonFlags(cmd *cobra.Command) (commonFlags, error) {
	var err error

	f := commonFlags{}

	if cmd.Flag("all-namespaces").Value.String() == "true" {
		f.allNamespaces = true
	}

	if cmd.Flag("oddities") != nil {
		if cmd.Flag("oddities").Value.String() == "true" {
			f.showOddities = true
		}
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
			f.sortList, err = splitAndFilterList(rawSortString, "ABCDEFGHIJKLMNOPQRSTUVWXYZ!%-")
			if err != nil {
				return commonFlags{}, err
			}
		}
	}

	if cmd.Flag("match") != nil {
		if len(cmd.Flag("match").Value.String()) > 0 {
			rawMatchString := cmd.Flag("match").Value.String()
			f.filterList, err = splitAndFilterList(rawMatchString, "ABCDEFGHIJKLMNOPQRSTUVWXYZ!%-0123456789<>=*?")
			if err != nil {
				return commonFlags{}, err
			}
		}
	}

	return f, nil
}

func splitAndFilterList(rawSortString string, filterString string) ([]string, error) {
	// based on a whitelist approach sort just removes invalid chars,
	// we cant check header names as we dont know them at this point
	var sortList []string
	var rawCase string

	rawSortList := strings.Split(rawSortString, ",")
	for i := 0; i < len(rawSortList); i++ {
		safeStr := ""
		rawItem := strings.TrimSpace(rawSortList[i])
		if len(rawItem) <= 0 {
			continue
		}

		// current used chars in headers are A-Z ! and % nothing else is needed
		// so pointless using regex
		rawCase = strings.ToUpper(rawItem)
		for _, v := range strings.Split(rawCase, "") {
			if strings.Contains(filterString, v) {
				safeStr += v
			}
		}

		if len(safeStr) != len(rawItem) {
			return []string{}, errors.New("invalid characters in column name")
		}
		sortList = append(sortList, safeStr)
	}

	return sortList, nil
}
