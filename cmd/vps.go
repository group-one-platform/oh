package cmd

import (
	"fmt"
	"github.com/edvin/oh/api"
	"github.com/edvin/oh/cache"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

var vpsCmd = &cobra.Command{
	Use:   "vps",
	Short: "Commands to manipulate your Virtual Servers",
	Long:  `Configure and control your Virtual Server instances`,
}

var validateSingleVpsIdArg = func(cmd *cobra.Command, args []string) error {
	switch len(args) {
	case 0:
		return fmt.Errorf("you must specify the VPS Id")
	case 1:
		return nil
	default:
		return fmt.Errorf("only one positional argument expected (the VPS ID), got %d", len(args))
	}
}

func init() {
	rootCmd.AddCommand(vpsCmd)
}

func completeVpsIds(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// only complete the first positional (<vps-id>)
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	vpsList, err := cache.Call(cache.KeyCloudServers, cache.DefaultTTL, func() ([]api.CloudServer, error) {
		return api.ListCloudServers()
	})

	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var comps []string
	for _, v := range vpsList {
		id := strconv.Itoa(v.Id)
		if strings.HasPrefix(id, toComplete) {
			comps = append(comps, fmt.Sprintf("%s\t%s", id, v.Name))
		}
	}

	return comps, cobra.ShellCompDirectiveNoFileComp
}
